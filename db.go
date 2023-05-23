package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var db *sql.DB
// Словарь Queries. Ключ - строка, значение - sql statement, иначе говоря запрос
var Queries map[string]*sql.Stmt

func connect() error {
	var e error
	
	Queries = make(map[string]*sql.Stmt)

	// Генерируем строку подключения со значениями из кфг файла
	db, e = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.PgHost, cfg.PgPort, cfg.PgUser, cfg.PgPass, cfg.PgDB))
	if e != nil {
		return e
	}

	return nil
}

func prepareQueries() {
	var e error

	Queries["Select#Films"], e = db.Prepare(`Select * from "Film" Order by "Id" desc`)
	Queries["Select#FilmsByDate"], e = db.Prepare(`Select * from "Film" where "Date"=$1`)
	Queries["Insert#Film"], e = db.Prepare(`Insert Into "Film" ("Title", "Description", "Image") values($1, $2, $3)`)
	Queries["Select#CurrentFilmId"], e = db.Prepare(`Select "Id" from "Film" Order by "Id"`)
	Queries["Select#FilmById"], e = db.Prepare(`Select "Title", "Description", "Image" from "Film" where "Id"=$1`)

	Queries["Select#DatePrice2D"], e = db.Prepare(`Select "Id", "Date", "Price", "FilmId" from "DatePrice2D" order by "FilmId" desc`)
	Queries["Select#DatePrice2DByDate"], e = db.Prepare(`Select "Id", "Date", "Price", "FilmId" from "DatePrice2D" where "Date"::text LIKE $1 `)
	Queries["Select#DatePrice2DByDateWithoutDuplicates"], e = db.Prepare(`SELECT MIN("Id") AS "Id", "FilmId" from "DatePrice2D" where "Date"::text LIKE $1 group by "FilmId" ORDER BY "FilmId" desc`)
	Queries["Insert#DatePrice2D"], e = db.Prepare(`Insert Into "DatePrice2D" ("Date", "Price", "FilmId") values($1, $2, $3)`)
	Queries["Select#DatePrice2DByIdAndDate"], e = db.Prepare(`Select "Date", "Price" from "DatePrice2D" where "FilmId"=$1 and  "Date"=$2`)

	Queries["Select#DatePrice3D"], e = db.Prepare(`Select "Id", "Date", "Price", "FilmId" from "DatePrice3D" order by "FilmId" desc`)
	Queries["Select#DatePrice3DByDate"], e = db.Prepare(`Select "Id", "Date", "Price", "FilmId" from "DatePrice3D" where "Date"::text LIKE $1`)
	Queries["Insert#DatePrice3D"], e = db.Prepare(`Insert Into "DatePrice3D" ("Date", "Price", "FilmId") values($1, $2, $3)`)
	Queries["Select#DatePrice3DByIdAndDate"], e = db.Prepare(`Select "Date", "Price" from "DatePrice3D" where "FilmId"=$1 and  "Date"=$2`)
	
	Queries["Select#Seats2DByIdAndDate"], e = db.Prepare(`Select "FilmId", "Date", "row", "column", "taken" from "AvailableSeat2D" where "FilmId"=$1 and  "Date"=$2`)
	Queries["Select#Seats3DByIdAndDate"], e = db.Prepare(`Select "FilmId", "Date", "row", "column", "taken" from "AvailableSeat3D" where "FilmId"=$1 and  "Date"=$2`)
	Queries["Insert#Seats2D"], e = db.Prepare(`Insert Into "AvailableSeat2D" ("FilmId", "Date", "row", "column", "taken") values($1, $2, $3, $4, $5)`)
	Queries["Insert#Seats3D"], e = db.Prepare(`Insert Into "AvailableSeat3D" ("FilmId", "Date", "row", "column", "taken") values($1, $2, $3, $4, $5)`)

	if e != nil {
		
		panic(e.Error())
	}
}

