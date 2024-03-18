package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	l "VK_app/internal/dbconn"
	st "VK_app/internal/structures"
	"VK_app/pkg/postgresql"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary Login
// @Tags auth
// @Description Login of a user or admin.
// @ID login
// @Accept json
// @Produce json
// @Param input body st.User true "login"
// @Success 200 {object} st.StatusOKMessage "user was successfully logged in"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Router /filmlibrary/login [post]
func Login(c *gin.Context) {
	var user st.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong data format or missing data"})
		return
	}

	var storedPassword string
	row := l.Db.QueryRow("SELECT password FROM filmsactors.users WHERE login = $1", user.Login)
	err := row.Scan(&storedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login is wrong"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password is wrong"})
		return
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":          user.Login,
		"hashedpassword": user.Password,
		"role":           user.Role,
		"exp":            time.Now().Add(time.Hour * 24).Unix(),
	})

	SignedToken, err := jwtToken.SignedString(st.Secret)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	c.Header("Authorization", SignedToken)
	c.JSON(http.StatusOK, gin.H{"message": "login was completed successfully"})
}

// RegisterUser godoc
// @Summary Register
// @Tags auth
// @Description Registration of a new user or admin.
// @ID register
// @Accept json
// @Produce json
// @Param input body st.User true "register"
// @Success 200 {object} st.StatusOKMessage "user was successfully registered"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Router /filmlibrary/registration [post]
func RegisterUser(c *gin.Context) {
	var user st.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong data format or missing data"})
		return
	}
	err := postgresql.AddUser(&user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't add user to database or user already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

// PostFilm godoc
// @Summary AddFilm
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Add a new film to the database.
// @ID add-film
// @Accept json
// @Produce json
// @Param input body st.Film true "Film object for adding"
// @Success 201 {object} st.StatusOKMessage "film was successfully added"
// @Failure 500 {object} st.StatusInternalServerErrorMessage "internal server error"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Router /filmlibrary/admin/films [post]
func PostFilm(c *gin.Context) {
	if _, ok := c.Get("isAuthorizedAdmin"); !ok {
		log.Println("Unauthorized admin")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized admin access denied"})
		return
	}
	var film st.Film
	if err := c.ShouldBindJSON(&film); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := postgresql.AddFilm(film)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

// PostActor godoc
// @Summary AddActor
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Add a new actor to the database.
// @ID add-actor
// @Accept json
// @Produce json
// @Param input body st.Actor true "Actor object for adding"
// @Success 201 {object} st.StatusOKMessage "actor was successfully added"
// @Failure 500 {object} st.StatusInternalServerErrorMessage "internal server error"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Failure 401 {object} st.StatusUnauthorizedMessage "unauthorized"
// @Router /filmlibrary/admin/actors [post]
func PostActor(c *gin.Context) {
	if _, ok := c.Get("isAuthorizedAdmin"); !ok {
		log.Println("Unauthorized admin")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized admin access denied"})
		return
	}
	var actor st.Actor
	if err := c.ShouldBindJSON(&actor); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := postgresql.AddActor(actor)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

// UpdateActor godoc
// @Summary UpdateActor
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Update an actor in the database.
// @ID update-actor
// @Accept json
// @Produce json
// @Param input body st.Actor true "Actor object for updating"
// @Success 200 {object} st.StatusOKMessage "actor was successfully updated"
// @Failure 500 {object} st.StatusInternalServerErrorMessage "internal server error"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Failure 401 {object} st.StatusUnauthorizedMessage "unauthorized"
// @Router /filmlibrary/admin/actor [put]
func UpdateActor(c *gin.Context) {
	if _, ok := c.Get("isAuthorizedAdmin"); !ok {
		log.Println("Unauthorized admin")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized admin access denied"})
		return
	}
	var actor st.Actor
	if err := c.ShouldBindJSON(&actor); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := postgresql.UpdateActor(actor)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// DeleteActor godoc
// @Summary DeleteActor
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Delete an actor from the database.
// @ID delete-actor
// @Accept json
// @Produce json
// @Param input body st.Actor true "Actor object for deleting"
// @Success 200 {object} st.StatusOKMessage "actor was successfully deleted"
// @Failure 500 {object} st.StatusInternalServerErrorMessage "internal server error"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Failure 401 {object} st.StatusUnauthorizedMessage "unauthorized"
// @Router /filmlibrary/admin/actor [delete]
func DeleteActor(c *gin.Context) {
	if _, ok := c.Get("isAuthorizedAdmin"); !ok {
		log.Println("Unauthorized admin")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized admin access denied"})
		return
	}
	var actor st.Actor
	if err := c.ShouldBindJSON(&actor); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := postgresql.DelActor(actor.Id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})

}

// UpdateFilm godoc
// @Summary UpdateFilm
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Update a film in the database.
// @ID update-film
// @Accept json
// @Produce json
// @Param input body st.Film true "Film object for updating"
// @Success 200 {object} st.StatusOKMessage "film was successfully updated"
// @Failure 500 {object} st.StatusInternalServerErrorMessage "internal server error"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Failure 401 {object} st.StatusUnauthorizedMessage "unauthorized"
// @Router /filmlibrary/admin/film [put]
func UpdateFilm(c *gin.Context) {
	if _, ok := c.Get("isAuthorizedAdmin"); !ok {
		log.Println("Unauthorized admin")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized admin access denied"})
		return
	}
	var film st.Film
	if err := c.ShouldBindJSON(&film); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := postgresql.UpdateFilm(film)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// DeleteFilm godoc
// @Summary DeleteFilm
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Delete a film from the database.
// @ID delete-film
// @Accept json
// @Produce json
// @Param input body st.Film true "Film object for deleting"
// @Success 200 {object} st.StatusOKMessage
// @Failure 500 {object} st.StatusInternalServerErrorMessage
// @Failure 400 {object} st.StatusBadRequestMessage
// @Failure 404 {object} st.StatusNotFoundMessage
// @Failure 401 {object} st.StatusUnauthorizedMessage
// @Router /filmlibrary/admin/film [delete]
func DeleteFilm(c *gin.Context) {
	if _, ok := c.Get("isAuthorizedAdmin"); !ok {
		log.Println("Unauthorized admin")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized admin access denied"})
		return
	}
	var film st.Film
	if err := c.ShouldBindJSON(&film); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := postgresql.DelFilm(film.Id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetSortedFilms godoc
// @Summary GetSortedFilms
// @Security ApiKeyAuth
// @Tags User Functions
// @Description Get films sorted by rating, name, or date.
// @ID get-sorted-films
// @Accept json
// @Produce json
// @Param input body st.KeySort true "key: rating, name, or date to be sorted by"
// @Success 200 {object} st.StatusOKMessage "ok"
// @Failure 500 {object} st.StatusInternalServerErrorMessage "internal server error"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Failure 401 {object} st.StatusUnauthorizedMessage "unauthorized"
// @Router /filmlibrary/filmssorted [post]
func GetSortedFilms(c *gin.Context) {
	if _, ok := c.Get("isAuthorizedUser"); !ok {
		log.Println("Unauthorized user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user access denied"})
		return
	}
	var films []st.Film
	var sortKey st.KeySort
	if err := c.ShouldBindJSON(&sortKey); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sortKey.Key != "key" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong key"})
		return
	}
	if sortKey.Value == "" || sortKey.Value == "rating" {
		films = postgresql.GetFilmsSortedRating()
	} else if sortKey.Value == "name" {
		films = postgresql.GetFilmsSortedName()
	} else {
		films = postgresql.GetFilmsSortedDate()
	}
	if films == nil {
		log.Println("No films found")
		c.JSON(http.StatusNotFound, gin.H{"error": "No films found"})
		return
	}
	c.JSON(http.StatusOK, films)
}

// GetFilmByPiece godoc
// @Summary GetFilmByPiece
// @Security ApiKeyAuth
// @Tags User Functions
// @Description Get films based on a JSON fragment, which contains a piece of the film name or actor name.
// @ID get-film-by-piece
// @Accept json
// @Produce json
// @Param input body st.JSONFragment true "JSON fragment with a piece of the film name or actor name to search for"
// @Success 200 {object} st.StatusOKMessage "ok"
// @Failure 500 {object} st.StatusInternalServerErrorMessage "internal server error"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Failure 401 {object} st.StatusUnauthorizedMessage "unauthorized"
// @Router /filmlibrary/filmspiece [post]
func GetFilmByPiece(c *gin.Context) {
	if _, ok := c.Get("isAuthorizedUser"); !ok {
		log.Println("Unauthorized user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user access denied"})
		return
	}
	var JSONInput st.JSONFragment
	var films []st.Film
	var buf bytes.Buffer
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &JSONInput); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if JSONInput.Key == "actor" {
		films = postgresql.GetFilmsPieceActor(JSONInput.Fragment)
	} else {
		films = postgresql.GetFilmsPieceFilm(JSONInput.Fragment)
	}
	if films == nil {
		log.Println("No films found")
		c.JSON(http.StatusNotFound, gin.H{"error": "No films found"})
		return
	}
	c.JSON(http.StatusOK, films)
}

// GetAllActors godoc
// @Summary GetAllActors
// @Security ApiKeyAuth
// @Tags User Functions
// @Description Get all actors
// @ID get-all-actors
// @Accept json
// @Produce json
// @Success 200 {object} st.StatusOKMessage "ok"
// @Failure 500 {object} st.StatusInternalServerErrorMessage "internal server error"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Failure 401 {object} st.StatusUnauthorizedMessage "unauthorized"
// @Router /filmlibrary/actors [get]
func GetAllActors(c *gin.Context) {
	if _, ok := c.Get("isAuthorizedUser"); !ok {
		log.Println("Unauthorized user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user access denied"})
		return
	}
	actors := postgresql.GetFilmsActor()
	c.JSON(http.StatusOK, actors)
}

// GetSortedFilmsAdmin godoc
// @Summary GetSortedFilmsAdmin
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Get films sorted by rating, name, or date.
// @ID get-sorted-films-admin
// @Accept json
// @Produce json
// @Param input body st.KeySort true "key: rating, name, or date to be sorted by"
// @Success 200 {object} st.StatusOKMessage "ok"
// @Failure 500 {object} st.StatusInternalServerErrorMessage "internal server error"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Failure 401 {object} st.StatusUnauthorizedMessage "unauthorized"
// @Router /filmlibrary/admin/filmssorted [post]
func GetSortedFilmsAdmin(c *gin.Context) {
	if _, ok := c.Get("isAuthorizedAdmin"); !ok {
		log.Println("Unauthorized admin")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized admin access denied"})
		return
	}
	var films []st.Film
	var buf bytes.Buffer
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	piece := buf.String()
	if piece == "" || piece == "rating" {
		films = postgresql.GetFilmsSortedRating()
	} else if piece == "name" {
		films = postgresql.GetFilmsSortedName()
	} else {
		films = postgresql.GetFilmsSortedDate()
	}
	if films == nil {
		log.Println("No films found")
		c.JSON(http.StatusNotFound, st.StatusNotFoundMessage{Message: "No films found"}) // gin.H{"error": "No films found"})
		return
	}
	c.JSON(http.StatusOK, films)
}

// GetFilmByPieceAdmin godoc
// @Summary GetFilmByPieceAdmin
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Get films based on a JSON fragment, which contains a piece of the film name or actor name.
// @ID get-film-by-piece-admin
// @Accept json
// @Produce json
// @Param input body st.JSONFragment true "JSON fragment containing a piece of the film name or actor name to be searched for"
// @Success 200 {object} st.StatusOKMessage "ok"
// @Failure 500 {object} st.StatusInternalServerErrorMessage "internal server error"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Failure 401 {object} st.StatusUnauthorizedMessage "unauthorized"
// @Router /filmlibrary/admin/filmspiece [post]
func GetFilmByPieceAdmin(c *gin.Context) {
	if _, ok := c.Get("isAuthorizedAdmin"); !ok {
		log.Println("Unauthorized admin")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized admin access denied"})
		return
	}
	var JSONInput st.JSONFragment
	var films []st.Film
	var buf bytes.Buffer
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &JSONInput); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if JSONInput.Key == "actor" {
		films = postgresql.GetFilmsPieceActor(JSONInput.Fragment)
	} else {
		films = postgresql.GetFilmsPieceFilm(JSONInput.Fragment)
	}
	if films == nil {
		log.Println("No films found")
		c.JSON(http.StatusNotFound, gin.H{"error": "No films found"})
		return
	}
	c.JSON(http.StatusOK, films)
}

// GetAllActorsAdmin godoc
// @Summary GetAllActors
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Get all actors
// @ID get-all-actors-admin
// @Accept json
// @Produce json
// @Success 200 {object} st.StatusOKMessage "ok"
// @Failure 500 {object} st.StatusInternalServerErrorMessage "internal server error"
// @Failure 400 {object} st.StatusBadRequestMessage "bad request"
// @Failure 404 {object} st.StatusNotFoundMessage "not found"
// @Failure 401 {object} st.StatusUnauthorizedMessage "unauthorized"
// @Router /filmlibrary/admin/actors [get]
func GetAllActorsAdmin(c *gin.Context) {
	if _, ok := c.Get("isAuthorizedAdmin"); !ok {
		log.Println("Unauthorized admin")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized admin access denied"})
		return
	}
	actors := postgresql.GetFilmsActor()
	c.JSON(http.StatusOK, actors)
}
