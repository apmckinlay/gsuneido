// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		fn = CompatibleUOMFactor.Calc
		Assert(fn(1, 'kg', 'g') is: 1000)
		Assert(fn(1, 'g', 'kg') is: 0.001)
		Assert(fn(1, 'kg', 'g', rate?:) is: 0.001)
		Assert(fn(1, 'g', 'kg', rate?:) is: 1000)
		Assert(fn(10, 'kg', 'kg') is: 10)

		Assert(fn(1, 'kg', 'km') is: false)
		Assert(fn(1, 'kg', 'km', rate?:) is: false)
		}

	Test_add()
		{
		fn = CompatibleUOMFactor.Add
		qty1 = Object(value: 5, uom: 'miles')
		qty2 = Object(value: 8, uom: 'miles')
		Assert(fn(qty1, qty2) is: #(value: 13, uom: 'miles'))

		qty1 = Object(value: 5, uom: 'kilometers')
		qty2 = Object(value: 8, uom: 'miles')
		Assert(fn(qty1, qty2) is: #(value: 11.11, uom: 'miles'))

		qty1 = Object(value: 10, uom: 'lb')
		qty2 = Object(value: 15, uom: 'm3')
		Assert(fn(qty1, qty2) is: false)
		}
	Test_inverse()
		{
		for u in Uom_Conversions
			Assert(u.etauom_inverse is: 1 / u.etauom_size)
		}
	}