package cql_test

import (
	"reflect"
	"strings"
	"testing"

	cql "github.com/johnnywidth/cql-calculator/cql-parser"
)

func TestParse(t *testing.T) {
	var testCases = []struct {
		s    string
		stmt *cql.Statement
	}{
		{
			s: `CREATE TABLE t.table_name 
				(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (id, email))`,
			stmt: &cql.Statement{
				TableName: "table_name",
				Colums: map[string]cql.Column{
					"id":          {Name: "id", Type: "int"},
					"email":       {Name: "email", Type: "text"},
					"name":        {Name: "name", Type: "text"},
					"status":      {Name: "status", Type: "tinyint"},
					"uploaded_at": {Name: "uploaded_at", Type: "timestamp"},
				},
				PK: map[string]cql.Column{
					"id": {Name: "id", Type: "int"},
				},
				CK: map[string]cql.Column{
					"email": {Name: "email", Type: "text"},
				},
				SK: map[string]cql.Column{
					"name": {Name: "name", Type: "text"},
				},
			},
		},
		{
			s: `CREATE TABLE IF NOT EXISTS t.table_name 
				(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (id, email))`,
			stmt: &cql.Statement{
				TableName: "table_name",
				Colums: map[string]cql.Column{
					"id":          {Name: "id", Type: "int"},
					"email":       {Name: "email", Type: "text"},
					"name":        {Name: "name", Type: "text"},
					"status":      {Name: "status", Type: "tinyint"},
					"uploaded_at": {Name: "uploaded_at", Type: "timestamp"},
				},
				PK: map[string]cql.Column{
					"id": {Name: "id", Type: "int"},
				},
				CK: map[string]cql.Column{
					"email": {Name: "email", Type: "text"},
				},
				SK: map[string]cql.Column{
					"name": {Name: "name", Type: "text"},
				},
			},
		},
		{
			s: `CREATE TABLE table_name 
				(id int, email text, name text, status tinyint, uploaded_at timestamp, PRIMARY KEY ((id), email))`,
			stmt: &cql.Statement{
				TableName: "table_name",
				Colums: map[string]cql.Column{
					"id":          {Name: "id", Type: "int"},
					"email":       {Name: "email", Type: "text"},
					"name":        {Name: "name", Type: "text"},
					"status":      {Name: "status", Type: "tinyint"},
					"uploaded_at": {Name: "uploaded_at", Type: "timestamp"},
				},
				PK: map[string]cql.Column{
					"id": {Name: "id", Type: "int"},
				},
				CK: map[string]cql.Column{
					"email": {Name: "email", Type: "text"},
				},
				SK: map[string]cql.Column{},
			},
		},
		{
			s: `CREATE TABLE table_name 
				(id int, email text, name text, status tinyint, uploaded_at timestamp, PRIMARY KEY ((id, email), name))`,
			stmt: &cql.Statement{
				TableName: "table_name",
				Colums: map[string]cql.Column{
					"id":          {Name: "id", Type: "int"},
					"email":       {Name: "email", Type: "text"},
					"name":        {Name: "name", Type: "text"},
					"status":      {Name: "status", Type: "tinyint"},
					"uploaded_at": {Name: "uploaded_at", Type: "timestamp"},
				},
				PK: map[string]cql.Column{
					"id":    {Name: "id", Type: "int"},
					"email": {Name: "email", Type: "text"},
				},
				CK: map[string]cql.Column{
					"name": {Name: "name", Type: "text"},
				},
				SK: map[string]cql.Column{},
			},
		},
		{
			s: `CREATE TABLE table_name 
				(id int, email text, name text, status tinyint, uploaded_at timestamp, PRIMARY KEY ((id, email)))`,
			stmt: &cql.Statement{
				TableName: "table_name",
				Colums: map[string]cql.Column{
					"id":          {Name: "id", Type: "int"},
					"email":       {Name: "email", Type: "text"},
					"name":        {Name: "name", Type: "text"},
					"status":      {Name: "status", Type: "tinyint"},
					"uploaded_at": {Name: "uploaded_at", Type: "timestamp"},
				},
				PK: map[string]cql.Column{
					"id":    {Name: "id", Type: "int"},
					"email": {Name: "email", Type: "text"},
				},
				CK: map[string]cql.Column{},
				SK: map[string]cql.Column{},
			},
		},
		{
			s: `CREATE TABLE table_name 
				(id int, email text, name text, status tinyint, uploaded_at timestamp, PRIMARY KEY (id))`,
			stmt: &cql.Statement{
				TableName: "table_name",
				Colums: map[string]cql.Column{
					"id":          {Name: "id", Type: "int"},
					"email":       {Name: "email", Type: "text"},
					"name":        {Name: "name", Type: "text"},
					"status":      {Name: "status", Type: "tinyint"},
					"uploaded_at": {Name: "uploaded_at", Type: "timestamp"},
				},
				PK: map[string]cql.Column{
					"id": {Name: "id", Type: "int"},
				},
				CK: map[string]cql.Column{},
				SK: map[string]cql.Column{},
			},
		},
		{
			s: `CREATE TABLE table_name (
				id int, set_text set<text>, map_int_text map<int,text>, list_int list<int>, 
				frozen_map map<text,frozen<list<text>>>, nested_map map<text,map<text,text>>, PRIMARY KEY (id))`,
			stmt: &cql.Statement{
				TableName: "table_name",
				Colums: map[string]cql.Column{
					"id":           {Name: "id", Type: "int"},
					"set_text":     {Name: "set_text", Type: "set<text>"},
					"map_int_text": {Name: "map_int_text", Type: "map<int,text>"},
					"list_int":     {Name: "list_int", Type: "list<int>"},
					"frozen_map":   {Name: "frozen_map", Type: "map<text,frozen<list<text>>>"},
					"nested_map":   {Name: "nested_map", Type: "map<text,map<text,text>>"},
				},
				PK: map[string]cql.Column{
					"id": {Name: "id", Type: "int"},
				},
				CK: map[string]cql.Column{},
				SK: map[string]cql.Column{},
			},
		},
	}

	for i, tc := range testCases {
		stmt, err := cql.NewParser(strings.NewReader(tc.s)).Parse()
		if !reflect.DeepEqual(tc.stmt, stmt) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tc.s, tc.stmt, stmt)
			t.Error(err)
		}
	}
}

