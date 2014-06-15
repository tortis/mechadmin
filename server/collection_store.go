package main

import "os"
import "log"
import "encoding/gob"
import "io"

/* System Collection UIDs. */
const ALL_SYS_COL uint32 = 0
const UNK_SYS_COL uint32 = 1
const OFF_SYS_COL uint32 = 2

type CollectionStore struct {
	collections	map[uint32]*Collection
	file		*os.File
}

func NewCollectionStore(filename string) *CollectionStore {
	s := &CollectionStore {collections: make(map[uint32]*Collection)}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if (err != nil) {
		log.Fatal("CollectionStore:", err)
	}
	s.file = f
	if err := s.loadCollections(); err != nil {
		log.Println("CollectionStore:", err)
	}
	return s
}

type collectionRecord struct {
	key uint32
	col *Collection
}

func (s *CollectionStore) loadCollections() error {
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}
	d := gob.NewDecoder(s.file)
	var err error
	for err == nil {
		var r collectionRecord
		if err = d.Decode(&r); err ==nil {
			s.collections[r.key] = r.col
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func (s *CollectionStore) saveCollection(key uint32, col *Collection) error {
	e := gob.NewEncoder(s.file)
	return e.Encode(collectionRecord{key, col})
}

func (s *CollectionStore) Get(key uint32) *Collection {
	return s.collections[key]
}

func (s *CollectionStore) Put(key uint32, col *Collection) {
	s.collections[key] = col
	s.saveCollection(key, col)
}
