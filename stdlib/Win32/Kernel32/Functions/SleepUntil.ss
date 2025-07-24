// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// - guarantee sleep up to a specified time, since Sleep might wake up earlier than expected
// - maximum sleep time is 24 hours, an exception will be thrown if exceeded. Due to
//	limitations of integer size there will be integer conversion issues at around 24 days
//	of sleep time in ms but 24 hours seems like a more practical limit for this function
function (until)
	{
	if until.MinusHours(Date()) > 24 /* = max sleep time in hours */
		throw "sleep time exceeded 24 hours, system date may not be correct"
	// make sure we only call Date() once per iteration and that Sleep is > 0
	while 0 < (ms = until.MinusSeconds(Date()) * 1000) /* = 1000 ms in one second */
		Thread.Sleep(ms)
	}