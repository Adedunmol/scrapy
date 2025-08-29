package api

import (
	"context"
	"encoding/json"
	"github.com/Adedunmol/scrapy/scrapy"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type JobRequest struct {
	//Email      string `json:"email"`
	SearchTerm string `json:"search_term"`
	Location   string `json:"location"`
}

func WriteJSONResponse(responseWriter http.ResponseWriter, data interface{}, statusCode int) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	if err := json.NewEncoder(responseWriter).Encode(data); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

func FetchJobsHandler(responseWriter http.ResponseWriter, request *http.Request) {

	var body JobRequest

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		return
	}

	if body.SearchTerm == "" {
		response := Response{
			Status:  "error",
			Message: "search term is missing",
		}

		WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
		return
	}

	// call the coordinator function to get all the jobs
	jobs := scrapy.Coordinator(context.Background(), false, body.SearchTerm, body.Location)

	response := Response{
		Status:  "success",
		Message: "scraped jobs gotten successfully",
		Data:    map[string]interface{}{"jobs": jobs},
	}

	WriteJSONResponse(responseWriter, response, http.StatusOK)
	return
}
