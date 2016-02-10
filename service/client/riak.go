package client

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