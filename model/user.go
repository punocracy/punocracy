package model

import "hash"

// PermissionLevel defines the level of privileges a user has
type PermissionLevel int

const (
	// Administrator is the user with the highest privilege
	Administrator PermissionLevel = iota + 1
	// Curator is the user that reviews user submitted phrases
	Curator
	// RegularUser is the user that can rate and submit phrases
	RegularUser
	// NonUser is self explanatory
	NonUser
)

// User is the data structure for a user
type User struct {
	userID    int             `bson:"userID"`
	username  string          `bson:"username"`
	password  hash.Hash       `bson:"passwordHash"`
	email     string          `bson:"email"`
	permLevel PermissionLevel `bson:"permLevel"`
}
