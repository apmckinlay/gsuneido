// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if "" isnt from = UserSettings.Get('email_from', '')
		return from
	return OptContribution('EmailFrom', function () { return '' })()
	}