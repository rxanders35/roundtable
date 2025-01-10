package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

func PathTransformer(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashstr := hex.EncodeToString(hash[:])
	blocksize := 5
	slicelen := len(hashstr) / blocksize
	paths := make([]string, slicelen)
	for i := 0; i < slicelen; i++ {
		from, to := i*blocksize, (i*blocksize + blocksize)
		paths[i] = hashstr[from:to]
	}

	return PathKey{
		Path:     strings.Join(paths, "/"),
		Original: hashstr,
	}
}

type Transform func(string) PathKey

var DefaultTransform = func(key string) string {
	return key
}

type PathKey struct {
	Path     string
	Original string
}

type StorageOpts struct {
	Transform Transform
}

type Storage struct {
	StorageOpts
}

func NewStorage(opts StorageOpts) *Storage {
	return &Storage{
		StorageOpts: opts,
	}
}

func (s *Storage) writeStream(key string, r io.Reader) error {
	path := s.Transform(key)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	fileBytes := md5.Sum(buf.Bytes())
	file := hex.EncodeToString(fileBytes[:])
	fullpath := path + "/" + file

	f, err := os.Create(fullpath)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, buf)
	if err != nil {
		return nil
	}

	log.Printf("Wrote %d bytes to disk: %s", n, fullpath)

	return nil
}
