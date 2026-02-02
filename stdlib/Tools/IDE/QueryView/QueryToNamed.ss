// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
// Used by QueryStrategyViewer (QueryView)
// to convert suitable queries to named argument form
// for QueryStrategy1 (Query1/Empty?/First/Last)
function(query)
	{
	if Client?()
		return ServerEval("QueryToNamed", query)
	ast = Query.Parse(query)
	if ast.type is 'table'
		return [ast.name]
	if ast.type isnt 'where'
		return [query]
	src = ast.source
	if src.type isnt 'table'
		return [query]
	args = [src.name]

	expr = ast.expr
	if expr.type isnt #Nary or expr.op isnt #And
		return [query]
	for i in ..expr.size
		{
		e = expr[i]
		try
			{
			if e.type isnt #Binary or e.op isnt #Is or
				e.lhs.type isnt #Ident or e.rhs.type isnt #Constant
				return [query]
			}
		catch (unused, 'member not found: "type"') // Build: future exe
			return [query] // this will be handled in a future exe; can remove then
		args[e.lhs.name] = e.rhs.value
		}
	return args
	}
