// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaTonumber
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		.Check(fn, [u, '1000 g'], [s, 'kg'], [n, 1])
		.Check(fn, [r, '1 g'], [s, 'kg'], [n, 1000])
		.Check(fn, [u, '1000 grams'], [s, 'kg'], [n, 1])
		.Check(fn, [u, '1000 gram'], [s, 'kgs'], [n, 1])
		.Check(fn, [u, '1000 grams'], [s, 'kgs'], [n, 1])
		.Check(fn, [u, '10 abc'], [s, 'abc'], [n, 10])

		.Check(fn, [u, '10 abc'], [s, FormulaTonumber.EmptyPlaceHolder], [n, 10])

		.CheckError(fn, [u, ''], [s, 'kg'], 'Invalid Value')
		.CheckError(fn, [u, 'kg'], [s, 'kg'], 'Invalid Value')
		.CheckError(fn, [u, '10'], [s, 'kg'],'Invalid Value')
		.CheckError(fn, [n, 0], [s, 'kg'],
			"TONUMBER Field must be a <Quantity> or <Rate>")
		.CheckError(fn, [u, '0 kg'], [n, 0], "TONUMBER Unit must be a <String>")
		.CheckError(fn, [r, '1000 g'], [s, ''], 'TONUMBER unit must not be empty')
		.CheckError(fn, [u, '0 kg'], [s, 'm'], 'Incompatible unit of measure')
		.CheckError(fn, [u, '10 abc'], [s, 'efg'], 'Incompatible unit of measure')
		}

	Test_Validate()
		{
		fn = FormulaTonumber.Validate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		Assert(fn([u], [s]) is: [n])
		Assert(fn([r], [s]) is: [n])
		Assert(fn([u, r], [s]) is: [n])

		Assert({ fn() } throws: "TONUMBER missing arguments")
		Assert({ fn([u]) } throws: "TONUMBER missing arguments")
		Assert({ fn([r]) } throws: "TONUMBER missing arguments")
		Assert({ fn([b], [n], [b], [n]) } throws: "TONUMBER too many arguments")
		Assert({ fn([b, u], [s]) }
			throws: "TONUMBER Field must be a <Quantity> or <Rate>")
		Assert({ fn([u], [s, n]) } throws: "TONUMBER Unit must be a <String>")
		}
	}
