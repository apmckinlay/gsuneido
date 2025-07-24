// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaAdd
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE
		d = FORMULATYPE.DATE

		.Check(fn, [n, 100], [n, 10], [n, 110])

		.Check(fn, [u, '1 kg'], [u, '100 g'], [u, '1.1 kg'])
		.Check(fn, [r, '1 kg'], [r, '100 g'], [r, '100001 kg'])
		.Check(fn, [u, '100 g'], [u, '1 kg'], [u, '1100 g'])
		.Check(fn, [r, '100 g'], [r, '1 kg'], [r, '100 g']) // .001 is omitted
		.Check(fn, [u, '100 abc'], [u, '10 abc'], [u, '110 abc'])
		.Check(fn, [r, '100 abc'], [r, '10 abc'], [r, '110 abc'])

		.Check(fn, [u, '10 hrs'], [u, '10 days'], [u, '250 hrs'])
		.Check(fn, [u, '10 months'], [u, '1 year'], [u, '22 months'])

		.CheckError(fn, [d, ''], [u, '10 days'], 'Invalid Value')
		.CheckError(fn, [d, #20180101], [u, ''], 'Invalid Value')
		.CheckError(fn, [d, #20180101], [u, 'days'], 'Invalid Value')
		.Check(fn, [d, #20180101], [u, '10 days'], [d, #20180111])
		.Check(fn, [d, #20180101], [u, '1 month'], [d, #20180201])
		.Check(fn, [d, #20180101], [u, '30 mins'], [d, #20180101.0030])
		.Check(fn, [d, #20180101], [u, '-30 mins'], [d, #20171231.2330])
		.Check(fn, [u, '-30 mins'], [d, #20180101], [d, #20171231.2330])

		.CheckError(fn, [s, '100'], [n, 10], 'Operation not supported')
		.CheckError(fn, [u, '100 abc'], [r, '10 abc'], 'Operation not supported')
		.CheckError(fn, [u, '100 abc'], [n, 10], 'Operation not supported')
		.CheckError(fn, [u, '100 abc'], [s, '10'], 'Operation not supported')
		.CheckError(fn, [r, '100 abc'], [n, 10], 'Operation not supported')
		.CheckError(fn, [r, '100 abc'], [s, '10'], 'Operation not supported')
		.CheckError(fn, [d, #20180101], [d, #20180102], 'Operation not supported')
		.CheckError(fn, [u, '10 abc'], [u, '2.5 efg'], 'Incompatible unit of measure')
		}
	}