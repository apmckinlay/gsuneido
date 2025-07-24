// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		rss = Rss2()
		rss.Channel(title: 'My Channel',
			link: 'http://abc.com/mychannel.html',
			description: 'my channel of stuff')
		rss.Image(url: 'http://suneido.com/image.png', title: 'my image')
		rss.TextInput(title: 'Search', description: 'Search this feed',
				name: 'query', link: 'http://suneido.com/search.cgi')
		rss.AddItem(title: 'My Item',
			description: 'my item of stuff', guid: 1)
		rss.AddItem(title: 'Second Item',
			description: 'more stuff', guid: 2)
		Assert(rss.ToString() like: .eg)
		}
	eg: '<?xml version="1.0" encoding="us-ascii"?>
<rss version="2.0">
<channel>
<title>My Channel</title>
<link>http://abc.com/mychannel.html</link>
<description>my channel of stuff</description>
<image>
<url>http://suneido.com/image.png</url>
<title>my image</title>
<link>http://abc.com/mychannel.html</link>
</image>
<textInput>
<title>Search</title>
<description>Search this feed</description>
<name>query</name>
<link>http://suneido.com/search.cgi</link>
</textInput>
<item>
<title>My Item</title>
<description>my item of stuff</description>
<guid>1</guid>
</item>
<item>
<title>Second Item</title>
<description>more stuff</description>
<guid>2</guid>
</item>
</channel>
</rss>'
	Test_one()
		{
		rss = Rss2()
		rss.Channel(title: 'My Channel',
			link: 'http://abc.com/mychannel.html',
			description: 'my channel of stuff')
		Assert({ rss.Channel() }
			throws: "Rss2: only one channel allowed")
		rss.Image(url: 'http://suneido.com/image.png', title: 'my image')
		Assert({ rss.Image() }
			throws: "Rss2: only one image allowed")
		rss.TextInput(title: 'Search', description: 'Search this feed',
				name: 'query', link: 'http://suneido.com/search.cgi')
		Assert({ rss.TextInput() }
			throws: "Rss2: only one textInput allowed")
	}
	Test_required()
		{
		rss = Rss2()
		Assert({ rss.Image(title: 'image') }
			throws: "Rss2: missing element: image/url")
		}
	Test_item_must_have_title_or_description()
		{
		rss = Rss2()
		Assert({ rss.AddItem() }
			throws: "Rss2: item must have either title or description")
		}
	}