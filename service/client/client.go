package client

// Client is the interface to the storer
type Client interface {
	GetDataTypes() (map[string]string, error)
	CreateBucketType(bucketName, dataType string) error
	DeleteBucketType(bucketName, bucketType string) error
	CreateUser(username, password string) error
	DeleteUser(username string) error
	GrantUserAccess(username, bucketName string) error
	RevokeUserAccess(username, bucketName string) error
	BucketStatus(bucketName string) error
}

// NilClient implements client interface doing nothing
type NilClient struct {
}

// NewNilClient Creates a new nil client
func NewNilClient() *NilClient {
	return &NilClient{}
}

func (n *NilClient) GetDataTypes() (map[string]string, error)             { return make(map[string]string), nil }
func (n *NilClient) CreateBucketType(bucketName, dataType string) error   { return nil }
func (n *NilClient) DeleteBucketType(bucketName, bucketType string) error { return nil }
func (n *NilClient) CreateUser(username, password string) error           { return nil }
func (n *NilClient) DeleteUser(username string) error                     { return nil }
func (n *NilClient) GrantUserAccess(username, bucketName string) error    { return nil }
func (n *NilClient) RevokeUserAccess(username, bucketName string) error   { return nil }
func (n *NilClient) BucketStatus(bucketName string) error                 { return nil }
