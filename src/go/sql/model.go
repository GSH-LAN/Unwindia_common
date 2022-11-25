// Package sql provides some utility functions for working with sql databases
package sql

// Table is an interface that should be implemented by all models which represents a table in the database
type Table interface {
	TableName() string
}
