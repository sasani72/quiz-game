package entity

type User struct {
	ID          uint
	Name        string
	PhoneNumber string
	// Password is always hashed
	Password string
}
