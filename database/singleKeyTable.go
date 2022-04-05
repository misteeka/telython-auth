package database

type SingleKeyTable struct {
	KeyTable *KeyTable
	key      string
}

func NewSingleKeyTable(name string, key string, params ...string) *SingleKeyTable {
	// params:
	// [0] dataSource
	// [1]
	return &SingleKeyTable{
		KeyTable: NewKeyTable(name, params...),
		key:      key,
	}
}

func columnSliceToString(columns ...string) string {
	result := ""
	for i := 0; i < len(columns); i++ {
		result += " `" + columns[i] + "`"
	}
	return result
}

func (table *SingleKeyTable) GetString(key interface{}, columns ...string) (string, bool, error) {
	return table.KeyTable.GetString(table.key, key, columns...)
}
func (table *SingleKeyTable) GetInt(key interface{}, columns ...string) (int, bool, error) {
	return table.KeyTable.GetInt(table.key, key, columns...)
}
func (table *SingleKeyTable) GetInt64(key interface{}, columns ...string) (int64, bool, error) {
	return table.KeyTable.GetInt64(table.key, key, columns...)
}
func (table *SingleKeyTable) GetFloat(key interface{}, columns ...string) (float64, bool, error) {
	return table.KeyTable.GetFloat(table.key, key, columns...)
}

func (table *SingleKeyTable) Get(key interface{}, columns []string, data []interface{}) (error, bool) {
	return table.KeyTable.Get(table.key, key, columns, data)
}
func (table *SingleKeyTable) Set(key interface{}, columns []string, data []interface{}) error {
	return table.KeyTable.Set(table.key, key, columns, data)
}

func (table *SingleKeyTable) SingleSet(key interface{}, column string, data interface{}) error {
	return table.KeyTable.Set(table.key, key, []string{column}, []interface{}{data})
}

func (table *SingleKeyTable) Put(key interface{}, columns []string, values []interface{}) error {
	return table.KeyTable.Put(append(columns, table.key), append(values, key))
}

func (table *SingleKeyTable) Remove(key interface{}) error {
	return table.KeyTable.Remove(table.key, key)
}
