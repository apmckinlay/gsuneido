// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
// NOTE: this is deliberately named so it will NOT be found by the test runner
// because it is slow and accesses the internet
Test
	{
	http: "http://httpbin.io"
	https: "https://httpbin.io"
	Test_Get()
		{
		for url in [.http, .https]
			{
			resp = HttpClient2("GET", url $ "/dump/request")
			Assert(resp.content like:
				"GET /dump/request HTTP/1.1
				Host: httpbin.io
				Accept-Encoding: gzip
				User-Agent: Suneido
				")
			}
		}
	Test_Timeout()
		{
		t = Date()
		Assert({ HttpClient2("GET", .https $ "/delay/10", timeout: 1) }
			throws: "deadline")
		Assert(Date().MinusSeconds(t).Round(0) is: 1)
		}
	Test_Reader()
		{
		s = "now is the time for all good men to come to the aid of their party"
		rdr = {|n|
			n = Min(n, 10) /*= small chunks */
			result = s[..n]
			s = s[n..]
			result
			}
		resp = HttpClient2("POST", .http $ "/dump/request", rdr,
			[content_length: s.Size()])
		Assert(resp.content has: s)
		}
	}