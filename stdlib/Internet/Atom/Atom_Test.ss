// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		a = new Atom
		a.Feed(title: 'My Feed', link: 'http:/mysite.com/myfeed.html',
			id: 1, updated: #20090425.123456.Format(Atom.DateFormat))
		a.AddEntry(title: 'Entry One', link: 'http:/mysite.com/myentry1.html',
			author: 'Andrew', content: 'my text',
			id: 2, updated: #20090425.123456.Format(Atom.DateFormat))
		Assert(.eg like: a.ToString())
		}
eg: '<?xml version="1.0" encoding="us-ascii"?>
<feed>
	<title>My Feed</title>
	<link href="http:/mysite.com/myfeed.html" rel="alternate" type="text/html" />
	<id>1</id>
	<updated>2009-04-25T12:34:56Z</updated>
<entry>
	<title>Entry One</title>
	<link href="http:/mysite.com/myentry1.html" rel="alternate" type="text/html" />
	<author><name>Andrew</name></author>
	<id>2</id>
	<updated>2009-04-25T12:34:56Z</updated>
	<content type="text">my text</content>
</entry>
</feed>
'
	}