package application

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/carbocation/interpose"
	_ "github.com/go-sql-driver/mysql"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	"github.com/alvarosness/punocracy/handlers"
	"github.com/alvarosness/punocracy/middlewares"
)

// New is the constructor for Application struct.
func New(config *viper.Viper) (*Application, error) {
	dsn := config.Get("dsn").(string)
	urlString := config.Get("mongoURL").(string)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(urlString))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Connect
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Check connection with ping
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	mongodb := client.Database("punocracy")

	cookieStoreSecret := config.Get("cookie_secret").(string)

	app := &Application{}
	app.config = config
	app.dsn = dsn
	app.db = db
	app.mongodb = mongodb
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, nil
}

// Application is the application object that runs HTTP server.
type Application struct {
	config       *viper.Viper
	dsn          string
	db           *sqlx.DB
	mongodb      *mongo.Database
	sessionStore sessions.Store
}

func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetDB(app.db))
	middle.Use(middlewares.SetMongo(app.mongodb))
	middle.Use(middlewares.SetSessionStore(app.sessionStore))
	middle.Use(middlewares.Logging())

	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *gorilla_mux.Router {
	MustLogin := middlewares.MustLogin

	router := gorilla_mux.NewRouter()

	router.Host("gopunned.me")
	router.NotFoundHandler = http.HandlerFunc(handlers.HandleNotFound)

	router.Handle("/now", http.HandlerFunc(handlers.GetHome)).Methods("GET")
	router.Handle("/now", http.HandlerFunc(handlers.PostHome)).Methods("POST")

	router.HandleFunc("/", handlers.HandleRoot).Methods("GET", "POST", "PUT", "DELETE")

	router.HandleFunc("/submit", handlers.GetSubmit).Methods("GET")
	router.HandleFunc("/submit", handlers.PostSubmit).Methods("POST")

	router.HandleFunc("/history", handlers.GetHistory).Methods("GET")
	router.HandleFunc("/history", handlers.PostHistory).Methods("POST")

	router.HandleFunc("/words/{letter}", handlers.GetWords).Methods("GET")

	router.HandleFunc("/queuerater", handlers.GetCurator).Methods("GET")
	router.HandleFunc("/queuerater", handlers.PostCurator).Methods("POST")

	router.HandleFunc("/about", handlers.GetAbout).Methods("GET")

	router.HandleFunc("/signup", handlers.GetSignup).Methods("GET")
	router.HandleFunc("/signup", handlers.PostSignup).Methods("POST")

	router.HandleFunc("/login", handlers.GetLogin).Methods("GET")
	router.HandleFunc("/login", handlers.PostLogin).Methods("POST")

	router.HandleFunc("/logout", handlers.GetLogout).Methods("GET")

	router.Handle("/users/{userID:[0-9]+}", MustLogin(http.HandlerFunc(handlers.PostPutDeleteUsersID))).Methods("POST", "PUT", "DELETE")

	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return router
}
