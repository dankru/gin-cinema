package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	rtime "github.com/ivahaev/russian-time"
	_ "github.com/lib/pq"
)

//added root for specific film, added redact mode for admin
func main() {
	e := connect()
	if e != nil{
		panic(e.Error())
	}
	
	prepareQueries()
	router := gin.Default()
	
	router.Use(static.Serve("/", static.LocalFile(cfg.Assets, false)))

	// создали куки с очень сложным приватным ключом
	store := sessions.NewCookieStore([]byte("TheMostSecretWord"))
	// Эту куку записываем в сам рутер с ключом сессии
	router.Use(sessions.Sessions("session", store))

	router.LoadHTMLGlob(cfg.HTML + "*.html")

	router.Static("assets", cfg.Assets)


	router.GET("/", getIndexPage)
	router.GET("/:date", getSortedIndexPage)


	films := router.Group("/films")
	{
		films.GET("/", selectCurrentId)
		films.GET("/:dimension/:id/:date/:price", getFilm)
		films.GET("/:dimension/:id/:date/:price/seats", getSeats)

		films.POST("/:dimension/:id/:date/:price", addSeats)

		films.POST("/", addFilm)
		films.POST("/datePrice2d", addDatePrice2d)
		films.POST("/datePrice3d", addDatePrice3d)
	}

	
	router.Run(cfg.ServerHost + ":" + cfg.ServerPort)
}
func getSeats(c *gin.Context){
	var s Seats
	id := c.Param("id")
	date := c.Param("date")
	dimension := c.Param("dimension")
	s.FilmId, _ = strconv.Atoi(id)
	s.Date = date

	if (dimension == "2d"){
		s.SelectSeats2DByIdAndDate()

	}	else {
		s.SelectSeats3DByIdAndDate()
	}
	c.JSON(200, gin.H{"Seats": s.Rows})
}
func addSeats(c *gin.Context) {
	var s Seats
	e := c.BindJSON(&s)
	id := c.Param("id")
	date := c.Param("date")
	s.FilmId, e = strconv.Atoi(id)
	s.Date = date



	if e != nil{
		fmt.Println(e)
	}
	fmt.Println("seat initialized with: ", s)
	dimension := c.Param("dimension")
	if dimension == "2d"{
		s.Initialize2D()

	} else{
		s.Initialize3D()
	}
}

