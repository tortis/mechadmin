package main

import "os"
import "io"
import "encoding/gob"
import "log"
import "math/rand"
import "github.com/tortis/mechadmin/types"

type ComputerStore struct {
	computers map[string]*Computer
	file      *os.File
}

type computerRecord struct {
	key  string
	comp *Computer
}

func NewComputerStore(filename string) *ComputerStore {
	cs := &ComputerStore{computers: make(map[string]*Computer, 0)}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("CollectionStore:", err)
	}
	cs.file = f
	if err := cs.loadComputers(); err != nil {
		log.Println("ComputerStore:", err)
	}
	return cs
}

func (s *ComputerStore) loadComputers() error {
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}
	d := gob.NewDecoder(s.file)
	var err error
	for err == nil {
		var r computerRecord
		if err = d.Decode(&r); err == nil {
			s.computers[r.key] = r.comp
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func (s *ComputerStore) saveComputer(name string, comp *Computer) error {
	e := gob.NewEncoder(s.file)
	return e.Encode(computerRecord{name, comp})
}

func (s *ComputerStore) Get(name string) *Computer {
	return s.computers[name]
}

func (s *ComputerStore) Put(name string, comp *Computer) error {
	s.computers[name] = comp
	return s.saveComputer(name, comp)
}

func (s *ComputerStore) Delete(name string) {
	//Also delete computer from associated collections.
	delete(s.computers, name)
}

func (s *ComputerStore) UpdateOrAddComputer(stat *types.Status) *Computer {
	_, present := s.computers[stat.CN]
	if present == true {
		s.computers[stat.CN].Info = *stat
		return s.computers[stat.CN]
	} else {
		c := &Computer{
			Info: *stat,
			UID:  rand.Uint32(),
		}
		s.Put(stat.CN, c)
		ColStore.Get(ALL_SYS_COL).AddComputer(stat.CN)
		return c
	}
}
