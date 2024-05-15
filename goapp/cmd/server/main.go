package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/flpnascto/otel-go/goapp/configs"
	"github.com/flpnascto/otel-go/goapp/internal/entity"
	"github.com/flpnascto/otel-go/goapp/internal/infra/api"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		cepQuery := parts[len(parts)-1]
		cep, err := entity.NewCep(cepQuery)
		if err != nil {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		city, err := api.FetchCepApi(cep)
		if err != nil {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}

		temp, err := api.FetchWeatherApi(*city, configs.WeatherApiKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]any{
			"city":   city,
			"temp_C": temp.TempC,
			"temp_F": temp.TempF,
			"temp_K": temp.TempK,
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}

		w.Write(responseBytes)
	})

	http.ListenAndServe(":8080", mux)
}
