// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_main()
		{
		fn = FormulaAnd
		n = FORMULATYPE.NUMBER
		b = FORMULATYPE.BOOLEAN

		.Check(fn, [b, ''], [b, ''], [b, false])
		.Check(fn, [b, ''], [b, false], [b, false])
		.Check(fn, [b, ''], [b, true], [b, false])
		.Check(fn, [b, false], [b, ''], [b, false])
		.Check(fn, [b, false], [b, false], [b, false])
		.Check(fn, [b, false], [b, true], [b, false])
		.Check(fn, [b, true], [b, ''], [b, false])
		.Check(fn, [b, true], [b, false], [b, false])
		.Check(fn, [b, true], [b, true], [b, true])

		.CheckError(fn, [b, ''], [n, ''], "Operation not supported")
		}
	}