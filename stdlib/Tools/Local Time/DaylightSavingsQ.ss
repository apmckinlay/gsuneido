// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
function (date)
	{
	// Is date within the standard Canada/U.S. daylight savings period?
	// second Sunday in March , 2 A.M. to first Sunday in November, 2 A.M
	// The answer will be wrong for up to two hours on the transition dates,
	// but this should not matter much since transition time is 2 A.M.

	secondSundayOfMarch = Date(date.Year() $ '0314', 'yyyyMMdd')
	dayOfWeek = secondSundayOfMarch.WeekDay(sun: 0)
	if dayOfWeek isnt 0
		secondSundayOfMarch = secondSundayOfMarch.Plus(days: -dayOfWeek)

	firstSundayInNovember = Date(date.Year() $ '1101', 'yyyyMMdd')
	dayOfWeek = firstSundayInNovember.WeekDay(sun: 0)
	if dayOfWeek isnt 0
		firstSundayInNovember = firstSundayInNovember.Plus(days: 7 - dayOfWeek)

	return date < secondSundayOfMarch or date >= firstSundayInNovember
	}