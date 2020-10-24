package apigateway

import (
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	api "github.com/contact-tracker/apiService/pkg/api/http"
)

func NewServer(port, checkInsPort, placesPort, usersPort string) (server *api.Server) {
	server = api.NewServer(port)
	r := server.Router
	r.Mount("/check-ins", newProxyRouter("http://localhost:"+checkInsPort, true))
	r.Mount("/places", newProxyRouter("http://localhost:"+placesPort, true))
	r.Mount("/users", newProxyRouter("http://localhost:"+usersPort, true))

	return
}

func newProxyRouter(targetHost string, useCors bool) http.Handler {
	targetUrl, err := url.Parse(targetHost)
	if err != nil {
		log.Fatalln(err)
	}
	if useCors {
		r := chi.NewRouter()
		corsHandler := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}).Handler
		r.Use(corsHandler)
		r.Mount("/", api.ProxyReq(targetUrl))
		return r
	}
	return api.ProxyReq(targetUrl)
}
