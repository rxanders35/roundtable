package main

import (
	"bytes"
	"testing"
)

func TestTransform(t *testing.T) {
	key := "bestpictures"
	pathkey := PathTransformer(key)
	expectedOG := "b2f8f1dd50fdeec113ac1b5066d1d3a10f70f1dc"
	expectedPath := "b2f8f/1dd50/fdeec/113ac/1b506/6d1d3/a10f7/0f1dc"
	if pathkey.Path != expectedPath {
		t.Errorf("have:\n %s, \nwant: %s", pathkey.Path, expectedPath)
	}
}

func TestStorage(t *testing.T) {
	opts := StorageOpts{
		Transform: PathTransformer,
	}
	s := NewStorage(opts)
	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("filename", data); err != nil {
		t.Error(err)
	}
}
