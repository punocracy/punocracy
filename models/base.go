package models

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

/*
    This type is to be derived from for each table in our database
*/
type Base struct {
	db    *sqlx.DB
	table string
	hasID bool
}

//not correct, base is just a begining. 
/*
func CreateBase ( iDb *sqlx.DB, iTable string, iHasID bool){
    b := &Base{ db = iDb, table = iTable, hadID = iHasID}
    return b
}
*/

/*
*   This function is used to check if a transaction has been commited, and a new one needs to be assigned
    This transaction system is designed such that, if no transaction is specified each function called is atomic, otherwise if the developer specifies a transaction, the functions can be run in the same transaction
*/
func (b *Base) newTransactionIfNeeded(tx *sqlx.Tx) (*sqlx.Tx, bool, error) {
	var err error
	wrapInSingleTransaction := false

	if tx != nil {
		return tx, wrapInSingleTransaction, nil
	}
        //new transaction started
	tx, err = b.db.Beginx()
	if err == nil {
		wrapInSingleTransaction = true
	}

	if err != nil {
		return nil, wrapInSingleTransaction, err
	}

	return tx, wrapInSingleTransaction, nil
}

/*
   Inserts one row into a table 
*/
func (b *Base) InsertIntoTable(tx *sqlx.Tx, data map[string]interface{}) (sql.Result, error) {
	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	qMarks := make([]string, 0)
	values := make([]interface{}, 0)
        
        //grows the list of key value pairs extracted from data and stored as seperate lists
        //keys represent the columns and qmarks are the bindvars which are replaced by values when tx.Exec is called
	for key, value := range data {
		keys = append(keys, key)
		qMarks = append(qMarks, "?")
		values = append(values, value)
	}

	query := fmt.Sprintf(
		"INSERT INTO %v (%v) VALUES (%v)",
		b.table,
		strings.Join(keys, ","),
		strings.Join(qMarks, ","))

	result, err := tx.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	return result, err
}

func (b *Base) UpdateFromTable(tx *sqlx.Tx, data map[string]interface{}, where string) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}
        //creates a new transaction if one has not been specified
        //returns nil on error
	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	keysWithQuestionMarks := make([]string, 0)
	values := make([]interface{}, 0)

	for key, value := range data {
		keysWithQuestionMark := fmt.Sprintf("%v=?", key)
		keysWithQuestionMarks = append(keysWithQuestionMarks, keysWithQuestionMark)
		values = append(values, value)
	}

	query := fmt.Sprintf(
		"UPDATE %v SET %v WHERE %v",
		b.table,
		strings.Join(keysWithQuestionMarks, ","),
		where)

	result, err = tx.Exec(query, values...)

	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	return result, err
}

func (b *Base) UpdateByID(tx *sqlx.Tx, data map[string]interface{}, id int64) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	keysWithQuestionMarks := make([]string, 0)
	values := make([]interface{}, 0)

	for key, value := range data {
		keysWithQuestionMark := fmt.Sprintf("%v=?", key)
		keysWithQuestionMarks = append(keysWithQuestionMarks, keysWithQuestionMark)
		values = append(values, value)
	}

	// Add id as part of values
	values = append(values, id)

	query := fmt.Sprintf(
		"UPDATE %v SET %v WHERE id=?",
		b.table,
		strings.Join(keysWithQuestionMarks, ","))

	result, err = tx.Exec(query, values...)

	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	return result, err
}

func (b *Base) UpdateByKeyValueString(tx *sqlx.Tx, data map[string]interface{}, key, value string) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	keysWithQuestionMarks := make([]string, 0)
	values := make([]interface{}, 0)

	for key, value := range data {
		keysWithQuestionMark := fmt.Sprintf("%v=?", key)
		keysWithQuestionMarks = append(keysWithQuestionMarks, keysWithQuestionMark)
		values = append(values, value)
	}

	// Add value as part of values
	values = append(values, value)

	query := fmt.Sprintf(
		"UPDATE %v SET %v WHERE %v=?",
		b.table,
		strings.Join(keysWithQuestionMarks, ","),
		key)

	result, err = tx.Exec(query, values...)

	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	return result, err
}

func (b *Base) DeleteFromTable(tx *sqlx.Tx, where string) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("DELETE FROM %v", b.table)

	if where != "" {
		query = query + " WHERE " + where
	}

	result, err = tx.Exec(query)

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	if err != nil {
		return nil, err
	}

	return result, err
}

func (b *Base) DeleteById(tx *sqlx.Tx, id int64) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("DELETE FROM %v WHERE id=?", b.table)

	result, err = tx.Exec(query, id)

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	if err != nil {
		return nil, err
	}

	return result, err
}
