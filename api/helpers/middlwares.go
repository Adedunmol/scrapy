package helpers

import (
	"context"
	"net/http"
	"strings"
)

var Unauthorized = Response{
	Status:  "error",
	Message: http.StatusText(http.StatusUnauthorized),
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {

		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			WriteJSONResponse(responseWriter, Unauthorized, http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authHeader, " ")

		if len(tokenString) != 2 {
			WriteJSONResponse(responseWriter, Unauthorized, http.StatusUnauthorized)
			return
		}

		data, err := DecodeToken(tokenString[1])
		if err != nil {
			WriteJSONResponse(responseWriter, Unauthorized, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(request.Context(), "email", data.Email)
		ctx = context.WithValue(ctx, "user_id", data.ID)

		newRequest := request.WithContext(ctx)
		next.ServeHTTP(responseWriter, newRequest)
	})
}
