// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
// see: http://www.vervest.org/htp/
class
	{
	UnixTime()
		{
		t = Date(.InternetFormat().Replace(' GMT', ''))
		days = t.NoTime().MinusDays(#19700101)
		seconds = t.MinusSeconds(t.NoTime()).Round(0)
		return days * 24.HoursInSeconds() + seconds /*= 24 hours a day*/
		}
	Header()
		{ // max-age is in seconds
		return Http('HEAD', 'http://www.google.com',
			header: Object(Pragma: 'no-cache',
				'Cache-Control': 'no-cache,no-store,max-age=0',
				Expires: Date().InternetFormat())).header
		}
	InternetFormat()
		{
		hdr = .Header()
		result = hdr.Extract('^Date: (.* GMT)$')
		if result is false
			SuneidoLog('INFO: Htp: header missing Date', params: hdr.Lines())
		return result
		}
	InternetFormatWithThrow()
		{
		if not String?(val = .InternetFormat())
			throw 'Expected date string but result was ' $ Display(val)
		return val
		}
	}