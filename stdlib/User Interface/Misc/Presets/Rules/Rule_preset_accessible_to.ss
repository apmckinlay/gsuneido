// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	return not Object?(.report_options) or
		not .report_options.GetDefault(#private?, false)
		? 'Everyone'
		: .user
	}
