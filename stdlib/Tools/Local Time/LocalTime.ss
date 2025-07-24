// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(to_state_prov, from_state_prov, to_timezone = '', from_timezone = '')
		{
		curdate = Date()
		if .cantConvert?(to_state_prov, from_state_prov)
			return curdate
		if to_timezone is ''
			{
			if not TimeZones.Member?(to_state_prov)
				return curdate
			else
				to_timezone = TimeZones[to_state_prov].zone
			}
		if from_timezone is ''
			{
			if not TimeZones.Member?(from_state_prov)
				return curdate
			else
				from_timezone = TimeZones[from_state_prov].zone
			}

		torec = TimeZones[to_state_prov]
		fromrec = TimeZones[from_state_prov]

		offset = .calc_offset(to_timezone, from_timezone)
		curdate = curdate.Plus(hours: offset.hours, minutes: offset.minutes)

		curdate = .calc_daylightsavings(curdate, torec, fromrec)

		return curdate
		}

	cantConvert?(to_state_prov, from_state_prov)
		{
		return to_state_prov is "" or from_state_prov is "" or
			not TimeZones.Member?(to_state_prov) or not TimeZones.Member?(from_state_prov)
		}

	calc_offset(to_timezone, from_timezone)
		{
		offset_to_zone = .offset_times(to_timezone)
		offset_to_from = .offset_times(from_timezone)

		hours = minutes = 0
		if offset_to_from.hours is offset_to_zone.hours and
			offset_to_from.minutes is offset_to_zone.minutes
			{
			hours = 0
			minutes = 0
			}
		else if offset_to_from.hours >= offset_to_zone.hours
			{
			hours = -(offset_to_from.hours - offset_to_zone.hours)
			minutes = -(offset_to_from.minutes - offset_to_zone.minutes)
			}
		else if offset_to_from.hours < offset_to_zone.hours
			{
			hours = (offset_to_from.hours - offset_to_zone.hours).Abs()
			minutes = offset_to_from.minutes + offset_to_zone.minutes.Abs()
			}
		return Object(:minutes, :hours)
		}

	zoneOffsets: (
		HST : -4
		AKST: -3
		PST : -2
		MST : -1
		CST :  0
		EST :  1
		AST :  2
		NST :  2)

	offset_times(zone) // based on SK as home time-zone
		{
		hours = .zoneOffsets.GetDefault(zone, 0)
		minutes = 0
		if zone is "NST"
			minutes = 30
		return Object(:minutes, :hours)
		}

	calc_daylightsavings(curdate, torec, fromrec)
		{
		try
			{
			if not DaylightSavings?(curdate) and not torec.Member?('nodaylight?') and
				fromrec.Member?('nodaylight?')
				return curdate.Plus(hours: 1)

			if not DaylightSavings?(curdate) and not fromrec.Member?('nodaylight?') and
				torec.Member?('nodaylight?')
				return curdate.Plus(hours: -1)
			}
		catch (e)
			{
			SuneidoLog('ERROR: (CAUGHT) ' $ e $ ' while trying to calc daylight saving',
				caughtMsg: 'Displayed to user there was a problem getting local time')
			return 'Problem getting local time'
			}

		return curdate
		}
	}