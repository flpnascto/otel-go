package web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/flpnascto/otel-go/goapp/internal/entity"
	"github.com/flpnascto/otel-go/goapp/internal/infra/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type TemplateData struct {
	RequestNameOTEL string
	OTELTracer      trace.Tracer
}

type Webserver struct {
	TemplateData *TemplateData
}

func NewServer(templateData *TemplateData) *Webserver {
	return &Webserver{
		TemplateData: templateData,
	}
}

func (we *Webserver) CreateServer() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))
	// promhttp
	router.Handle("/metrics", promhttp.Handler())
	router.Get("/{cep}", we.HandleRequest)
	return router
}

func (h *Webserver) HandleRequest(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, spanHandler := h.TemplateData.OTELTracer.Start(ctx, "HANDLER REQUEST "+h.TemplateData.RequestNameOTEL)
	defer spanHandler.End()

	cepQuery := chi.URLParam(r, "cep")
	cep, err := entity.NewCep(cepQuery)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	ctx, spanCepApi := h.TemplateData.OTELTracer.Start(ctx, "FETCH CEP API "+h.TemplateData.RequestNameOTEL)
	city, err := api.FetchCepApi(cep)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	spanCepApi.End()

	_, spanWeatherApi := h.TemplateData.OTELTracer.Start(ctx, "FETCH WEATHER API "+h.TemplateData.RequestNameOTEL)
	temp, err := api.FetchWeatherApi(*city, viper.GetString("WEATHER_API_KEY"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	spanWeatherApi.End()

	response := map[string]any{
		"city":   city,
		"temp_C": temp.TempC,
		"temp_F": temp.TempF,
		"temp_K": temp.TempK,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
