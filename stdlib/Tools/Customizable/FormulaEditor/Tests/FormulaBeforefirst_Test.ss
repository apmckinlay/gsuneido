// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaBeforefirst
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN

		.Check(fn, [s, ''], [s, ''], [s, ''])
		.Check(fn, [s, 'test'], [s, ''], [s, ''])
		.Check(fn, [s, ''], [s, 'test'], [s, ''])

		.Check(fn, [s, 'Test'], [s, ' '], [s, 'Test'])
		.Check(fn, [s, 'ttt Test'], [s, ' '], [s, 'ttt'])
		.Check(fn, [s, 'ttt Test test'], [s, ' '], [s, 'ttt'])
		.Check(fn, [s, 'ttt Test test'], [s, 'test'], [s, 'ttt Test '])
		.Check(fn, [s, 'ttt Test test'], [s, 'bbb'], [s, 'ttt Test test'])

		.CheckError(fn, [b, true], [s, 'test'],
			"Formula: BEFOREFIRST text must be a <String>")
		.CheckError(fn, [s, 'test'], [n, 10],
			"Formula: BEFOREFIRST delimiter must be a <String>")
		.CheckError(fn, [s, 0], [s, ' '],
			'Formula: BEFOREFIRST failed to extract substring')
		}
	}
