// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_UrlEncodeValues()
		{
		Assert(AmazonAWS.UrlEncodeValues([param1: "SAMPLE_PARAMS"])
			is: 'param1=SAMPLE_PARAMS')
		Assert(AmazonAWS.UrlEncodeValues(#('abc\x9a123\x80'))
			is: "abc%C5%A1123%E2%82%AC")
		Assert(AmazonAWS.UrlEncodeValues(#('b', 'a', 'c'))
			is: "a&b&c")
		}
	}