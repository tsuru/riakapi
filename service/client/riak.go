package client

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"

	"github.com/Sirupsen/logrus"
	riak "github.com/basho/riak-go-client"

	"gitlab.qdqmedia.com/shared-projects/riakapi/config"
)

// Riak admin cmd fmts
const (
	//NOTE: Bucket types need to be created before setting up the service
	createBucketTypeCmd   = `riak-admin bucket-type create %s '{"props":{"datatype":"%s"}}'`
	activateBucketTypeCmd = `riak-admin bucket-type activate %s`
)

// Riak is the entrypoint for riak client
type Riak struct {
	//SSHConnection SSH connection (for riak-admin manage operations)
	SSHClient *ssh.Client

	// RiakClient riak lowlevel client (for riak bucket operations)
	RiakClient *riak.Client
}

// NewRiak creates a riak client and the ssh connection
func NewRiak(cfg *config.ServiceConfig) *Riak {

	// Create riak client
	rClient, err := riak.NewClient(&riak.NewClientOptions{
		RemoteAddresses: []string{fmt.Sprintf("%s:%d", cfg.RiakHost, cfg.RiakPort)},
	})

	if err != nil {
		logrus.Fatalf("Error connecting to riak: %v", err)
	}

	// Create SSH connection
	sshConfig := &ssh.ClientConfig{
		User: cfg.SSHUser,
		Auth: cfg.SSHAuthMethods,
	}
	addr := fmt.Sprintf("%s:%d", cfg.SSHHost, cfg.SSHPort)
	sClient, err := ssh.Dial("tcp", addr, sshConfig)

	if err != nil {
		logrus.Fatalf("Error connecting with ssh: %v", err)
	}

	return &Riak{
		RiakClient: rClient,
		SSHClient:  sClient,
	}
}

// GetBucketTypes Gets Riak plans
func (c *Riak) GetBucketTypes() ([]map[string]string, error) {
	var r []map[string]string

	for k, v := range bucketTypes {
		r = append(r, map[string]string{
			"name":        k,
			"description": v,
		})
	}
	return r, nil
}

// CreateBucket Creates a bucket on riak
func (c *Riak) CreateBucket(bucketName, bucketType string) error {

	// Select the correct data type and create
	var cmd riak.Command
	var err error
	switch bucketType {
	case BucketTypeCounter:
		cmd, err = riak.NewUpdateCounterCommandBuilder().
			WithBucketType(bucketType).
			WithBucket(bucketName).
			Build()
	case BucketTypeSet:
		cmd, err = riak.NewUpdateSetCommandBuilder().
			WithBucketType(bucketType).
			WithBucket(bucketName).
			Build()
	case BucketTypeMap:
		cmd, err = riak.NewUpdateMapCommandBuilder().
			WithBucketType(bucketType).
			WithBucket(bucketName).
			Build()
	default:
		return errors.New("Not valid bucket data type")
	}
	if err != nil {
		return fmt.Errorf("Could not create bucket type: %v", err)
	}

	err = c.RiakClient.Execute(cmd)

	if err != nil {
		return fmt.Errorf("Could not create bucket type: %v", err)
	}

	// Set props on bucket type
	propsCmd, err := riak.NewStoreBucketTypePropsCommandBuilder().
		WithBucketType(bucketType).
		WithAllowMult(true).
		Build()
	if err != nil {
		return fmt.Errorf("Could not set props on bucket type: %v", err)
	}

	err = c.RiakClient.Execute(propsCmd)

	if err != nil {
		return fmt.Errorf("Could not set props on bucket type: %v", err)
	}

	logrus.Infof("Created bucket '%s' on bucket type '%s'", bucketName, bucketType)
	return nil
}

func (c *Riak) DeleteBucket(bucketName, bucketType string) error {
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
