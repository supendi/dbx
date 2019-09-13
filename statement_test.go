package dbx

import (
	"context"
	"testing"
)

func TestAddParameter(t *testing.T) {
	statement := NewStatement(context.Background(), "SELECT * FROM unit_testx WHERE id=:id")
	statement.AddParameter("id", 1)

	if len(statement.Parameters) != 1 {
		t.Error("Parameters must have 1 element")
	}
}
