// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		_report = Object(Params: #(choosefields: ''))
		Assert(AccessCurrentPrint.AccessCurrentPrint_fields(#()) is: #())

		_report = Object(Params: #(choosefields: 'a,,b'))
		Assert(AccessCurrentPrint.AccessCurrentPrint_fields(#()) is: #(a, b))

		_report = Object(Params: #(choosefields: 'a,,b'))
		Assert(AccessCurrentPrint.AccessCurrentPrint_fields(#(a: field_a))
			is: #(field_a, b))
		}
	}