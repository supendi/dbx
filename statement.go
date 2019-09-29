package dbx

type SqlParameter struct {
	Name  string
	Value string
}

type Statement struct {
	SQL        string
	Parameters map[string]interface{}
}

//AddParameter add new parameter to sql statement
func (me *Statement) AddParameter(name string, value interface{}) {
	me.Parameters[name] = value
}

func NewStatement(sql string, params ...*SqlParameter) *Statement {
	var statement = &Statement{
		SQL: sql,
	}
	statement.Parameters = make(map[string]interface{})

	for _, param := range params {
		statement.Parameters[param.Name] = param.Value
	}
	return statement
}
