// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		router = RackRouter([['get', '/prefix', { }],
			['get', '/file', {}]])
		f = router.RackRouter_find_route
		Assert(f(#(method: 'GET', path: 'foo')) is: false)
		Assert(f(#(method: 'GET', path: 'x/prefix')) is: false)
		Assert(f(#(method: 'PUT', path: '/prefix')) is: false)
		Assert(f(#(method: 'GET', path: '/prefix')) isnt: false)
		Assert(f(#(method: 'GET', path: '/prefix/more')) isnt: false)

		Assert(f(#(method: 'GET', path: '/fileSomethingElse')) is: false)
		Assert(f(#(method: 'GET', path: '/file/seperated')) isnt: false)
		}
	}