// Package vdb implements a local vector database using SQLite and the
// sqlite-vec extension. This package uses the github.com/ollama/ollama/api
// module directly, but when creating a [VectorDb], the model name must be
// specified.

package vdb

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
	return nil, nil
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
