// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.scan('')
		.scan('a b c')
		.scan('Func(a, b: c) { for x in y z += 123 }')
		}
	scan(text)
		{
		scan = Scanner(text)
		scanwc = ScannerWithContext(text)
		ob = Object()
		ahead = false
		while scanwc isnt token = scanwc.Next()
			{
			Assert(.next(scan) is: token)
			Assert(scan.Type() is: scanwc.Type())
			Assert(scan.Keyword?() is: scanwc.Keyword?())
			Assert(scan.Position() is: scanwc.Position())
			if ob.Size() > 0
				{
				Assert(ob.Last() is: scanwc.Prev())
				Assert(ahead is: token)
				}
			if ob.Size() > 1
				Assert(ob[ob.Size() - 2] is: scanwc.Prev2())
			ob.Add(token)
			ahead = scanwc.Ahead()
			}
		Assert(scan.Next() is: scan)
		}
	next(scan)
		{
		do
			token = scan.Next()
			while token isnt scan and token.White?()
		return token
		}
	}