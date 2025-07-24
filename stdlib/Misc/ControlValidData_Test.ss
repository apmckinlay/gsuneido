// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		field = .MakeDatadict(baseClass: 'Field_boolean', control: 'CheckBox')
		fn = ControlValidData?
		rec = []

		rec[field] = true
		Assert(fn(rec, field))
		Assert(not fn(rec, field, value: #()))
		Assert(not fn(rec, field, value: 'abc'))

		rec[field] = 'abc'
		Assert(not fn(rec, field))
		Assert(fn(rec, field, value: false))
		}
	}