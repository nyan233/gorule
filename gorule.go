package gorule

type Argument struct {
	Name string
	Val  interface{}
}

func ExecuteSimpleBoolExpr(expr string, args ...Argument) (bool, error) {
	parseExpr, err := gExprBuilder.BuildExpr(expr)
	if err != nil {
		return false, err
	}
	argsMap := make(map[string]interface{}, len(args))
	for _, arg := range args {
		argsMap[arg.Name] = arg.Val
	}
	return boolEval(parseExpr, argsMap), nil
}
