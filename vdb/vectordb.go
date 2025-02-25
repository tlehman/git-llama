// Package vdb implements a local vector database using SQLite and the
// sqlite-vec extension. This package uses the github.com/ollama/ollama/api
// module directly, but when creating a [VectorDb], the model name must be
// specified.

package vdb

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// VectorDatabase represents a file-backed SQLite db with sqlite-vec extension that
// will store all the embeddings for the git repo
type VectorDatabase struct {
	filename  string
	modelname string
	dimension int
}

// Vector is a wrapper around a slice of float64s, this enables vector addition with
// methods like v1.Add(v2)
type Vector struct {
	values []float64 // len(values) is the dimensions
}

// Open will attempt to find the SQLite db at filename, and if that fails, then create it,
// and if the creation fails, it will return an error
func Open(filename string, modelname string) (*VectorDatabase, error) {
	// Check if the database file already exists
	_, err := os.Stat(filename)
	if err == nil {
		// If it exists, attempt to open and read from it
		vecdb, err := openFromExistingFile(filename, modelname)
		if err != nil {
			return nil, fmt.Errorf("failed to open existing db: %s\n", err)
		}
		return vecdb, nil
	}

	// If not found, create the SQLite database using the go-sqlite3 library
	vecdb, err := createNewDatabase(filename, modelname)
	if err != nil {
		return nil, fmt.Errorf("failed to create new db: %s", err)
	}
	return vecdb, nil
}

func openFromExistingFile(filename string, modelname string) (*VectorDatabase, error) {
	// Use the go-sqlite3 library to open and read from the existing database file
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	//var vecdb VectorDatabase, then .Scan(&vecdb.dimension)
	var version string
	err = db.QueryRow(`SELECT sqlite_version();`).Scan(&version)
	if err != nil {
		return nil, fmt.Errorf("failed to get sqlite version: %s\n", err)
	}
	fmt.Printf("sqlite version = %s\n", version)

	return nil, nil
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

	// Open a transaction.
	err = db.Ping()
	if err != nil {
		return nil, err
	}

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

func (vectordb *VectorDatabase) Get(id string) *Vector {
	return nil
}

func (vectordb *VectorDatabase) Insert(id string, input string) error {
	return nil
}

func (vectordb *VectorDatabase) Update(id string, input string) error {
	return nil
}
