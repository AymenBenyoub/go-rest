package server

import (
	
	"net/http"
	"rest/internal/handlers"
)

func NewRouter( userHandler handlers.UserHandler, postHandler handlers.PostHandler) http.Handler {

	router := http.NewServeMux()

	userHandler.RegisterRoutes(router)
	postHandler.RegisterRoutes(router)

	return router
}
