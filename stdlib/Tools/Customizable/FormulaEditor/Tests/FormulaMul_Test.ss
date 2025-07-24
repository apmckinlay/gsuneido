// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBase_Test
	{
	Test_one()
		{
		fn = FormulaMul
		n = FORMULATYPE.NUMBER
		s = FORMULATYPE.STRING
		u = FORMULATYPE.UOM
		r = FORMULATYPE.UOM_RATE

		.Check(fn, [n, 2.5], [n, 10], [n, 25])

		.Check(fn, [n, ''], [n, 10], [n, 0])
		.Check(fn, [n, 2.5], [u, '10 kg'], [u, '25 kg'])
		.Check(fn, [u, '10 kg'], [n, 2.5], [u, '25 kg'])
		.Check(fn, [u, '10 kg'], [r, '2.5 kg'], [n, 25])
		.Check(fn, [u, '10 kg'], [r, '2.5 g'], [n, 25000])
		.Check(fn, [r, '2.5 g'], [u, '10 kg'], [n, 25000])
		.Check(fn, [u, '10 kgs'], [r, '2.5 g'], [n, 25000])
		.Check(fn, [u, '10 kg'], [r, '2.5 grams'], [n, 25000])
		.Check(fn, [u, '10 kgs'], [r, '2.5 grams'], [n, 25000])

		.Check(fn, [u, '10 abc'], [r, '2.5 abc'], [n, 25])
		.Check(fn, [u, '10 abc'], [n, 2.5], [u, '25 abc'])
		.Check(fn, [u, '10 abc'], [n, ''], [u, '0 abc'])

		.CheckError(fn, [s, ''], [n, 10], 'Operation not supported')
		.CheckError(fn, [s, '123'], [n, 10], 'Operation not supported')
		.CheckError(fn, [u, ''], [n, 10], 'Invalid Value')
		.CheckError(fn, [n, 10], [u, ''], 'Invalid Value')
		.CheckError(fn, [u, ''], [r, "10 kg"], 'Invalid Value')
		.CheckError(fn, [u, '10 kg'], [r, ""], 'Invalid Value')
		.CheckError(fn, [u, '10 abc'], [u, '2.5 efg'], 'Operation not supported')
		.CheckError(fn, [u, '10 abc'], [r, '2.5 efg'], 'Incompatible unit of measure')
		}

	Test_ConvertValue()
		{
		fn = FormulaBase.ConvertValue
		fromValue = 10
		fromUom = 'kg'
		toUom = 'lb'
		rate? = false
		Assert(fn(fromValue, fromUom, toUom, rate?).Round(2) is: 22.05)
		fromUom = 'kgs'
		Assert(fn(fromValue, fromUom, toUom, rate?).Round(2) is: 22.05)
		toUom = 'lbs'
		Assert(fn(fromValue, fromUom, toUom, rate?).Round(2) is: 22.05)
		toUom = 'gram'
		Assert(fn(fromValue, fromUom, toUom, rate?) is: 10000)
		toUom = 'grams'
		Assert(fn(fromValue, fromUom, toUom, rate?) is: 10000)
		rate? = true
		Assert(fn(fromValue, fromUom, toUom, rate?) is: 0.01)
		toUom = 'liter'
		Assert(fn(fromValue, fromUom, toUom, rate?) is: false)
		toUom = 'liters'
		Assert(fn(fromValue, fromUom, toUom, rate?) is: false)
		toUom = ''
		Assert(fn(fromValue, fromUom, toUom, rate?) is: false)
		toUom = 'abc'
		Assert(fn(fromValue, fromUom, toUom, rate?) is: false)
		}
	}
