package postgresql

import (
	"fmt"
	"log"

	"VK_app/internal/structures"

	l "VK_app/internal/dbconn"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// AddUser adds a new user to the database.
//
// It takes a pointer to a structures.User struct as a parameter.
// It returns an error if there was an issue adding the user.

func AddUser(u *structures.User) error {
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		log.Printf("addUser: %v\n", err)
		return err
	}
	u.Password = string(hashedpassword) //hashedpassword
	fmt.Println(u.Password)
	_, err = l.Db.Exec("INSERT INTO users (login,password,role) VALUES ($1, $2, $3)", u.Login, u.Password, u.Role)
	if err != nil {
		log.Printf("addUser: %v\n", err)
		return err
	}
	log.Println("user was appended successfully")
	return nil
}

// GetFilmsSortedRating returns a slice of structures.Film sorted by rating in descending order.
//
// No parameters.
// Returns a slice of structures.Film.
func GetFilmsSortedRating() []structures.Film {
	var films []structures.Film
	rows, err := l.Db.Query("SELECT * FROM films ORDER BY rating DESC")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		film := structures.Film{}
		err := rows.Scan(&film.Id, &film.Name, &film.Description, &film.Date, &film.Rating)
		if err != nil {
			log.Println(err)
			continue
		}
		films = append(films, film)
	}
	return films
}

// GetFilmsSortedDate retrieves and returns films sorted by date.
//
// No parameters.
// Returns a slice of structures.Film.
func GetFilmsSortedDate() []structures.Film {
	var films []structures.Film
	rows, err := l.Db.Query("SELECT * FROM films ORDER BY date DESC")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		film := structures.Film{}
		err := rows.Scan(&film.Id, &film.Name, &film.Description, &film.Date, &film.Rating)
		if err != nil {
			log.Println(err)
			continue
		}
		films = append(films, film)
	}
	return films
}

// GetFilmsSortedName retrieves and returns films sorted by name in descending order.
//
// No parameters.
// Returns a slice of structures.Film.
func GetFilmsSortedName() []structures.Film {
	var films []structures.Film
	rows, err := l.Db.Query("SELECT * FROM films ORDER BY name DESC")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		film := structures.Film{}
		err := rows.Scan(&film.Id, &film.Name, &film.Description, &film.Date, &film.Rating)
		if err != nil {
			log.Println(err)
			continue
		}
		films = append(films, film)
	}
	return films
}

// GetFilmsPieceActor retrieves films related to a specific actor piece.
//
// Parameter:
// piece string - the actor piece to search for in the database.
// []structures.Film - a slice of structures.Film containing the retrieved films.
func GetFilmsPieceActor(piece string) []structures.Film {
	var films []structures.Film
	rows, err := l.Db.Query("SELECT * FROM films WHERE id IN (SELECT film_id FROM actorsfilms WHERE actor_id IN (SELECT id FROM actors WHERE name LIKE $1))", "%"+piece+"%")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		film := structures.Film{}
		err := rows.Scan(&film.Id, &film.Name, &film.Description, &film.Date, &film.Rating)
		if err != nil {
			log.Println(err)
			continue
		}
		films = append(films, film)
	}
	return films
}

// GetFilmsPieceFilm retrieves films containing a specific piece in their name.
//
// piece - the string to search for in film names.
// []structures.Film - a slice of Film structures containing the matching films.
func GetFilmsPieceFilm(piece string) []structures.Film {
	var films []structures.Film
	rows, err := l.Db.Query("SELECT * FROM films WHERE name LIKE $1", "%"+piece+"%")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		film := structures.Film{}
		err := rows.Scan(&film.Id, &film.Name, &film.Description, &film.Date, &film.Rating)
		if err != nil {
			log.Println(err)
			continue
		}
		films = append(films, film)
	}
	return films
}

// GetFilmsActor retrieves a list of actors along with their films from the database.
//
// No parameters.
// Returns a slice of structures.ActorResponse.
func GetFilmsActor() []structures.ActorResponse {
	var actors = map[int]structures.ActorResponse{}
	var actorsResponse []structures.ActorResponse
	var idsactors []int
	rows, err := l.Db.Query("SELECT * FROM actors")
	if err != nil {
		log.Println(err)
		return actorsResponse
	}
	for rows.Next() {
		actor := structures.ActorResponse{}
		err := rows.Scan(&actor.Id, &actor.Name, &actor.Surname, &actor.FatherName, &actor.BirthDate, &actor.Sex)
		if err != nil {
			log.Println(err)
			continue
		}
		actors[actor.Id] = actor
		idsactors = append(idsactors, actor.Id)
	}
	for _, id := range idsactors {
		var films []structures.FilmResponse
		rows, err = l.Db.Query("SELECT id, name FROM films WHERE id IN (SELECT film_id FROM actorsfilms WHERE actor_id=$1)", id)
		if err != nil {
			log.Println(err)
			return actorsResponse
		}
		for rows.Next() {
			film := structures.FilmResponse{}
			err := rows.Scan(&film.Id, &film.Name)
			if err != nil {
				log.Println(err)
				continue
			}
			films = append(films, film)
		}
		actor := actors[id]
		actor.Films = films
		actors[id] = actor
	}
	for _, value := range actors {
		actorsResponse = append(actorsResponse, value)
	}
	return actorsResponse
}

