// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaStartswith
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN

		.Check(fn, [s, ''], [s, ''], [b, true])
		.Check(fn, [s, ''], [s, 'test'], [b, false])
		.Check(fn, [s, 'Test'], [s, 'test'], [b, true])
		.Check(fn, [s, 'Test ttt'], [s, 'test'], [b, true])
		.Check(fn, [s, 'ttt Test'], [s, 'test'], [b, false])
		.Check(fn, [s, 'T(est'], [s, 't('], [b, true])

		.CheckError(fn, [b, true], [s, 'test'],
			"Formula: STARTSWITH Text must be a <String>, <Quantity> or <Rate>")
		.CheckError(fn, [s, 'test'], [n, 10],
			"Formula: STARTSWITH Substring must be a <String>, <Quantity> or <Rate>")
		}

	Test_Validate()
		{
		fn = FormulaStartswith.Validate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		Assert(fn([s], [s]) is: [b])
		Assert(fn([s, u, r], [s, u, r]) is: [b])

		Assert({ fn([s]) } throws: "STARTSWITH missing arguments")
		Assert({ fn([s], [s], [s]) } throws: "STARTSWITH too many arguments")
		Assert({ fn([s, n], [s]) }
			throws: "STARTSWITH Text must be a <String>, <Quantity> or <Rate>")
		Assert({ fn([s], [s, n]) }
			throws: "STARTSWITH Substring must be a <String>, <Quantity> or <Rate>")
		}
	}