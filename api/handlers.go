package api

type JobRequest struct {
	//Email      string `json:"email"`
	SearchTerm string `json:"search_term"`
	Location   string `json:"location"`
}

//func FetchJobsHandler(responseWriter http.ResponseWriter, request *http.Request) {
//	var body JobRequest
//
//	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
//		response := helpers.Response{
//			Status:  "error",
//			Message: "error decoding body",
//		}
//
//		helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
//		return
//	}
//
//	if body.SearchTerm == "" {
//		response := helpers.Response{
//			Status:  "error",
//			Message: "search term is missing",
//		}
//
//		helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
//		return
//	}
//
//	// call the coordinator function to get all the jobs
//	jobs := scrapy.Coordinator(context.Background(), false, body.SearchTerm, body.Location)
//
//	response := helpers.Response{
//		Status:  "success",
//		Message: "scraped jobs gotten successfully",
//		Data:    map[string]interface{}{"jobs": jobs},
//	}
//
//	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
//	return
//}
