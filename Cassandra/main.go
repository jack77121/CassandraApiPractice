package Cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

// Session w/ uppercase variable can be used in other package
var Session *gocql.Session

func init() {
	var err error

	cluster := gocql.NewCluster("127.0.0.1", "127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5")
	cluster.Keyspace = "testgoapi"
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("cassandra init done")
}
