package main

import "os"
import "log"
import "encoding/gob"
import "io"
import "errors"

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
	f, err := os.OpenFile(s.fname, os.O_RDONLY, 0644)
	defer f.Close()

	if err != nil {
		return err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	d := gob.NewDecoder(f)
	err = nil
	for err == nil {
		var r collectionRecord
		if err = d.Decode(&r); err == nil {
			s.collections[r.Key] = r.Col
			s.collections[r.Key].sub = NewSubscription()
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func (s *CollectionStore) saveCollections() error {
	f, err := os.OpenFile(s.fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("WARNING: Unable to open store file. Some data may be lost.", err)
		return err
	}
	defer f.Close()
	e := gob.NewEncoder(f)
	for k, v := range s.collections {
		err := e.Encode(collectionRecord{Key: k, Col: v})
		if err != nil {
			log.Println("WARNING: Failed to write a collection to the store. Some data may be lost.", err)
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
