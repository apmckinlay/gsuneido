// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaDate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		d = FORMULATYPE.DATE

		.Check(fn, [s, '2018-05-11'], [s, 'yyyy-dd-MM'], [d, #20181105])
		.Check(fn, [s, 't'], [s, ''], [d, Date().NoTime()])
		.Check(fn, [s, 'h'], [s, ''], [d, Date().NoTime().EndOfMonth()])

		.Check(fn, [s, '05/11/2018'], [s, 'MMddyyyy'], [d, #20180511])
		.Check(fn, [s, '05/11/2018'], [s, 'ddMMyyyy'], [d, #20181105])
		.Check(fn, [s, 'May 11, 2018'], [s, 'MMddyyyy'], [d, #20180511])

		.CheckError(fn, [s, ''], [s, 'MMddyyyy'], 'DATE cannot convert')
		.CheckError(fn, [s, 'a'], [s, 'MMddyyyy'], 'DATE cannot convert')
		.CheckError(fn, [s, '05/11/2018'], [s, 'a'],
			'DATE cannot convert "05/11/2018" to a Date')
		.CheckError(fn, [s, '05/11/2018'], [s, ''],
			'DATE cannot convert "05/11/2018" to a Date')
		.CheckError(fn, [n, 1], [s, ''], "DATE Date must be a <String>")
		.CheckError(fn, [s, ''], [n, 1], "DATE Format must be a <String>")
		}

	Test_Validate()
		{
		fn = FormulaDate.Validate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		d = FORMULATYPE.DATE

		Assert(fn([s], [s]) is: [d])

		Assert({ fn([s]) } throws: "DATE missing arguments")
		Assert({ fn([s], [s], [s]) } throws: "DATE too many arguments")
		Assert({ fn([s, n], [s]) } throws: "DATE Date must be a <String>")
		Assert({ fn([s], [s, n]) } throws: "DATE Format must be a <String>")
		}
	}
