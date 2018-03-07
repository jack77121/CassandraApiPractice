package Users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gocql/gocql"
	"github.com/jack77121/CassandraApiPractice/Cassandra"
)

// Get User api
func Get(w http.ResponseWriter, r *http.Request) {
	var userList []User
	m := map[string]interface{}{}

	query := "SELECT * FROM testgoapi.users"
	iterable := Cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		userList = append(userList, User{
			ID:        m["id"].(gocql.UUID),
			FirstName: m["firstname"].(string),
			LastName:  m["lastname"].(string),
			Email:     m["email"].(string),
			Age:       m["age"].(int),
			City:      m["city"].(string),
		})
		m = map[string]interface{}{}
	}
	json.NewEncoder(w).Encode(AllUsersResponse{Users: userList})
}

// GetOne : retrive one user info. by the UUID
func GetOne(w http.ResponseWriter, r *http.Request) {
	var user User
	var errs []string
	m := map[string]interface{}{}
	found := false

	vars := mux.Vars(r)
	uuid, err := gocql.ParseUUID(vars["user_uuid"])

	if err != nil {
		errs = append(errs, err.Error())
	} else {
		query := "SELECT * FROM users WHERE id=? LIMIT 1"
		iter := Cassandra.Session.Query(query, uuid).Consistency(gocql.Two).Iter()
		for iter.MapScan(m) {
			found = true
			user = User{
				ID:        m["id"].(gocql.UUID),
				FirstName: m["firstname"].(string),
				LastName:  m["lastname"].(string),
				Email:     m["email"].(string),
				Age:       m["age"].(int),
				City:      m["city"].(string),
			}
		}
		if found == false {
			fmt.Println("id not found")
		}

		if found == true {
			json.NewEncoder(w).Encode(GetUserResponse{User: user})
		} else {
			json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})

		}

	}
}

// Enrich - first name +  last name
func Enrich(uuids []gocql.UUID) map[string]string {
	if len(uuids) != 0 {
		names := map[string]string{}
		m := map[string]interface{}{}

		query := "SELECT id, firstname, lastname FROM users WHERE id IN ?"
		iterate := Cassandra.Session.Query(query, uuids).Iter()
		for iterate.MapScan(m) {
			firstname := m["firstname"].(string)
			lastname := m["lastname"].(string)
			names[m["id"].(string)] = fmt.Sprintf("%s %s", firstname, lastname)
			m = map[string]interface{}{}
		}
		return names
	}
	return map[string]string{}
}
