package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"

	. "github.com/mssql_example/sql_example"
)

const strVERSION string = "0.18 compiled on 8/11/2015"

// sqltest is a small application for demonstrating/testing/learning about SQL database connectivity from Go
func main() {

	// Flags
	ptrVersion := flag.Bool("version", false, "Display program version")
	ptrDeleteIt := flag.Bool("deletedb", false, "Delete the database")
	ptrServer := flag.String("server", "localhost", "Server to connect to")
	ptrUser := flag.String("username", "sa", "Username for authenticating to database; if you use a backslash, it must be escaped or in quotes")
	ptrPass := flag.String("password", "pass@word1", "Password for database connection")
	ptrDBName := flag.String("dbname", "test_db", "Database name")

	flag.Parse()

	// Does the user just want the version of the application?
	if *ptrVersion == true {
		fmt.Println("Version " + strVERSION)
		os.Exit(0)
	}

	// Open connection to the database server; this doesn't verify anything until you
	// perform an operation (such as a ping).
	db, err := sql.Open("mssql", "server="+*ptrServer+";user id="+*ptrUser+";password="+*ptrPass)
	if err != nil {
		fmt.Println("From Open() attempt: " + err.Error())
	}

	// When main() is done, this should close the connections
	defer db.Close()

	// Does the user want to delete the database?
	if *ptrDeleteIt == true {
		boolDBExist, err := CheckDB(db, *ptrDBName)
		if err != nil {
			fmt.Println("Error running CheckDB: " + err.Error())
			os.Exit(1)
		}
		if boolDBExist {
			fmt.Println("(sqltest) Deleting database " + *ptrDBName + "...")
			DropDB(db, *ptrDBName)
			os.Exit(0)
		} else {

			// Database doesn't seem to exist...
			fmt.Println("(sqltest) Database " + *ptrDBName + " doesn't appear to exist...")
			os.Exit(1)

		}
	}

	// Let's start the tests...
	fmt.Println("********************************")

	// Is the database running?
	strResult := PingServer(db)
	fmt.Println("(sqltest) Ping of Server Result Was: " + strResult)

	fmt.Println("********************************")

	// Does the database exist?
	boolDBExist, err := CheckDB(db, *ptrDBName)
	if err != nil {
		fmt.Println("(sqltest) Error running CheckDB: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("(sqltest) Database Existence Check: " + strconv.FormatBool(boolDBExist))

	fmt.Println("********************************")

	// If it doesn't exist, let's create the base database
	if !boolDBExist {

		CreateDBAndTable(db, *ptrDBName)
		fmt.Println("********************************")

	}

	// Enter a test record
	boolDBExist, err = CheckDB(db, *ptrDBName)
	if err != nil {
		fmt.Println("(sqltest) CheckDB() error: " + err.Error())
		os.Exit(1)
	}

	if boolDBExist == true {

		err := AddToContent(db, *ptrDBName, "Bob", 1437506592, "Hello!")
		if err != nil {
			fmt.Println("(sqltest) Error adding line to content: " + err.Error())
			os.Exit(1)
		}

		err = AddToContent(db, *ptrDBName, "user", 1437506648, "Now testing memory")
		if err != nil {
			fmt.Println("(sqltest) Error adding line to content: " + err.Error())
			os.Exit(1)
		}

		err = AddToContent(db, *ptrDBName, "user", 1437503394, "test, text!")
		if err != nil {
			fmt.Println("(sqltest) Error adding line to content: " + err.Error())
			os.Exit(1)
		}

		err = AddToContent(db, *ptrDBName, "Bob", 1437506592, "Hope this works!")
		if err != nil {
			fmt.Println("(sqltest) Error adding line to content: " + err.Error())
			os.Exit(1)
		}

	}

	fmt.Println("(sqltest) Completed entering test records.")

	fmt.Println("********************************")

	fmt.Println("(sqltest) Deleting records from a particular source.")

	// Delete from a source
	int64Deleted, err := RemoveFromContentBySource(db, "user")
	if err != nil {
		fmt.Println("(sqltest) Error deleting records by source: " + err.Error())
		os.Exit(1)
	} else {

		// How many records were removed?
		fmt.Println("Removed " + strconv.FormatInt(int64Deleted, 10) + " records")
		fmt.Println("********************************")

	}

	// Get the content
	slcstrSource, slcint64Timestamp, slcstrContent, err := GetContent(db)
	if err != nil {
		fmt.Println("(sqltest) Error getting content: " + err.Error())
	}

	// Now read the contents
	for i := range slcstrContent {

		fmt.Println("Entry " + strconv.Itoa(i) + ": " + strconv.FormatInt(slcint64Timestamp[i], 10) + ", from " + slcstrSource[i] + ": " + slcstrContent[i])

	}

}