func TestParse_Errors(t *testing.T) {
	var testCases = []string{
		`CREATE TABLE t. 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (id, email))`,
		`CREATE TABLE IF t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (id, email))`,
		`CREATE TABLE IF NOT t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (id, email))`,
		`CREATE TABLE NOT t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (id, email))`,
		`CREATE TABLE EXISTS t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (id, email))`,
		`CREATE TABLE IF EXISTS t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (id, email))`,
		`CREATE t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (id, email))`,
		`TABLE table_name 
			(id int, email text, name text, status tinyint, uploaded_at timestamp, PRIMARY KEY ((id), email))`,
		`table_name 
			(id int, email text, name text, status tinyint, uploaded_at timestamp, PRIMARY KEY ((id, email), name))`,
		`CREATE TABLE 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (id, email))`,
		`CREATE TABLE t.table_name 
			id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (id, email))`,
		`CREATE TABLE t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY id, email))`,
		`CREATE TABLE t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY (id, email))`,
		`CREATE TABLE t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (, email))`,
		`CREATE TABLE t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, (id, email))`,
		`CREATE TABLE t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY ((), email))`,
		`CREATE TABLE t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY ((id name), email))`,
		`CREATE TABLE t.table_name 
			(id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (id, email,))`,
	}

	for _, tc := range testCases {
		_, err := cql.NewParser(strings.NewReader(tc)).Parse()
		if err == nil {
			t.Fatalf("Should be err for %s, got nil", tc)
		}
	}
}
