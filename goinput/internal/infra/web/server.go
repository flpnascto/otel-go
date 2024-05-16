package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type MicroserviceClimateResponse struct {
	City  string  `json:"city"`
	TempC float32 `json:"temp_C"`
	TempF float32 `json:"temp_F"`
	TempK float32 `json:"temp_K"`
}

type TemplateData struct {
	RequestNameOTEL string
	ExternalCallURL string
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
	router.Post("/input", we.HandleRequest)
	return router
}

func (h *Webserver) HandleRequest(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, spanHandler := h.TemplateData.OTELTracer.Start(ctx, "HANDLER REQUEST "+h.TemplateData.RequestNameOTEL)
	defer spanHandler.End()

	var cepBody struct {
		Cep string `json:"cep"`
	}
	err := json.NewDecoder(r.Body).Decode(&cepBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	var req *http.Request

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, h.TemplateData.ExternalCallURL+"/"+cepBody.Cep, nil)
	if err != nil {
		http.Error(w, "Invalid ExternalCallMethod", http.StatusBadRequest)
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var response MicroserviceClimateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		http.Error(w, "Erro 2 "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Erro 3 "+err.Error(), http.StatusInternalServerError)
		return
	}

}
