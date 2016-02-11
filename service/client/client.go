/*Package client will respond to tsuru events creating by default one
bucket type for each data type (counter, set and map). Having these bucket types
for each new service instance the client will create a bucket, and for each app
binding to this instance the riak client will create a user and an ACL to access
to this bucket.
*/
package client

const (
	//BucketTypeCounter is a counter data type bucket type
	BucketTypeCounter = "tsuru-counter"
	//BucketTypeSet is a set data type bucket type
	BucketTypeSet = "tsuru-set"
	//BucketTypeMap is a map data type bucket type
	BucketTypeMap = "tsuru-map"
)

// bucketTypes lists all the bucket types available
var bucketTypes = map[string]string{
	BucketTypeCounter: "Bucket type of counter data type",
	BucketTypeSet:     "Bucket type of set data type",
	BucketTypeMap:     "Bucket type of map data type",
}

// Client is the interface to the storer
type Client interface {
	GetBucketTypes() ([]map[string]string, error)
	CreateBucket(bucketName, bucketType string) error
	DeleteBucket(bucketName, bucketType string) error
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

func (c *Nil) GetBucketTypes() ([]map[string]string, error)       { return []map[string]string{}, nil }
func (c *Nil) CreateBucket(bucketName, bucketType string) error   { return nil }
func (c *Nil) DeleteBucket(bucketName, bucketType string) error   { return nil }
func (c *Nil) CreateUser(username, password string) error         { return nil }
func (c *Nil) DeleteUser(username string) error                   { return nil }
func (c *Nil) GrantUserAccess(username, bucketName string) error  { return nil }
func (c *Nil) RevokeUserAccess(username, bucketName string) error { return nil }
func (c *Nil) BucketStatus(bucketName string) error               { return nil }
