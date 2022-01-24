package sql

import "github.com/jmoiron/sqlx"

type sqlClient struct {
	*sqlx.DB
	modelFieldCache map[string]string
}

func New(driverName, dataSourceName string) (*sqlClient, error) {
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return &sqlClient{
		db,
		make(map[string]string),
	}, nil
}
