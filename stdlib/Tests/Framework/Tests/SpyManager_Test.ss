// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	fn1: `A(.a, .B = 1, _c = false, ._d = {})
		{
		return 1 + .b()
		}`
	fn1WithSpy: `A(.a, .B = 1, _c = false, ._d = {})
		{
		res = SpyManager().Spy(1, Object(:a, :b, :c, :d))
		if res.action is 'return' { return res.value }
		if res.action is 'throw' { throw res.value }
		return 1 + .b()
		}`
	fn2: `b() { return 1 }`
	fn2WithSpy: `b() {
		res = SpyManager().Spy(2, Object())
		if res.action is 'return' { return res.value }
		if res.action is 'throw' { throw res.value }
		return 1 }`
	fn3: `C(@args)
		{
		return args
		}`
	fn3WithSpy: `C(@args)
		{
		res = SpyManager().Spy(3, Object(:args))
		if res.action is 'return' { return res.value }
		if res.action is 'throw' { throw res.value }
		return args
		}`

	Test_buildCodeWithTdop()
		{
		fn = SpyManager.SpyManager_buildCodeWithTdop
		source = `class
			{
			C: class
				{
				` $ .fn1 $ `
				` $ .fn3 $ `
				}
			` $ .fn2 $ `
			}`
		expected = `class
			{
			C: class
				{
				` $ .fn1WithSpy $ `
				` $ .fn3WithSpy $ `
				}
			` $ .fn2WithSpy $ `
			}`
		code = fn(source, #(
			(Id: 1, Paths: ('C', 'A'), Params: "(a,b=1,_c=false,_d=[])"),
			(Id: 2, Paths: ('b'), Params: "()"),
			(Id: 3, Paths: ('C', 'C'), Params: "(@args)")))
		Assert(.removeEmptyLine(code) like: expected)

		Assert({ fn(source, #((Paths: ('D')))) }
			throws: "SpyOn cannot find specified method - D")
		Assert({ fn(source, #((Paths: ('C', 'D')))) }
			throws: "SpyOn cannot find specified method - C.D")
		}

	Test_buildCodeWithClassHelp()
		{
		fn = SpyManager.SpyManager_buildCodeWithClassHelp
		source = .fn1.Replace('^A', 'function ')
		expected = .fn1WithSpy.Replace('^A', 'function ')
		code = fn(source,
			#((Id: 1, Paths: (), Params: "(a,b=1,_c=false,_d=[])", Method?: false)))
		Assert(.removeEmptyLine(code) like: expected)

		source = `class
			{
			` $ .fn1 $ `
			` $ .fn2 $ `
			` $ .fn3 $ `
			}`
		expected = `class
			{
			` $ .fn1WithSpy $ `
			` $ .fn2WithSpy $ `
			` $ .fn3WithSpy $ `
			}`
		code = fn(source, #(
			(Id: 1, Paths: ('A'), Params: "(a,b=1,_c=false,_d=[])", Method?:),
			(Id: 2, Paths: ('b'), Params: "()", Method?:),
			(Id: 3, Paths: ('C'), Params: "(@args)", Method?:)))
		Assert(.removeEmptyLine(code) like: expected)

		}

	removeEmptyLine(s)
		{
		return s.Lines().RemoveIf(#Blank?).Join('\r\n')
		}
	}