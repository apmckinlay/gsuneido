// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaNeg
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE
		d = FORMULATYPE.DATE
		b = FORMULATYPE.BOOLEAN

		.Check(fn, [n, 100], [n, -100])
		.Check(fn, [u, '100 unit'], [u, '-100 unit'])
		.Check(fn, [r, '-100 unit'], [r, '100 unit'])

		.Check(fn, [n, ''], [n, 0])

		.CheckError(fn, [s, '100'], 'Operation not supported: "- <String>"')
		.CheckError(fn, [d, #20170101], 'Operation not supported: "- <Date>"')
		.CheckError(fn, [b, true], 'Operation not supported: "- <Boolean>"')
		}
	}