func (m *Film) Select() error {
	stmt, ok := Queries["Select#Films"]
	if !ok {
		return errors.New("Chosen query doesn't exist")
	}
	
	rows, e := stmt.Query()
	if e != nil{
		return e
	}
	for rows.Next() {
		e = rows.Scan(&m.Id, &m.Title, &m.Description, &m.Image)

		if e != nil {	
			return e
		}

		m.Rows = append(m.Rows, Film{Id: m.Id, Title: m.Title, Description: m.Description, Image: m.Image})
	}
	fmt.Println("Rows=", m.Rows)

	return nil
}
func (m *Film) Add() error {
	stmt, ok := Queries["Insert#Film"]
	if !ok {
		return errors.New("Films query doesn't exist")
	}
	_ = stmt.QueryRow(m.Title, m.Description, m.Image	)
	fmt.Println("After querying object state is= ", m)

	return nil
}
func (m *Film) selectCurrentId() error {
	stmt, ok := Queries["Select#CurrentFilmId"]
	if !ok {
		return errors.New("Films query doesn't exist")
	}
	r, e := stmt.Query()
	for r.Next() {
		e = r.Scan(&m.Id)

		if e != nil {	
			return e
		}

		m.Rows = append(m.Rows, Film{Id: m.Id})
	}
	fmt.Println("m.Rows =", m.Rows)
	return nil
}
func (m *Film) SelectById() error {
	stmt, ok := Queries["Select#FilmById"]
	if !ok {
		return errors.New("Films query doesn't exist")
	}
	fmt.Println("SelectBy Id chosen id=", m.Id)
	r := stmt.QueryRow(m.Id)
	e := r.Scan(&m.Title, &m.Description, &m.Image)
	if e != nil { 
		fmt.Println(e.Error())
		return errors.New("Invalid login or password")
	}
	return nil
}

func (m *DatePrice2D) Select() error {
	stmt, ok := Queries["Select#DatePrice2D"]
	if !ok {
		return errors.New("Chosen query doesn't exist")
	}
	rows, e := stmt.Query()
	if e != nil{
		return e
	}
	for rows.Next() {
		e = rows.Scan(&m.Id, &m.Date, &m.Price, &m.FilmId)
		s := strings.Split(m.Date, "T")

		if e != nil {	
			fmt.Println(e)
		}
		
		m.Rows = append(m.Rows, DatePrice2D{Id: m.Id, Date: m.Date, Price: m.Price, FilmId: m.FilmId, Time: s[1]})
	}

	return nil
}
func (m *DatePrice2D) SelectByDate() error {
	stmt, ok := Queries["Select#DatePrice2DByDate"]
	if !ok {
		return errors.New("Chosen query doesn't exist")
	}
	rows, e := stmt.Query(m.Date+"%")
	if e != nil{
		return e
	}
	for rows.Next() {
		e = rows.Scan(&m.Id, &m.Date, &m.Price, &m.FilmId)
		s := strings.Split(m.Date, "T")

		if e != nil {	
			fmt.Println(e)
		}
		
		m.Rows = append(m.Rows, DatePrice2D{Id: m.Id, Date: m.Date, Price: m.Price, FilmId: m.FilmId, Time: s[1]})
	}
	return nil
}
func (m *DatePrice2D) SelectByDateWithoutDuplicates() error {
	stmt, ok := Queries["Select#DatePrice2DByDateWithoutDuplicates"]
	if !ok {
		return errors.New("Chosen query doesn't exist")
	}
	rows, e := stmt.Query(m.Date+"%")
	if e != nil{
		return e
	}
	for rows.Next() {
		e = rows.Scan(&m.Id, &m.FilmId)

		if e != nil {	
			fmt.Println(e)
		}
		
		m.Rows = append(m.Rows, DatePrice2D{Id: m.Id, FilmId: m.FilmId})
	}
	fmt.Println("Found dp2d non-duplicates=", m.Rows)
	return nil
}

func (m *DatePrice2D) Add() error {
	stmt, ok := Queries["Insert#DatePrice2D"]
	if !ok {
		return errors.New("datePrice2d query doesn't exist")
	}
	_ = stmt.QueryRow(m.Date, m.Price, m.FilmId)
	fmt.Println("After querying object state is= ", m)

	return nil
}
func (m *DatePrice2D) SelectByIdAndDate() error {
	stmt, ok := Queries["Select#DatePrice2DByIdAndDate"]
	if !ok {
		return errors.New("Films query doesn't exist")
	}
	fmt.Println("chosen filmId=", m.FilmId)
	fmt.Println("chosen Date=", m.Date)

	r := stmt.QueryRow(m.FilmId, m.Date)
	e := r.Scan(&m.Date, &m.Price)
	if e != nil { 
		fmt.Println(e.Error())
		return errors.New("Invalid login or password")
	}
	s := strings.Split(m.Date, "T")
	m.Time = s[1]
	fmt.Println("found datePrice =", m)
	return nil
}

