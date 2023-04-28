//go:generate specgen-golang service-go --server {{server.value}} --spec-file spec.yaml --module-name the-service --generate-path ./spec --services-path ./services --swagger-path docs/swagger.yaml

package main

import (
	"flag"
	"net/http"
	"the-service/services"
	"the-service/services/v2"
	"the-service/spec"
	{{#server.vestigo}}
	"github.com/husobee/vestigo"
	{{/server.vestigo}}
	{{#server.httprouter}}
	"github.com/julienschmidt/httprouter"
	{{/server.httprouter}}
	{{#server.chi}}
	"github.com/go-chi/chi/v5"
	{{#cors.value}}
	"github.com/go-chi/cors"
	{{/cors.value}}
	{{/server.chi}}
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func main() {
	port := flag.String("port", "8081", "port number")
	flag.Parse()

	decimal.MarshalJSONWithoutQuotes = true

	{{#server.vestigo}}
	router := vestigo.NewRouter()
	{{/server.vestigo}}
	{{#server.httprouter}}
	router := httprouter.New()
	{{/server.httprouter}}
	{{#server.chi}}
	router := chi.NewRouter()
	{{/server.chi}}

	{{#cors.value}}
	{{#server.vestigo}}
	router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin: []string{"*", "*"},
	})
	{{/server.vestigo}}
	{{#server.httprouter}}
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Set("Access-Control-Allow-Origin", "*")
		}
	})
	{{/server.httprouter}}
	{{#server.chi}}
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*", "*"},
	}))
	{{/server.chi}}
	{{/cors.value}}

	echoServiceV2 := &v2.EchoService{}
	echoService := &services.EchoService{}
	checkService := &services.CheckService{}

	spec.AddRoutes(router, echoServiceV2, echoService, checkService)

	{{#server.vestigo}}
	router.Get("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))).ServeHTTP)
	{{/server.vestigo}}
	{{#server.httprouter}}
	router.ServeFiles("/docs/*filepath", http.Dir("./docs"))
	{{/server.httprouter}}
	{{#server.chi}}
	router.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))
	{{/server.chi}}

	log.Infof("Starting service on port: %s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, router))
}