func getIndexPage(c *gin.Context) {
	s := sessions.Default(c)
	admin := false 
	login := false

	role := s.Get("MySecretKey")

	if role == true {
		admin = true
		login = true
	}	else if role == false {
		login = true
	}

	TodayTime := time.Now()
	
	var d2 DatePrice2D
	date := c.Param("date")
	d2.Date = strings.Split(date, " ")[0]
	d2.Select()

	var d3 DatePrice3D
	d3.Date = strings.Split(TodayTime.String(), " ")[0]
	d3.Select()

	// Сделать добавление последовательно по разным Idшникам
	var m Film
	m.Id = d2.FilmId
	m.Select()
	// m.Rows = append(m.Rows, Film{Id: m.Id, Title: m.Title, Description: m.Description, Image: m.Image})
	// for _, element := range d2.Rows{
	// 	fmt.Println("element=", element)
	// 	m.SelectById()
	// 	if m.Id != element.FilmId{
	// 		m.Rows = append(m.Rows, Film{Id: m.Id, Title: m.Title, Description: m.Description, Image: m.Image})
	// 	}
	// }
	// -----------------------------------

	fmt.Println("founded films=", m.Rows)

	reg:= regexp.MustCompile("декабря|января|февраля|марта|апреля|мая|июня|июля|августа|сентября|октября|ноября")

	TomorrowTime := TodayTime.AddDate(0,0,1)
	DATomorrowTime := TodayTime.AddDate(0,0,2)
	FourthTime := TodayTime.AddDate(0,0,3)
	FifthTime := TodayTime.AddDate(0,0,4)

	russianTimeToday := rtime.Time(TodayTime)
	russianTimeTomorrow := rtime.Time(TomorrowTime)
	russianTimeDATomorrow := rtime.Time(DATomorrowTime)
	russianFourthTime := rtime.Time(FourthTime)
	russianFifthTime := rtime.Time(FifthTime)

	TodayDay := russianTimeToday.Weekday().String()
	TomorrowDay := russianTimeTomorrow.Weekday().String()
	DATomorrowDay := russianTimeDATomorrow.Weekday().String()
	FourthTimeDay := russianFourthTime.Weekday().String()
	russianFifthDay := russianFifthTime.Weekday().String()

	TodayDOM := strconv.FormatInt(int64(TodayTime.Day()), 10)
	TomorrowDOM := strconv.FormatInt(int64(TomorrowTime.Day()), 10)
	DATomorrowDOM := strconv.FormatInt(int64(DATomorrowTime.Day()), 10)
	FourthDOM := strconv.FormatInt(int64(FourthTime.Day()), 10)
	FifthDOM := strconv.FormatInt(int64(FifthTime.Day()), 10)

	Tabs := [5]string{strings.Split(TodayTime.String(), " ")[0], strings.Split(TomorrowTime.String()," ")[0], strings.Split(DATomorrowTime.String(), " ")[0],
	strings.Split(FourthTime.String(), " ")[0], strings.Split(FifthTime.String(), " ")[0]}

	Days := [5]string{TodayDOM +" "+ reg.FindAllString(russianTimeToday.String(), -1)[0], 
				TomorrowDOM +" "+ reg.FindAllString(russianTimeTomorrow.String(), -1)[0],
				DATomorrowDOM +" "+ reg.FindAllString(russianTimeDATomorrow.String(), -1)[0], 
				FourthDOM +" "+ reg.FindAllString(russianFourthTime.String(), -1)[0], 
				FifthDOM +" "+ reg.FindAllString(russianFifthTime.String(), -1)[0]}	

	WeekDays := [5]string{TodayDay, TomorrowDay, DATomorrowDay, FourthTimeDay, russianFifthDay}
	
	var Content Cinema
	Content.DatePrice2DRows = d2.Rows
	Content.DatePrice3DRows = d3.Rows
	Content.Available = true
	Content.FilmRows = m.Rows
	c.HTML(200, "index.html", gin.H{
			"Admin": admin,
			"isLogin": login,
			"Content": Content,
			"WeekDays": WeekDays,
			"Days": Days,
			"Tabs": Tabs,
	})
}
func getSortedIndexPage(c *gin.Context) {
	s := sessions.Default(c)
	admin := false 
	login := false

	role := s.Get("MySecretKey")

	if role == true {
		admin = true
		login = true
	}	else if role == false {
		login = true
	}

	TodayTime := time.Now()
	
	var d2 DatePrice2D
	date := c.Param("date")
	d2.Date = strings.Split(date, " ")[0]
	d2.Select()
	fmt.Println("d2.rows=", d2.Rows)

	var nonDuplicates DatePrice2D
	nonDuplicates.Date = strings.Split(date, " ")[0]
	nonDuplicates.SelectByDateWithoutDuplicates()

	var d3 DatePrice3D
	d3.Date = strings.Split(date, " ")[0]
	d3.SelectByDate()

	// Сделать добавление последовательно по разным Idшникам
	var m Film
	for _, element := range nonDuplicates.Rows{
		m.Id = element.FilmId
		m.SelectById()
		m.Rows = append(m.Rows, Film{Id: m.Id, Title: m.Title, Description: m.Description, Image: m.Image})
	}

	reg:= regexp.MustCompile("декабря|января|февраля|марта|апреля|мая|июня|июля|августа|сентября|октября|ноября")

	TomorrowTime := TodayTime.AddDate(0,0,1)
	DATomorrowTime := TodayTime.AddDate(0,0,2)
	FourthTime := TodayTime.AddDate(0,0,3)
	FifthTime := TodayTime.AddDate(0,0,4)

	russianTimeToday := rtime.Time(TodayTime)
	russianTimeTomorrow := rtime.Time(TomorrowTime)
	russianTimeDATomorrow := rtime.Time(DATomorrowTime)
	russianFourthTime := rtime.Time(FourthTime)
	russianFifthTime := rtime.Time(FifthTime)

	TodayDay := russianTimeToday.Weekday().String()
	TomorrowDay := russianTimeTomorrow.Weekday().String()
	DATomorrowDay := russianTimeDATomorrow.Weekday().String()
	FourthTimeDay := russianFourthTime.Weekday().String()
	russianFifthDay := russianFifthTime.Weekday().String()

	TodayDOM := strconv.FormatInt(int64(TodayTime.Day()), 10)
	TomorrowDOM := strconv.FormatInt(int64(TomorrowTime.Day()), 10)
	DATomorrowDOM := strconv.FormatInt(int64(DATomorrowTime.Day()), 10)
	FourthDOM := strconv.FormatInt(int64(FourthTime.Day()), 10)
	FifthDOM := strconv.FormatInt(int64(FifthTime.Day()), 10)

	Tabs := [5]string{strings.Split(TodayTime.String(), " ")[0], strings.Split(TomorrowTime.String()," ")[0], strings.Split(DATomorrowTime.String(), " ")[0],
	strings.Split(FourthTime.String(), " ")[0], strings.Split(FifthTime.String(), " ")[0]}

	Days := [5]string{TodayDOM +" "+ reg.FindAllString(russianTimeToday.String(), -1)[0], 
				TomorrowDOM +" "+ reg.FindAllString(russianTimeTomorrow.String(), -1)[0],
				DATomorrowDOM +" "+ reg.FindAllString(russianTimeDATomorrow.String(), -1)[0], 
				FourthDOM +" "+ reg.FindAllString(russianFourthTime.String(), -1)[0], 
				FifthDOM +" "+ reg.FindAllString(russianFifthTime.String(), -1)[0]}	

	WeekDays := [5]string{TodayDay, TomorrowDay, DATomorrowDay, FourthTimeDay, russianFifthDay}
	



	var Content Cinema
	Content.DatePrice2DRows = d2.Rows
	Content.DatePrice3DRows = d3.Rows
	Content.Available = true
	Content.FilmRows = m.Rows
	c.HTML(200, "index.html", gin.H{
			"Admin": admin,
			"isLogin": login,
			"Content": Content,
			"WeekDays": WeekDays,
			"Days": Days,
			"Tabs": Tabs,
	})
}
func getFilm(c *gin.Context) {
	s := sessions.Default(c)
	admin := false 
	login := false

	role := s.Get("MySecretKey")
	var isLogin = false
	if role == true {
		admin = true
		isLogin = true
	}	else if role == false {
		isLogin = true
	}

	var err error
	var f Film
	id := c.Param("id")
	// price := c.Param("Price")
	f.Id, _ = strconv.Atoi(id)
	err = f.SelectById()
	if err != nil {
		fmt.Println(err.Error())
	}

	date := c.Param("date")

	dimension := c.Param("dimension")
	rows := [...]int{1,2,3,4,5,6,7,8,9,10}
	columns := [...]int{1,2,3,4,5,6,7,8}


	var seats Seats
	seats.FilmId = f.Id
	seats.Date = date

	if dimension == "2d"{
		var d2 DatePrice2D
		d2.Date = date
		d2.FilmId = f.Id
		d2.SelectByIdAndDate()
		
		
		fmt.Println("seats state is=", s)
		seats.SelectSeats2DByIdAndDate()

		
		fmt.Println("seats=", s)
		c.HTML(200, "film.html", gin.H {
			"admin": admin,
			"isLogin": isLogin,
			"login": login,
			"Film": f,
			"DatePrice": d2,
			"Seats": seats.Rows,
			"Rows": rows,
			"Columns": columns,
		})
	} else{
		var d3 DatePrice3D
		d3.Date = date
		d3.FilmId = f.Id
		d3.SelectByIdAndDate()

		seats.SelectSeats3DByIdAndDate()
		
		c.HTML(200, "film.html", gin.H {
			"admin": admin,
			"isLogin": isLogin,
			"login": login,
			"Film": f,
			"DatePrice": d3,
			"Seats": seats.Rows,
			"Rows": rows,
			"Columns": columns,
		})
	}
}
func addFilm(c *gin.Context) {
	var f Film
	e := c.BindJSON(&f)
	f.Add()
	fmt.Println("f=", f)
	fmt.Println("e=", e)
}
func addDatePrice2d(c *gin.Context) {
	var d2 DatePrice2D
	e := c.BindJSON(&d2)
	fmt.Println("error=", e)
	d2.Add()
}
func addDatePrice3d(c *gin.Context) {
	var d3 DatePrice3D
	e := c.BindJSON(&d3)
	fmt.Println("d3=", d3)
	fmt.Println("e=", e)
	d3.Add()
}
func selectCurrentId(c *gin.Context) {
	var f Film
	f.selectCurrentId()
	c.JSON(200, f.Id)
}