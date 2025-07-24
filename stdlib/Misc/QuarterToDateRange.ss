// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function (quarter, year = false)
	{
	qtrRanges = #(
		1: (start: '0101', end: '0331'),
		2: (start: '0401', end: '0630'),
		3: (start: '0701', end: '0930'),
		4: (start: '1001', end: '1231'))

	range = qtrRanges[Number(quarter)]
	year = year is false ? '#' $ Date().Year() : '#' $ year
	return Object(start: Date(year $ range.start), end: Date(year $ range.end))
	}