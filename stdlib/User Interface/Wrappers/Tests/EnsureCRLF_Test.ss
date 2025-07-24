// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		fn = EnsureCRLF
		Assert(fn('') is: '')
		Assert(fn('a') is: 'a')
		Assert(fn('abc') is: 'abc')
		Assert(fn('\r') is: '\r\n')
		Assert(fn('\n') is: '\r\n')
		Assert(fn('\r\n') is: '\r\n')

		Assert(fn('\r\r') is: '\r\n\r\n')
		Assert(fn('\n\n') is: '\r\n\r\n')
		Assert(fn('\n\r') is: '\r\n\r\n')

		Assert(fn('\r\taaa\nbbbb\r\nccccc') is: '\r\n\taaa\r\nbbbb\r\nccccc')
		a = 'a'.Repeat(400)
		b = 'b'.Repeat(1000)
		Assert(fn('\n\t' $ a $ '\r' $ b $ '\nccccc')
			is: '\r\n\t' $ a $ '\r\n' $ b $ '\r\nccccc')
		Assert(fn('\r\n\t' $ a $ '\r\n' $ b $ '\r\n')
			is: '\r\n\t' $ a $ '\r\n' $ b $ '\r\n')
		}
	}