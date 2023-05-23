package main

type Genre struct {
	Name string
	Rows []Genre
}

type Setting struct {
	ServerHost string
	ServerPort string
	PgHost     string
	PgPort     string
	PgUser     string
	PgDB       string
	PgPass     string
	Image      string
	Data       string
	Assets     string
	HTML       string
}

type User struct {
	Id       string `json:"-"`
	Login    string `json:"Login"`
	Password string `json:"Password"`
	Admin    bool   `json:"Admin"`
	Rows     []User
}

type Film struct {
	Id          int
	Title       string
	Description string
	Image       string
	Rows        []Film
}

type Cinema struct {
	Id              int `json:"-"`
	DatePrice2DRows []DatePrice2D
	DatePrice3DRows []DatePrice3D
	Available       bool
	FilmRows        []Film
	Rows            []Cinema
}

type DatePrice2D struct {
	Id     int     `json:"-"`
	Date   string  `json:"Date"`
	Price  float32 `json:"Price"`
	FilmId int
	Time   string
	Rows   []DatePrice2D
}
type DatePrice3D struct {
	Id     int `json:"-"`
	Date   string
	Price  float32
	FilmId int
	Time   string
	Rows   []DatePrice3D
}

type Seats struct {
	Id     int    `json:"-"`
	FilmId int    `json:"FilmId"`
	Date   string `json:"Date"`
	Row    int    `json:"row"`
	Column int    `json:"column"`
	Taken  bool   `json:"taken"`
	Rows   []Seats
}