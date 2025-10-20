package model

type Role int

const (
	Guest Role = iota
	Client
	Astronomer
)
