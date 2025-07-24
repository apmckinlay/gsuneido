// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaConvert
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		.Check(fn, [u, '1000 g'], [s, 'kg'], [u, '1 kg'])
		.Check(fn, [u, '1000 grams'], [s, 'kg'], [u, '1 kg'])
		.Check(fn, [u, '1000 gram'], [s, 'kgs'], [u, '1 kgs'])
		.Check(fn, [u, '1000 grams'], [s, 'kgs'], [u, '1 kgs'])
		.Check(fn, [r, '1 g'], [s, 'kg'], [r, '1000 kg'])
		.Check(fn, [u, '10 abc'], [s, 'abc'], [u, '10 abc'])

		.CheckError(fn, [u, ''], [s, 'kg'], 'Invalid Value')
		.CheckError(fn, [u, 'kg'], [s, 'kg'], 'Invalid Value')
		.CheckError(fn, [u, '10'], [s, 'kg'],'Invalid Value')
		.CheckError(fn, [r, '1000 g'], [s, ''], 'CONVERT unit must not be empty')
		.CheckError(fn, [r, '1000 g'], [s, FormulaTonumber.EmptyPlaceHolder],
			'Incompatible unit of measure')
		.CheckError(fn, [n, 0], [s, 'kg'], "CONVERT Field must be a <Quantity> or <Rate>")
		.CheckError(fn, [u, '0 kg'], [n, 0], "CONVERT Unit must be a <String>")
		.CheckError(fn, [u, '0 kg'], [s, 'm'], 'Incompatible unit of measure')
		.CheckError(fn, [u, '10 abc'], [s, 'efg'], 'Incompatible unit of measure')
		}

	Test_Validate()
		{
		fn = FormulaConvert.Validate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		Assert(fn([u], [s]) is: [u])
		Assert(fn([r], [s]) is: [r])
		Assert(fn([u, r], [s]) is: [u, r])

		Assert({ fn() } throws: "CONVERT missing arguments")
		Assert({ fn([u]) } throws: "CONVERT missing arguments")
		Assert({ fn([r]) } throws: "CONVERT missing arguments")
		Assert({ fn([b], [n], [b], [n]) } throws: "CONVERT too many arguments")
		Assert({ fn([b, u], [s]) } throws: "CONVERT Field must be a <Quantity> or <Rate>")
		Assert({ fn([u], [s, n]) } throws: "CONVERT Unit must be a <String>")
		}
	}
