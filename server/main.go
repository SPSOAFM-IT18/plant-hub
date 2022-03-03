package main

import (
	"log"
	"net/http"

	sens "github.com/SPSOAFM-IT18/dmp-plant-hub/sensors"

	"github.com/SPSOAFM-IT18/dmp-plant-hub/env"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	plg "github.com/99designs/gqlgen/graphql/playground"
	"github.com/SPSOAFM-IT18/dmp-plant-hub/database"
	"github.com/SPSOAFM-IT18/dmp-plant-hub/graph"
	"github.com/SPSOAFM-IT18/dmp-plant-hub/graph/generated"
	"github.com/SPSOAFM-IT18/dmp-plant-hub/rest/router"
	seq "github.com/SPSOAFM-IT18/dmp-plant-hub/sequences"
	"github.com/go-chi/chi"
	webs "github.com/gorilla/websocket"
)

func main() {
	cMoist := make(chan float64)
	cTemp := make(chan float64)
	cHum := make(chan float64)
	cPumpState := make(chan bool)

	sensei := sens.Init()

	var db = database.Connect()

	go seq.MeasurementSequence(sensei, cMoist, cTemp, cHum, cPumpState)
	go seq.SaveOnFourHoursPeriod(db, cMoist, cTemp, cHum)
	go seq.Controller(db, sensei, cMoist, cPumpState)

	gqlRouter := chi.NewRouter()
	restRouter := router.Router()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	// router.Use(cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://4.2.0.126:3000", "http://localhost:3000", "http://4.2.0.126", "http://4.2.0.225:5000"},
	// 	AllowCredentials: true,
	// 	Debug:            true,
	// }).Handler)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: webs.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return r.Host == "http://4.2.0.126:3000"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	gqlRouter.Handle("/graphql", plg.Handler("GraphQL playground", "/query"))
	gqlRouter.Handle("/query", srv)
	/*http.Handle("/graphql", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)*/

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", env.Process("GO_GQL_API_PORT"))
	go log.Fatal(http.ListenAndServe(":"+env.Process("GO_GQL_API_PORT"), gqlRouter))
	log.Fatal(http.ListenAndServe(":"+env.Process("GO_REST_API_PORT"), restRouter))

}
