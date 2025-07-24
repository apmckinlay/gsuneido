// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function(e)
	{
	if e.Trim().Prefix?('curl:') or
		e.Lower().Has?("500 internal server error") or
		e.Has?("(empty header)") // from Http.ResponseCode
		return true
	code = e.Extract(' \d\d\d ')
	if code isnt false and code.Trim().Prefix?('5')
		return true
	return false
	}