// DelActor deletes information about an actor from the database.
//
// It takes an integer parameter 'id' and returns an error.
func DelActor(id int) error {
	_, err := l.Db.Exec("DELETE FROM actors WHERE id=$1", id)
	if err != nil {
		log.Println("problem with deleting information about actor", err)
		return err
	}
	return nil
}

// DelFilm deletes the film with the given ID.
//
// Parameter(s):
//
//	id int - the ID of the film to be deleted.
//
// Return type(s):
//
//	error - an error, if any.
func DelFilm(id int) error {
	_, err := l.Db.Exec("DELETE FROM films WHERE id=$1", id)
	if err != nil {
		log.Println("problem with deleting information about actor", err)
		return err
	}
	return nil
}

// UpdateFilm updates film information in the database.
//
// Takes a structures.Film object as input.
// Returns an error.
func UpdateFilm(film structures.Film) error {
	args := []interface{}{}
	counter := 1
	query := "UPDATE films SET"
	if film.Name != "" {
		query += fmt.Sprintf(" name=$%d,", counter)
		args = append(args, film.Name)
		counter++
	}
	if film.Description != "" {
		query += fmt.Sprintf(" description=$%d,", counter)
		args = append(args, film.Description)
		counter++
	}
	if film.Date != "" {
		query += fmt.Sprintf(" date=$%d,", counter)
		args = append(args, film.Date)
		counter++
	}
	if film.Rating != 0 {
		query += fmt.Sprintf(" rating=$%d,", counter)
		args = append(args, film.Rating)
		counter++
	}
	if counter == 1 {
		log.Println("no fields to update")
		return nil
	}

	query = query[:len(query)-1]
	query += fmt.Sprintf(" WHERE id=$%d", counter)
	args = append(args, film.Id)
	_, err := l.Db.Exec(query, args...)
	if err != nil {
		log.Println("problem with updating information about film", err)
		return err
	}
	return nil
}

// UpdateActor updates the information of an actor in the actors table.
//
// The actor parameter is the structure containing the updated actor information.
// It should have at least one field set to update the corresponding record in the database.
//
// The function returns an error if there was a problem with updating the actor information.
func UpdateActor(actor structures.Actor) error {
	args := []interface{}{}
	counter := 1
	query := "UPDATE actors SET"
	if actor.Name != "" {
		query += fmt.Sprintf(" name=$%d,", counter)
		args = append(args, actor.Name)
		counter++
	}
	if actor.Surname != "" {
		query += fmt.Sprintf(" surname=$%d,", counter)
		args = append(args, actor.Surname)
		counter++
	}
	if actor.FatherName != "" {
		query += fmt.Sprintf(" fathername=$%d,", counter)
		args = append(args, actor.FatherName)
		counter++
	}
	if actor.BirthDate != "" {
		query += fmt.Sprintf(" birthdate=$%d,", counter)
		args = append(args, actor.BirthDate)
		counter++
	}
	if actor.Sex != "" {
		query += fmt.Sprintf(" sex=$%d,", counter)
		args = append(args, actor.Sex)
		counter++
	}
	if counter == 1 {
		log.Println("no fields to update")
		return nil
	}
	query = query[:len(query)-1]
	query += fmt.Sprintf(" WHERE id=$%d", counter)
	args = append(args, actor.Id)
	_, err := l.Db.Exec(query, args...)
	if err != nil {
		log.Println("problem with updating information about actor", err)
		return err
	}
	return nil
}

// AddActor adds an actor to the database.
//
// It takes a structures.Actor as a parameter and returns an error.
func AddActor(actor structures.Actor) error {
	_, err := l.Db.Exec("INSERT INTO actors (name, surname, fathername, birthdate, sex) VALUES ($1, $2, $3, $4, $5)", actor.Name, actor.Surname, actor.FatherName, actor.BirthDate, actor.Sex)
	if err != nil {
		log.Println("problem with adding information about actor", err)
		return err
	}
	return nil
}

// CheckActor checks the existence of an actor in the database.
//
// It takes a structures.Actor as a parameter and returns an error.
func CheckActor(id int) error {
	var name string
	err := l.Db.QueryRow("SELECT name FROM actors WHERE id = $1", id).Scan(&name)
	if err != nil {
		log.Println("problem with checking information about actor", err)
		return err
	}
	return nil
}

// CheckFilm checks theexistsence of a film in the database.
//
// Parameter: film of type structures.Film.
// Returns an error.
func CheckFilm(id int) error {
	var name string
	err := l.Db.QueryRow("SELECT name FROM films WHERE id = $1", id).Scan(&name)
	if err != nil {
		log.Println("problem with checking information about film", err)
		return err
	}
	return nil
}

func AddActorFilm(actorFilm structures.ActorFilm) error {
	_, err := l.Db.Exec("INSERT INTO actorsfilms (actor_id, film_id) VALUES ($1, $2)", actorFilm.ActorID, actorFilm.FilmID)
	if err != nil {
		log.Println("problem with adding information about actor", err)
		return err
	}
	return nil
}

// AddFilm adds a film to the database.
//
// Parameter: film structures.Film
// Return type: error
func AddFilm(film structures.Film) error {
	_, err := l.Db.Exec("INSERT INTO films (name, description, date, rating) VALUES ($1, $2, $3, $4)", film.Name, film.Description, film.Date, film.Rating)
	if err != nil {
		log.Println("problem with adding information about film", err)
		return err
	}
	return nil
}
