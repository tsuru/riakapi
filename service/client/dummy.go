package client

import (
	"errors"
)

// Dummy is the entrypoint for riak dummy client
type Dummy struct {
	*Riak
}

// NewDummy creates a dummy client, useful for testing
func NewDummy() *Dummy {
	return &Dummy{&Riak{host: "127.0.0.1", port: 0}}

}

func (c *Dummy) CreateBucketType(bucketName, dataType string) error {
	return errors.New("Not implemented")
}
func (c *Dummy) DeleteBucketType(bucketName, bucketType string) error {
	return errors.New("Not implemented")
}
func (c *Dummy) CreateUser(username, password string) error {
	return errors.New("Not implemented")
}
func (c *Dummy) DeleteUser(username string) error {
	return errors.New("Not implemented")
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
