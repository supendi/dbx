package dbx

//SQLParameter represent the sql parameter
type SQLParameter struct {
	Name  string
	Value string
}

//Statement represent the SQL statement
type Statement struct {
	SQL        string
	Parameters map[string]interface{}
}

//AddParameter add new parameter to sql statement
func (me *Statement) AddParameter(name string, value interface{}) {
	me.Parameters[name] = value
}

//NewStatement returns new SQL statement instance
func NewStatement(sql string, params ...*SQLParameter) *Statement {
	var statement = &Statement{
		SQL: sql,
	}
	statement.Parameters = make(map[string]interface{})

	for _, param := range params {
		statement.Parameters[param.Name] = param.Value
	}
	return statement
}
