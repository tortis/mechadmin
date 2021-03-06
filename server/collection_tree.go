package main

import "math/rand"
import "strconv"
import "log"
import "encoding/gob"
import "os"

type CollectionTree struct {
	Name         string
	UID          uint32
	Children     []*CollectionTree // For OUs
	IsSystem     bool
	IsCollection bool
}

func NewCollectionTree() *CollectionTree {
	t := &CollectionTree{
		Name:         `Collections`,
		UID:          rand.Uint32(),
		Children:     make([]*CollectionTree, 0, 0),
		IsSystem:     true,
		IsCollection: false,
	}
	/* Add the system collections to the root OU. */
	allSys := &CollectionTree{
		Name:         `All Systems`,
		UID:          ALL_SYS_COL,
		Children:     make([]*CollectionTree, 0, 0),
		IsSystem:     true,
		IsCollection: true,
	}
	err := ColStore.Put(ALL_SYS_COL, NewCollection(allSys.Name, ALL_SYS_COL))
	if err == nil {
		t.Children = append(t.Children, allSys)
	}

	unkSys := &CollectionTree{
		Name:         `Unknown Systems`,
		UID:          UNK_SYS_COL,
		Children:     make([]*CollectionTree, 0, 0),
		IsSystem:     true,
		IsCollection: true,
	}
	err = ColStore.Put(UNK_SYS_COL, NewCollection(unkSys.Name, UNK_SYS_COL))
	if err == nil {
		t.Children = append(t.Children, unkSys)
	}

	offSys := &CollectionTree{
		Name:         `Offline Systems`,
		UID:          OFF_SYS_COL,
		Children:     make([]*CollectionTree, 0, 0),
		IsSystem:     true,
		IsCollection: true,
	}
	err = ColStore.Put(OFF_SYS_COL, NewCollection(offSys.Name, OFF_SYS_COL))
	if err == nil {
		t.Children = append(t.Children, offSys)
	}
	return t
}

func LoadTree(filename string) *CollectionTree {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Println("Could not open the collection tree file. Creating new tree.")
		return NewCollectionTree()
	}
	defer f.Close()
	d := gob.NewDecoder(f)
	var r *CollectionTree
	err = d.Decode(&r)
	if err != nil {
		log.Println("Failed to read collection tree file. Creating new tree.")
		return NewCollectionTree()
	}
	return r
}

func (t *CollectionTree) Save(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("Could not save collection tree to disk. Some data may be lost.", err)
		return err
	}
	defer f.Close()
	e := gob.NewEncoder(f)
	return e.Encode(t)
}

func (t *CollectionTree) NewOU(name string) *CollectionTree {
	ou := &CollectionTree{
		Name:         name,
		UID:          rand.Uint32(),
		Children:     make([]*CollectionTree, 0, 10),
		IsSystem:     false,
		IsCollection: false,
	}
	t.Children = append(t.Children, ou)
	return ou
}

func (t *CollectionTree) NewCol(name string) *CollectionTree {
	col := &CollectionTree{
		Name:         name,
		UID:          rand.Uint32(),
		IsSystem:     false,
		IsCollection: true,
	}
	t.Children = append(t.Children, col)
	ColStore.Put(col.UID, NewCollection(col.Name, col.UID))
	return col
}

func (t *CollectionTree) generateHtmlHlp() string {
	var html string
	for _, v := range t.Children {
		if v.IsCollection {
			html += `<li class="file" onclick="requestCol(this);" id="` + strconv.FormatUint(uint64(v.UID), 10) + `"><a href="#">` + v.Name + `</a></li>`
		} else {
			html += `<li><label for="` + strconv.FormatUint(uint64(v.UID), 10) + `">` + v.Name + `</label> <input type="checkbox" id="` + strconv.FormatUint(uint64(v.UID), 10) + `" /><ol>`
			if v.Children != nil {
				html += v.generateHtmlHlp()
			}
			html += `</ol></li>`
		}
	}
	return html
}

func (t *CollectionTree) GenerateHTML() string {
	var html string
	html += `<ol class="tree">`
	html += t.generateHtmlHlp()
	html += `</ol>`
	return html
}
