// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
JsTranslate // to get output stuff
	{
	CallClass(ast, outerName = false)
		{
		if Type(ast) isnt 'AstNode'
			{
			.Value(ast, isConst:)
			return
			}

		switch (ast.type)
			{
		case 'Function':
			JsTranslateFunction(ast, outerName)
		case 'Object', 'Record':
			JsTranslateObject(ast)
		case 'Class':
			JsTranslateClass(ast, outerName)
		case 'Nary':
			.Value(AstFoldExpr(ast), isConst:)
		default:
			.Value(ast.value)
			}
		}
	}
