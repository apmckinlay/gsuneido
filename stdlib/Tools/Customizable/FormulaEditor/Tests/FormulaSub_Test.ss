// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaSub
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE
		d = FORMULATYPE.DATE

		.Check(fn, [n, 100], [n, 10], [n, 90])

		.Check(fn, [u, '1 kg'], [u, '100 g'], [u, '.9 kg'])
		.Check(fn, [r, '100 kg'], [r, '1 g'], [r, '-900 kg'])
		.Check(fn, [u, '100 g'], [u, '1 kg'], [u, '-900 g'])
		.Check(fn, [r, '1 g'], [r, '100 kg'], [r, '.9 g'])
		.Check(fn, [u, '100 abc'], [u, '10 abc'], [u, '90 abc'])
		.Check(fn, [u, '10 abc'], [u, '100 abc'], [u, '-90 abc'])
		.Check(fn, [r, '100 abc'], [r, '10 abc'], [r, '90 abc'])

		.CheckError(fn, [d, #20180101.0102], [d, ''], 'Invalid Value')
		.Check(fn, [d, #20180101.0102], [d, #20180101.0101], [u, "60 secs"])
		.Check(fn, [d, #20180101.0101], [d, #20180101.0101], [u, "0 secs"])
		.Check(fn, [d, #20180101.0101], [d, #20180101.0102], [u, "-60 secs"])
		.Check(fn, [d, #20180101.0202], [d, #20180101.0101], [u, "3660 secs"])
		.Check(fn, [d, #21010103.0102], [d, #20180101.0101], [u, "2619388800 secs"])

		.Check(fn, [d, #20180102], [u, '1 day'], [d, #20180101])
		.Check(fn, [d, #20180102], [u, '-1 day'], [d, #20180103])

		.CheckError(fn, [s, '100'], [n, 10], 'Operation not supported')
		.CheckError(fn, [u, '-1 day'], [d, #20180102], 'Operation not supported')
		.CheckError(fn, [u, '100 abc'], [r, '10 abc'], 'Operation not supported')
		.CheckError(fn, [u, '100 abc'], [n, 10], 'Operation not supported')
		.CheckError(fn, [u, '100 abc'], [s, '10'], 'Operation not supported')
		.CheckError(fn, [r, '100 abc'], [n, 10], 'Operation not supported')
		.CheckError(fn, [r, '100 abc'], [s, '10'], 'Operation not supported')
		.CheckError(fn, [u, '10 abc'], [u, '2.5 efg'], 'Incompatible unit of measure')
		}
	}