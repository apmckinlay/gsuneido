// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaUom
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		.Check(fn, [u, '1000 g'], [s, 'g'])
		.Check(fn, [r, '1000 g'], [s, 'g'])
		.Check(fn, [u, ' abc'], [s, 'abc'])
		.Check(fn, [u, '1000'], [s, ''])

		.CheckError(fn, [n, 0], "UOM Field must be a <Quantity> or <Rate>")
		}

	Test_Validate()
		{
		fn = FormulaUom.Validate
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		b = FORMULATYPE.BOOLEAN
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		Assert(fn([u]) is: [s])
		Assert(fn([r]) is: [s])
		Assert(fn([u, r]) is: [s])

		Assert({ fn() } throws: "UOM missing arguments")
		Assert({ fn([u], [n]) } throws: "UOM too many arguments")
		Assert({ fn([b, u]) } throws: "UOM Field must be a <Quantity> or <Rate>")
		}
	}