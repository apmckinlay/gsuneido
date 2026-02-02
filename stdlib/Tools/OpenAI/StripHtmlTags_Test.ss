// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		fn = StripHtmlTags
		Assert(fn("<p>Hello <b>world</b>!</p>") is: "Hello world!")
		Assert(fn('Test<br/>line') is: "Test\nline")
		Assert(fn('<div class="test"><span>Hi</span> <em>there</em></div>')
			is: "Hi there")
		Assert(fn("<ul>
	<li>One</li>
	<li>Two</li>
	</ul>") is: "One\nTwo")
		Assert(fn("<h1>Title</h1><p>Para.</p>") is: "Title\nPara.")
		Assert(fn("No tags here!") is: "No tags here!")
		Assert(fn("") is: "")
		}
	}