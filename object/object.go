package object

import "fmt"

type ObjectType string

const (
	INTEGER      = "INTEGER"
	BOOLEAN      = "BOOLEAN"
	NULL         = "NULL"
	RETURN_VALUE = "RETURN_VALUE"
	ERROR        = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (i *Integer) Type() ObjectType { return INTEGER }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN }

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL }

func (n *Null) Inspect() string { return "null" }

// ReturnValue 返り値を上流に返していくオブジェクト。
type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE }

func (r *ReturnValue) Inspect() string { return r.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR }

func (e *Error) Inspect() string { return "ERROR: " + e.Message }
