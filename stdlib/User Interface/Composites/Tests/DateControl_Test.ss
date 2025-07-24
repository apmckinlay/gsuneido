// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		datetime = Date()
		Assert(DateControl.FormatValue(datetime, showTime:) is: datetime.ShortDateTime())
		Assert(DateControl.FormatValue(#20081001.135055250, showTime:)
			is: #20081001.135055250.ShortDateTime())
		}

	Test_validData()
		{
		Assert(DateControl.ValidData?(""))
		Assert(DateControl.ValidData?("fred") is: false)
		Assert(DateControl.ValidData?(Date().NoTime()))
		Assert(DateControl.ValidData?("", mandatory:) is: false)
		}

	Test_validDateCode?()
		{
		valid = DateControl.DateControl_validDateCode?
		Assert(valid('') is: false)
		Assert(valid('t'))
		Assert(valid('t+'))
		Assert(valid('m'))
		Assert(valid('m-'))
		Assert(valid('Y'))
		Assert(valid('m--'))
		Assert(valid('m++'))
		Assert(valid('ph'))
		Assert(valid('pm'))
		Assert(valid('pm++++'))
		Assert(valid('ph---'))

		Assert(valid('ym++') is: false)
		Assert(valid('m++m') is: false)
		Assert(valid('x') is: false)

		Assert(valid('t', showTime:))
		Assert(valid('t9', showTime:))
		Assert(valid('t20', showTime:))
		Assert(valid('t150', showTime:))
		Assert(valid('t-0150', showTime:))
		Assert(valid('t-150', showTime:))
		Assert(valid('t-1850', showTime:))
		Assert(valid('M-1850', showTime:))
		Assert(valid('y--1850', showTime:))
		Assert(valid('y+1850', showTime:))

		Assert(valid('t-170', showTime:) is: false)
		Assert(valid('t-2570', showTime:) is: false)
		Assert(valid('x-120', showTime:) is: false)
		Assert(valid('m-99', showTime:) is: false)

		text = 't' $ '+'.Repeat(99)
		Assert(valid(text, showTime:))
		Assert(valid(text $ '0100', showTime:))
		Assert(valid(text $ '01', showTime:))
		Assert(valid(text $ '00100', showTime:) is: false)
		Assert(valid(text $ '0100') is: false)
		Assert(valid(text $ 'a' $ '+'.Repeat(99)) is: false)
		Assert(valid(text $ 'a1111', showTime:) is: false)
		Assert(valid(text $ 'a111', showTime:) is: false)
		Assert(valid(text $ '111a', showTime:) is: false)
		Assert(valid(text $ '$\/*', showTime:) is: false)
		Assert(valid('today') is: false)
		}

	Test_ConvertToDate()
		{
		c = DateControl.ConvertToDate
		date = Date().NoTime()
		Assert(c('t', convertDateCodes?:) is: date)
		Assert(c('t++', convertDateCodes?:) is: date.Plus(days: 2))
		Assert(c('T++', convertDateCodes?:) is: date.Plus(days: 2))
		Assert(c('t0000', convertDateCodes?:, showTime:) is: date.Plus(milliseconds: 1))

		before = Date() // IF there is no time displayed, 'when' it was entered is used
		Assert(c('t', convertDateCodes?:, showTime:) between: Object(before, Date()))

		Assert(c('t2359', convertDateCodes?:, showTime:)
			is: date.Plus(hours: 23, minutes: 59))
		Assert(c('t2750', convertDateCodes?:, showTime:) is: false)
		Assert(c('r++', convertDateCodes?: false) is: 'r++')

		year = Date().Year()
		Assert(c('r', convertDateCodes?:) is: Date(year $ '1231').NoTime())
		Assert(c('y', convertDateCodes?:) is: Date(year $ '0101').NoTime())

		Assert(c('05062001', convertDateCodes?:, format: 'MMddyyyy') is: #20010506)
		Assert(c('05062001dddd', convertDateCodes?:, format: 'MMddyyyy') is: false)

		date = Date().NoTime()
		Assert(c('ph', convertDateCodes?:) is: date.Replace(day: 1).Plus(days: -1))
		Assert(c('pm', convertDateCodes?:) is: date.Replace(day: 1).Plus(months: -1))
		}

	Test_largeConvertStr()
		{
		c =  DateControl.ConvertToDate
		date = Date().Plus(days: 99).NoTime()
		text = 't' $ '+'.Repeat(99)
		Assert(c(text, convertDateCodes?:) is: date)
		Assert(c(text $ '1100', convertDateCodes?:, showTime:) is: date.Plus(hours: 11))

		date = Date().Minus(days: 99).NoTime()
		text = 't'.RightFill(100, '-')
		Assert(c(text, convertDateCodes?:) is: date)
		Assert(c(text $ '1100', convertDateCodes?:, showTime:) is: date.Plus(hours: 11))
		}

	Test_ConvertToDateNoConvert()
		{
		c =  DateControl.ConvertToDate

		Assert(c('t',   convertDateCodes?: false) is: 't')
		Assert(c('t++', convertDateCodes?: false) is: 't++')
		Assert(c('T++', convertDateCodes?: false) is: 'T++')
		Assert(c('r++', convertDateCodes?: false) is: 'r++')
		Assert(c('Y--', convertDateCodes?: false) is: 'Y--')
		Assert(c('h+-', convertDateCodes?: false) is: 'h+-')
		Assert(c('m-++-', convertDateCodes?: false) is: 'm-++-')

		Assert(c('i++', convertDateCodes?: false) is: false)
		}
	}