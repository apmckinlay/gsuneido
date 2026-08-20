// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
// object and record constants (not constructor calls)
SuJsTranslate
	{
	CallClass(ast)
		{
		Assert(ast.type in ('Object', 'Record'))
		t = ast.type.Lower()

		if ast.size is 0
			{
			.Print('su.empty_' $ t)
			return
			}

		.Print('su.mk' $ t.Capitalize() $ '(')
		first = true
		firstNamed = true
		for (i = 0; i < ast.size; i++)
			{
			if (first) first = false; else .Print(', ')
			if ast[i].named is false
				SuJsTranslateConstant(ast[i].value)
			else
				{
				if firstNamed
					{
					.Print('null, ')
					firstNamed = false
					}
				SuJsTranslateConstant(ast[i].key)
				.Print(', ')
				SuJsTranslateConstant(ast[i].value)
				}
			}
		.Print(')')
		}
	}
