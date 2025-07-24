// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
function (expr, expected = '**notfalse**', exception = false)
	{
	if expected is 'throws'
		return Catch({ expr.Eval() }).Has?(exception) // Eval is okay here
	else if expected is '**notfalse**'
		return expr.Eval() isnt false // Eval is okay here
	else
		return expr.Eval() is expected.Compile() // Eval is okay here
	}