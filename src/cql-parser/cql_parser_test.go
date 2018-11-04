package cql_test

import (
	"reflect"
	"strings"
	"testing"

	cql "github.com/johnnywidth/cql-calculator/src/cql-parser"
)

func TestParse(t *testing.T) {
	var testCase = []struct {
		s    string
		stmt *cql.Statement
		err  string
	}{
		{
			s: `CREATE TABLE t.video (video_id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (video_id, email))`,
			stmt: &cql.Statement{
				TableName: "video",
				Colums: map[string]cql.Column{
					"video_id":    {Name: "video_id", Type: "int"},
					"email":       {Name: "email", Type: "text"},
					"name":        {Name: "name", Type: "text"},
					"status":      {Name: "status", Type: "tinyint"},
					"uploaded_at": {Name: "uploaded_at", Type: "timestamp"},
				},
				PK: map[string]cql.Column{
					"video_id": {Name: "video_id", Type: "int"},
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
			s: `CREATE TABLE video (video_id int, email text, name text, status tinyint, uploaded_at timestamp, PRIMARY KEY ((video_id), email))`,
			stmt: &cql.Statement{
				TableName: "video",
				Colums: map[string]cql.Column{
					"video_id":    {Name: "video_id", Type: "int"},
					"email":       {Name: "email", Type: "text"},
					"name":        {Name: "name", Type: "text"},
					"status":      {Name: "status", Type: "tinyint"},
					"uploaded_at": {Name: "uploaded_at", Type: "timestamp"},
				},
				PK: map[string]cql.Column{
					"video_id": {Name: "video_id", Type: "int"},
				},
				CK: map[string]cql.Column{
					"email": {Name: "email", Type: "text"},
				},
				SK: map[string]cql.Column{},
			},
		},
		{
			s: `CREATE TABLE video (video_id int, email text, name text, status tinyint, uploaded_at timestamp, PRIMARY KEY ((video_id, email), name))`,
			stmt: &cql.Statement{
				TableName: "video",
				Colums: map[string]cql.Column{
					"video_id":    {Name: "video_id", Type: "int"},
					"email":       {Name: "email", Type: "text"},
					"name":        {Name: "name", Type: "text"},
					"status":      {Name: "status", Type: "tinyint"},
					"uploaded_at": {Name: "uploaded_at", Type: "timestamp"},
				},
				PK: map[string]cql.Column{
					"video_id": {Name: "video_id", Type: "int"},
					"email":    {Name: "email", Type: "text"},
				},
				CK: map[string]cql.Column{
					"name": {Name: "name", Type: "text"},
				},
				SK: map[string]cql.Column{},
			},
		},
		{
			s: `CREATE TABLE video (video_id int, email text, name text, status tinyint, uploaded_at timestamp, PRIMARY KEY ((video_id, email)))`,
			stmt: &cql.Statement{
				TableName: "video",
				Colums: map[string]cql.Column{
					"video_id":    {Name: "video_id", Type: "int"},
					"email":       {Name: "email", Type: "text"},
					"name":        {Name: "name", Type: "text"},
					"status":      {Name: "status", Type: "tinyint"},
					"uploaded_at": {Name: "uploaded_at", Type: "timestamp"},
				},
				PK: map[string]cql.Column{
					"video_id": {Name: "video_id", Type: "int"},
					"email":    {Name: "email", Type: "text"},
				},
				CK: map[string]cql.Column{},
				SK: map[string]cql.Column{},
			},
		},
		{
			s: `CREATE TABLE video (video_id int, email text, name text, status tinyint, uploaded_at timestamp, PRIMARY KEY (video_id))`,
			stmt: &cql.Statement{
				TableName: "video",
				Colums: map[string]cql.Column{
					"video_id":    {Name: "video_id", Type: "int"},
					"email":       {Name: "email", Type: "text"},
					"name":        {Name: "name", Type: "text"},
					"status":      {Name: "status", Type: "tinyint"},
					"uploaded_at": {Name: "uploaded_at", Type: "timestamp"},
				},
				PK: map[string]cql.Column{
					"video_id": {Name: "video_id", Type: "int"},
				},
				CK: map[string]cql.Column{},
				SK: map[string]cql.Column{},
			},
		},
	}

	for i, tt := range testCase {
		stmt, err := cql.NewParser(strings.NewReader(tt.s)).Parse()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.stmt, stmt)
		}
	}
}

func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
