// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaMatches
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		.Check(fn, [s, ''], [s, ''], [b, true])
		.Check(fn, [s, ''], [s, 'test'], [b, false])
		.Check(fn, [s, 'Test'], [s, 'test'], [b, false])
		.Check(fn, [s, 'Test'], [s, '(?i)test'], [b, true])
		.Check(fn, [s, 'Tes(t) ttt'], [s, '^(?i)(?q)tes(t)(?-q)'], [b, true])
		.Check(fn, [s, 'ttt Tes(t)'], [s, '(?i)(?q)tes(t)(?-q)$'], [b, true])

		.Check(fn, [u, '1 unit'], [s, 'unit'], [b, true])
		.Check(fn, [r, '1 unit'], [s, 'test'], [b, false])
		.Check(fn, [s, 'ttt 1 unit ttt'], [u, '1 unit'], [b, true])
		.Check(fn, [s, 'ttt 2 unit ttt'], [r, '1 unit'], [b, false])

		.CheckError(fn, [s, 'Test'], [s, 't('],
			'Formula: MATCHES failed to search "t(" from "Test"')
		.CheckError(fn, [b, true], [s, 'test'],
			"Formula: MATCHES Text must be a <String>, <Quantity> or <Rate>")
		.CheckError(fn, [s, 'test'], [n, 10],
			"Formula: MATCHES Substring must be a <String>, <Quantity> or <Rate>")
		}

	Test_Validate()
		{
		fn = FormulaMatches.Validate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		Assert(fn([s], [s]) is: [b])
		Assert(fn([s, u, r], [s, u, r]) is: [b])

		Assert({ fn([s]) } throws: "MATCHES missing arguments")
		Assert({ fn([s], [s], [s]) } throws: "MATCHES too many arguments")
		Assert({ fn([s, n], [s]) }
			throws: "MATCHES Text must be a <String>, <Quantity> or <Rate>")
		Assert({ fn([s], [s, n]) }
			throws: "MATCHES Substring must be a <String>, <Quantity> or <Rate>")
		}
	}