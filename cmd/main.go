package main

import (
	"net/http"

	l "VK_app/internal/dbconn"
	h "VK_app/pkg/handlers"

	logger "VK_app/pkg/logger"

	_ "VK_app/api/docs"
	middle "VK_app/pkg/middleware"
	"log"

	gin "github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

//swag init --parseDependency --parseInternal -d ./,./handlers && go run main.go - to start

// @title VK Film library
// @version 1.0
// @description VK_app_film_library project
// @host localhost:8080
// @BasePath /
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
	log.SetOutput(logger.LogFile)
	swaggerRouter := gin.Default()
	swaggerRouter.Use(gin.Logger())
	swaggerRouter.Use(gin.Logger())
	swaggerRouter.Use(gin.Recovery())

	swaggerRouter.POST("/filmlibrary/registration", h.RegisterUser)
	swaggerRouter.POST("/filmlibrary/login", h.Login)

	UserGroup := swaggerRouter.Group("/filmlibrary")
	UserGroup.Use(middle.CheckToken)
	UserGroup.POST("/filmssorted", h.GetSortedFilms)
	UserGroup.POST("/filmspiece", h.GetFilmByPiece)
	UserGroup.GET("/actors", h.GetAllActors)

	AdminGroup := swaggerRouter.Group("/filmlibrary/admin")
	AdminGroup.Use(middle.CheckTokenAdmin)
	AdminGroup.DELETE("/film", h.DeleteFilm)
	AdminGroup.PUT("/film", h.UpdateFilm)
	AdminGroup.DELETE("/actor", h.DeleteActor)
	AdminGroup.PUT("/actor", h.UpdateActor)
	AdminGroup.POST("/actors", h.PostActor)
	AdminGroup.POST("/films", h.PostFilm)
	AdminGroup.POST("/filmssorted", h.GetSortedFilmsAdmin)
	AdminGroup.POST("/filmspiece", h.GetFilmByPieceAdmin)
	AdminGroup.GET("/actors", h.GetAllActorsAdmin)

	swaggerRouter.GET("/docs", func(c *gin.Context) { c.Redirect(http.StatusFound, "swagger/index.html") })
	swaggerRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := http.ListenAndServe(":8080", swaggerRouter); err != nil {
		log.Printf("Failed to start server: %s\n", err.Error())
		return
	}
}
