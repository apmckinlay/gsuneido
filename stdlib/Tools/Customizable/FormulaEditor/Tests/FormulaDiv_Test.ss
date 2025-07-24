// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaDiv
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		.Check(fn, [n, ''], [n, 10], [n, 0])

		.CheckError(fn, [u, ''], [n, 10], 'Invalid Value')
		.CheckError(fn, [n, 10], [u, ''], 'Invalid Value')
		.CheckError(fn, [u, ''], [u, "10 kg"], 'Invalid Value')
		.CheckError(fn, [u, '10 kg'], [u, ""], 'Invalid Value')

		.Check(fn, [n, 2.5], [n, 10], [n, .25])

		.Check(fn, [n, 2.5], [u, '10 kg'], [r, '.25 kg'])
		.Check(fn, [u, '10 kg'], [n, 2.5], [u, '4 kg'])
		.Check(fn, [u, '10 kg'], [u, '2.5 kg'], [n, 4])
		.Check(fn, [r, '2.5 g'], [r, '10 kg'], [n, 250])

		.Check(fn, [u, '10 abc'], [u, '2.5 abc'], [n, 4])
		.Check(fn, [u, '10 abc'], [n, 2.5], [u, '4 abc'])

		.CheckError(fn, [s, ''], [n, 10], 'Operation not supported')
		.CheckError(fn, [s, '123'], [n, 10], 'Operation not supported')
		.CheckError(fn, [u, '10 abc'], [u, '2.5 efg'], 'Incompatible unit of measure')
		.CheckError(fn, [u, '10 abc'], [r, '2.5 efg'], 'Operation not supported')
		}
	}