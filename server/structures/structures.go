package structures

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type Film struct {
	Id          int `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Rating      float32 `json:"rating"`
}

type Actor struct {
	Id         int `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	FatherName string `json:"fathername"`
	BirthDate  string `json:"birthdate"`
	Sex        string `json:"sex"`
}

type FilmResponse struct {
	Id   int `json:"id"`
	Name string `json:"name"`
}

type ActorResponse struct {
	Id         int `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	FatherName string `json:"fathername"`
	BirthDate  string `json:"birthdate"`
	Sex        string `json:"sex"`
	Films      []FilmResponse `json:"films"`
}

type JSONFragment struct {
	Key      string `json:"key"`
	Fragment string `json:"fragment"`
}

var (
	Secret = []byte("gBElG5NThZSyeBysksiwusdbqlnwkqhrbv10481u592g")
)
