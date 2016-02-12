package client

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	riak "github.com/basho/riak-go-client"
	"golang.org/x/crypto/ssh"

	"gitlab.qdqmedia.com/shared-projects/riakapi/config"
)

// Riak admin cmd fmts
const (
	checkBucketTypePresentCmd = `sudo riak-admin bucket-type list | grep -e "%s"`
	createBucketTypeCmd       = `sudo riak-admin bucket-type create %s '{"props":{"datatype":"%s", "allow_mult": true}}'`
	activateBucketTypeCmd     = `sudo riak-admin bucket-type activate %s`
)

// This will hold the added instances on tsuru
const (
	RiakInstancesInfoBucket = "tsuru-instances"
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

	for k, v := range BucketTypes {
		r = append(r, map[string]string{
			"name":        k,
			"description": v,
		})
	}
	return r, nil
}

// CreateBucket Creates a bucket on riak
func (c *Riak) CreateBucket(bucketName, bucketType string) error {

	// Check valid bucketType
	if _, ok := BucketTypes[bucketType]; !ok {
		logrus.Errorf("%s is not a valid bucket type", bucketType)
		return errors.New("Not valid bucket type")
	}

	// First ensure the data types are createed (with Riak-admin)
	if err := c.ensureBucketTypePresent(bucketType); err != nil {
		logrus.Errorf("Could not ensure bucket type '%s' presence: %v", bucketType, err)
		return err
	}

	// Second create bucket on the bucket type
	if err := c.ensureBucketPresent(bucketName, bucketType); err != nil {
		logrus.Errorf("Could not create bucket '%s' on bucket '%s': %v", bucketName, bucketType, err)
		return err
	}

	// Third save the location of the created bucket (store its datatype)
	if err := c.saveBucketLocation(bucketName, bucketType); err != nil {
		logrus.Errorf("Could not save bucket '%s' location: %v", bucketName, err)
		return err
	}
	c.GetBucketType(bucketName)
	logrus.Infof("Bucket '%s' of bucket type '%s' ready", bucketName, bucketType)
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

// GetBucketType returns the bucket type based on the bucket name
func (c *Riak) GetBucketType(bucketName string) string {
	var bucketType string

	// Create command
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucket(RiakInstancesInfoBucket).
		WithKey(bucketName).
		Build()
	if err != nil {
		return ""
	}

	if err = c.RiakClient.Execute(cmd); err != nil {
		return ""
	}

	fvc, ok := cmd.(*riak.FetchValueCommand)
	if !ok {
		return ""
	}
	if len(fvc.Response.Values) == 0 {
		return ""
	}

	bucketType = string(fvc.Response.Values[0].Value)
	logrus.Debugf("Retrieved Bucket type '%s' from bucket name '%s'", bucketType, bucketName)

	return bucketType
}

//ensureBucketTypePresent checks bucket type present and if not will create adn activate it
func (c *Riak) ensureBucketTypePresent(bucketType string) error {
	// Sessions are channels on the same connection
	var session *ssh.Session
	var err error

	// Check bucket type is present
	logrus.Debugf("Check bucket type '%s' is created", bucketType)
	cmd := fmt.Sprintf(checkBucketTypePresentCmd, bucketType)

	session, _ = c.SSHClient.NewSession()
	err = session.Run(cmd)
	session.Close()

	// If error will need to create the bucket
	if err != nil {
		n, _ := NameBucketTypeMapping[bucketType]
		cmd = fmt.Sprintf(createBucketTypeCmd, bucketType, n)
		session, _ = c.SSHClient.NewSession()
		err = session.Run(cmd)
		session.Close()
		if err != nil {
			return fmt.Errorf("Could not create bucket type: %v", err)
		}
		logrus.Debugf("Bucket type '%s' created", bucketType)
	} else {
		logrus.Debugf("Bucket type '%s' already present", bucketType)
	}
	// Activate always
	cmd = fmt.Sprintf(activateBucketTypeCmd, bucketType)
	session, _ = c.SSHClient.NewSession()
	err = session.Run(cmd)
	session.Close()
	if err != nil {
		return fmt.Errorf("Failed activating bucket type: %s", bucketType)
	}
	logrus.Debugf("Bucket type '%s' activated", bucketType)
	return nil
}

func (c *Riak) ensureBucketPresent(bucketName, bucketType string) error {
	// Select the correct data type and create the bucket
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
	}
	if err != nil {
		return fmt.Errorf("Could not create bucket type: %v", err)
	}

	if err = c.RiakClient.Execute(cmd); err != nil {
		return fmt.Errorf("Could not create bucket type: %v", err)
	}

	// Create bucket
	propsCmd, err := riak.NewStoreBucketTypePropsCommandBuilder().
		WithBucketType(bucketType).
		WithAllowMult(true).
		Build()
	if err != nil {
		return fmt.Errorf("Could not set props on bucket type: %v", err)
	}

	if err = c.RiakClient.Execute(propsCmd); err != nil {
		return fmt.Errorf("Could not set props on bucket type: %v", err)
	}

	logrus.Debugf("Bucket '%s' of bucket type '%s' created", bucketName, bucketType)
	return nil
}

// saveBucketLocation will save the location (default bucket type) of the bucketname
// in key->value form: bucketName->bucketType this is used so we can reach the
// bucket when we don't have the bucketType.
func (c *Riak) saveBucketLocation(bucketNameKey, bucketTypeValue string) error {
	// Our value to store
	obj := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte(bucketTypeValue),
	}

	// Create command
	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucket(RiakInstancesInfoBucket).
		WithKey(bucketNameKey).
		WithContent(obj).
		Build()
	if err != nil {
		return fmt.Errorf("Could not store bucket location: %v", err)
	}

	if err := c.RiakClient.Execute(cmd); err != nil {
		return fmt.Errorf("Could not store bucket location: %v", err)
	}
	logrus.Debugf("Bucket '%s' location stored", bucketNameKey)
	return nil
}
