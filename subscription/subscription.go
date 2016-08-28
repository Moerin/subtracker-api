package subscription

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/yonmey/subtracker-api/lib/dbc"
	"github.com/yonmey/subtracker-api/lib/errorHandler"
)

// Subscription type
type Subscription struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Duration int    `db:"duration"`
}

// GetAll Gets all the subscriptions
func GetAll(c echo.Context) error {
	db, err := dbc.Connect()
	errorHandler.CheckErr(err)

	subscriptions := []Subscription{}
	err = db.Select(&subscriptions, "SELECT * FROM subscriptions")
	errorHandler.CheckErr(err)

	defer db.Close()

	return c.JSON(http.StatusOK, &subscriptions)
}

// GetOne Get one subscription
func GetOne(c echo.Context) error {
	db, err := dbc.Connect()
	errorHandler.CheckErr(err)

	var ID = c.Param("id")
	subscription := Subscription{}
	err = db.Get(&subscription, "SELECT * FROM subscriptions WHERE id =$1", ID)
	errorHandler.CheckErr(err)

	defer db.Close()

	return c.JSON(http.StatusOK, &subscription)
}

// Add Add one subscription
func Add(c echo.Context) error {
	db, err := dbc.Connect()
	errorHandler.CheckErr(err)

	subscriptionName := c.FormValue("name")
	subscriptionPeriod := c.FormValue("duration")

	stmt, err := db.Preparex("INSERT INTO subscriptions(name, duration) VALUES (?, ?)")
	errorHandler.CheckErr(err)

	_, err = stmt.Exec(subscriptionName, subscriptionPeriod)
	errorHandler.CheckErr(err)
	defer db.Close()

	return c.String(http.StatusCreated, "Subscription created")
}

// Delete Deletes a subscription
func Delete(c echo.Context) error {
	db, err := dbc.Connect()
	errorHandler.CheckErr(err)

	stmt, err := db.Prepare("DELETE FROM subscriptions WHERE id = ?")
	errorHandler.CheckErr(err)

	_, err = stmt.Exec(c.Param("id"))
	errorHandler.CheckErr(err)
	defer db.Close()

	return c.String(http.StatusCreated, "Subscription deleted")
}

// Update Updates a subscription
func Update(c echo.Context) error {
	db, err := dbc.Connect()
	errorHandler.CheckErr(err)

	var updateColumns string
	for k := range c.FormParams() {
		updateColumns += k + " = \"" + c.FormValue(k) + "\", "
	}

	updateColumns = strings.TrimRight(updateColumns, ", ")

	stmt, err := db.Preparex("UPDATE subscriptions SET " + updateColumns + " WHERE id = ?")
	errorHandler.CheckErr(err)

	_, err = stmt.Exec(c.Param("id"))
	errorHandler.CheckErr(err)
	defer db.Close()

	return c.String(http.StatusCreated, "Subscription Updated")
}
