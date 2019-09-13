package dbx

import "context"

type Statement struct {
	context    context.Context
	SQL        string
	Parameters map[string]interface{}
}

//AddParameter add new parameter to sql statement
func (me *Statement) AddParameter(name string, value interface{}) {
	me.Parameters[name] = value
}

func NewStatement(context context.Context, sql string) *Statement {
	var statement = &Statement{
		context: context,
		SQL:     sql,
	}
	statement.Parameters = make(map[string]interface{})
	return statement
}
