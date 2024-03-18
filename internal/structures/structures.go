package structures

//swagger:model
type User struct {
	Login    string `json:"login" example:"john_doe"`
	Password string `json:"password" example:"psjfb10"`
	Role     int    `json:"role" example:"1"`
}

//swagger:model
type Film struct {
	Id          int     `json:"id" example:"3"`
	Name        string  `json:"name" example:"Затмение"`
	Description string  `json:"description" example:"Описание фильма"`
	Date        string  `json:"date" example:"20161125"`
	Rating      float32 `json:"rating" example:"5.8"`
}

//swagger:model
type Actor struct {
	Id         int    `json:"id" example:"4"`
	Name       string `json:"name" example:"Сергей"`
	Surname    string `json:"surname" example:"Бурунов"`
	FatherName string `json:"fathername" example:"Александрович"`
	BirthDate  string `json:"birthdate" example:"19770306"`
	Sex        string `json:"sex" example:"m"`
}

//swagger:model
type FilmResponse struct {
	Id   int    `json:"id" example:"3"`
	Name string `json:"name" example:"Затмение"`
}

//swagger:model
type ActorResponse struct {
	Id         int            `json:"id" example:"4"`
	Name       string         `json:"name" example:"Сергей"`
	Surname    string         `json:"surname" example:"Бурунов"`
	FatherName string         `json:"fathername" example:"Александрович"`
	BirthDate  string         `json:"birthdate" example:"19770306"`
	Sex        string         `json:"sex" example:"m"`
	Films      []FilmResponse `json:"films" example:"[{\"id\":3,\"name\":\"Затмение\"}]"`
}

//swagger:model
type JSONFragment struct {
	Key      string `json:"key" example:"actor"`
	Fragment string `json:"fragment" example:"иану"`
}

//swagger:model
type KeySort struct {
	Key   string `json:"key" example:"key"`
	Value string `json:"value" example:"name"`
}

//swagger:model
type StatusOKMessage struct {
	Message string `json:"status" example:"ok"`
}

//swagger:model
type StatusBadRequestMessage struct {
	Message string `json:"message" example:"Error: bad request"`
}

//swagger:model
type StatusUnauthorizedMessage struct {
	Message string `json:"message" example:"Error: unauthorized"`
}

//swagger:model
type StatusNotFoundMessage struct {
	Message string `json:"message" example:"Error: not found"`
}

//swagger:model
type StatusInternalServerErrorMessage struct {
	Message string `json:"message" example:"Error: internal server error"`
}

//swagger:model
var (
	Secret = []byte("gBElG5NThZSyeBysksiwusdbqlnwkqhrbv10481u592g")
)
