// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		_x = 0
		_y = 1
		.f()
		Assert(_x is: 0)
		Assert(_y is: 1)
		Assert({ _z } throws: "uninitialized" )
		}
	f()
		{
		Assert(_x is: 0)
		Assert(_y is: 1)
		_y = 2
		Assert({ _z } throws: "uninitialized" )
		_z = 3
		Assert(_y is: 2)
		.g()
		Assert(_x is: 0)
		Assert(_y is: 2)
		Assert(_z is: 3)
		}
	g()
		{
		Assert(_x is: 0)
		Assert(_y is: 2)
		Assert(_z is: 3)
		}

	Test_blocks()
		{
		_x = 0
		function () { _x = 1; b = { return }; function (b) { b(); _x = 5 }(b) }()
		Assert(_x is 0)
		}

	Test_gsuneido_bug()
		{
		Assert(.a())
		}
	a()
		{
		_singleLine = true
		return .b()
		}
	b()
		{
		if Random(100) is false
			_singleLine = 123 // in locals but uninit
		return .c()
		}
	c(_singleLine = false)
		{
		return singleLine
		}
	}