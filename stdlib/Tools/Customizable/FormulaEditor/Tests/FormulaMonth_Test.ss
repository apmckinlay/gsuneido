// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaMonth
		d = FORMULATYPE.DATE
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING

		.Check(fn, [d, #20180809], [n, 8])

		.CheckError(fn, [d, ''], 'Formula: Invalid Value')
		.CheckError(fn, [s, ''], "Formula: MONTH Field must be a <Date>")
		}

	Test_Validate()
		{
		fn = FormulaMonth.Validate
		d = FORMULATYPE.DATE
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING

		Assert(fn([d]) is: [n])

		Assert({ fn() } throws: "MONTH missing arguments")
		Assert({ fn([d], [d]) } throws: "MONTH too many arguments")
		Assert({ fn([s, d]) } throws: "MONTH Field must be a <Date>")
		}
	}