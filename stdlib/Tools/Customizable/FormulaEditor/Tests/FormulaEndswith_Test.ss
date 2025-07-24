// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaEndswith
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		.Check(fn, [s, ''], [s, ''], [b, true])
		.Check(fn, [s, ''], [s, 'test'], [b, false])
		.Check(fn, [s, 'Test'], [s, 'test'], [b, true])
		.Check(fn, [s, 'ttt Test'], [s, 'test'], [b, true])
		.Check(fn, [s, 'Test ttt'], [s, 'test'], [b, false])
		.Check(fn, [s, 'T(est'], [s, '(est'], [b, true])

		.Check(fn, [u, '1 unit'], [s, 'unit'], [b, true])
		.Check(fn, [r, '1 unit'], [s, 'test'], [b, false])
		.Check(fn, [s, 'ttt 1 unit ttt'], [u, '1 unit'], [b, false])
		.Check(fn, [s, 'ttt 2 unit ttt'], [r, '1 unit'], [b, false])

		.CheckError(fn, [b, true], [s, 'test'],
			"Formula: ENDSWITH Text must be a <String>, <Quantity> or <Rate>")
		.CheckError(fn, [s, 'test'], [n, 10],
			"Formula: ENDSWITH Substring must be a <String>, <Quantity> or <Rate>")
		}

	Test_Validate()
		{
		fn = FormulaEndswith.Validate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		Assert(fn([s], [s]) is: [b])
		Assert(fn([s, u, r], [s, u, r]) is: [b])

		Assert({ fn([s]) } throws: "ENDSWITH missing arguments")
		Assert({ fn([s], [s], [s]) } throws: "ENDSWITH too many arguments")
		Assert({ fn([s, n], [s]) }
			throws: "ENDSWITH Text must be a <String>, <Quantity> or <Rate>")
		Assert({ fn([s], [s, n]) }
			throws: "ENDSWITH Substring must be a <String>, <Quantity> or <Rate>")
		}
	}