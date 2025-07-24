// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	stdReports = Object()
	for list in GetContributions('StandardScheduleReports')
		stdReports.Merge(list())
	return stdReports
	}