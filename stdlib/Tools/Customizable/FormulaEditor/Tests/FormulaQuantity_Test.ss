// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaQuantity
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		u = FORMULATYPE.UOM

		.Check(fn, [n, ''], [s, 'unit'], [u, ''])
		.Check(fn, [n, 1], [s, 'unit'], [u, '1 unit'])

		.CheckError(fn, [s, 'a'], [s, 'unit'], "QUANTITY Value must be a <Number>")
		.CheckError(fn, [n, 1], [n, 1], "QUANTITY Unit must be a <string>")
		}

	Test_Validate()
		{
		fn = FormulaQuantity.Validate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		u = FORMULATYPE.UOM

		Assert(fn([n], [s]) is: [u])

		Assert({ fn() } throws: "QUANTITY missing arguments")
		Assert({ fn([n], [s], [s], [n]) } throws: "QUANTITY too many arguments")
		Assert({ fn([s, n], [s]) } throws: "QUANTITY Value must be a <Number>")
		Assert({ fn([n], [s, n]) } throws: "QUANTITY Unit must be a <string>")
		}
	}