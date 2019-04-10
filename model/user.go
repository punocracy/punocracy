package model

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
	UserID    int             `bson:"userID"`
	Username  string          `bson:"username"`
	Password  string          `bson:"passwordHash"`
	Email     string          `bson:"email"`
	PermLevel PermissionLevel `bson:"permLevel"`
}
