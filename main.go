package main

import (
    "log"
    "net/http"
    "database/sql"
    
    // Echo framework
    "github.com/labstack/echo"
    "github.com/labstack/echo/engine/standard"
    
    // SQLite3
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    e := echo.New()
    
    // Routing
    e.GET("/subscriptions", listSubscriptions)
    e.GET("/subscription/:id", getSubscription)
    
    // Server
    e.Run(standard.New(":8080"))
}

func listSubscriptions(c echo.Context) error {
    // Open database connection
    db, err := sql.Open("sqlite3", "subscriptions")
    checkErr(err)
    
    var (
        subscriptionId int
        subscriptionName string
        subscriptionPeriod string
    )
    
    rows, err := db.Query("SELECT * FROM subscriptions")
    checkErr(err)
    
    defer rows.Close()
    
    var subscriptions interface{}
    for rows.Next() {
        err = rows.Scan(&subscriptionId, &subscriptionName, &subscriptionPeriod)
        checkErr(err)
        
        /*
         * Todo
         * - Create an array
         * - Create a structure
         */
        subscriptions = map[string]interface{}{
            "id": subscriptionId,
            "name": subscriptionName,
            "period": subscriptionPeriod,
        }
    }
    
    return c.JSON(http.StatusOK, subscriptions)
}

func getSubscription(c echo.Context) error {
    // Open database connection
    db, err := sql.Open("sqlite3", "subscriptions")
    checkErr(err)
    
    var (
        subscriptionId int
        subscriptionName string
        subscriptionPeriod string
    )
    
    err = db.QueryRow("SELECT * FROM subscriptions WHERE id = ?", c.Param("id")).Scan(&subscriptionId, &subscriptionName, &subscriptionPeriod)
    checkErr(err)
    defer db.Close()
    
    subscription := map[string]interface{}{
        "id": subscriptionId,
        "name": subscriptionName,
        "period": subscriptionPeriod,
    }
   
    return c.JSON(http.StatusOK, subscription)
}

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}