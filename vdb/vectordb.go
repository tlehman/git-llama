// Package vdb implements a local vector database using SQLite and the
// sqlite-vec extension. This package uses the github.com/ollama/ollama/api
// module directly, but when creating a [VectorDb], the model name must be
// specified.

package vdb

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"math"
	"regexp"
	"strings"

	_ "github.com/asg017/sqlite-vec-go-bindings/ncruces"
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

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

func (vectordb *VectorDatabase) TableName() string {
	cleanedModelName := nonAlphanumericRegex.ReplaceAllString(vectordb.modelname, "")
	return fmt.Sprintf("vec_%s", cleanedModelName)
}

// Vector is a wrapper around a slice of float64s, this enables vector addition with
// methods like v1.Add(v2)
type Vector struct {
	Values []float32 // len(values) is the dimensions
}

func (v *Vector) String() string {
	if len(v.Values) == 0 {
		return ""
	}

	// Convert each float32 to string
	strValues := make([]string, len(v.Values))
	for i, val := range v.Values {
		strValues[i] = fmt.Sprintf("%.3f", val) // %g gives compact float representation
	}

	// Join with commas (no extra space for SQLite compatibility)
	return strings.Join(strValues, ",")
}

func (v *Vector) Add(u *Vector) *Vector {
	if len(v.Values) != len(u.Values) || len(v.Values) == 0 {
		return nil
	}
	var values []float32 = make([]float32, len(v.Values))
	for i, _ := range v.Values {
		values[i] = v.Values[i] + u.Values[i]
	}
	return &Vector{Values: values}
}

func (v *Vector) Sub(u *Vector) *Vector {
	if len(v.Values) != len(u.Values) || len(v.Values) == 0 {
		return nil
	}
	var values []float32 = make([]float32, len(v.Values))
	for i, _ := range v.Values {
		values[i] = v.Values[i] - u.Values[i]
	}
	return &Vector{Values: values}
}

func (v *Vector) Norm() float32 {
	var norm float64
	for i, _ := range v.Values {
		vsquared := v.Values[i] * v.Values[i]
		if vsquared != 0.0 {
			fmt.Printf("v^2 = %v\n", vsquared)
		}
		norm += float64(vsquared)
	}
	return float32(math.Sqrt(norm))
}

func (v *Vector) Equals(u *Vector) bool {
	if len(v.Values) != len(u.Values) {
		return false
	}
	for i, vi := range v.Values {
		if vi != u.Values[i] {
			return false
		}
	}
	return true
}

// Open will attempt to find the SQLite db at filename, and if that fails, then create it,
// and if the creation fails, it will return an error
func Open(filename string, modelname string) (*VectorDatabase, error) {
	var vecdb *VectorDatabase = &VectorDatabase{
		filename:  filename,
		modelname: modelname,
	}
	// open the SQL DB on the VectorDatabase
	db, err := sqlite3.Open(filename)
	if err != nil {
		fmt.Errorf("failed opening sqlite3 in memory mode: %s\n", err)
		return nil, err
	}
	vecdb.DB = db

	return vecdb, nil
}

// CreateTableIdempotent takes the dimension, this was to decouple the vdb from the ollm package
func (vectordb *VectorDatabase) CreateTableIdempotent(dim int) error {
	vectordb.dimension = dim

	sql := fmt.Sprintf(
		"CREATE VIRTUAL TABLE IF NOT EXISTS %s USING vec0(id text UNIQUE, embedding float[%d]);",
		vectordb.TableName(),
		vectordb.dimension,
	)
	err := vectordb.DB.Exec(sql)
	return err
}

func (vectordb *VectorDatabase) Get(id string) *Vector {
	var vec Vector
	sql := fmt.Sprintf(
		"SELECT embedding FROM %s WHERE id = '%s';",
		vectordb.TableName(),
		id,
	)
	stmt, _, err := vectordb.DB.Prepare(sql)
	if err != nil {
		fmt.Printf("failed preparing SQL in Get(%s): %s", id, err)
		return nil
	}
	defer stmt.Close()
	if stmt.Step() {
		var values []float32
		columnRaw := stmt.ColumnRawBlob(0)
		if columnRaw == nil {
			fmt.Printf("no embedding found for id %s: ", id)
			return nil
		}
		// TODO make a Vector function that does this conversion and then put a unit test around it
		// float32 is 4 bytes, so the length of the return values must be len(columnRaw)/4
		values = make([]float32, len(columnRaw)/4)
		for i := 0; i < len(columnRaw)/4; i++ {
			bits := binary.LittleEndian.Uint32(columnRaw[i*4 : (i+1)*4])
			values[i] = math.Float32frombits(bits)
		}
		return &Vector{Values: values}
	}

	return &vec
}

func (vectordb *VectorDatabase) Insert(id string, embedding *Vector) error {
	tx := vectordb.DB.Begin()
	sql := fmt.Sprintf(
		"INSERT INTO %s(id, embedding) values ('%s', '[%s]');",
		vectordb.TableName(),
		id,
		embedding.String(),
	)
	err := vectordb.DB.Exec(sql)
	if err != nil {
		fmt.Printf("failed executing sql, rolling back tx: %s\n", err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (vectordb *VectorDatabase) Update(id string, input string) error {
	return nil
}

func (vectordb *VectorDatabase) Close() error {
	return vectordb.DB.Close()
}
