package client

// Client is the interface to the storer
type Client interface {
	GetDataTypes() ([]map[string]string, error)
	CreateBucketType(bucketName, dataType string) error
	DeleteBucketType(bucketName, bucketType string) error
	CreateUser(username, password string) error
	DeleteUser(username string) error
	GrantUserAccess(username, bucketName string) error
	RevokeUserAccess(username, bucketName string) error
	BucketStatus(bucketName string) error
}

// Nil implements client interface doing nothing
type Nil struct {
}

// NewNil Creates a new nil client
func NewNil() *Nil {
	return &Nil{}
}

func (c *Nil) GetDataTypes() ([]map[string]string, error)           { return []map[string]string{}, nil }
func (c *Nil) CreateBucketType(bucketName, dataType string) error   { return nil }
func (c *Nil) DeleteBucketType(bucketName, bucketType string) error { return nil }
func (c *Nil) CreateUser(username, password string) error           { return nil }
func (c *Nil) DeleteUser(username string) error                     { return nil }
func (c *Nil) GrantUserAccess(username, bucketName string) error    { return nil }
func (c *Nil) RevokeUserAccess(username, bucketName string) error   { return nil }
func (c *Nil) BucketStatus(bucketName string) error                 { return nil }
