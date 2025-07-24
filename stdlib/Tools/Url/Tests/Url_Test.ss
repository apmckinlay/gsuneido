// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Decode()
		{
		Assert(Url.Decode("") is: "")
		Assert(Url.Decode("hello world") is: "hello world")
		Assert(Url.Decode("%21hello+there+world%21") is: "!hello there world!")
		}

	Test_Split()
		{
		for scheme in #('', 'http')
			for user in #('', 'bob', 'sue:password')
				for host in #('', 'ibm.com', '127.0.0.1')
					for port in #('', 80)
						for path in #('', '/dir', '/dir/sub')
							for query in #('', 'age=1')
								for frag in #('', 'frag')
									.test(scheme, user, host, port, path, query, frag)
		Assert(Url.Split('mailto:bob@acme.com')
			is: #(scheme: 'mailto', user: 'bob', host: 'acme.com'))
		}

	test(scheme, user, host, port, path, query, frag)
		{
		if path is ''
			frag = ''
		x = Url.Split(Opt(scheme, '://') $ Opt(user, '@') $ host $ Opt(':', port) $
			path $ Opt('?', query) $ Opt('#', frag))
		Assert(x.GetDefault(#scheme, '') is: scheme)
		Assert(x.GetDefault(#host, '') is: host)
		Assert(x.GetDefault(#port, '') is: port)
		Assert(x.GetDefault(#path, '') is: path $ Opt('?', query))
		Assert(x.GetDefault(#basepath, '') is: path)
		Assert(x.GetDefault(#query, '') is: query)
		Assert(x.GetDefault(#fragment, '') is: frag)
		}

	Test_extractScheme()
		{
		ob = Object()
		scheme = Url.Url_extractScheme

		Assert(scheme("", ob) is: "")
		Assert(ob hasntMember: "scheme")
		Assert(scheme("www.url.com", ob) is: "www.url.com")
		Assert(ob hasntMember: "scheme")
		Assert(scheme("mailto:bob@hi.com", ob) is: "bob@hi.com")
		Assert(ob.scheme is: "mailto")
		Assert(scheme("http://www.test.com", ob) is: "www.test.com")
		Assert(ob.scheme is: "http")
		}

	Test_Encode()
		{
		Assert(Url.Encode('') is: '')
		Assert(Url.Encode('abc') is: 'abc')
		Assert(Url.Encode('hello world') is: 'hello+world')
		Assert(Url.Encode('a,b') is: 'a%2Cb')
		Assert(Url.Encode('\x05') is: '%05')
		Assert(Url.Encode('\xee') is: '%EE')
		Assert(Url.Encode('http://appserver:8080/Wiki?ProgrammingDepartment#1') is:
			'http://appserver:8080/Wiki?ProgrammingDepartment#1')

		Assert(Url.Encode("http://www.website.com/bobs'site",
			#(value: '15', query: '"table"'))
			is:	'http://www.website.com/bobs%27site?query=%22table%22&value=15')

		Assert(Url.Encode("<ab>", #()) is: "%3Cab%3E")
		}

	Test_EncodePreservePath()
		{
		Assert(Url.EncodePreservePath('') is: '')
		Assert(Url.EncodePreservePath('/') is: '/')
		Assert(Url.EncodePreservePath('a/bc') is: 'a/bc')
		Assert(Url.EncodePreservePath('a/bc/') is: 'a/bc/')
		Assert(Url.EncodePreservePath('hello/world') is: 'hello/world')
		Assert(Url.EncodePreservePath('hello/world/') is: 'hello/world/')
		Assert(Url.EncodePreservePath('start/a,b') is: 'start/a%2Cb')
		Assert(Url.EncodePreservePath('start/a,b/') is: 'start/a%2Cb/')
		Assert(Url.EncodePreservePath("<a/b>") is: "%3Ca/b%3E")
		Assert(Url.EncodePreservePath("<a/b>/") is: "%3Ca/b%3E/")
		Assert(Url.EncodePreservePath('start&/a,b') is: 'start%26/a%2Cb')
		Assert(Url.EncodePreservePath('start&/a,b/') is: 'start%26/a%2Cb/')
		}

	decodeEncodeTestCases: (
		("123", #(123))
		("0", #(0))
		("0123", #("0123"))
		("1&abc&n=3&s=xyz", #(1, 'abc', n: 3, s: 'xyz'))
		("=", #('': ''))
		("hello", #("hello"))
		("hello&world", #("hello", "world"))
		("age=25&text=hello%20there", #(age: 25, text: "hello there"))
		("hello%26world", #("hello&world"))
		("0026188e8211", #("0026188e8211"))
		("123456", #(123456))
		("003048617785", #("003048617785"))
		('one%26two&x%26y=1%3D2', #('one&two', 'x&y': '1=2'))
		('x%3Dy=123', #('x=y': 123))
		('a=1&b=2', #(a: 1, b: 2))
		('a=1&b=qq', #(a: 1, b: qq))
		('a=1&b=%2B', #(a: 1, b: '+'))
		)

	Test_SplitQuery()
		{
		for x in .decodeEncodeTestCases
			Assert(Url.SplitQuery(x[0]) is: x[1])
		}

	Test_BuildQuery()
		{
		for x in .decodeEncodeTestCases
			Assert(Url.BuildQuery(x[1]) is: "?" $ x[0])

		Assert(Url.BuildQuery(#()) is: "")
		}

	Test_EncodeValues()
		{
		for x in .decodeEncodeTestCases
			Assert(Url.EncodeValues(x[1]) is: x[0])

		Assert(Url.EncodeValues(#()) is: "")
		}

	Test_EncodeQueryValue()
		{
		encode = Url.EncodeQueryValue
		Assert(encode('') is: '')
		Assert(encode('abc') is: 'abc')
		Assert(encode('abc=') is: 'abc%3D')
		Assert(encode('abc==') is: 'abc%3D%3D')
		Assert(encode('=abc') is: '%3Dabc')
		Assert(encode('==') is: '%3D%3D')
		Assert(encode('hello world') is: 'hello%20world')
		Assert(encode('-h-e_l.l~o world-') is: '-h-e_l.l~o%20world-')
		Assert(encode('a,b') is: 'a%2Cb')
		Assert(encode('\x05') is: '%05')
		Assert(encode('\xee') is: '%EE')
		Assert(encode("Test\r\n Line") is: "Test%0D%0A%20Line")
		Assert(encode("Test=#:@?/&;!\"$%'()+,<>[\\]^`Line") is:
			"Test%3D%23%3A%40%3F%2F%26%3B%21%22%24%25%27%28%29%2B%2C%3C%3E%5B%5C%5D%5E" $
				"%60Line")
		}
	}