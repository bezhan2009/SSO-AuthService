package models

type User struct {
	ID       uint
	Name     string
	Email    string
	PassHash []byte
	IsAdmin  bool
}
