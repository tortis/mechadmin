package main

import "math/rand"
import "strconv"
import "io/ioutil"
import "encoding/json"

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
		Children:     make([]*CollectionTree, 0, 10),
		IsSystem:     true,
		IsCollection: false,
	}
	/* Add the system collections to the root OU. */
	allSys := &CollectionTree{
		Name:         `All Systems`,
		UID:          ALL_SYS_COL,
		IsSystem:     true,
		IsCollection: true,
	}
	ColStore.Put(ALL_SYS_COL, NewCollection(allSys.Name, ALL_SYS_COL))
	t.Children = append(t.Children, allSys)

	unkSys := &CollectionTree{
		Name:         `Unknown Systems`,
		UID:          UNK_SYS_COL,
		IsSystem:     true,
		IsCollection: true,
	}
	ColStore.Put(UNK_SYS_COL, NewCollection(unkSys.Name, UNK_SYS_COL))
	t.Children = append(t.Children, unkSys)

	offSys := &CollectionTree{
		Name:         `Offline Systems`,
		UID:          OFF_SYS_COL,
		IsSystem:     true,
		IsCollection: true,
	}
	ColStore.Put(OFF_SYS_COL, NewCollection(offSys.Name, OFF_SYS_COL))
	t.Children = append(t.Children, offSys)

	return t
}

func LoadTree(filename string) (*CollectionTree, error) {
	jsonByte, err := ioutil.ReadFile(filename)
	if err != nil {
		t := NewCollectionTree()
		println(filename+" was not found. Creating new collection tree.")
		return t, nil
	} else {
		println(filename+" found. Loading collection tree.")
		err = json.Unmarshal(jsonByte, root)
		return nil, err
	}
}

func (t *CollectionTree) Save(filename string) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
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
