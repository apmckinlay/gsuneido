// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getAddress()
		{
		func = HttpLinkControl.MergePrefix
		Assert(func('') is: 'http://')
		Assert(func('www.google.ca') is: 'http://www.google.ca')
		Assert(func('http://www.google.ca')	is: 'http://www.google.ca')
		Assert(func('https://www.google.ca') is: 'https://www.google.ca')
		}
	}