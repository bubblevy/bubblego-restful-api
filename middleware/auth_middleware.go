package middleware

import (
	"bubblevy/restful-api/helper"
	"bubblevy/restful-api/model/web"
	"net/http"
)

type authMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *authMiddleware {
	return &authMiddleware{Handler: handler}
}

func (middleware *authMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("API-Key") == "BUBBLEKEY" {
		// ok
		middleware.Handler.ServeHTTP(writer, request)
	} else {
		//error api key
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		webResponse := web.WebResponse{
			Code:    http.StatusUnauthorized,
			Error:   true,
			Message: "Invalid API key. Please provide a valid key.",
		}

		helper.WriteToResponseBody(writer, webResponse)
	}
}
