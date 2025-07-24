// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// this is tests for built-in methods
// Dates_Test is for methods defined in Dates
// SuJsWebTest
Test
	{
	Test_false()
		{
		Assert(Date(false) is: false)
		}
	Test_Minus()
		{
		d1 = #19000101
		d2 = #20000303
		n = d2.MinusDays(d1)
		Assert(d1.Plus(days: n) is: d2)

		d1 = Date('9:45')
		d2 = Date('10:30')
		n = d2.MinusSeconds(d1)
		Assert(n is: 45 * 60)
		Assert(d1.Plus(seconds: n) is: d2)

		d1 = Date('9:45')
		d2 = Date('10:30')
		n = d2.MinusSeconds(d1) / 60
		Assert(n is: 45)
		Assert(d1.Plus(minutes: n) is: d2)
		}
	Test_Parse()
		{
		Assert(Date("Feb 14 2000") is: #20000214)
		Assert(Date("14 Feb 2000") is: #20000214)
		Assert(Date("Feb 14 00") is: #20000214)
		Assert(Date("Aug 14 00") is: #20000814)
		Assert(Date("Feb 14 01") is: #20010214)
		Assert(Date("July 14 00") is: #20000714)
		Assert(Date("Feb 14 01") is: #20010214)
		Assert(Date("0/2/14") is: #20000214)
		Assert(Date("60/4/25") is: #19600425)
		Assert(Date("4/25/60") is: #19600425)
		Assert(Date("7-8-9 3:4:5am", "yyMMdd") is: #20070809.030405)
		Assert(Date("3:04pm 7-8-9", "yyMMdd") is: #20070809.1504)
		Assert(Date("14 Feb 20", "yyMMdd") is: #20140220)
		Assert(Date("14 Feb 99") is: #19990214)
		Assert(Date("9:45").Format("Hmm") is: "945")
		Assert(Date("9:45pm").Format("Hmm") is: "2145")
		Assert(Date("15:45").Format("Hmm") is: "1545")
		Assert(Date("3am").Format("Hmm") is: "300")
		Assert(Date("3pm").Format("Hmm") is: "1500")
		Assert(Date("000214") is: #20000214)
		Assert(Date("20000214", "yyMMdd") is: #20000214)
		Assert(Date('Wed Jul 21 00:45:55 2003\n') is: #20030721.004555)
		Assert(Date('201406011957') is: false)
		Assert(Date('#201406011957') is: false)

		Assert(Date('#11111111.1') is: false)
		Assert(Date('#20179999') is: false)
		Assert(Date('#201799-1') is: false)
		Assert(Date('a '.Repeat(33)) is: false)
		}
	Test_Normalize()
		{
		// month -> year
		Assert(Date(year: 2000, month: 0).Format("yyyyMM") is: "199912")
		Assert(Date(year: 2000, month: 13).Format("yyyyMM") is: "200101")
		// day -> month
		Assert(Date(year: 2000, month: 2, day: 0).Format("MMdd") is: "0131")
		Assert(Date(year: 2000, month: 2, day: 30).Format("MMdd") is: "0301")
		// hour -> day
		Assert(Date(day: 15, hour: -1).Format("ddHH") is: "1423")
		Assert(Date(day: 15, hour: 24).Format("ddHH") is: "1600")
		// minute -> hour
		Assert(Date(hour: 6, minute: -1).Format("HHmm") is: "0559")
		Assert(Date(hour: 6, minute: 60).Format("HHmm") is: "0700")
		// second -> minute
		Assert(Date(minute: 6, second: -1).Format("mmss") is: "0559")
		Assert(Date(minute: 6, second: 60).Format("mmss") is: "0700")
		}
	Test_FormatEn()
		{
		d = Date(year: 1999, month: 1, day: 2, hour: 3, minute: 4, second: 5)
		Assert(d.FormatEn("yyyy/MM/dd") is: "1999/01/02")
		Assert(d.FormatEn("yy/M/d") is: "99/1/2")
		Assert(d.FormatEn("hh:mm:ss AA") is: "03:04:05 AM")
		Assert(d.FormatEn("ddd d MMM yyyy H:mm:ss") is: "Sat 2 Jan 1999 3:04:05")
		d = Date(year: 2000, month: 11, day: 22, hour: 23, minute: 44, second: 55)
		Assert(d.FormatEn("yyyy/MM/dd") is: "2000/11/22")
		Assert(d.FormatEn("yy/M/d") is: "00/11/22")
		Assert(d.FormatEn("h:m:s AA") is: "11:44:55 PM")
		Assert(d.FormatEn("ddd d MMM yyyy H:mm:ss") is: "Wed 22 Nov 2000 23:44:55")
		}
	Test_Plus()
		{
		d1 = Date(year: 1999, month: 1, day: 1,
			hour: 0, minute: 0, second: 0, millisecond: 0)
		d2 = Date(year: 2000, month: 2, day: 2,
			hour: 1, minute: 1, second: 1, millisecond: 0)
		// no normalization
		d3 = d1.Plus(years: 1, months: 1, days: 1, hours: 1, minutes: 1, seconds: 1)
		Assert(d2 is: d3)
		Assert(d1.Plus(years: 1).Plus(years: -1) is: d1)
		Assert(d1.Plus(months: 1).Plus(months: -1) is: d1)
		Assert(d1.Plus(days: 1).Plus(days: -1) is: d1)
		Assert(d1.Plus(hours: 1).Plus(hours: -1) is: d1)
		Assert(d1.Plus(minutes: 1).Plus(minutes: -1) is: d1)
		Assert(d1.Plus(seconds: 1).Plus(seconds: -1) is: d1)
		// need normalization
		Assert(d1.Plus(months: -1).Plus(months: 1) is: d1)
		Assert(d1.Plus(months: 20).Plus(months: -20) is: d1)
		Assert(d1.Plus(months: -20).Plus(months: 20) is: d1)
		Assert(d1.Plus(days: -1).Plus(days: 1) is: d1)
		Assert(d1.Plus(days: 99).Plus(days: -99) is: d1)
		Assert(d1.Plus(days: -99).Plus(days: 99) is: d1)
		Assert(d1.Plus(hours: -1).Plus(hours: 1) is: d1)
		Assert(d1.Plus(hours: 999).Plus(hours: -999) is: d1)
		Assert(d1.Plus(hours: -999).Plus(hours: 999) is: d1)
		Assert(d1.Plus(minutes: -1).Plus(minutes: 1) is: d1)
		Assert(d1.Plus(minutes: 999).Plus(minutes: -999) is: d1)
		Assert(d1.Plus(minutes: -999).Plus(minutes: 999) is: d1)
		Assert(d1.Plus(seconds: -1).Plus(seconds: 1) is: d1)
		Assert(d1.Plus(seconds: 999).Plus(seconds: -999) is: d1)
		Assert(d1.Plus(seconds: -999).Plus(seconds: 999) is: d1)

		Assert(Display(d1.Plus(seconds: -1)) is: "#19981231.235959")
		}
	Test_Time()
		{
		Assert(Date('0:00 am').Format("HH:mm") is: "00:00")
		Assert(Date('12:00 am').Format("HH:mm") is: "00:00")
		Assert(Date('12:05 am').Format("HH:mm") is: "00:05")
		Assert(Date('12:10 am').Format("HH:mm") is: "00:10")
		Assert(Date('12:20 am').Format("HH:mm") is: "00:20")
		Assert(Date('12:30 am').Format("HH:mm") is: "00:30")
		Assert(Date('13:00 am') is: false)
		Assert(Date('1:00 am').Format("HH:mm") is: "01:00")
		Assert(Date('11:45 am').Format("HH:mm") is: "11:45")
		Assert(Date('12:00 pm').Format("HH:mm") is: "12:00")
		Assert(Date('12:30 pm').Format("HH:mm") is: "12:30")
		Assert(Date('1:00 pm').Format("HH:mm") is: "13:00")
		Assert(Date('13:00 pm').Format("HH:mm") is: "13:00")
		Assert(Date('11:45 pm').Format("HH:mm") is: "23:45")

		Assert(#20020118.1220.Format('hh:mm AA') is: "12:20 PM")
		Assert(Date("2002/01/18 12:05 pm") is: #20020118.1205)
		}
	Test_WeekDay()
		{
		Assert(#20030331.WeekDay() is: 1)
		Assert(#20030331.WeekDay(0) is: 1)
		Assert(#20030331.WeekDay(1) is: 0)
		Assert(#20030331.WeekDay('mon') is: 0)
		Assert(#20030331.WeekDay('Monday') is: 0)
		}
	Test_month_day()
		{
		date = Date().NoTime()
		month = date.Month()
		day = date.Day()
		Assert(Date(month $ '/' $ day, "yMd") is: date)
		Assert(Date(day $ '/' $ month, "dMy") is: date)
		}
	Test_invalid()
		{
		Assert({ Date.End().Plus(years: 1) } throws: 'bad date')
		}

	Test_parse_default_year()
		{
		d = Date()
		y = d.Year()
		m = d.Month()
		data = .parse_data[m]
		for (m = 1; m <= 12; ++m)
			Assert(Date(.parse_data[m][0] $ ' 1').Year() is: y + data[m],
				msg: "for " $ .parse_data[m][0] $ ' 1')
		}
	parse_data: (
		1:  (Jan, 0,  0,  0,  0,  0,  0,  0, -1, -1, -1, -1, -1)
		2:  (Feb, 0,  0,  0,  0,  0,  0,  0,  0, -1, -1, -1, -1)
		3:  (Mar, 0,  0,  0,  0,  0,  0,  0,  0,  0, -1, -1, -1)
		4:  (Apr, 0,  0,  0,  0,  0,  0,  0,  0,  0,  0, -1, -1)
		5:  (May, 0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0, -1)
		6:  (Jun, 0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0)
		7:  (Jul, 1,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0)
		8:  (Aug, 1,  1,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0)
		9:  (Sep, 1,  1,  1,  0,  0,  0,  0,  0,  0,  0,  0,  0)
		10: (Oct, 1,  1,  1,  1,  0,  0,  0,  0,  0,  0,  0,  0)
		11: (Nov, 1,  1,  1,  1,  1,  0,  0,  0,  0,  0,  0,  0)
		12: (Dec, 1,  1,  1,  1,  1,  1,  0,  0,  0,  0,  0,  0)
		)
	}