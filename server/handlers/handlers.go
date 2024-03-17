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
// @Router /login [post]

// Login handles the user login process by decoding the request body, checking credentials, and generating a JWT token for authorization.
//
// Parameters:
//   - w: http.ResponseWriter to send responses.
//   - r: http.Request containing the user login information.
//
// Returns nothing.
func Login(w http.ResponseWriter, r *http.Request) {
	var user st.User
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// ищем пользователя в базе данных
	var storedPassword string
	// забираем пароль и id пользователя из базы данных
	row := l.Db.QueryRow("SELECT password FROM filmsactors.users WHERE login = $1", user.Login)
	err = row.Scan(&storedPassword)
	// если пользователь не найден
	if err != nil {
		log.Println(err)
		http.Error(w, "Login is wrong", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		log.Println(err)
		http.Error(w, "Password is", http.StatusInternalServerError)
		return
	}
	user.Password = storedPassword
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":          user.Login,
		"hashedpassword": user.Password,
		"role":           user.Role,
		"exp":            time.Now().Add(time.Hour * 24).Unix(),
	})

	SignedToken, err := jwtToken.SignedString(st.Secret)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", SignedToken)
	w.WriteHeader(http.StatusOK)
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
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user st.User
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = postgresql.AddUser(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
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
func PostFilm(w http.ResponseWriter, r *http.Request) {
	var film st.Film
	var buf bytes.Buffer
	// read the request body
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// unmarshal the request body into a Film struct
	if err = json.Unmarshal(buf.Bytes(), &film); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// add the Film to the database
	err = postgresql.AddFilm(film)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
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
func PostActor(w http.ResponseWriter, r *http.Request) {
	var actor st.Actor
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &actor); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = postgresql.AddActor(actor)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
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
func UpdateActor(w http.ResponseWriter, r *http.Request) {
	var actor st.Actor
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &actor); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = postgresql.UpdateActor(actor)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
func DeleteActor(w http.ResponseWriter, r *http.Request) {
	var actor st.Actor
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &actor); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = postgresql.DelActor(actor.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
func UpdateFilm(w http.ResponseWriter, r *http.Request) {
	var film st.Film
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &film); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = postgresql.UpdateFilm(film)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
func DeleteFilm(w http.ResponseWriter, r *http.Request) {
	var film st.Film
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &film); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = postgresql.DelFilm(film.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
func GetSortedFilms(w http.ResponseWriter, r *http.Request) {
	var films []st.Film
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, "No films found", http.StatusNotFound)
		return
	}
	filmsJSON, err := json.Marshal(films)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(filmsJSON)
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
func GetFilmByPiece(w http.ResponseWriter, r *http.Request) {
	var JSONInput st.JSONFragment
	var films []st.Film
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &JSONInput); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if JSONInput.Key == "actor" {
		films = postgresql.GetFilmsPieceActor(JSONInput.Fragment)
	} else {
		films = postgresql.GetFilmsPieceFilm(JSONInput.Fragment)
	}
	if films == nil {
		log.Println("No films found")
		http.Error(w, "No films found", http.StatusNotFound)
	}
	filmsJSON, err := json.Marshal(films)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(filmsJSON)
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
func GetAllActors(w http.ResponseWriter, r *http.Request) {
	actors := postgresql.GetFilmsActor()
	actorsJSON, err := json.Marshal(actors)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(actorsJSON)
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
func GetSortedFilmsAdmin(w http.ResponseWriter, r *http.Request) {
	var films []st.Film
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, "No films found", http.StatusNotFound)
		return
	}
	filmsJSON, err := json.Marshal(films)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(filmsJSON)
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
func GetFilmByPieceAdmin(w http.ResponseWriter, r *http.Request) {
	var JSONInput st.JSONFragment
	var films []st.Film
	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &JSONInput); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if JSONInput.Key == "actor" {
		films = postgresql.GetFilmsPieceActor(JSONInput.Fragment)
	} else {
		films = postgresql.GetFilmsPieceFilm(JSONInput.Fragment)
	}
	if films == nil {
		log.Println("No films found")
		http.Error(w, "No films found", http.StatusNotFound)
	}
	filmsJSON, err := json.Marshal(films)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(filmsJSON)
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
func GetAllActorsAdmin(w http.ResponseWriter, r *http.Request) {
	actors := postgresql.GetFilmsActor()
	actorsJSON, err := json.Marshal(actors)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(actorsJSON)
}
