package model

type Role int

const (
	Client Role = iota
	Astronomer
	Admin
)
