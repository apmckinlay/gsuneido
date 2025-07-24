// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// this is not Memoize'd
	// since we should be skipping at a higher level if the code hasn't changed
	CallClass(code, fromStdlib?)
		{
		try
			ast = Suneido.Parse(code)
		catch
			return #() // ignore errors caused by editing code
		refs = Object().Set_default(0)
		stdNames = Qc_stdNames(fromStdlib?)
		.Traverse(ast, refs, stdNames)
		return refs
		}
	Traverse(ast, refs, stdNames, includeConstant? = false, extraFn = false)
		{
		if Type(ast) isnt 'AstNode'
			return
		switch ast.type
			{
		case #Ident:
			name = ast.name
			if name.Capitalized?() and not stdNames.Member?(name) and
				(includeConstant? or not name.Upper?())
				++refs[name]
		case #Class:
			name = ast.base
			if name.Capitalized?() and not stdNames.Member?(name)
				++refs[name]
		default:
			}
		if extraFn isnt false
			extraFn(ast, refs)

		for (i = 0; false isnt c = ast.children[i]; ++i)
			.Traverse(c, refs, stdNames, :includeConstant?, :extraFn) // recursive
		}
	}