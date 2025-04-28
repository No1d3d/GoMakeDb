package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run gomakedb <path_to_db_descriptor_txt_file>")
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Validate first line
	if !scanner.Scan() || scanner.Text() != "#!GoMakeDB" {
		log.Fatal("First line must be #!GoMakeDB")
	}

	// Read database path and name
	var dbDir, dbName string

	for scanner.Scan() {
		line := cleanLine(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		dbDir = line
		break
	}

	for scanner.Scan() {
		line := cleanLine(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		dbName = line
		break
	}

	if dbDir == "" || dbName == "" {
		log.Fatal("Database path or name is missing")
	}

	// Ensure directory exists
	err = os.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	dbPath := filepath.Join(dbDir, dbName)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	fmt.Println("Database created:", dbPath)

	var tableName string
	var columns []string

	// Parse tables
	for scanner.Scan() {
		line := cleanLine(scanner.Text())
		if line == "" {
			continue
		}

		if line == "---" {
			if tableName != "" {
				createTable(db, tableName, columns)
			}
			// Start a new table
			tableName = ""
			columns = []string{}
			continue
		}

		if tableName == "" {
			tableName = line
		} else {
			columns = append(columns, line)
		}
	}

	// Create last table if exists
	if tableName != "" {
		createTable(db, tableName, columns)
	}

	fmt.Println("Database setup completed successfully.")
}

func cleanLine(line string) string {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "#") {
		return ""
	}
	return line
}

func createTable(db *sql.DB, tableName string, columns []string) {
	if len(columns) == 0 {
		log.Fatalf("No columns specified for table %s", tableName)
	}

	columnDefs := append([]string{}, columns...)

	createStmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, strings.Join(columnDefs, ", "))
	_, err := db.Exec(createStmt)
	if err != nil {
		log.Fatalf("Failed to create table %s: %v", tableName, err)
	}

	fmt.Println("Created table:", tableName)
}
