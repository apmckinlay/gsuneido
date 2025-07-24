// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_format()
		{
		fn = CalendarWebPage.CalendarWebPage_format

		// no events
		expected = `({"test", "events":[]})`
		Assert(fn('"test", ', #()) is: expected)

		// single day event
		events = #(
			#(start_date: #20241111, title: "Test Event", type: "Holidays",
				date: #20241111, subtype: ""))
		expected = `({"test", "events":[{"type":"Holidays","subtype":"",` $
			`"start_date":"2024-11-11","multidays":1,"span":1,"title":"Test Event"}]})`
		Assert(fn('"test", ', events) is: expected)

		// multiday event
		events = #(#(start_date: #20241112, title: "Test Multiday", multidays: 3,
			type: "Events", end:, desc: "Shop Work", date: #20241112,
			subtype: "inshop", completed: 0, span: 3))
		expected = `({"test", "events":[{"type":"Events","subtype":"inshop",` $
			`"start_date":"2024-11-12","multidays":3,"span":3,` $
			`"end_date":"2024-11-14","completed":0,"end":true,` $
			`"title":"Test Multiday","desc":"Shop Work"}]})`
		Assert(fn('"test", ', events) is: expected)
		}

	Test_handleInvalidChar()
		{
		f = CalendarWebPage.CalendarWebPage_handleInvalidChar
		Assert(f("abc") is: "abc")
		Assert(f("ab\\c") is: "ab\\\\c")
		Assert(f('ab"cd"ef') is: 'ab\\"cd\\"ef')
		Assert(f('abc\r\n') is: 'abc ')
		Assert(f('abc\r\n\n') is: 'abc  ')
		Assert(f('abc\r\nefg') is: 'abc efg')
		Assert(f('ab\\c\r\n') is: 'ab\\\\c ')
		Assert(f('ab"c\r\n') is: 'ab\\"c ')
		}

	Test_date_span()
		{
		f = CalendarWebPage.CalendarWebPage_span
		Assert(f(#20090601, 1) is: 1)
		Assert(f(#20090601, 3) is: 3)
		Assert(f(#20090601, 5) is: 5)
		Assert(f(#20090601, 6) is: 6)
		Assert(f(#20090601, 7) is: 6)
		Assert(f(#20090601, 8) is: 6)
		Assert(f(#20090601, 100) is: 6)

		Assert(f(#20090531, 1) is: 1)
		Assert(f(#20090531, 3) is: 3)
		Assert(f(#20090531, 7) is: 7)
		Assert(f(#20090531, 100) is: 7)

		Assert(f(#20090606, 1) is: 1)
		Assert(f(#20090606, 2) is: 1)
		Assert(f(#20090606, 3) is: 1)
		Assert(f(#20090606, 20) is: 1)

		Assert(f(#20091228, 4) is: 4)
		}

	from: #20090531
	to: #20090628
	Test_add_events()
		{
		events = Object()
		e = Object(type:'test', date:#20090601, title:'this is a test')
		CalendarWebPage.CalendarWebPage_add_events(e, .from, .to, events)
		Assert(events is: Object(Object(type: 'test', subtype: '',
			title: 'this is a test',  date: #20090601, start_date:#20090601)))

		events = Object()
		e = Object(type:'test', date:#20090601, title:'this is a test',
			multidays: 2)
		CalendarWebPage.CalendarWebPage_add_events(e, .from, .to, events)
		Assert(events is: Object(Object(type: 'test', subtype: '',
			title: 'this is a test',  date: #20090601,	start_date:#20090601,
			multidays: 2, span: 2,	completed: 0, end: true)))

		events = Object()
		e = Object(type:'test', date:#20090601, title:'this is a test',
			multidays: 6)
		CalendarWebPage.CalendarWebPage_add_events(e, .from, .to, events)
		Assert(events is: Object(Object(type: 'test', subtype: '',
			date: #20090601, start_date:#20090601, title:'this is a test',
			multidays: 6, span: 6, completed: 0, end: true)))

		events = Object()
		e = Object(type:'test', date:#20090529, title:'this is a test',
			multidays: 6)
		CalendarWebPage.CalendarWebPage_add_events(e, .from, .to, events)
		Assert(events is: Object(Object(type: 'test', subtype: '', date: #20090531,
			start_date:#20090529, title:'this is a test', multidays: 6, span: 4,
			completed: 2, end: true)))
		}

	Test_add_events_MoreThanAWeek()
		{
		events = Object()
		e = Object(type:'test', date:#20090601, title:'this is a test',
			multidays: 7)
		CalendarWebPage.CalendarWebPage_add_events(e, .from, .to, events)
		Assert(events is: Object(
			Object(type: 'test', subtype: '', date: #20090601, start_date:#20090601,
				title:'this is a test', multidays: 7, span: 6, completed: 0, end: false),
			Object(type: 'test', subtype: '', date: #20090607, start_date:#20090601,
				title:'this is a test', multidays: 7, span: 1, completed: 6, end: true)))

		events = Object()
		e = Object(type:'test', date:#20090601, title:'this is a test',
			multidays: 18)
		CalendarWebPage.CalendarWebPage_add_events(e, .from, .to, events)
		Assert(events is: Object(
			Object(type: 'test', subtype: '', date: #20090601, start_date:#20090601,
				title:'this is a test', multidays: 18, span: 6, completed: 0, end: false),
			Object(type: 'test', subtype: '', date: #20090607, start_date:#20090601,
				title:'this is a test', multidays: 18, span: 7, completed: 6, end: false),
			Object(type: 'test', subtype: '', date: #20090614, start_date:#20090601,
				title:'this is a test', multidays: 18, span: 5, completed: 13,
				end: true)))

		events = Object()
		e = Object(type:'test', date:#20090529, title:'this is a test',
			multidays: 18)
		CalendarWebPage.CalendarWebPage_add_events(e, .from, .to, events)
		Assert(events is: Object(
			Object(type: 'test', subtype: '', date: #20090531, start_date:#20090529,
				title:'this is a test', multidays: 18, span: 7, completed: 2, end: false),
			Object(type: 'test', subtype: '', date: #20090607, start_date:#20090529,
				title:'this is a test', multidays: 18, span: 7, completed: 9, end: false),
			Object(type: 'test', subtype: '', date: #20090614, start_date:#20090529,
				title:'this is a test', multidays: 18, span: 2, completed: 16,
				end: true)))

		events = Object()
		e = Object(type:'test', date:#20090615, title:'this is a test',
			multidays: 30)
		CalendarWebPage.CalendarWebPage_add_events(e, .from, .to, events)
		Assert(events is: Object(
			Object(type: 'test', subtype: '', date: #20090615, start_date:#20090615,
				title:'this is a test', multidays: 30, span: 6, completed: 0, end: false),
			Object(type: 'test', subtype: '', date: #20090621, start_date:#20090615,
				title:'this is a test', multidays: 30, span: 7, completed: 6, end: false)
			))
		}

	Test_add_events_LargeDateRange()
		{
		events = Object()
		e = Object(type:'test', date:#20090520, title:'this is a test',
			multidays: 100)
		CalendarWebPage.CalendarWebPage_add_events(e, .from, .to, events)
		Assert(events is: Object(
			Object(type: 'test', subtype: '', date: #20090531, start_date:#20090520,
				title:'this is a test', multidays: 100, span: 7, completed: 11,
				end: false),
			Object(type: 'test', subtype: '', date: #20090607, start_date:#20090520,
				title:'this is a test', multidays: 100, span: 7, completed: 18,
				end: false),
			Object(type: 'test', subtype: '', date: #20090614, start_date:#20090520,
				title:'this is a test', multidays: 100, span: 7, completed: 25,
				end: false),
			Object(type: 'test', subtype: '', date: #20090621, start_date:#20090520,
				title:'this is a test', multidays: 100, span: 7, completed: 32,
				end: false)
			))
		}

	Test_add_events_NoEndDate()
		{
		events = Object()
		e = Object(type:'test', date:#20090520, title:'this is a test',
			multidays: Date.End().MinusDays(#20090520))
		CalendarWebPage.CalendarWebPage_add_events(e, .from, .to, events)
		Assert(events is: Object(
			Object(type: 'test', subtype: '', date: #20090531, start_date:#20090520,
				title:'this is a test', multidays: 361816, span: 7, completed: 11,
				end: false),
			Object(type: 'test', subtype: '', date: #20090607, start_date:#20090520,
				title:'this is a test', multidays: 361816, span: 7, completed: 18,
				end: false),
			Object(type: 'test', subtype: '', date: #20090614, start_date:#20090520,
				title:'this is a test', multidays: 361816, span: 7, completed: 25,
				end: false),
			Object(type: 'test', subtype: '', date: #20090621, start_date:#20090520,
				title:'this is a test', multidays: 361816, span: 7, completed: 32,
				end: false)
			))

		events = Object()
		e = Object(type:'test', date: #20090520.1230, title:'this is a test',
			multidays: Date.End().MinusDays(#20090520))
		CalendarWebPage.CalendarWebPage_add_events(e, .from, .to, events)
		Assert(events is: Object(
			Object(type: 'test', subtype: '', date: #20090531, start_date:#20090520,
				title:'this is a test', multidays: 361816, span: 7, completed: 11,
				end: false),
			Object(type: 'test', subtype: '', date: #20090607, start_date:#20090520,
				title:'this is a test', multidays: 361816, span: 7, completed: 18,
				end: false),
			Object(type: 'test', subtype: '', date: #20090614, start_date:#20090520,
				title:'this is a test', multidays: 361816, span: 7, completed: 25,
				end: false),
			Object(type: 'test', subtype: '', date: #20090621, start_date:#20090520,
				title:'this is a test', multidays: 361816, span: 7, completed: 32,
				end: false)
			))
		}

	Test_parseEventSources()
		{
		f = CalendarWebPage.CalendarWebPage_parse_src
		str = 'Holiday$*$Staff$*$Staff__3 month anniversary$*$Staff__Birthday$*$' $
			'Staff__year anniversary'
		ob = #(Holiday: #(), Staff: #("3 month anniversary", "Birthday",
			"year anniversary"))
		Assert(f(str) is: ob)

		str = 'Holiday$*$Staff__3 month anniversary$*$Staff__Birthday$*$' $
			'Staff__year anniversary'
		ob = #(Holiday: #(), Staff: #("3 month anniversary", "Birthday",
			"year anniversary"))
		Assert(f(str) is: ob)

		str = 'Staff__3 month anniversary$*$Staff__Birthday$*$Staff__year anniversary'
		ob = #(Staff: #("3 month anniversary", "Birthday", "year anniversary"))
		Assert(f(str) is: ob)

		str = 'Holiday$*$Staff$*$Staff__3 - month anniversary$*$Staff__1 , Birthday'
		ob = #(Holiday: #(), Staff: #("3 - month anniversary", "1 , Birthday"))
		Assert(f(str) is: ob)
		}
	}
