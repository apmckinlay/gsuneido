// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if .admin? is true or .curUser is .user
		return true
	private? = Object?(.report_options)
		? .report_options.GetDefault(#private?, false)
		: false
	return not private?
	}
