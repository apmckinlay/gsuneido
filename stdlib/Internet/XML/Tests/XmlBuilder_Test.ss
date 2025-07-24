// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.check('') { }
		.check('<p>hello</p>') { .p('hello') }
		.check('<p>hello</p>') { .p { 'hello' } }
		.check('<p>&lt;&amp;&gt;</p>') { .p('<&>') }
		.check('<p>&lt;&amp;&gt;</p>') { .p { '<&>' } }
		.check('<img src="pic.jpg" />') { .img(src: 'pic.jpg') }
		.check('<img src="&lt;&amp;&gt;" />') { .img(src: '<&>') }
		.check('<html><head><title>My Title</title></head><body>' $
			'<h1>First</h1><p>The first paragraph</p>' $
			'<h1>Second</h1><p>A second paragraph</p></body></html>')
			{
			.html
				{
				.head { .title { 'My Title' } }
				.body
					{
					.h1 { 'First' }
					.p { 'The first paragraph' }
					.h1 { 'Second' }
					.p { 'A second paragraph' }
					}
				}
			}
		.check('<p>hello<b>world</b></p>')
			{ .p { .S('hello'); .b('world') } }
		.check('<?xml version="1.0" encoding="us-ascii"?>')
			{ .Instruct() }
		.check('<?xml version="1.1" encoding="UTF-8"?>')
			{ .Instruct(version: '1.1', encoding: 'UTF-8') }
		.check('<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" ' $
			'"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">')
			{ .Declare() }
		.check('<!DECLARE symbol "string">')
			{ .Declare('DECLARE symbol', 'string') }
		.check('<!-- comment -->')
			{ .Comment('comment') }
		.check('<TestUpper>Test</TestUpper>')
			{ .TestUpper() { 'Test' } }
		}
	check(s, block)
		{
		b = XmlBuilder()
		b.Eval(block)
		Assert(b.ToString() is: s)
		}
	}