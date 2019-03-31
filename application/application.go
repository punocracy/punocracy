package application

import (
	"net/http"

	"github.com/carbocation/interpose"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"

	"github.com/alvarosness/punocracy/handlers"
	"github.com/alvarosness/punocracy/middlewares"
)

// New is the constructor for Application struct.
func New(config *viper.Viper) (*Application, error) {
	cookieStoreSecret := config.Get("cookie_secret").(string)

	app := &Application{}
	app.config = config
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, nil
}

// Application is the application object that runs HTTP server.
type Application struct {
	config       *viper.Viper
	sessionStore sessions.Store
}

func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetSessionStore(app.sessionStore))

	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *gorilla_mux.Router {
	router := gorilla_mux.NewRouter()

	router.Handle("/", http.HandlerFunc(handlers.GetHome)).Methods("GET")
	router.Handle("/login", http.HandlerFunc(handlers.GetLogin)).Methods("GET")
	router.Handle("/login", http.HandlerFunc(handlers.PostLogin)).Methods("POST")
	router.Handle("/register", http.HandlerFunc(handlers.GetRegister)).Methods("GET")
	router.Handle("/register", http.HandlerFunc(handlers.PostRegister)).Methods("POST")
	router.Handle("/now", http.HandlerFunc(handlers.GetNow)).Methods("GET")
	router.Handle("/now", http.HandlerFunc(handlers.PostNow)).Methods("POST")

	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return router
}
