// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	old_fmt: false
	Setup()
		{
		.old_fmt = Settings.Get('ShortDateFormat')
		}
	Test_main()
		{
		Settings.Set('ShortDateFormat', "yyyy-MM-dd")

		fn = FormulaCat
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE
		d = FORMULATYPE.DATE

		.Check(fn, [d, #20180101], [s, ' test'], [s, '2018-01-01 test'])
		.CheckError(fn, [d, ''], [s, ' test'], 'Invalid Value')
		.Check(fn, [r, '1 abc'], [s, ' test'], [s, '1 abc test'])
		.Check(fn, [u, '1 abc'], [s, ' test'], [s, '1 abc test'])
		.Check(fn, [n, 123], [s, ' test'], [s, '123 test'])
		.Check(fn, [n, ''], [s, ' test'], [s, ' test'])
		}
	Teardown()
		{
		if .old_fmt isnt false
			Settings.Set('ShortDateFormat', .old_fmt)
		super.Teardown()
		}
	}