// Package vdb implements a local vector database using SQLite and the
// sqlite-vec extension. This package uses the github.com/ollama/ollama/api
// module directly, but when creating a [VectorDb], the model name must be
// specified.

package vdb

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"

	_ "github.com/asg017/sqlite-vec-go-bindings/ncruces"
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
)

// VectorDatabase represents a file-backed SQLite db with sqlite-vec extension that
// will store all the embeddings for the git repo
type VectorDatabase struct {
	filename  string
	modelname string
	dimension int
	DB        *sqlite3.Conn
}

// Vector is a wrapper around a slice of float64s, this enables vector addition with
// methods like v1.Add(v2)
type Vector struct {
	values []float64 // len(values) is the dimensions
}

// Open will attempt to find the SQLite db at filename, and if that fails, then create it,
// and if the creation fails, it will return an error
func Open(filename string, modelname string) (*VectorDatabase, error) {
	var vecdb *VectorDatabase
	// Check if the database file already exists
	_, err := os.Stat(filename)
	if err == nil {
		// If it exists, attempt to open and read from it
		vecdb, err = openFromExistingFile(filename, modelname)
		if err != nil {
			return nil, fmt.Errorf("failed to open existing db: %s\n", err)
		}
	} else {
		// Otherwise, create the SQLite database using the go-sqlite3 library
		vecdb, err = createNewDatabase(filename, modelname)
		if err != nil {
			return nil, fmt.Errorf("failed to create new db: %s", err)
		}
	}

	// open the SQL DB on the VectorDatabase
	vecdb.DB, err = sqlite3.Open(":memory:")
	if err != nil {
		fmt.Printf("failed opening sqlite3 in memory mode: %s\n", err)
		return nil, err
	}

	// check the sqlite_version and the vec_version
	stmt, _, err := vecdb.DB.Prepare(`SELECT sqlite_version(), vec_version()`)
	if err != nil {
		fmt.Printf("failed getting vec_version(): %s\n", err)
		return nil, err
	}
	stmt.Step()
	fmt.Printf("sqlite_version() = %s, vec_version() = %s\n", stmt.ColumnText(0), stmt.ColumnText(1))

	return vecdb, nil
}

func (vectordb *VectorDatabase) Get(id string) *Vector {
	var vec Vector
	// query the vector for that id
	return &vec
}

func (vectordb *VectorDatabase) Insert(id string, value string) error {
	return nil
}

func (vectordb *VectorDatabase) Update(id string, input string) error {
	return nil
}

func (vectordb *VectorDatabase) Close() error {
	return vectordb.DB.Close()
}

func openFromExistingFile(filename string, modelname string) (*VectorDatabase, error) {
	var vecdb VectorDatabase

	// Use the go-sqlite3 library to open and read from the existing database file
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var version string
	err = db.QueryRow(`SELECT sqlite_version();`).Scan(&version)
	if err != nil {
		return nil, fmt.Errorf("failed to get sqlite version: %s\n", err)
	}
	fmt.Printf("sqlite version = %s\n", version)

	return &vecdb, nil
}

func createNewDatabase(filename string, modelname string) (*VectorDatabase, error) {
	// Use the go-sqlite3 library to create a new database file and add it as an attachment.
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", filename))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Open a transaction.
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create the git_embeddings table

	// Insert the new model info into our table
	//_, err = db.Exec(`INSERT INTO data VALUES (''::TEXT, $modelname, $dimension, 0)`)

	// Commit the changes (db.Commit was not found)

	if err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	defaultdim := 768

	return &VectorDatabase{
		filename:  filename,
		modelname: modelname,
		dimension: defaultdim,
	}, nil
}
