package main

import (
    "fmt"
    "log"
    "net/http"
    "database/sql"
    
    // Echo framework
    "github.com/labstack/echo"
    "github.com/labstack/echo/engine/standard"
    
    // SQLite3
    _ "github.com/mattn/go-sqlite3"
)

const serverPort = ":8080"

func main() {
    e := echo.New()
    
    // Routing
    // Subscription consulting
    e.GET("/subscriptions", listSubscriptions)
    e.GET("/subscription/:id", getSubscription)
    
    // CRUD
    e.POST("/subscription/add", addSubscription)
    e.POST("/subscription/delete", deleteSubscription)
    e.POST("/subscription/update/:id", updateSubscription)
    
    // Server
    fmt.Println("Serving on port", serverPort)
    e.Run(standard.New(serverPort))
}

func listSubscriptions(c echo.Context) error {
    // Open database connection
    db, err := sql.Open("sqlite3", "subscriptions")
    checkErr(err)
    
    rows, err := db.Query("SELECT * FROM subscriptions")
    checkErr(err)
    
    defer rows.Close()
    
    var subscriptions Subscription
    for rows.Next() {
        err = rows.Scan(&subscriptions.Id, &subscriptions.Name, &subscriptions.Duration)
        checkErr(err)
    } 
    
    return c.JSON(http.StatusOK, subscriptions)
}

func getSubscription(c echo.Context) error {
    // Open database connection
    db, err := sql.Open("sqlite3", "subscriptions")
    checkErr(err)
    
    var subscription Subscription
    err = db.QueryRow("SELECT * FROM subscriptions WHERE id = ?", c.Param("id")).Scan(&subscription.Id, &subscription.Name, &subscription.Duration)
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
    db, err := sql.Open("sqlite3", "subscriptions") 
    checkErr(err)
    
    subscriptionName := c.FormValue("name")
    subscriptionPeriod := c.FormValue("period")

    stmt, err := db.Prepare("INSERT INTO subscriptions(name, duration) VALUES (?, ?)")
    checkErr(err)  
    
    _, err = stmt.Exec(subscriptionName, subscriptionPeriod)
    checkErr(err)
    defer db.Close() 
    
    return c.String(http.StatusCreated, "Subscription created")
}

func deleteSubscription(c echo.Context) error {
    db, err := sql.Open("sqlite3", "subscriptions") 
    checkErr(err)
    
    stmt, err := db.Prepare("DELETE FROM subscriptions WHERE id = ?")
    checkErr(err)  
    
    _, err = stmt.Exec(c.FormValue("id"))
    checkErr(err)
    defer db.Close() 
    
    return c.String(http.StatusCreated, "Subscription deleted")
}

// TODO
func updateSubscription(c echo.Context) error {
    db, err := sql.Open("sqlite3", "subscriptions") 
    checkErr(err)
    
    stmt, err := db.Prepare("UPDATE subscriptions SET ?  = ? WHERE id = ?")
    checkErr(err)  
    
    _, err = stmt.Exec(c.Param("id"))
    checkErr(err)
    defer db.Close() 
    
    return c.String(http.StatusCreated, "Subscription deleted")
}

// Todo: add other fields
type Subscription struct {
    Id int `json:id`
    Name string `json:name`
    Duration string `json:duration`
}