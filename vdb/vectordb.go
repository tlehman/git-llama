// Package vdb implements a local vector database using SQLite and the
// sqlite-vec extension. This package uses the github.com/ollama/ollama/api
// module directly, but when creating a [VectorDb], the model name must be
// specified.

package vdb

import (
	_ "embed"
	"fmt"

	_ "github.com/asg017/sqlite-vec-go-bindings/ncruces"
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
)

// VectorDatabase represents a file-backed SQLite db with sqlite-vec extension that
// will store all the embeddings for the git repo. The modelname is stored because it
// is necessary know since it is necessary to work with the vectors later. Different models
// have a vectors of different dimension.
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
	var vecdb *VectorDatabase = &VectorDatabase{
		filename:  filename,
		modelname: modelname,
	}
	// open the SQL DB on the VectorDatabase
	db, err := sqlite3.Open("/Users/tobi/.git-llama.db")
	if err != nil {
		fmt.Errorf("failed opening sqlite3 in memory mode: %s\n", err)
		return nil, err
	}
	vecdb.DB = db

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
