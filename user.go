package main

import (
	"time"
)

type User struct {
	ID        string
	Firstname string
	Lastname  string
	Email     string
	Age       uint
	Created   time.Time
}
