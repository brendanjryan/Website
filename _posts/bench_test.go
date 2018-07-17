package main

// import (
// 	"bytes"
// 	"fmt"
// 	"strings"
// 	"testing"
// 	"time"
// )

// func fmtStatStr(stat string, tags map[string]string) string {
// 	parts := []string{}
// 	for k, v := range tags {
// 		if v != "" {
// 			parts = append(parts, fmt.Sprintf("%s:%s", k, v))
// 		}
// 	}

// 	return fmt.Sprintf("%s|%s", stat, strings.Join(parts, ","))
// }

// func TestfmtStatStr(t *testing.T){

// 	tests := []struct{
// 		msg string
// 		stat string
// 		tags map[string]string
// 		exp string
// 	}{
// 		{
// 			msg:  "empty string and tags",
// 			stat: "",
// 			tags: nil,
// 			exp:  "",
// 		},
// 		{
// 			msg:  "empty tags",
// 			stat: "handler.latency",
// 			tags: nil,
// 			exp:  "foo",
// 		},
// 		{
// 			msg:  "with tags",
// 			stat: "handler.latency",
// 			tags: map[string]string{
// 				"host": "aws_",
// 				"service": "users",
// 			},
// 			exp:  "foo",
// 		},
// 	}

// 	for _  , tt := range tests {
// 		t.Run(tt.msg, func(t *testing.T) {
// 			assert.Equal(t, tt.exp, fmtStatStr(tt.stat, tt.tags))
// 		})
// 	}
// }

// // User represents a user in our application.
// type User struct {
// 	Id int64
// 	FirstName string
// 	LastName string
// }

// // UserMutator defines a type which is able to preform CRUD operations on User data.
// type UserMutator interface {
// 	Fetch(int64) (*User, error)
// 	Create(*User) (*User, error)
// 	Update(int64, *User) (*User, error)
// 	Delete(int64) error
// }

// // UserManager handles datastore operations for 'user' data.
// type UserManager struct {
// 	// database tables, caches, etc...
// }

// // Fetch retrieves a single user model keyed by id.
// func (*UserManager) Fetch(id int64) (*User, error) {
// 	return nil, nil
// }

// // Create instantiates a new user model
// func (*UserManager) Create(u *User) (*User, error) {
// 	return nil, nil
// }

// // Update updates the user keyed by the provided id to a new state.
// func (*UserManager) Update(id int64, u *User) (*User, error) {
// 	return nil, nil
// }

// // Delete deletes the user keyed by the provided id.
// func (*UserManager) Delete(id int64) error {
// 	return nil
// }

// // UserManagerCreator wraps a UserManager and provides additional
// // functionality for setting up and tearing down this manager's underlying
// // datastores.
// type UserManagerCreator struct {
// 	UserManager
// }

// // Setup instantiates resources used by this manager.
// func (u *UserManagerCreator) Setup() error {
// 	return nil
// }

// // Teardown cleans up resources used by this manager.
// func (u *UserManagerCreator) Teardown() error {
// 	return nil
// }

// var (
// 	timeMultiplier = 1.0
// 	batchMaxAge = time.Second
// 	batchSize = 25
// )

// func flushBatches(items <- chan string) {

// 	var batch []string
// 	select {

// 	case i := <- items:
// 		batch = append(batch, i)

// 		if len(batch) >= batchSize {
// 			log.Println("flushing batch of events")
// 			batch = ni
// 		}

// 	case <-time.Tick(time.Duration(timeMultiplier) * batchMaxAge):
// 		log.Println("flushing batch of stale events")
// 		batch = nil
// 	}
// }

// func main() {

// 	countChan := make(chan int)

// 	for c := range countChan {
// 		log.Println(c)
// 	}

// }

// var test = struct {
// 	stat string
// 	tags map[string]string
// }{
// 	"handler.sample",
// 	map[string]string{
// 		"os":     "ios",
// 		"locale": "en-US",
// 	},
// }

// func Benchmark_fmtStatStrBefore(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		fmtStatStr(test.stat, test.tags)
// 	}
// }

// func Benchmark_fmtStatStrAfter(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		fmtStatStrAfter(test.stat, test.tags)
// 	}
// }
