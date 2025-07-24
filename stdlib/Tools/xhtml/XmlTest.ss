// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(Xml('br') is: "<br />")
		Assert(Xml('p', "hello world") is: "<p>hello world</p>")
		Assert(Xml('p', #("hello world")) is: "<p>hello world</p>")
		Assert(Xml('p', "hello world", font: 'Arial', size: 14)
			is: '<p font="Arial" size="14">hello world</p>')
		Assert(Xml('p', #("hello world", font: 'Arial', size: 14))
			is: '<p font="Arial" size="14">hello world</p>')
		Assert(Xml("?works" data: "test.wks") is: '<?works data="test.wks"?>')
		Assert(Xml("!--", "this is a comment") is: "<!--this is a comment-->")
		}
	}

