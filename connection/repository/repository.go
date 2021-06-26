package repository

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

type Repository struct {
	DatabaseDriver *neo4j.Driver
}
