package database

import (
	"database/sql"
	"errors"
	"fmt"
	"main/log"
	"strconv"
)

type KeyTable struct {
	name   string
	driver *sql.DB
}

func NewKeyTable(name string, params ...string) *KeyTable {
	// params:
	// [0] dataSource
	// [1]
	return &KeyTable{
		name:   name,
		driver: db,
	}
}

func Serialize(value interface{}) string {
	switch array := value.(type) {
	case map[string]interface{}:
		var result string
		keys := make([]string, 0, len(array))
		values := make([]interface{}, 0, len(array))
		for i := 0; i < len(array); i++ {
			if i == len(array)-1 {
				result += fmt.Sprintf("%s:%v", keys[i], values[i])
			} else {
				result += fmt.Sprintf("%s:%v,", keys[i], values[i])
			}
		}
		return result
	case []interface{}:
		var result string
		for i := 0; i < len(array); i++ {
			if i == len(array)-1 {
				result += fmt.Sprintf("%v", array[i])
			} else {
				result += fmt.Sprintf("%v,", array[i])
			}
		}
		return result
	default:
		return ""
	}
}
func DeserializeMap(serializedData string) {
	var result map[string]interface{}
	for i := 0; i < len(serializedData); i++ {
		key := ""
		value := ""
		result[key] = value
	}
}
func DeserializeSlice(serializedData string) {
	for i := 0; i < len(serializedData); i++ {

	}

}

func value(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case []interface{}: // Serialize s
		return fmt.Sprintf("'%v'", v[0])
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'E', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (table *KeyTable) GetString(keyName interface{}, key interface{}, columns ...string) (string, bool, error) {
	var result string
	err, found := table.Get(keyName, key, columns, []interface{}{&result})
	if err != nil {
		return "", found, err
	}
	return result, found, nil
}
func (table *KeyTable) GetInt(keyName interface{}, key interface{}, columns ...string) (int, bool, error) {
	var result int
	err, found := table.Get(keyName, key, columns, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (table *KeyTable) GetInt64(keyName interface{}, key interface{}, columns ...string) (int64, bool, error) {
	var result int64
	err, found := table.Get(keyName, key, columns, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}
func (table *KeyTable) GetFloat(keyName interface{}, key interface{}, columns ...string) (float64, bool, error) {
	var result float64
	err, found := table.Get(keyName, key, columns, []interface{}{&result})
	if err != nil {
		return 0, found, err
	}
	return result, found, nil
}

func (table *KeyTable) Get(keyName interface{}, key interface{}, columns []string, data []interface{}) (error, bool) {
	err, found := table.get(keyName, key, columns, data)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
	}
	return err, found
}
func (table *KeyTable) Set(keyName interface{}, key interface{}, columns []string, data []interface{}) error {
	err := table.set(keyName, key, columns, data)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
	}
	return err
}
func (table *KeyTable) Put(columns []string, values []interface{}) error {
	err := table.put(columns, values)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
	}
	return err
}
func (table *KeyTable) Remove(keyName interface{}, key interface{}) error {
	err := table.remove(keyName, key)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
	}
	return err
}

func (table *KeyTable) get(keyName interface{}, key interface{}, columns []string, data []interface{}) (error, bool) {
	query := fmt.Sprintf("SELECT %s FROM `%s` WHERE `%v` = %s;", columnSliceToString(columns...), table.name, keyName, value(key))
	rows, err := table.driver.Query(query)
	if err != nil {
		return err, false
	}
	if rows.Next() {
		err := rows.Scan(data...)
		if err != nil {
			return err, true
		}
		rows.Close()
	} else {
		rows.Close()
		return nil, false
	}
	return nil, true
}
func (table *KeyTable) put(columns []string, values []interface{}) error {
	// `%s` = ?
	if len(columns) != len(values) {
		return errors.New("keyTable.Put : len(columns) != len(data) ")
	}
	columnsString := ""
	valuesString := ""
	for i := 0; i < len(columns); i++ {
		if i == len(columns)-1 {
			columnsString += fmt.Sprintf("`%s`", columns[i])
		} else {
			columnsString += fmt.Sprintf("`%s`, ", columns[i])
		}
	}
	for i := 0; i < len(values); i++ {
		if i == len(values)-1 {
			valuesString += fmt.Sprintf("%s", value(values[i]))
		} else {
			valuesString += fmt.Sprintf("%s, ", value(values[i]))
		}
	}
	query := fmt.Sprintf("INSERT INTO `%s` (%s) values (%s);", table.name, columnsString, valuesString)
	log.InfoLogger.Println(query)
	_, err := table.driver.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (table *KeyTable) set(keyName interface{}, key interface{}, columns []string, values []interface{}) error {
	// `%s` = ?
	if len(columns) != len(values) {
		return errors.New("keyTable.Set : len(columns) != len(values) ")
	}
	s := ""
	for i := 0; i < len(columns); i++ {
		if i == len(columns)-1 {
			s += fmt.Sprintf("`%s` = %s", columns[i], value(values[i]))
		} else {
			s += fmt.Sprintf("`%s` = %s, ", columns[i], value(values[i]))
		}
	}
	query := fmt.Sprintf("UPDATE `%s` SET %s WHERE `%s` = %s;", table.name, s, keyName, value(key))
	_, err := table.driver.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (table *KeyTable) remove(keyName interface{}, key interface{}) error {
	query := fmt.Sprintf("DELETE FROM `%s` WHERE `%s` = %s;", table.name, keyName, value(key))
	_, err := table.driver.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
