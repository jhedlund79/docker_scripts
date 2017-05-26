package sql_example


import (
"database/sql"
_ "github.com/denisenkom/go-mssqldb"
"strconv"
)

// PingServer uses a passed database handle to check if the database server works
func PingServer(db *sql.DB) string {

	err := db.Ping()
	if err != nil {
		return ("From Ping() Attempt: " + err.Error())
	}

	return ("Database Ping Worked...")

}

// CheckDB checks if the database "strDBName" exists on the MSSQL database engine.
func CheckDB(db *sql.DB, strDBName string) (bool, error) {

	// Does the database exist?
	result, err := db.Query("SELECT db_id('" + strDBName + "')")
	defer result.Close()
	if err != nil {
		return false, err
	}

	for result.Next() {
		var s sql.NullString
		err := result.Scan(&s)
		if err != nil {
			return false, err
		}

		// Check result
		if s.Valid {
			return true, nil
		} else {
			return false, nil
		}
	}

	// This return() should never be hit...
	return false, err
}

// CreateDBAndTable creates a new content database on the SQL Server along with
// the necessary tables. Keep in mind the user credentials that opened the database
// connection with sql.Open must have at least dbcreator rights to the database. The
// table (testtable) will have columns source (nvarchar), timestamp (bigint), and
// content (nvarchar).
func CreateDBAndTable(db *sql.DB, strDBName string) error {

	// Create the database
	_, err := db.Exec("CREATE DATABASE [" + strDBName + "]")
	if err != nil {
		return (err)
	}

	// Let's turn off AutoClose
	_, err = db.Exec("ALTER DATABASE [" + strDBName + "] SET AUTO_CLOSE OFF;")
	if err != nil {
		return (err)
	}

	// Create the tables
	_, err = db.Exec("USE " + strDBName + "; CREATE TABLE testtable (source nvarchar(100) NOT NULL, timestamp bigint NOT NULL, content nvarchar(4000) NOT NULL)")
	if err != nil {
		return (err)
	}

	return nil

}

// DropDB deletes the database strDBName.
func DropDB(db *sql.DB, strDBName string) error {

	// Drop the database
	_, err := db.Exec("DROP DATABASE [" + strDBName + "]")

	if err != nil {
		return err
	}

	return nil

}

// AddToContent adds new content to the database.
func AddToContent(db *sql.DB, strDBName string, strSource string, int64Timestamp int64, strContent string) error {

	// Add a record entry
	_, err := db.Exec("USE " + strDBName + "; INSERT INTO testtable (source, timestamp, content) VALUES ('" + strSource + "','" + strconv.FormatInt(int64Timestamp, 10) + "','" + strContent + "');")
	if err != nil {
		return err
	}

	return nil

}

// RemoveFromContentBySource removes a record from the database with source strSource. The
// int64 returned is a message indicating the number of rows affected.
func RemoveFromContentBySource(db *sql.DB, strSource string) (int64, error) {

	// Remove entries containing the source...
	result, err := db.Exec("DELETE FROM testtable WHERE source=$1;", strSource)
	if err != nil {
		return 0, err
	}

	// What was the result?
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil

}

// Query the content in the database and return the source (string), timestamp (int64), and
// content (string) as slices
func GetContent(db *sql.DB) ([]string, []int64, []string, error) {

	var slcstrContent []string
	var slcint64Timestamp []int64
	var slcstrSource []string

	// Run the query
	rows, err := db.Query("SELECT source, timestamp, content FROM testtable")
	if err != nil {
		return slcstrSource, slcint64Timestamp, slcstrContent, err
	}
	defer rows.Close()

	for rows.Next() {

		// Holding variables for the content in the columns
		var source, content string
		var timestamp int64

		// Get the results of the query
		err := rows.Scan(&source, &timestamp, &content)
		if err != nil {
			return slcstrSource, slcint64Timestamp, slcstrContent, err
		}

		// Append them into the slices that will eventually be returned to the caller
		slcstrSource = append(slcstrSource, source)
		slcstrContent = append(slcstrContent, content)
		slcint64Timestamp = append(slcint64Timestamp, timestamp)
	}

	return slcstrSource, slcint64Timestamp, slcstrContent, nil

}