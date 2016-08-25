package main

import (
	//"database/sql"
	"fmt"
	"log"
	"net/http"

	// Echo framework
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"

	// sqlx database wrapper
	"github.com/jmoiron/sqlx"

	// SQLite3
	_ "github.com/mattn/go-sqlite3"
)

const serverPort = ":8080"

func init() {
	schema := `CREATE TABLE IF NOT EXISTS subscriptions (
    id INTEGER PRIMARY KEY,
    name TEXT NULL,
    duration INTEGER NULL);`

	db, err := dbConnect()
	checkErr(err)

	_, err = db.Exec(schema)
	checkErr(err)
}

func dbConnect() (*sqlx.DB, error) {
	return sqlx.Connect("sqlite3", "subscriptions")
}

func main() {
	e := echo.New()

	// Routing
	// Subscription consulting
	e.GET("/subscriptions", listSubscriptions)
	e.GET("/subscription/:id", getSubscription)

	// CRUD
	e.POST("/subscription/add", addSubscription)
	e.DELETE("/subscription/delete/:id", deleteSubscription)
	e.PUT("/subscription/update/:id", updateSubscription)

	// Server
	fmt.Println("Serving on port", serverPort)
	e.Run(standard.New(serverPort))
}

func listSubscriptions(c echo.Context) error {
	db, err := dbConnect()
	checkErr(err)

	subscriptions := []Subscription{}
	err = db.Select(&subscriptions, "SELECT * FROM subscriptions")
	checkErr(err)

	defer db.Close()

	return c.JSON(http.StatusOK, subscriptions)
}

func getSubscription(c echo.Context) error {
	db, err := dbConnect()
	checkErr(err)

	var ID = c.Param("id")
	subscription := Subscription{}
	err = db.Get(&subscription, "SELECT * FROM subscriptions WHERE id =$1", ID)
	checkErr(err)

	defer db.Close()

	return c.JSON(http.StatusOK, &subscription)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func addSubscription(c echo.Context) error {
	db, err := dbConnect()
	checkErr(err)

	subscriptionName := c.FormValue("name")
	subscriptionPeriod := c.FormValue("duration")

	stmt, err := db.Preparex("INSERT INTO subscriptions(name, duration) VALUES (?, ?)")
	checkErr(err)

	_, err = stmt.Exec(subscriptionName, subscriptionPeriod)
	checkErr(err)
	defer db.Close()

	return c.String(http.StatusCreated, "Subscription created")
}

func deleteSubscription(c echo.Context) error {
	db, err := dbConnect()
	checkErr(err)

	stmt, err := db.Prepare("DELETE FROM subscriptions WHERE id = ?")
	checkErr(err)

	_, err = stmt.Exec(c.Param("id"))
	checkErr(err)
	defer db.Close()

	return c.String(http.StatusCreated, "Subscription deleted")
}

func updateSubscription(c echo.Context) error {
	db, err := dbConnect()
	checkErr(err)

	fmt.Println(c.FormParams())

	// TODO
	index := len(c.FormParams())
	var columns = "UPDATE subscriptions SET "
	i := 0
	for k, v := range c.FormParams() {
		if i == (index - 1) {
			columns += k + " = " + "\"" + v[0] + "\""
		} else {
			columns += k + " = " + "\"" + v[0] + "\"" + ", "
		}
		fmt.Println(k, v)
		i++
	}

	columns += " WHERE id = ?"

	fmt.Println(columns)
	stmt, err := db.Preparex(columns)
	checkErr(err)

	_, err = stmt.Exec(c.Param("id"))
	checkErr(err)
	defer db.Close()

	return c.String(http.StatusCreated, "Subscription Updated")
}

// Todo: add other fields
type Subscription struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Duration int    `db:"duration"`
}
