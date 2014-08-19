package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/steveyen/gkvlite"
	"reflect"
	"sync"
)

type Rank int

type User struct {
	Name string
	Rank
}

type Ladder []*User

func (l Ladder) Len() int      { return len(l) }
func (l Ladder) Swap(a, b int) { l[a], l[b] = l[b], l[a] }
func (l Ladder) Less(a, b int) bool {
	return l[a].Rank < l[b].Rank
}

type Gladder struct {
	*sync.Mutex
	*gkvlite.Store
}

func NewGladder(store *gkvlite.Store) *Gladder {
	return &Gladder{new(sync.Mutex), store}
}
func (g *Gladder) CreateUser(name string, initialRank Rank) error {
	g.Lock()
	defer g.Unlock()
	user := User{
		Name: name,
		Rank: initialRank,
	}
	return g.setObject("users", name, user)
}

func (g *Gladder) GetUser(name string) (*User, error) {
	g.Lock()
	defer g.Unlock()
	user := &User{}
	err := g.getObject("users", name, user)
	return user, err
}

func (g *Gladder) SaveUser(user *User) error {
	g.Lock()
	defer g.Unlock()
	return g.setObject("users", user.Name, user)
}

func (g *Gladder) GetUsers() (Ladder, error) {
	g.Lock()
	defer g.Unlock()
	var users []*User
	err := g.getObjects("users", &users)
	return Ladder(users), err
}

func (g *Gladder) setObject(cName, key string, obj interface{}) error {
	collection := GetOrCreateCollection(g.Store, cName)
	oEnc, err := GobEncode(obj)
	if err != nil {
		return err
	}
	err = collection.Set([]byte(key), oEnc)
	if err != nil {
		return err
	}
	err = collection.Write()
	if err != nil {
		return err
	}
	return g.Store.Flush()
}

func (g *Gladder) getObject(cName, key string, obj interface{}) error {
	collection := GetOrCreateCollection(g.Store, cName)
	value, err := collection.Get([]byte(key))
	if value == nil {
		return nil
	}
	if err != nil {
		return err
	}
	return GobDecode(value, obj)
}

func (g *Gladder) getObjects(cName string, outSlicePtr interface{}) error {
	collection := GetOrCreateCollection(g.Store, cName)
	svp := reflect.ValueOf(outSlicePtr)
	if svp.Kind() != reflect.Ptr {
		return fmt.Errorf("outSlicePtr must be a pointer to a slice of pointers")
	}
	sv := reflect.Indirect(svp)
	if sv.Kind() != reflect.Slice {
		return fmt.Errorf("outSlicePtr must be a pointer to a slice of pointers")
	}
	elemType := sv.Type().Elem().Elem()
	var visitErr error
	collection.VisitItemsAscend(nil, true, func(i *gkvlite.Item) bool {
		elem := reflect.New(elemType)
		buf := bytes.NewBuffer(i.Val)
		err := gob.NewDecoder(buf).DecodeValue(elem)
		if err != nil {
			visitErr = err
			return false
		}
		sv.Set(reflect.Append(sv, elem))
		return true
	})
	return visitErr
}

func GobEncode(e interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(e)
	return buf.Bytes(), err
}

func GobDecode(data []byte, e interface{}) error {
	buf := bytes.NewBuffer(data)
	return gob.NewDecoder(buf).Decode(e)
}

func GetOrCreateCollection(store *gkvlite.Store, cName string) *gkvlite.Collection {
	collection := store.GetCollection(cName)
	if collection == nil {
		collection = store.SetCollection(cName, bytes.Compare)
	}
	return collection
}