func (m *DatePrice3D) Select() error {
	stmt, ok := Queries["Select#DatePrice3D"]
	if !ok {
		return errors.New("Chosen query doesn't exist")
	}
	rows, e := stmt.Query()
	if e != nil{
		return e
	}
	for rows.Next() {
		e = rows.Scan(&m.Id, &m.Date, &m.Price, &m.FilmId)
		s := strings.Split(m.Date, "T")
		if e != nil {	
			fmt.Println(e)
		}
		m.Rows = append(m.Rows, DatePrice3D{Id: m.Id, Date: m.Date, Price: m.Price, FilmId: m.FilmId, Time: s[1]})
	}
	fmt.Println("DatePrice3D=", m.Rows)

	return nil
}
func (m *DatePrice3D) SelectByDate() error {
	stmt, ok := Queries["Select#DatePrice2DByDate"]
	if !ok {
		return errors.New("Chosen query doesn't exist")
	}
	rows, e := stmt.Query(m.Date+"%")
	if e != nil{
		return e
	}
	for rows.Next() {
		e = rows.Scan(&m.Id, &m.Date, &m.Price, &m.FilmId)
		s := strings.Split(m.Date, "T")

		if e != nil {	
			fmt.Println(e)
		}
		
		m.Rows = append(m.Rows, DatePrice3D{Id: m.Id, Date: m.Date, Price: m.Price, FilmId: m.FilmId, Time: s[1]})
	}
	return nil
}
func (m *DatePrice3D) SelectByIdAndDate() error {
	stmt, ok := Queries["Select#DatePrice3DByIdAndDate"]
	if !ok {
		return errors.New("Films query doesn't exist")
	}
	fmt.Println("chosen filmId=", m.FilmId)
	fmt.Println("date=", m.Date)
	r := stmt.QueryRow(m.FilmId, m.Date)
	e := r.Scan(&m.Date, &m.Price)
	if e != nil { 
		fmt.Println(e.Error())
		return errors.New("Invalid login or password")
	}

	s := strings.Split(m.Date, "T")
	m.Time = s[1]
	fmt.Println("found datePrice =", m)
	return nil
}
func (m *DatePrice3D) Add() error {
	stmt, ok := Queries["Insert#DatePrice3D"]
	if !ok {
		return errors.New("datePrice3d query doesn't exist")
	}

	_ = stmt.QueryRow(m.Date, m.Price, m.FilmId)
	fmt.Println("After querying object state is= ", m)

	return nil
}

func (m *Seats) Initialize2D() error {
	stmt, ok := Queries["Insert#Seats2D"]
	if !ok {
		return errors.New("Seats query doesn't exist")
	}
	_ = stmt.QueryRow(m.FilmId, m.Date, m.Row, m.Column, m.Taken)
	fmt.Println("After querying object state is= ", m)

	return nil
}

func (m *Seats) Initialize3D() error {
	stmt, ok := Queries["Insert#Seats3D"]
	if !ok {
		return errors.New("Seats query doesn't exist")
	}
	_ = stmt.QueryRow(m.FilmId, m.Date, m.Row, m.Column, m.Taken)
	fmt.Println("After querying object state is= ", m)

	return nil
}

func (m *Seats) SelectSeats2DByIdAndDate() error {
	stmt, ok := Queries["Select#Seats2DByIdAndDate"]
	if !ok {
		return errors.New("Seats query doesn't exist")
	}
	rows, e := stmt.Query(m.FilmId, m.Date)
	for rows.Next() {
		e = rows.Scan(&m.FilmId, &m.Date, &m.Row, &m.Column, &m.Taken)

		if e != nil {	
			return e
		}

		m.Rows = append(m.Rows, Seats{FilmId: m.FilmId, Date: m.Date, Row: m.Row, Column: m.Column, Taken: m.Taken})
	}
	return nil
}

func (m *Seats) SelectSeats3DByIdAndDate() error {
	stmt, ok := Queries["Select#Seats3DByIdAndDate"]
	if !ok {
		return errors.New("Seats query doesn't exist")
	}
	rows, e := stmt.Query(m.FilmId, m.Date)
	for rows.Next() {
		e = rows.Scan(&m.FilmId, &m.Date, &m.Row, &m.Column, &m.Taken)

		if e != nil {	
			return e
		}

		m.Rows = append(m.Rows, Seats{FilmId: m.FilmId, Date: m.Date, Row: m.Row, Column: m.Column, Taken: m.Taken})
	}
	return nil
}