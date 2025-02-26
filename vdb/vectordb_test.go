package vdb

import (
	"os"
	"testing"
)

const dbfilename = "/tmp/git-llama-test.db"
const modelname = "llama3.2"

var vectordb *VectorDatabase

// create test vector database
func setup() {
	vectordb, _ = Open(dbfilename, modelname)
}

// destroy test vector database
func teardown() {
	vectordb.Close()
	os.Remove(dbfilename)
}

func TestInsert(t *testing.T) {
	setup()
	defer teardown()

	err := vectordb.Insert("foo", "bar")
	if err != nil {
		t.Fatalf("failed inserting into db: %s", err)
	}
	fooVector := vectordb.Get("foo")
	if fooVector == nil {
		t.Fatal("vector was supposed to be non-nil")
	}
}
