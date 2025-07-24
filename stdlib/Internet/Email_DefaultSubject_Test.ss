// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_BuildSubject()
		{
		fn = Email_DefaultSubject.BuildSubject

		Assert(fn("", #(), []) is: '')

		Assert(fn("", #(city, date), []) is: '')

		Assert(fn("hello world", #(city, date), []) is: 'hello world')

		Assert(fn("hello <City> at <Date>", #(city, date), []) is: 'hello  at')

		rec = [city: "Saskatoon", date: #20010504]
		dateStr = #20010504.ShortDate()
		Assert(fn("hello <City> at <Date>", #(city, date), rec)
			is: 'hello Saskatoon at ' $ dateStr)

		Assert(fn("hello <City> \r\nat <Date>", #(city, date), rec)
			is: 'hello Saskatoon  at ' $ dateStr)

		Assert(fn("hello <City <Date>> at <Date> <world>   ", #(city, date), rec)
			is: 'hello <City ' $ dateStr $ '> at ' $ dateStr $ ' <world>')
		}

	Test_AddEmailSubjectInfo()
		{
		fn = Email_DefaultSubject.AddEmailSubjectInfo

		rpt = Object(Params: Object())
		fn(rpt, 'type', #(field: v1), #('field'))
		Assert(rpt is: #(Params: (
			EmailSubject: (data: (field: v1), type: 'type', cols: ('field')))))

		fn(rpt, 'type', #(field: v1), #('field'))
		Assert(rpt is: #(Params: (
			EmailSubject: (data: (field: v1), type: 'type', cols: ('field')))))

		fn(rpt, 'type', #(field: v2), #('field'))
		Assert(rpt is: #(Params: (NoEmailSubject:)))

		fn(rpt, 'type', #(field: v3), #('field'))
		Assert(rpt is: #(Params: (NoEmailSubject:)))
		}
	}