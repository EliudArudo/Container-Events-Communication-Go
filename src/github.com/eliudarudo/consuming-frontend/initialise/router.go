package initialise

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eliudarudo/consuming-frontend/controllers"
	"github.com/eliudarudo/consuming-frontend/env"
	"github.com/eliudarudo/consuming-frontend/logs"
	"github.com/gorilla/mux"
)

var routerFilename = "initialise/router.go"

// App defines the Mux router object
type App struct {
	Router *mux.Router
}

func setUpRouter() {
	portString := fmt.Sprintf(":%d", env.Port)

	app := &App{}
	app.initialize()

	logMessage := fmt.Sprintf("Starting server on port: %v\n", env.Port)
	logs.StatusFileMessageLogging("SUCCESS", routerFilename, "initialiseRouter", logMessage)

	app.run(portString)
}

func (a *App) initialize() {
	a.Router = mux.NewRouter()
	a.setRouters()
	a.Router.NotFoundHandler = a.handleRequest(controllers.RouterHandler404)
}

func (a *App) setRouters() {
	a.get("/", a.handleRequest(controllers.IndexController))

	a.post("/task", a.handleRequest(controllers.RequestRouteController))
}

func (a *App) get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type requestHandlerFunction func(w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler requestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}
