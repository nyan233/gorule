package gorule

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

/*
	布尔表达式求值器
*/

type Visitor func(args map[string]interface{}) interface{}

func boolEval(expr ast.Expr, args map[string]interface{}) bool {
	switch expr.(type) {
	case *ast.BinaryExpr:
		binExpr := expr.(*ast.BinaryExpr)
		if binExpr.Op == token.LAND {
			return boolEval(binExpr.X, args) && boolEval(binExpr.Y, args)
		} else if binExpr.Op == token.LOR {
			return boolEval(binExpr.X, args) || boolEval(binExpr.Y, args)
		} else {
			return boolFactorEval(binExpr, args)
		}
	case *ast.ParenExpr:
		parenExpr := expr.(*ast.ParenExpr)
		return boolEval(parenExpr.X.(*ast.BinaryExpr), args)
	default:
		return reflect.Indirect(reflect.ValueOf(visit(expr, args)(args))).Bool() == true
	}
}

func visit(expr ast.Expr, args map[string]interface{}) Visitor {
	switch expr.(type) {
	case *ast.BinaryExpr:
		return func(args map[string]interface{}) interface{} {
			return boolEval(expr, args)
		}
	case *ast.ParenExpr:
		return func(args map[string]interface{}) interface{} {
			parenExpr := expr.(*ast.ParenExpr)
			return boolEval(parenExpr.X, args)
		}
	case *ast.IndexExpr:
		indexExpr := expr.(*ast.IndexExpr)
		return func(args map[string]interface{}) interface{} {
			rawVal := visit(indexExpr.X, args)(args)
			rawValIv := reflect.Indirect(reflect.ValueOf(rawVal))
			if !rawValIv.IsValid() {
				panic(fmt.Sprintf("array is not valid"))
			}
			indexVal := visit(indexExpr.Index, args)(args)
			indexIv := reflect.Indirect(reflect.ValueOf(indexVal))
			if !indexIv.IsValid() {
				panic(fmt.Sprintf("index is not valid"))
			}
			switch rawValIv.Kind() {
			case reflect.Slice:
				return rawValIv.Index(int(indexIv.Int())).Interface()
			case reflect.Map:
				return rawValIv.MapIndex(indexIv).Interface()
			default:
				panic(fmt.Sprintf("IndexExpr : unknown visit type : %s", rawValIv.Type()))
			}
		}
	case *ast.CallExpr:
		callExpr := expr.(*ast.CallExpr)
		return func(args map[string]interface{}) interface{} {
			funIdent, ok := callExpr.Fun.(*ast.Ident)
			if !ok {
				panic("func name only is static")
			}
			funcName := funIdent.Name
			switch funcName {
			case "len":
				if len(callExpr.Args) == 0 {
					panic("args len is zero for call len")
				} else if len(callExpr.Args) > 1 {
					panic("args len overflow for call len")
				}
				return callInnerLen(visit(callExpr.Args[0], args)(args))
			case "cap":
				if len(callExpr.Args) == 0 {
					panic("args len is zero for call cap")
				} else if len(callExpr.Args) > 1 {
					panic("args len overflow for call cap")
				}
				return callInnerCap(visit(callExpr.Args[0], args)(args))
			default:
				panic(fmt.Sprintf("no support func name : %s", funcName))
			}
		}
	case *ast.SelectorExpr:
		selectorExpr := expr.(*ast.SelectorExpr)
		ident, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			panic(fmt.Sprintf("selector X is not ast.Ident : %s", reflect.TypeOf(selectorExpr.X)))
		}
		argName := ident.Name
		val, ok := args[argName]
		if !ok {
			panic(fmt.Sprintf("argument is not found : %s", argName))
		}
		return func(args map[string]interface{}) interface{} {
			iV := reflect.Indirect(reflect.ValueOf(val))
			field := iV.FieldByName(selectorExpr.Sel.Name)
			if !field.IsValid() {
				panic(fmt.Sprintf("%s.%s field is not found", argName, ident.Name))
			}
			return field.Interface()
		}
	case *ast.Ident:
		ident := expr.(*ast.Ident)
		return func(args map[string]interface{}) interface{} {
			return ident.Obj.Data
		}
	case *ast.BasicLit:
		bl := expr.(*ast.BasicLit)
		return func(args map[string]interface{}) interface{} {
			switch bl.Kind {
			case token.STRING:
				val2 := bl.Value[1 : len(bl.Value)-1]
				return strings.ReplaceAll(val2, "\\", "")
			case token.INT:
				var base int
				if len(bl.Value) > 2 {
					if bl.Value[0] == '0' {
						base = 8
					}
					switch bl.Value[:2] {
					case "0b":
						base = 2
					case "0x":
						base = 16
					default:
						base = 10
					}

				} else {
					base = 10
				}
				parseInt, err := strconv.ParseInt(bl.Value, base, 64)
				if err != nil {
					panic(err)
				}
				return parseInt
			case token.FLOAT:
				parseFloat, err := strconv.ParseFloat(bl.Value, 64)
				if err != nil {
					panic(err)
				}
				return parseFloat
			case token.CHAR:
				return bl.Value[1 : len(bl.Value)-1]
			default:
				panic(fmt.Sprintf("unknown kind : %s", bl.Kind))
			}
		}
	default:
		panic(fmt.Sprintf("visit unknown expr type : %s", reflect.TypeOf(expr)))
	}
}

func boolFactorEval(binExpr *ast.BinaryExpr, args map[string]interface{}) bool {
	actual := visit(binExpr.X, args)(args)
	expected := visit(binExpr.Y, args)(args)
	baseType := reflect.TypeOf(actual)
	switch baseType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare[int64](binExpr.Op, reflect.ValueOf(expected).Int(), reflect.ValueOf(actual).Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return compare[uint64](binExpr.Op, reflect.ValueOf(expected).Uint(), reflect.ValueOf(actual).Uint())
	case reflect.Float32, reflect.Float64:
		return compare[float64](binExpr.Op, reflect.ValueOf(expected).Float(), reflect.ValueOf(actual).Float())
	case reflect.String:
		return compare[string](binExpr.Op, expected.(string), actual.(string))
	case reflect.Bool:
		switch binExpr.Op {
		case token.LAND:
			return reflect.ValueOf(expected).Bool() && reflect.ValueOf(actual).Bool()
		case token.LOR:
			return reflect.ValueOf(expected).Bool() || reflect.ValueOf(actual).Bool()
		default:
			return reflect.ValueOf(expected).Bool() == reflect.ValueOf(actual).Bool()
		}
	default:
		panic(fmt.Sprintf("unknown type : %s", baseType.Kind()))
	}
}

type CompareType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~string | ~float32 | ~float64
}

func compare[T CompareType](op token.Token, expected, actual T) bool {
	switch op {
	case token.EQL:
		// ==
		return expected == actual
	case token.LSS:
		// <
		return expected < actual
	case token.GTR:
		// >
		return expected > actual
	case token.NEQ:
		// !=
		return expected != actual
	case token.LEQ:
		// <=
		return expected <= actual
	case token.GEQ:
		// >=
		return expected >= actual
	default:
		return false
	}
}
