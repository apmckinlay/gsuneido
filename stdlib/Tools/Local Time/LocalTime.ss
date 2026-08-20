// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
class
	{
	zoneOffsets: #(
		HST : -10
		AKST: -9
		PST : -8
		MST : -7
		CST : -6
		EST : -5
		AST : -4
		NST : -3.5)
	hourMinutesConversion: 60
	CallClass(stateProv, timezone)
		{
		if .cantConvert?(stateProv)
			return .defaultTime()

		gmtTime = .gmtTime()
		info = .localTimeInfo(stateProv, timezone)
		curTime = gmtTime.Plus(minutes: (info.offSetOverride) * .hourMinutesConversion)
		if not info.nodaylightAdjustment? and .DaylightSavings?(gmtTime, timezone)
			curTime = curTime.Plus(hours: 1)
		return curTime
		}

	gmtTime()
		{
		return Date().GMTime()
		}

	defaultTime()
		{
		return Date()
		}

	cantConvert?(stateProv)
		{
		return stateProv is "" or not TimeZones.Member?(stateProv)
		}

	localTimeInfo(stateProv, timezone)
		{
		offSetOverride = 0
		nodaylightAdjustment? = false
		zoneRec = .timeZoneRec(stateProv)
		if zoneRec.zone is timezone
			{
			offSetOverride = zoneRec.Member?('offSetOverride')
				? zoneRec.offSetOverride
				: .zoneOffsets[timezone]
			nodaylightAdjustment? = zoneRec.Member?('nodaylightAdjustment?')
			}
		else
			{
			offSetOverride = .hasOverride?(zoneRec, timezone)
					? zoneRec.exceptions[timezone].offSetOverride
					: .zoneOffsets.GetDefault(timezone,
						.zoneOffsets[TimeZones[stateProv].zone])
			if zoneRec.Member?('exceptions') and
				Object?(zoneRec.exceptions) and
				zoneRec.exceptions.Member?(timezone)
				nodaylightAdjustment? = Object?(zoneRec.exceptions[timezone]) and
					zoneRec.exceptions[timezone].Member?('nodaylightAdjustment?')
			}

		return Object(:offSetOverride, :nodaylightAdjustment?)
		}

	timeZoneRec(stateProv)
		{
		return TimeZones[stateProv]
		}

	hasOverride?(zoneRec, timezone)
		{
		return zoneRec.Member?('exceptions') and
			Object?(zoneRec.exceptions) and
			zoneRec.exceptions.Member?(timezone) and
			Object?(zoneRec.exceptions[timezone]) and
			zoneRec.exceptions[timezone].Member?('offSetOverride')
		}

	/*
	This is to replace the existing DaylightSavings? so that LocalTime has all it needs to
	calculate local time. The daylight saving calculation can be accurate to seconds
	NOTE: ALWAYS pass in GMT time for the first argument (date)
	*/
	DaylightSavings?(date, timezone)
		{
		secondSundayOfMarch = Date(date.Year() $ '0314', 'yyyyMMdd')
		firstSundayInNovember = Date(date.Year() $ '1101', 'yyyyMMdd')
		Assert(Date?(secondSundayOfMarch) and Date?(firstSundayInNovember))

		dayOfWeek = secondSundayOfMarch.WeekDay()
		if dayOfWeek isnt 0
			secondSundayOfMarch = secondSundayOfMarch.Plus(days: -dayOfWeek)
		dayOfWeek = firstSundayInNovember.WeekDay()
		if dayOfWeek isnt 0
			firstSundayInNovember = firstSundayInNovember.Plus(days: 7 - dayOfWeek)

		dateInfo = Object(year: date.Year(), month: date.Month(), date: date.Day())
		if .transitionDay?(dateInfo, 3 /*= March*/, secondSundayOfMarch)
			return date >= .transitionDate(dateInfo, .forwardMap, timezone)
		else if .transitionDay?(dateInfo, 11 /*= November*/, firstSundayInNovember)
			return date < .transitionDate(dateInfo, .rollBackMap, timezone)

		return date >= secondSundayOfMarch and date < firstSundayInNovember
		}

	transitionDay?(dateInfo, month, tranDate)
		{
		return dateInfo.month is month and dateInfo.date is tranDate.Day()
		}

	transitionDate(dateInfo, map, timezone)
		{
		return Date('#' $ Display(dateInfo.year) $
			Display(dateInfo.month).LeftFill(2, '0') $
			Display(dateInfo.date).LeftFill(2, '0') $ '.' $
			map[timezone])
		}

	forwardMap: #(
		HST:	'1200',
		AKST:	'1100',
		PST:	'1000',
		MST: 	'0900',
		CST:	'0800',
		EST:	'0700',
		AST:	'0600',
		NST:	'0530'
		)

	rollBackMap: #(
		HST:	'1200',
		AKST:	'1000',
		PST:	'0900',
		MST: 	'0800',
		CST:	'0700',
		EST:	'0600',
		AST:	'0500',
		NST:	'0430'
		)
	}
