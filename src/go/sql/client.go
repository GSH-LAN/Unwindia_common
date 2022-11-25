package sql

import "github.com/jmoiron/sqlx"

type SqlClient struct {
	*sqlx.DB
	modelFieldCache map[string]string
}

func New(driverName, dataSourceName string) (*SqlClient, error) {
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return &SqlClient{
		db,
		make(map[string]string),
	}, nil
}
