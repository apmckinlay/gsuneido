// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaYear
		d = FORMULATYPE.DATE
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING

		.Check(fn, [d, #20180809], [n, 2018])

		.CheckError(fn, [d, ''], 'Formula: Invalid Value')
		.CheckError(fn, [s, ''], "Formula: YEAR Field must be a <Date>")
		}

	Test_Validate()
		{
		fn = FormulaYear.Validate
		d = FORMULATYPE.DATE
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING

		Assert(fn([d]) is: [n])

		Assert({ fn() } throws: "YEAR missing arguments")
		Assert({ fn([d], [d]) } throws: "YEAR too many arguments")
		Assert({ fn([s, d]) } throws: "YEAR Field must be a <Date>")
		}
	}