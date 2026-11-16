package api

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	BirthDay string `json:"birth_day"`
}

type UserInput struct {
	Name     string `json:"name"`
	BirthDay string `json:"birth_day"`
}

type ListFilter struct {
	Limit  int
	Offset int
	Name   string
	MinAge *int
	MaxAge *int
}

type Account struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
