package main

import (
	"net/http"

	l "VK_app/server/dbconn"
	h "VK_app/server/handlers"

	logger "VK_app/server/logger"
	middle "VK_app/server/middlewear"

	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/files" // swagger embed files
	_ "VK_app/docs"
	gin "github.com/gin-gonic/gin"
)

// @title VK Film library
// @version 1.0
// @description VK_app_film_library project

// @host localhost:8080
// @BasePath /filmlibrary

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @security ApiKeyAuth

// @securityDefinitions.apikey AdminKeyAuth
// @in header
// @name Authorization
// @security AdminKeyAuth
func init() {
	l.Db = l.Connection()
}

// main is the entry point of the program.
//
// It initializes the database connection, sets up the logger, and configures the HTTP router.
// The main function then starts the HTTP server and listens for incoming requests.
// It returns an error if there is an issue starting the server.
func main() {
	defer l.Db.Close()
	logger.LogFile = logger.LoggerInit()
	defer logger.LogFile.Close()
	loggerOut := log.New(logger.LogFile, "", log.LstdFlags)
	log.SetOutput(logger.LogFile)
	r := chi.NewRouter()
	swaggerRouter:=gin.Default()
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: loggerOut}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/filmlibrary/registration", h.RegisterUser)
	r.Post("/filmlibrary/login", h.Login)
	r.Group(func(rUser chi.Router) {
		rUser.Use(middle.CheckToken)
		rUser.Post("/filmlibrary/filmssorted", h.GetSortedFilms)
		rUser.Post("/filmlibrary/filmspiece", h.GetFilmByPiece)
		rUser.Get("/filmlibrary/actors", h.GetAllActors)
	})
	r.Group(func(rAdm chi.Router) {
		rAdm.Use(middle.CheckTokenAdmin)
		rAdm.Delete("/filmlibrary/admin/film", h.DeleteFilm)
		rAdm.Put("/filmlibrary/admin/film", h.UpdateFilm)
		rAdm.Delete("/filmlibrary/admin/actor", h.DeleteActor)
		rAdm.Put("/filmlibrary/admin/actor", h.UpdateActor)
		rAdm.Post("/filmlibrary/admin/actors", h.PostActor)
		rAdm.Post("/filmlibrary/admin/films", h.PostFilm)
		rAdm.Post("/filmlibrary/admin/filmssorted", h.GetSortedFilmsAdmin)
		rAdm.Post("/filmlibrary/admin/filmspiece", h.GetFilmByPieceAdmin)
		rAdm.Get("/filmlibrary/admin/actors", h.GetAllActorsAdmin)
	})
	swaggerRouter.GET("filmlibrary/swagger", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Printf("Failed to start server: %s\n", err.Error())
		return
	}
}
