// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		response = `HTTP/1.1 301 Moved Permanently
Location: http://www.google.com/
Content-Type: text/html; charset=UTF-8
Date: Mon, 10 Feb 2020 14:50:28 GMT
Expires: Wed, 11 Mar 2020 14:50:28 GMT
Cache-Control: public, max-age=2592000
Server: gws
Content-Length: 219
X-XSS-Protection: 0
X-Frame-Options: SAMEORIGIN
Connection: close`
		response2 = response.Replace(`Date: Mon, 10 Feb 2020 14:50:28 GMT`,
			`date: Mon, 10 Feb 2020 14:50:28 GMT`)
		response3 = response.Replace(`Date: Mon, 10 Feb 2020 14:50:28 GMT`,
			`date: Mon, 10 Feb 2020 14:50:28 gmt`)
		.SpyOn(Http).Return(
			[content: "", header: response],
			[content: "", header: response],
			[content: "", header: response],
			[content: "", header: 'invalid'],
			[content: "", header: 'invalid'],
			[content: "", header: response2],
			[content: "", header: response2],
			[content: "", header: response2],
			[content: "", header: response3],
			[content: "", header: response3])
		Assert(Htp.InternetFormat() is: 'Mon, 10 Feb 2020 14:50:28 GMT')
		Assert(Htp.InternetFormatWithThrow() is: 'Mon, 10 Feb 2020 14:50:28 GMT')
		Assert(Htp.UnixTime() is: 1581346228)

		Assert(Htp.InternetFormat() is: false)
		Assert({ Htp.InternetFormatWithThrow() }
			throws: 'Expected date string but result was')

		Assert(Htp.InternetFormat() is: 'Mon, 10 Feb 2020 14:50:28 GMT')
		Assert(Htp.InternetFormatWithThrow() is: 'Mon, 10 Feb 2020 14:50:28 GMT')
		Assert(Htp.UnixTime() is: 1581346228)

		Assert(Htp.InternetFormat() is: false)
		Assert({ Htp.InternetFormatWithThrow() }
			throws: 'Expected date string but result was')
		}
	}