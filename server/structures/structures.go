package structures

//swagger:model
type User struct {
	Login    string `json:"login" example:"johndoe"`
	Password string `json:"password" example:"psjfb10"`
	Role     int    `json:"role" example:"1"`
}

//swagger:model
type Film struct {
	Id          int     `json:"id" example:"7"`
	Name        string  `json:"name" example:"Интерстеллар"`
	Description string  `json:"description" example:"Описание фильма"`
	Date        string  `json:"date" example:"20170301"`
	Rating      float32 `json:"rating" example:"8.9"`
}

//swagger:model
type Actor struct {
	Id         int    `json:"id" example:"7"`
	Name       string `json:"name" example:"Александр"`
	Surname    string `json:"surname" example:"Пушкин"`
	FatherName string `json:"fathername" example:"Сергеевич"`
	BirthDate  string `json:"birthdate" example:"17990606"`
	Sex        string `json:"sex" example:"m"`
}

//swagger:model
type FilmResponse struct {
	Id   int    `json:"id" example:"7"`
	Name string `json:"name" example:"Интерстеллар"`
}

//swagger:model
type ActorResponse struct {
	Id         int            `json:"id" example:"7"`
	Name       string         `json:"name" example:"Александр"`
	Surname    string         `json:"surname" example:"Пушкин"`
	FatherName string         `json:"fathername" example:"Сергеевич"`
	BirthDate  string         `json:"birthdate" example:"17990606"`
	Sex        string         `json:"sex" example:"m"`
	Films      []FilmResponse `json:"films" example:"[{\"id\":7,\"name\":\"Интерстеллар\"}]"`
}

//swagger:model
type JSONFragment struct {
	Key      string `json:"key" example:"actor"`
	Fragment string `json:"fragment" example:"иану"`
}

//swagger:model
var (
	Secret = []byte("gBElG5NThZSyeBysksiwusdbqlnwkqhrbv10481u592g")
)
