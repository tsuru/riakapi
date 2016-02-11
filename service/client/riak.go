package client

import "errors"

const (
	RiakDataTypeFlag     = "flag"
	RiakDataTypeRegister = "register"
	RiakDataTypeCounter  = "counter"
	RiakDataTypeSet      = "set"
	RiakDataTypeMap      = "map"
)

// DataTypes lists all the datatypes available for bucket types on redis
var dataTypes = map[string]string{
	RiakDataTypeFlag:     "Bucket type of flag data type",
	RiakDataTypeRegister: "Bucket type of register data type",
	RiakDataTypeCounter:  "Bucket type of counter data type",
	RiakDataTypeSet:      "Bucket type of set data type",
	RiakDataTypeMap:      "Bucket type of map data type",
}

// Riak is the entrypoint for riak client
type Riak struct {
	host string
	port int
}

// NewRiak creates a riak client
func NewRiak(host string, port int) *Riak {
	return &Riak{
		host: host,
		port: port,
	}
}

// GetDataTypes Gets Riak plans
func (c *Riak) GetDataTypes() ([]map[string]string, error) {
	var r []map[string]string

	for k, v := range dataTypes {
		r = append(r, map[string]string{
			"name":        k,
			"description": v,
		})
	}
	return r, nil
}

func (c *Riak) CreateBucketType(bucketName, dataType string) error {
	return errors.New("Not implemented")
}
func (c *Riak) DeleteBucketType(bucketName, bucketType string) error {
	return errors.New("Not implemented")
}
func (c *Riak) CreateUser(username, password string) error {
	return errors.New("Not implemented")
}
func (c *Riak) DeleteUser(username string) error {
	return errors.New("Not implemented")
}
func (c *Riak) GrantUserAccess(username, bucketName string) error {
	return errors.New("Not implemented")
}
func (c *Riak) RevokeUserAccess(username, bucketName string) error {
	return errors.New("Not implemented")
}
func (c *Riak) BucketStatus(bucketName string) error {
	return errors.New("Not implemented")
}
