package sql

import (
	"fmt"
	"reflect"

	"github.com/GSH-LAN/Unwindia_common/src/go/helper"
)

func (d *sqlClient) getFieldsFromModel(model Table) (tableName string, fieldList string, err error) {
	if tabler, ok := model.(Table); ok {
		tableName = tabler.TableName()
	}
	fieldList, err = d.getFieldsFromModelWithTablename(model, tableName)
	return
}

func (d *sqlClient) getFieldsFromModelWithTablename(model interface{}, tableName string) (string, error) {

	if queryString, ok := d.modelFieldCache[tableName]; ok {
		// log.Debugf("Retrieved query for %s from cache: %s", tableName, queryString)
		return queryString, nil
	}

	var columns []string
	var queryString string

	if err := d.Select(&columns, fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s'", tableName)); err != nil {
		return "", err
	}

	reflectedModel := reflect.ValueOf(model)

	for i := 0; i < reflectedModel.NumField(); i++ {
		typeField := reflectedModel.Type().Field(i)
		if val, ok := typeField.Tag.Lookup("db"); ok {
			if helper.StringSliceContains(columns, val) {
				if len(queryString) > 0 {
					queryString = queryString + ", " + val
				} else {
					queryString = val
				}
			}
		}
	}

	d.modelFieldCache[tableName] = queryString

	return queryString, nil
}
