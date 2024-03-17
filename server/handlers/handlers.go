package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	l "VK_app/server/dbconn"
	"VK_app/server/postgresql"
	st "VK_app/server/structures"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Login
// @Tags auth
// @Description Login of a user or admin.
// @ID login
// @Accept json
// @Produce json
// @Param input body st.User true "login"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /filmlibrary/login [post]

func Login(c *gin.Context) {
	var user st.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// @Summary Register
// @Tags auth
// @Description Registration of a new user or admin.
// @ID register
// @Accept json
// @Produce json
// @Param input body st.User true "register"
// @Success 201 {string} string "created"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /registration [post]

// RegisterUser handles the registration of a user based on the incoming request.
//
// Parameters:
//   - w: http.ResponseWriter
//   - r: *http.Request
//
// Return type: void
func RegisterUser(c *gin.Context) {
	var user st.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := postgresql.AddUser(&user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

// @Summary AddFilm
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Add a new film to the database.
// @ID add-film
// @Accept json
// @Produce json
// @Param input body st.Film true "Film object for adding"
// @Success 201 {string} string "created"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /admin/films [post]

// PostFilm handles the HTTP POST request for adding a film to the database.
//
// It takes a http.ResponseWriter and a *http.Request as parameters,
// unmarshals the request body into a Film struct, and adds the Film to the
// database. If there is an error, it sends an appropriate HTTP response
// with an error message.
func PostFilm(c *gin.Context) {
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

// @Summary AddActor
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Add a new actor to the database.
// @ID add-actor
// @Accept json
// @Produce json
// @Param input body st.Actor true "Actor object for adding"
// @Success 201 {string} string "created"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /admin/actors [post]

// PostActor handles the POST request for adding an actor to the database.
//
// It takes in the http.ResponseWriter and http.Request as parameters.
// It does not return anything.
func PostActor(c *gin.Context) {
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

// @Summary UpdateActor
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Update an actor in the database.
// @ID update-actor
// @Accept json
// @Produce json
// @Param input body st.Actor true "Actor object for updating"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /admin/actor [put]

// UpdateActor updates the actor using the data from the HTTP request body.
//
// Parameters: w http.ResponseWriter, r *http.Request
// Return type: None

func UpdateActor(c *gin.Context) {
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
}

// @Summary DeleteActor
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Delete an actor from the database.
// @ID delete-actor
// @Accept json
// @Produce json
// @Param input body st.Actor true "Actor object for deleting"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /admin/actor [delete]

// DeleteActor deletes an actor based on the provided request body.
//
// Parameters:
//   - w: http.ResponseWriter
//   - r: http.Request

func DeleteActor(c *gin.Context) {
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
}

// @Summary UpdateFilm
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Update a film in the database.
// @ID update-film
// @Accept json
// @Produce json
// @Param input body st.Film true "Film object for updating"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /admin/film [put]

// UpdateFilm updates a film based on the request body.
//
// Parameters: w http.ResponseWriter, r *http.Request
// Return type: None

func UpdateFilm(c *gin.Context) {
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
}

// @Summary DeleteFilm
// @Security AdminKeyAuth
// @Tags Admin Functions
// @Description Delete a film from the database.
// @ID delete-film
// @Accept json
// @Produce json
// @Param input body st.Film true "Film object for deleting"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /admin/film [delete]

// DeleteFilm deletes a film using the data from the request body.
//
// Parameters:
//
//	w http.ResponseWriter - the response writer
//	r *http.Request - the incoming request

func DeleteFilm(c *gin.Context) {
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
}

// @Summary GetSortedFilms
// @Security ApiKeyAuth
// @Tags User Functions
// @Description Get films sorted by rating, name, or date.
// @ID get-sorted-films
// @Accept string
// @Produce json
// @Param input body string true "Rating, name, or date to be sorted by"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /filmssorted [post]

// GetSortedFilms retrieves and returns sorted films based on the given criteria.
//
// Parameters: w http.ResponseWriter, r *http.Request
// Return type: None
func GetSortedFilms(c *gin.Context) {
	var films []st.Film
	var buf bytes.Buffer
	// читаем тело запроса
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
		c.JSON(http.StatusNotFound, gin.H{"error": "No films found"})
		return
	}
	c.JSON(http.StatusOK, films)
}

// @Summary GetFilmByPiece
// @Security ApiKeyAuth
// @Tags User Functions
// @Description Get films based on a JSON fragment, which contains a piece of the film name or actor name.
// @ID get-film-by-piece
// @Accept json
// @Produce json
// @Param input body st.JSONFragment true "JSON fragment with a piece of the film name or actor name to search for"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /filmspiece [post]

// GetFilmByPiece is a Go function to retrieve films based on a JSON fragment.
// It takes a http.ResponseWriter and *http.Request as parameters and does not return anything.
func GetFilmByPiece(c *gin.Context) {
	var JSONInput st.JSONFragment
	var films []st.Film
	var buf bytes.Buffer
	// читаем тело запроса
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

// @Summary GetAllActors
// @Security ApiKeyAuth
// @Tags User Functions
// @Description Get all actors
// @ID get-all-actors
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /actors [get]

// GetAllActors retrieves all actors from the PostgreSQL database and returns them as a JSON response.
//
// Parameters:
// - w: an http.ResponseWriter object used to write the JSON response.
// - r: an http.Request object representing the incoming request.
//
// Returns: None.
func GetAllActors(c *gin.Context) {
	actors := postgresql.GetFilmsActor()
	c.JSON(http.StatusOK, actors)
}

// @Summary GetSortedFilmsAdmin
// @Security AdminKeyAuth
// @Tags User Functions
// @Description Get films sorted by rating, name, or date.
// @ID get-sorted-films
// @Accept string
// @Produce json
// @Param input body string true "Rating, name, or date to be sorted by"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /admin/filmssorted [post]

// GetSortedFilms retrieves and returns sorted films based on the given criteria.
//
// Parameters: w http.ResponseWriter, r *http.Request
// Return type: None

func GetSortedFilmsAdmin(c *gin.Context) {
	var films []st.Film
	var buf bytes.Buffer
	// читаем тело запроса
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
		c.JSON(http.StatusNotFound, gin.H{"error": "No films found"})
		return
	}
	c.JSON(http.StatusOK, films)
}

// @Summary GetFilmByPieceAdmin
// @Security AdminKeyAuth
// @Tags User Functions
// @Description Get films based on a JSON fragment, which contains a piece of the film name or actor name.
// @ID get-film-by-piece
// @Accept json
// @Produce json
// @Param input body st.JSONFragment true "JSON fragment containing a piece of the film name or actor name to be searched for"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /admin/filmspiece [post]

// GetFilmByPiece is a Go function to retrieve films based on a JSON fragment.
// It takes a http.ResponseWriter and *http.Request as parameters and does not return anything.
func GetFilmByPieceAdmin(c *gin.Context) {
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

// @Summary GetAllActors
// @Security AdminKeyAuth
// @Tags User Functions
// @Description Get all actors
// @ID get-all-actors
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /admin/actors [get]

// GetAllActors retrieves all actors from the PostgreSQL database and returns them as a JSON response.
//
// Parameters:
// - w: an http.ResponseWriter object used to write the JSON response.
// - r: an http.Request object representing the incoming request.
//
// Returns: None.
func GetAllActorsAdmin(c *gin.Context) {
	actors := postgresql.GetFilmsActor()
	c.JSON(http.StatusOK, actors)
}
