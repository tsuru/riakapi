package client

import (
	"errors"
	"sync"

	"github.com/Sirupsen/logrus"

	"gitlab.qdqmedia.com/shared-projects/riakapi/utils"
)

type UserProps struct {
	Username string
	Password string
	ACL      []string // bucket names wich can access
}

// Dummy is the entrypoint for riak dummy client
type Dummy struct {
	*Riak

	// Our custom database (on the instance to allow parallel tests)
	Buckets map[string]string
	Users   map[string]*UserProps

	bucketsMutex *sync.Mutex
	usersMutex   *sync.Mutex
}

// NewDummy creates a dummy client, useful for testing
func NewDummy() *Dummy {
	return &Dummy{
		Riak:         &Riak{},
		Buckets:      map[string]string{},
		Users:        map[string]*UserProps{},
		bucketsMutex: &sync.Mutex{},
		usersMutex:   &sync.Mutex{},
	}
}

func (c *Dummy) GetBucketType(bucketName string) string {
	return c.Buckets[bucketName]
}

func (c *Dummy) GetBucketTypes() ([]map[string]string, error) {
	var r []map[string]string

	for k, v := range BucketTypes {
		r = append(r, map[string]string{
			"name":        k,
			"description": v,
		})
	}
	return r, nil

}

func (c *Dummy) CreateBucket(bucketName, bucketType string) error {
	// Check bucket type
	if _, ok := BucketTypes[bucketType]; !ok {
		return errors.New("Not valid bucket type")
	}

	c.bucketsMutex.Lock()
	defer c.bucketsMutex.Unlock()
	if _, ok := c.Buckets[bucketName]; !ok {
		c.Buckets[bucketName] = bucketType
		logrus.Infof("Bucket '%s' of type '%s' created", bucketName, bucketType)
		return nil
	}
	return errors.New("Bucket already declared")
}
func (c *Dummy) DeleteBucket(bucketName, bucketType string) error {
	c.bucketsMutex.Lock()
	defer c.bucketsMutex.Unlock()
	if _, ok := c.Buckets[bucketName]; ok {
		delete(c.Buckets, bucketName)
		return nil
	}
	return nil
}

func (c *Dummy) EnsureUserPresent(word string) (user, pass string, err error) {
	c.usersMutex.Lock()
	defer c.usersMutex.Unlock()
	//TODO: use salt
	user = utils.GenerateUsername(word)
	var props *UserProps
	var ok bool
	if props, ok = c.Users[user]; !ok {
		pass = word // handy when using on the tests
		c.Users[user] = &UserProps{
			Username: user,
			Password: pass,
			ACL:      []string{},
		}
		return
	}
	pass = props.Password
	return
}
func (c *Dummy) GrantUserAccess(username, bucketName string) error {
	c.usersMutex.Lock()
	defer c.usersMutex.Unlock()
	if user, ok := c.Users[username]; ok {
		// Check if present already (performance on dummy doesn't matter)
		for _, a := range user.ACL {
			if a == bucketName {
				return nil
			}
		}
		user.ACL = append(user.ACL, bucketName)
		return nil
	}
	return errors.New("Not present user")
}

func (c *Dummy) DeleteUser(username string) error {
	c.usersMutex.Lock()
	defer c.usersMutex.Unlock()

	if _, ok := c.Users[username]; ok {
		delete(c.Users, username)
		return nil
	}
	return errors.New("Theres no user to delete")
}
func (c *Dummy) RevokeUserAccess(username, bucketName string) error {
	return errors.New("Not implemented")
}
func (c *Dummy) BucketStatus(bucketName string) error {
	return errors.New("Not implemented")
}
