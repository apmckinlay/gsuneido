// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(MapBing.BuildUrl('100 Broadway', '', 'San Diego', 'CA', '') is:
			"http://www.bing.com/maps/?v=2&where1=100 Broadway , San Diego, CA&encType=1")

		Assert(MapBing.BuildUrl('100 Broadway', '', 'San Diego', 'CA',
			"42.217777778N,83.278888889W") is:
			"http://www.bing.com/maps/?v=2&where1=42.217777778N,83.278888889W&encType=1")

		Assert(MapBing.BuildUrl('', '', '', '', '', '90210') is:
			"http://www.bing.com/maps/?v=2&where1=90210&encType=1")
		}
	}