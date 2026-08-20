// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		_time = #20260720.105122397
		_gmtTime = _time
		_timeZoneRec = TimeZones.Copy()
		_defaultTime = Date()
		cl = LocalTime
			{
			LocalTime_gmtTime()
				{
				return _gmtTime
				}
			LocalTime_timeZoneRec(stateProv)
				{
				return _timeZoneRec[stateProv]
				}
			LocalTime_defaultTime()
				{
				return _defaultTime
				}
			}
		func = cl

		// no/wrong stateProv passed in
		stateProv = ''
		timezone = 'PST'
		Assert(func(stateProv, timezone) is: _defaultTime)
		stateProv = 'AA'
		Assert(func(stateProv, timezone) is: _defaultTime)

		// BC normal case with DST (but no switch)
		stateProv = 'BC'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))
		// AB with DST (but no switch)
		stateProv = 'AB'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))
		// AZ with DST (but no switch)
		stateProv = 'AZ'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		// HI with DST (but no switch)
		stateProv = 'HI'
		timezone = 'HST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -10))
		// SK with DST (but no switch)
		stateProv = 'SK'
		timezone = 'CST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))
		// SON with DST (but no switch)
		stateProv = 'SON'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		// YT with DST (but no switch)
		stateProv = 'SON'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))

		// AK with DST
		stateProv = 'AK'
		timezone = 'AKST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -8))
		// OR with DST
		stateProv = 'OR'
		timezone = 'PST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		timezone = 'HST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -9))
		timezone = 'AKST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -8))
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))
		timezone = 'CST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -5))
		timezone = 'EST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -4))
		timezone = 'AST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -3))
		timezone = 'NST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(minutes: -2.5 * 60))

		// NT with DST
		stateProv = 'NT'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))
		// ID with DST
		stateProv = 'ID'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))

		// MB with DST
		stateProv = 'MB'
		timezone = 'CST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -5))
		// MN with DST
		stateProv = 'MN'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -5))

		// SD with DST
		stateProv = 'SD'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))
		// SD with CST
		timezone = 'CST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -5))

		// SD with DST
		stateProv = 'ON'
		timezone = 'EST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -4))

		// NB with AST
		stateProv = 'NB'
		timezone = 'AST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -3))

		// NL with NST
		stateProv = 'NL'
		timezone = 'NST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(minutes: -2.5 * 60))

		stateProv = 'BC'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))

		// no DST after new gmtTime
		_gmtTime = _gmtTime.Plus(months: 4)

		// no/wrong stateProv passed in
		stateProv = ''
		timezone = 'PST'
		Assert(func(stateProv, timezone) is: _defaultTime)
		stateProv = 'AA'
		Assert(func(stateProv, timezone) is: _defaultTime)

		// BC normal case without DST (but no switch)
		stateProv = 'BC'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		// AB without DST (but no switch)
		stateProv = 'AB'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))
		// AZ without DST (but no switch)
		stateProv = 'AZ'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		// HI without DST (but no switch)
		stateProv = 'HI'
		timezone = 'HST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -10))
		// SK without DST (but no switch)
		stateProv = 'SK'
		timezone = 'CST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))
		// SON without DST (but no switch)
		stateProv = 'SON'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		// YT without DST (but no switch)
		stateProv = 'SON'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))

		// AK without DST
		stateProv = 'AK'
		timezone = 'AKST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -9))
		// OR without DST
		stateProv = 'OR'
		timezone = 'PST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -8))
		timezone = 'HST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -10))
		timezone = 'AKST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -9))
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		timezone = 'CST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))
		timezone = 'EST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -5))
		timezone = 'AST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -4))
		timezone = 'NST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(minutes: -3.5 * 60))

		// NT without DST
		stateProv = 'NT'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		// ID without DST
		stateProv = 'ID'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))

		// MB without DST
		stateProv = 'MB'
		timezone = 'CST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))
		// MN without DST
		stateProv = 'MN'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))

		// SD without DST
		stateProv = 'SD'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		// SD without CST
		timezone = 'CST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))

		// SD without DST
		stateProv = 'ON'
		timezone = 'EST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -5))

		// NB without AST
		stateProv = 'NB'
		timezone = 'AST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -4))

		// NL without NST
		stateProv = 'NL'
		timezone = 'NST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(minutes: -3.5 * 60))

		stateProv = 'BC'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))

		// more exception cases; no DST
		stateProv = 'BC'
		timezone = 'MST'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		_timeZoneRec.BC = [zone: "PST", nodaylightAdjustment?:, offSetOverride: -7,
			exceptions: [MST: []]]
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		_timeZoneRec.BC = [zone: "PST", nodaylightAdjustment?:, offSetOverride: -7,
			exceptions: [MST: [nodaylightAdjustment?:]]]
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		_timeZoneRec.BC = [zone: "PST", nodaylightAdjustment?:, offSetOverride: -7,
			exceptions: [MST: [nodaylightAdjustment?:, offSetOverride: -6]]]
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))

		// back to DST
		_gmtTime = _time
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))
		_timeZoneRec.BC = [zone: "PST", nodaylightAdjustment?:, offSetOverride: -7,
			exceptions: [MST: []]]
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))
		_timeZoneRec.BC = [zone: "PST", nodaylightAdjustment?:, offSetOverride: -7,
			exceptions: [MST: [nodaylightAdjustment?:, offSetOverride: -7]]]
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		_timeZoneRec.BC = [zone: "PST", nodaylightAdjustment?:, offSetOverride: -7,
			exceptions: [MST: [offSetOverride: -7]]]
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -6))

		// test invalid timezone input
		timezone = 'invalid'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -7))
		stateProv = 'MB'
		Assert(func(stateProv, timezone) is: _gmtTime.Plus(hours: -5))
		}

	Test_DaylightSavings?()
		{
		func = LocalTime.DaylightSavings?
		timezone = 'CST'

		// the following is from DaylightSavings_Test with return value reverted
		Assert(func(#20190101, timezone) is: false)
		Assert(func(#20190301, timezone) is: false)
		Assert(func(#20190309, timezone) is: false)
		Assert(func(#20190310, timezone) is: false)
		Assert(func(#20190310.0800, timezone))
		Assert(func(#20190310.0801, timezone))
		Assert(func(#20190310.1001, timezone))
		Assert(func(#20190615, timezone))
		Assert(func(#20190701, timezone))
		Assert(func(#20191102, timezone))
		Assert(func(#20191103, timezone))
		Assert(func(#20191103.065900, timezone))
		Assert(func(#20191103.070000, timezone) is: false)
		Assert(func(#20191103.100000, timezone) is: false)
		Assert(func(#20191201, timezone) is: false)
		Assert(func(#20191231, timezone) is: false)

		Assert(func(#20200101, timezone) is: false)
		Assert(func(#20200301, timezone) is: false)
		Assert(func(#20200307, timezone) is: false)

		Assert(func(#20200308, timezone) is: false)
		Assert(func(#20200308.075900, timezone) is: false)
		Assert(func(#20200308.080000, timezone))
		Assert(func(#20200615, timezone))
		Assert(func(#20200701, timezone))
		Assert(func(#20201031, timezone))
		Assert(func(#20201101, timezone))
		Assert(func(#20201101.065900, timezone))
		Assert(func(#20201101.070000, timezone) is: false)
		Assert(func(#20201201, timezone) is: false)
		Assert(func(#20201231, timezone) is: false)
		// the above is from DaylightSavings_Test with return value reverted

		// more edge cases for March spring forward
		gmtTime = #20260308.010000
		Assert(func(gmtTime, timezone) is: false)
		gmtTime = #20260308.075900
		Assert(func(gmtTime, timezone) is: false)
		gmtTime = #20260308.075959
		Assert(func(gmtTime, timezone) is: false)
		gmtTime = #20260308.080000
		Assert(func(gmtTime, timezone))
		gmtTime = #20260308.080001
		Assert(func(gmtTime, timezone))
		gmtTime = #20260308.100000
		Assert(func(gmtTime, timezone))
		gmtTime = #20260309.100000
		Assert(func(gmtTime, timezone))

		// more edge cases for November roll back
		gmtTime = #20261101.010000
		Assert(func(gmtTime, timezone))
		gmtTime = #20261101.065900
		Assert(func(gmtTime, timezone))
		gmtTime = #20261101.065959
		Assert(func(gmtTime, timezone))
		gmtTime = #20261101.070000
		Assert(func(gmtTime, timezone) is: false)
		gmtTime = #20261101.070001
		Assert(func(gmtTime, timezone) is: false)
		gmtTime = #20261101.070100
		Assert(func(gmtTime, timezone) is: false)
		gmtTime = #20261101.090000
		Assert(func(gmtTime, timezone) is: false)
		}
	}
