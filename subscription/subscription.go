package subscription

// Subscription type
type Subscription struct {
    Id int `db:"id"`
    Name string `db:"name"`
    Duration int `db:"duration"`
}