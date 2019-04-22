package models

import (
	_"github.com/go-sql-driver/mysql"
	"testing"
        "github.com/jmoiron/sqlx"
)

func newUserForTest(t *testing.T) *User{
    db, err := sqlx.Connect("mysql","root:*password_here@/punocracy")
    if err != nil {
        t.Errorf("database didnt connect. Error: %v", err)
    }

    return NewUser(db)
}

func TestDeleteUser(t *testing.T){
    u := newUserForTest(t) 
    rowToDelete := UserRow{ID: 10}
    u.DeleteUser(nil,rowToDelete)
}
/*
func TestUserCRUD(t *testing.T) {
	u := newUserForTest(t)

	// Signup
	userRow, err := u.Signup(nil, newEmailForTest(), "abc123", "abc123")
	if err != nil {
		t.Errorf("Signing up user should work. Error: %v", err)
	}
	if userRow == nil {
		t.Fatal("Signing up user should work.")
	}
	if userRow.ID <= 0 {
		t.Fatal("Signing up user should work.")
	}

	// DELETE FROM users WHERE id=...
	_, err = u.DeleteById(nil, userRow.ID)
	if err != nil {
		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
	}

}
*/
