package main

import "os"
import "log"
import "encoding/gob"
import "io"
import "errors"
import "strconv"

type CollectionStore struct {
	collections map[uint32]*Collection
	fname       string
}

func NewCollectionStore(filename string) *CollectionStore {
	s := &CollectionStore{collections: make(map[uint32]*Collection), fname: filename}
	if err := s.loadCollections(); err != nil {
		log.Println("CollectionStore-loadCollections:", err)
	}
	return s
}

type collectionRecord struct {
	Key uint32
	Col *Collection
}

func (s *CollectionStore) loadCollections() error {
	f, err := os.OpenFile(s.fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("CollectionStore-OpenFile:", err)
	}

	println("Loading collections from store.")
	defer f.Close()
	if _, err := f.Seek(0, 0); err != nil {
		println("Failed to seek 0,0 in file.")
		return err
	}
	d := gob.NewDecoder(f)
	err = nil
	count := 0
	for err == nil {
		var r collectionRecord
		if err = d.Decode(&r); err == nil {
			s.collections[r.Key] = r.Col
			count++
		}
	}
	println("Loaded " + strconv.Itoa(count) + " records.")
	if err == io.EOF {
		return nil
	}
	println("There was a problem loading records from the store: " + err.Error())
	return err
}

func (s *CollectionStore) saveCollections() error {
	f, err := os.OpenFile(s.fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("CollectionStore-saveCollections-OpenFile", err)
	}
	defer f.Close()
	e := gob.NewEncoder(f)
	for k, v := range s.collections {
		err := e.Encode(collectionRecord{Key: k, Col: v})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *CollectionStore) Get(key uint32) *Collection {
	return s.collections[key]
}

func (s *CollectionStore) Put(key uint32, col *Collection) error {
	if _, present := s.collections[key]; present {
		return errors.New("A collection with that key already exists. Collection was not added to store.")
	}
	s.collections[key] = col
	return nil
}
