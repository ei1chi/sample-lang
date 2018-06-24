package eval

import (
	"fmt"

	"github.com/ei1chi/sample-lang/ast"
	"github.com/ei1chi/sample-lang/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Stmts)

	case *ast.BlockStmt:
		return evalBlockStmt(node.Stmts)

	case *ast.ReturnStmt:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.IfExpr:
		return evalIfExpr(node)

	case *ast.ExprStmt:
		return Eval(node.Expr)

	case *ast.IntLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBooleanObject(node.Value)

	case *ast.PrefixExpr:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(node.Operator, right)

	case *ast.InfixExpr:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpr(node.Operator, left, right)
	}

	return nil
}

func evalProgram(stmts []ast.Stmt) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)

		// BlockStmts もしくは ProgramStmts の中で、
		// retrun 文が来たら中断して上流に返す
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStmt(stmts []ast.Stmt) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)

		// BlockStmts の中で、
		// retrun 文が来たら中断して上流に返す
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE || rt == object.ERROR {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpr(ope string, right object.Object) object.Object {
	switch ope {
	case "!":
		return evalBangOperatorExpr(right)
	case "-":
		return evalMinusPrefixOperatorExpr(right)
	default:
		return newError("unknown operator: %s%s", ope, right.Type())
	}
}

func evalInfixExpr(ope string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpr(ope, left, right)
	case ope == "==":
		return nativeBooleanObject(left == right)
	case ope == "!=":
		return nativeBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), ope, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), ope, right.Type())
	}
}

func nativeBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalBangOperatorExpr(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpr(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIntegerInfixExpr(ope string, left, right object.Object) object.Object {
	lval := left.(*object.Integer).Value
	rval := right.(*object.Integer).Value

	switch ope {
	case "+":
		return &object.Integer{Value: lval + rval}
	case "-":
		return &object.Integer{Value: lval - rval}
	case "*":
		return &object.Integer{Value: lval * rval}
	case "/":
		return &object.Integer{Value: lval / rval}

	case "<":
		return nativeBooleanObject(lval < rval)
	case ">":
		return nativeBooleanObject(lval > rval)
	case "==":
		return nativeBooleanObject(lval == rval)
	case "!=":
		return nativeBooleanObject(lval != rval)
	}
	return newError("unknown operator: %s %s %s", left.Type(), ope, right.Type())
}

func evalIfExpr(ie *ast.IfExpr) object.Object {
	cond := Eval(ie.Cond)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(ie.Cons)
	} else if ie.Alt != nil {
		return Eval(ie.Alt)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}
