package sql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type TestTable struct {
	Id    uint   `db:"id" `
	Value string `db:"value" `
}

func (TestTable) TableName() string {
	return "test_table"
}
func Test_sqlClient_getFieldsFromModel(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := sqlmock.NewRows([]string{"column_name"})
	rows.AddRow("id")
	rows.AddRow("value")

	mock.ExpectQuery("SELECT column_name FROM information_schema.columns WHERE table_name = 'test_table'").WillReturnRows(rows)
	mock.ExpectQuery("SELECT column_name FROM information_schema.columns WHERE table_name = 'unknown_table'").WillReturnRows(sqlmock.NewRows([]string{"column_name"}))

	fieldCache := make(map[string]string)

	c := sqlClient{
		sqlxDB,
		fieldCache,
	}

	type args struct {
		model     interface{}
		tableName string
		cacheLen  int
	}
	tests := []struct {
		name      string
		sqlClient sqlClient
		args      args
		want      string
		wantErr   bool
	}{
		{
			name:      "Test get fields success",
			sqlClient: c,
			args: args{
				model:     TestTable{},
				tableName: TestTable{}.TableName(),
				cacheLen:  1,
			},
			want:    "id, value",
			wantErr: false,
		},
		// TODO: check that result really comes from cache (not only that cache size didn't increase...)
		{
			name:      "Test get fields success from cache",
			sqlClient: c,
			args: args{
				model:     TestTable{},
				tableName: TestTable{}.TableName(),
				cacheLen:  1,
			},
			want:    "id, value",
			wantErr: false,
		},
		{
			name:      "Test get fields empty",
			sqlClient: c,
			args: args{
				model:     TestTable{},
				tableName: "unknown_table",
				cacheLen:  2,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sqlClient.getFieldsFromModelWithTablename(tt.args.model, tt.args.tableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlClient.getFieldsFromModelWithTablename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sqlClient.getFieldsFromModelWithTablename() = %v, want %v", got, tt.want)
			}
			if tt.args.cacheLen != len(tt.sqlClient.modelFieldCache) {
				t.Errorf("sqlClient.getFieldsFromModelWithTablename().modelFieldCache = %v, want %v", tt.sqlClient.modelFieldCache, tt.sqlClient.modelFieldCache)
			}
		})
	}
}

func Test_sqlClient_getFieldsFromModelWithoutTablename(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := sqlmock.NewRows([]string{"column_name"})
	rows.AddRow("id")
	rows.AddRow("value")

	mock.ExpectQuery("SELECT column_name FROM information_schema.columns WHERE table_name = 'test_table'").WillReturnRows(rows)
	mock.ExpectQuery("SELECT column_name FROM information_schema.columns WHERE table_name = 'unknown_table'").WillReturnRows(sqlmock.NewRows([]string{"column_name"}))

	fieldCache := make(map[string]string)

	type fields struct {
		DB              *sqlx.DB
		modelFieldCache map[string]string
	}
	type args struct {
		model Table
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantTableName string
		wantFieldList string
		wantErr       bool
	}{
		{
			name: "Test get fields success",
			fields: fields{
				DB:              sqlxDB,
				modelFieldCache: fieldCache,
			},
			args: args{
				model: TestTable{},
			},
			wantTableName: "test_table",
			wantFieldList: "id, value",
			wantErr:       false,
		},
		// TODO: Add test for cached results
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &sqlClient{
				DB:              tt.fields.DB,
				modelFieldCache: tt.fields.modelFieldCache,
			}
			gotTableName, gotFieldList, err := d.getFieldsFromModel(tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlClient.getFieldsFromModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTableName != tt.wantTableName {
				t.Errorf("sqlClient.getFieldsFromModel() gotTableName = %v, want %v", gotTableName, tt.wantTableName)
			}
			if gotFieldList != tt.wantFieldList {
				t.Errorf("sqlClient.getFieldsFromModel() gotFieldList = %v, want %v", gotFieldList, tt.wantFieldList)
			}
		})
	}
}
