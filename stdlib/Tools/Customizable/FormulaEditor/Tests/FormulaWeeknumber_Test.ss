// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaWeeknumber
		d = FORMULATYPE.DATE
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING

		.Check(fn, [d, #20180809], [n, 32]) /*= WeekNumber*/

		.CheckError(fn, [d, ''], 'Formula: Invalid Value')
		.CheckError(fn, [s, ''], "Formula: WEEKNUMBER Field must be a <Date>")
		}

	Test_Validate()
		{
		fn = FormulaWeeknumber.Validate
		d = FORMULATYPE.DATE
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING

		Assert(fn([d]) is: [n])

		Assert({ fn() } throws: "WEEKNUMBER missing arguments")
		Assert({ fn([d], [d]) } throws: "WEEKNUMBER too many arguments")
		Assert({ fn([s, d]) } throws: "WEEKNUMBER Field must be a <Date>")
		}
	}