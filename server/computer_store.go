package main

import "os"
import "io"
import "encoding/gob"
import "log"
import "math/rand"
import "github.com/tortis/mechadmin/types"
import "errors"

type ComputerStore struct {
	computers map[string]*Computer
	fname     string
}

type computerRecord struct {
	Key  string
	Comp *Computer
}

func NewComputerStore(filename string) *ComputerStore {
	cs := &ComputerStore{computers: make(map[string]*Computer, 0), fname: filename}
	if err := cs.loadComputers(); err != nil {
		log.Println("ComputerStore:", err)
	}
	return cs
}

func (s *ComputerStore) loadComputers() error {
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
		var r computerRecord
		if err = d.Decode(&r); err == nil {
			s.computers[r.Key] = r.Comp
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func (s *ComputerStore) saveComputers() error {
	f, err := os.OpenFile(s.fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("WARNING: Unable to save computers to store. Some data may have been lost.", err)
		return err
	}
	defer f.Close()
	e := gob.NewEncoder(f)
	for k, v := range s.computers {
		err := e.Encode(computerRecord{Key: k, Comp: v})
		if err != nil {
			log.Println("Failed to save a computer record to the store. Some data may havve been lost.", err)
		}
	}
	return nil
}

func (s *ComputerStore) Get(MAC string) *Computer {
	return s.computers[MAC]
}

func (s *ComputerStore) Put(comp *Computer) error {
	MAC := comp.Info.MAC
	if _, present := s.computers[MAC]; present {
		return errors.New("A computer with that key already exists. Computer was not added to store.")
	}
	s.computers[MAC] = comp
	return nil
}

func (s *ComputerStore) Delete(MAC string) {
	//Also delete computer from associated collections.
	delete(s.computers, MAC)
}

func (s *ComputerStore) UpdateOrAddComputer(stat *types.Status) *Computer {
	_, present := s.computers[stat.MAC]
	if present == true {
		s.computers[stat.MAC].Info = *stat
		return s.computers[stat.MAC]
	} else {
		c := &Computer{
			Info: *stat,
			UID:  rand.Uint32(),
		}
		s.Put(c)
		ColStore.Get(ALL_SYS_COL).AddComputer(stat.MAC)
		return c
	}
}
