// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
SuJsTranslate // to get output stuff
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
			SuJsTranslateFunction(ast, outerName)
		case 'Object', 'Record':
			SuJsTranslateObject(ast)
		case 'Class':
			SuJsTranslateClass(ast, outerName)
		case 'Nary':
			.Value(AstFoldExpr(ast), isConst:)
		default:
			.Value(ast.value)
			}
		}
	}
