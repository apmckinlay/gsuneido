// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(H('') is: '')
		Assert(H('now & then') is: 'now &amp; then')
		img = HtmlString('<img src="image.png"/>')
		Assert(H(img) is: img.Value())
		}
	}