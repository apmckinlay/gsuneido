// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaRate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		r = FORMULATYPE.UOM_RATE

		.Check(fn, [n, ''], [s, 'unit'], [r, ''])
		.Check(fn, [n, 1], [s, 'unit'], [r, '1 unit'])

		.CheckError(fn, [s, 'a'], [s, 'unit'], "RATE Value must be a <Number>")
		.CheckError(fn, [n, 1], [n, 1], "RATE Unit must be a <string>")
		}

	Test_Validate()
		{
		fn = FormulaRate.Validate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		r = FORMULATYPE.UOM_RATE

		Assert(fn([n], [s]) is: [r])

		Assert({ fn() } throws: "RATE missing arguments")
		Assert({ fn([n], [s], [s], [n]) } throws: "RATE too many arguments")
		Assert({ fn([s, n], [s]) } throws: "RATE Value must be a <Number>")
		Assert({ fn([n], [s, n]) } throws: "RATE Unit must be a <string>")
		}
	}