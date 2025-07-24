// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_main()
		{
		fn = FormulaNot
		n = FORMULATYPE.NUMBER
		b = FORMULATYPE.BOOLEAN

		.Check(fn, [b, ''], [b, true])
		.Check(fn, [b, false], [b, true])
		.Check(fn, [b, true], [b, false])

		.CheckError(fn, [n, ''], "Operation not supported")
		}
	}