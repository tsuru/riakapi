package client

import (
	"errors"
	"sync"

	"github.com/Sirupsen/logrus"
)

var buckets map[string]string
var users map[string]userProps

var bucketsMutex = &sync.Mutex{}
var usersMutex = &sync.Mutex{}

type userProps struct {
	username string
	password string
	ACL      []string // bucket names wich can access
}

// Dummy is the entrypoint for riak dummy client
type Dummy struct {
	*Riak
}

// NewDummy creates a dummy client, useful for testing
func NewDummy() *Dummy {
	// init once only
	var once sync.Once
	once.Do(func() {
		buckets = make(map[string]string)
		users = make(map[string]userProps)
	})

	return &Dummy{&Riak{host: "127.0.0.1", port: 0}}
}

// Used for tests
func (c *Dummy) Flush() {
	buckets = make(map[string]string)
	users = make(map[string]userProps)
}

func (c *Dummy) CreateBucketType(bucketName, dataType string) error {
	// Check bucket type
	if _, ok := dataTypes[dataType]; !ok {
		return errors.New("Not valid bucket data type")
	}

	bucketsMutex.Lock()
	defer bucketsMutex.Unlock()
	if _, ok := buckets[bucketName]; !ok {
		buckets[bucketName] = dataType
		logrus.Infof("Bucket '%s' of type '%s' created", bucketName, dataType)
		return nil
	}
	return errors.New("Bucket type already declared")
}
func (c *Dummy) DeleteBucketType(bucketName, bucketType string) error {
	bucketsMutex.Lock()
	defer bucketsMutex.Unlock()
	if _, ok := buckets[bucketName]; ok {
		delete(buckets, bucketName)
		return nil
	}
	return nil
}
func (c *Dummy) CreateUser(username, password string) error {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	if _, ok := users[username]; !ok {
		users[username] = userProps{
			username: username,
			password: password,
			ACL:      []string{},
		}
		return nil
	}
	return errors.New("User already present")
}
func (c *Dummy) DeleteUser(username string) error {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	if _, ok := users[username]; ok {
		delete(users, username)
		return nil
	}
	return errors.New("Theres no user to delete")
}
func (c *Dummy) GrantUserAccess(username, bucketName string) error {
	return errors.New("Not implemented")
}
func (c *Dummy) RevokeUserAccess(username, bucketName string) error {
	return errors.New("Not implemented")
}
func (c *Dummy) BucketStatus(bucketName string) error {
	return errors.New("Not implemented")
}
