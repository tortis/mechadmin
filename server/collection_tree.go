package main

import "math/rand"
import "strconv"
import "io/ioutil"
import "encoding/json"
import "os"
import "os/signal"
import "log"

type CollectionTree struct {
	Name         string
	UID          uint32
	Children     []*CollectionTree // For OUs
	IsSystem     bool
	IsCollection bool
}

var root *CollectionTree = &CollectionTree{}

func init() {
	ColStore = NewCollectionStore("colstore.gob")

	// Load root tree from file or create a new one.
	if err := loadRootTree(); err != nil {
		log.Fatal("CollectionTree:", err)
	}

	// Register to recieve the interrup signal so that
	// the root tree can be saved before the program exits.
	// The root tree should still be saved at a regular interval.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if saveRootTree() != nil {
			println("\nError! The root collection tree could not be saved before the program stopped. Some data may have been lost.")
			os.Exit(1)
		} else {
			println("\nSuccessfully saved root collection tree.")
			os.Exit(0)
		}
	}()
}

func loadRootTree() error {
	jsonByte, err := ioutil.ReadFile("root.json")
	if err != nil {
		root = NewCollectionTree()
		println("root.json was not found. Creating new collection tree.")
		return nil
	} else {
		println("root.json found. Loading collection tree.")
		err = json.Unmarshal(jsonByte, root)
		return err
	}
}

func saveRootTree() error {
	b, err := json.Marshal(root)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("root.json", b, 0644)
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
