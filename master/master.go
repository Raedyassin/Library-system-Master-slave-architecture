package main

import (
	"database/sql"
	"fmt"
	"log"

	"bufio"
	"net"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "library",
		AllowNativePasswords: true,
	} //root@127.0.0.1:3306

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end connect

	// Listen for incoming connections. create socket
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	log.Println("Listening on " + "localhost: 8080")

	for {
		// Listen for an incoming connection. is a synchronous function is block
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Handle connections in a new goroutine.
		go handleRequest(conn, db)
	}

}
func handleRequest(conn net.Conn, db *sql.DB) {
	///This line of code is used in Go to read a string from a network connection (conn) until it encounters a newline character ('\n')
	buffer, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Println("Network Issues.")
		conn.Close()
		return
	}

	// Process the command
	sqlCommand := strings.TrimSpace(string(buffer))
	parts := strings.Split(sqlCommand, " ")

	// switchCondition := ;
	switch strings.ToUpper(parts[0]) {

	case "INSERT":
		_, err := db.Exec(sqlCommand)
		if err != nil { //"INSERT INTO items (name, description) VALUES (?, ?)
			conn.Write([]byte("Error inserting data. or Insert Format error.\n"))
		} else {
			conn.Write([]byte("Insert successful.\n"))
		}

	case "SELECT":
		rows, err := db.Query(sqlCommand)
		if err != nil {
			conn.Write([]byte("Error reading data1. Or format error \n"))
			return
		}
		defer rows.Close()
		// return array of columns name
		columns, err := rows.Columns()
		if err != nil {
			log.Fatal(err)
		}

		// create 2 slice for read data from scan function
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		for rows.Next() {
			rowResponse := ""
			// Fetch values for the current row into the 'values' slice
			err := rows.Scan(valuePtrs...)
			if err != nil {
				conn.Write([]byte("Error reading data2.\n"))
				return
			}

			for i := 0; i < len(columns); i++ {
				// Convert value to string before concatenating
				stringValue := ""
				if byteValue, ok := values[i].([]byte); ok {
					stringValue = string(byteValue)
				} else {
					stringValue = fmt.Sprintf("%v", values[i])
				}
				rowResponse += fmt.Sprintf("%s: %s, ", columns[i], stringValue)
			}
			// Remove the trailing comma and space from the end
			rowResponse = strings.TrimSuffix(rowResponse, ", ")
			// Send the row response to the client
			conn.Write([]byte(rowResponse + "###"))
		}
		conn.Write([]byte("\n"))

	case "UPDATE":
		_, err := db.Exec(sqlCommand) // "UPDATE items SET name = ?, description = ? WHERE id = ?"
		if err != nil {
			conn.Write([]byte("Error updating data. Or format error\n"))
		} else {
			conn.Write([]byte("Update successful.\n"))
		}

	case "DELETE":
		_, err := db.Exec(sqlCommand) //"DELETE FROM items WHERE id = ?", id
		if err != nil {
			conn.Write([]byte("Error deleting data. Or format error\n"))
		} else {
			conn.Write([]byte("Delete successful.\n"))
		}

	default:
		conn.Write([]byte("Invalid command.\n"))
	}

	handleRequest(conn, db) // Loop back to handle more requests
}
