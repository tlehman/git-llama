package vdb

import (
	"os"
	"testing"
)

const dbfilename = "/tmp/git-llama-test.db"
const modelname = "llama3.2"

var vectordb *VectorDatabase

const dim = 3

// create test vector database
func setup() {
	vectordb, _ = Open(dbfilename, modelname)
	vectordb.CreateTableIdempotent(dim)
}

// destroy test vector database
func teardown() {
	vectordb.Close()
	os.Remove(dbfilename)
}

func TestInsert(t *testing.T) {
	setup()
	defer teardown()

	v := &Vector{Values: []float32{0, 1, -1}}
	err := vectordb.Insert("foo", v)
	if err != nil {
		t.Fatalf("failed inserting into db: %s", err)
	}
	fooVector := vectordb.Get("foo")
	if fooVector == nil {
		t.Fatal("vector was supposed to be non-nil")
	}
}

func TestGet(t *testing.T) {
	setup()
	defer teardown()
	v := &Vector{Values: []float32{0, 1, -1}}
	_ = vectordb.Insert("foo", v)
	u := vectordb.Get("foo")

	// TODO verify that u == v
	if !v.Equals(u) {
		t.Fatalf("u = %v\nv = %v\n", u, v)
	}
}
