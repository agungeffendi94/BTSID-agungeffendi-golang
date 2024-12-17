package helpers

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func mapBytesToString(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		if b, ok := v.([]byte); ok {
			m[k] = string(b)
		}
	}
	return m
}

func DatabaseQueryRows(db *sqlx.DB, query string, args ...interface{}) []map[string]interface{} {

	var datarows []map[string]interface{}
	rows, err := db.Queryx(query, args...)

	if err != nil {
		fmt.Println("Query Error", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			results := make(map[string]interface{})
			err = rows.MapScan(results)
			// fmt.Println(err)
			datarows = append(datarows, mapBytesToString(results))
		}
	}
	return datarows
}

func DatabaseQuerySingleRow(db *sqlx.DB, query string, args ...interface{}) map[string]interface{} {

	result := make(map[string]interface{})

	rows, err := db.Queryx(query, args...)

	if err != nil {
		fmt.Println("Query Error", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			results := make(map[string]interface{})
			err = rows.MapScan(results)
			return mapBytesToString(results)
		}
	}

	return result
}

func DatabaseQueryRowsE(db *sqlx.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {

	var datarows []map[string]interface{}
	rows, err := db.Queryx(query, args...)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			results := make(map[string]interface{})
			err = rows.MapScan(results)
			// fmt.Println(err)
			datarows = append(datarows, mapBytesToString(results))
		}
	}
	return datarows, err
}

func DatabaseQuerySingleRowE(db *sqlx.DB, query string, args ...interface{}) (map[string]interface{}, error) {

	result := make(map[string]interface{})

	rows, err := db.Queryx(query, args...)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			results := make(map[string]interface{})
			err = rows.MapScan(results)
			return mapBytesToString(results), err
		}
	}

	return result, err
}
