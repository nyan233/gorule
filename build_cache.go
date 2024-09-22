package gorule

import (
	"go/ast"
	"go/parser"
	"reflect"
	"sync"
	"unsafe"
)

var (
	gExprBuilder = newBuildCache()
)

type exprBuilder struct {
	rwmu       sync.RWMutex
	buildCache map[uintptr]cacheItem
}

type cacheItem struct {
	val     ast.Expr
	rawExpr string
}

func newBuildCache() *exprBuilder {
	return &exprBuilder{
		buildCache: make(map[uintptr]cacheItem, 32),
	}
}

func (bc *exprBuilder) BuildExpr(expr string) (ast.Expr, error) {
	header := (*reflect.StringHeader)(unsafe.Pointer(&expr))
	bc.rwmu.RLock()
	val, ok := bc.buildCache[header.Data]
	if ok && val.rawExpr == expr {
		bc.rwmu.RUnlock()
		return val.val, nil
	}
	bc.rwmu.RUnlock()
	bc.rwmu.Lock()
	defer bc.rwmu.Unlock()
	astVal, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, err
	}
	bc.buildCache[header.Data] = cacheItem{
		val:     astVal,
		rawExpr: expr,
	}
	return astVal, nil
}
