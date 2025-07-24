// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(XmlParser('') is: false)
		.cases.Each
			{ .check(XmlParser(it), it) }
		.check(XmlParser('<img\r\n\tsrc="pic.jpg" />'), '<img src="pic.jpg" />')
		}
	cases:
		(
		'<html>hello</html>'
		'<img src="pic.jpg" />'
		'<body>
			<h1>First</h1>
			<p>The first paragraph</p>
			<h1>Second</h1>
			<p>A second paragraph</p>
		</body>'
		'<html>
			<head>
				<title>My Title</title>
			</head>
			<body>
				<h1>First</h1>
				<p>The first paragraph</p>
				<h1>Second</h1>
				<p>A second paragraph</p>
			</body>
		</html>'
		)

	Test_path()
		{
		xml = XmlParser(.cases[3])
		.check(xml.body[0], .cases[2])
		.check(xml.body[0][0], '<h1>First</h1>')
		.check(xml.body[0][3], '<p>A second paragraph</p>')
		.check(xml.body.h1, '<h1>First</h1><h1>Second</h1>')
		Assert(xml.Name() is: 'html')
		Assert(xml.Attributes() is: #())
		Assert(xml.Children().Map(#Name) is: #(head, body))
		Assert(xml.head.title.Text() is: 'My Title')
		Assert(xml.body.h1.Text() is: 'FirstSecond')

		Assert(XmlParser('<img src="pic.jpg" align="right" />').Attributes()
			is: #(src: "pic.jpg", align: "right"))
		}

	check(xml, s)
		{
		Assert(xml.ToString().Tr('\t\r\n') is: s.Tr('\t\r\n'))
		}

	Test_invalidXml()
		{
		Assert({ XmlParser('test invalid xml message') } throws: 'Invalid xml format')
		}
	}