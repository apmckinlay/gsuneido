// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_main()
		{
		fn = FormulaIf
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN

		.Check(fn, [b, true], [n, 1], [s, 'a'], [n, 1])
		.Check(fn, [b, false], [n, 1], [s, 'a'], [s, 'a'])
		.Check(fn, [b, ''], [n, 1], [s, 'a'], [s, 'a'])
		.Check(fn, [b, false], [n, 1], [b, false], [n, 2], [b, false], [n, 3], [s, 'a'],
			[s, 'a'])
		.Check(fn, [b, false], [n, 1], [b, true], [n, 2], [b, true], [n, 3], [s, 'a'],
			[n, 2])

		.Check(fn, [b, false], [b, false], [b, true], [b, true])
		.Check(fn, [b, ''], [b, false], [b, true], [b, true])
		.Check(fn, [b, true], [b, false], [b, true], [b, false])

		.CheckError(fn, "IF must have at least 3 arguments")
		.CheckError(fn, [b, true], "IF must have at least 3 arguments")
		.CheckError(fn, [b, true], [n, 1], "IF must have at least 3 arguments")
		.CheckError(fn, [b, true], [n, 1], [b, true], [n, 2],
			"IF must have odd number of arguments")
		.CheckError(fn, [n, 1], [n, 2], [n, 3], "IF condition must be a <Boolean>")
		.CheckError(fn, [b, false], [n, 2], [n, 3], [n, 4], [n, 5],
			"IF condition must be a <Boolean>")

		// Test control flow
		Assert(fn({[type: b, value: true]}, {[type: n, value: 1]}, { throw 'Error' })
			is: [type: n, value: 1])
		Assert(fn({[type: b, value: false]}, {[type: n, value: 1]}, {[type: n, value: 2]})
			is: [type: n, value: 2])
		Assert({ fn({[type: b, value: false]}, {[type: n, value: 1]}, { throw 'Error' }) }
			throws: 'Error')

		Assert({ fn({[type: b, value: false]}, {[type: n, value: 1]}) }
			throws: 'IF must have at least 3 arguments')
		Assert({ fn({[type: b, value: false]}, {[type: n, value: 1]},
			{[type: b, value: true]}, {[type: n, value: 2]}) }
			throws: 'IF must have odd number of arguments')
		Assert({ fn({[type: s, value: false]}, {[type: n, value: 1]}, { throw 'Error' }) }
			throws: 'IF condition must be a <Boolean>')
		}

	Test_Validate()
		{
		fn = FormulaIf.Validate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN

		Assert(fn([b], [n], [n]) is: [n])
		Assert(fn([b], [n], [s]) is: [s, n])
		Assert(fn([b], [n, s], [b]) is: [b, n, s])
		Assert(fn([b], [n], [b, s]) is: [b, s, n])

		Assert({ fn([b], [n]) } throws: "IF must have at least 3 arguments")
		Assert({ fn([b], [n], [b], [n]) } throws: "IF must have odd number of arguments")
		Assert({ fn([b, n], [n], [b]) } throws: "IF condition must be a <Boolean>")
		}
	}
