package model

import "hash"

// PermissionLevel defines the level of privileges a user has
type PermissionLevel int

const (
	// Administrator is the user with the highest privilege
	Administrator PermissionLevel = iota
	// Curator is the user that reviews user submitted phrases
	Curator
	// RegularUser is the user that can rate and submit phrases
	RegularUser
	// NonUser is self explanatory
	NonUser
)

// User is the data structure for a user
type User struct {
	userID    int
	username  string
	password  hash.Hash
	email     string
	permLevel PermissionLevel
}
