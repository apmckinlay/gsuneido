// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaAfterlast
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN

		.Check(fn, [s, ''], [s, ''], [s, ''])
		.Check(fn, [s, 'test'], [s, ''], [s, ''])
		.Check(fn, [s, ''], [s, 'test'], [s, ''])

		.Check(fn, [s, 'Test'], [s, ' '], [s, ''])
		.Check(fn, [s, 'ttt Test'], [s, ' '], [s, 'Test'])
		.Check(fn, [s, 'ttt Test test'], [s, ' '], [s, 'test'])
		.Check(fn, [s, 'ttt Test test'], [s, 'ttt'], [s, ' Test test'])
		.Check(fn, [s, 'ttt Test test'], [s, 'bbb'], [s, ''])

		.CheckError(fn, [b, true], [s, 'test'],
			"Formula: AFTERLAST text must be a <String>")
		.CheckError(fn, [s, 'test'], [n, 10],
			"Formula: AFTERLAST delimiter must be a <String>")
		.CheckError(fn, [s, 0], [s, ' '],
			'Formula: AFTERLAST failed to extract substring')
		}
	}
