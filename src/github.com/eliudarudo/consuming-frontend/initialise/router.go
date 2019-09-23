package initialise

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eliudarudo/consuming-frontend/controllers"
	"github.com/gorilla/mux"
)

// App to create Router Instance
type App struct {
	Router *mux.Router
}

func initialiseRouter() {
	port := 8080
	portString := fmt.Sprintf(":%d", port)

	app := &App{}
	app.initialize()
	fmt.Printf("Starting server on port: %v\n", port)
	app.Run(portString)
}

func (a *App) initialize() {
	a.Router = mux.NewRouter()
	a.setRouters()
	a.Router.NotFoundHandler = a.handleRequest(controllers.RespondError)
}

// setRouters sets all the required routers
func (a *App) setRouters() {
	a.Get("/", a.handleRequest(controllers.HTTPIndex))

	a.Post("/task", a.handleRequest(controllers.HTTPPostTaskHandler))
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

// RequestHandlerFunction does something
type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}
