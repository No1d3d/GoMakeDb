# GoMakeDB
A utility written in Go that reads through a special MakeDB file and Creates a DB described in this file

## GoMakeDB Specification
First line must be:

*#!GoMakeDB*

### Comments:

Lines starting with # after the first one are treated as comments and ignored.

Empty lines are ignored.

### Database setup:

The first two meaningful lines (non-comment, non-empty):

Relative path to where the database will be stored (example: ./database)

Name of the database file (example: example.db)

### Table creation:

A table definition starts with:

*---*

Then:

Table name (single line)

Each column: one line per column, with column name and column type separated by space.

### Multiple tables:

Each new table starts again with ---.

### Example
```bash
#!GoMakeDB

./database
example.db
---
person
id int NOT NULL
lastName varchar(255) NOT NULL
firstName varchar(255)
age int
---
socials
id int NOT NULL
person_id int NOT NULL
contact_info varchar(255) NOT NULL
status int NOT NULL CHECK (status IN (0, 1))
```

### Running

To run the application, create a file with a description of your database
!!!(Do not forget to start it with #!GoMakeDB)!!!

After that, run

```bash
go build -o gomakedb main.go
go run gomakedb YourDbDescriptionFile.txt
```