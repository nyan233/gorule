package gorule

import "testing"

func TestNumberEval(t *testing.T) {
	ExecuteSimpleBoolExpr("60 + 60 * 24")
}
