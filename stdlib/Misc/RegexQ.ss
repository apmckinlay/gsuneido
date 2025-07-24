// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
function (regex)
	{
	if not String?(regex)
		return false
	try
		{
		unused = "" =~ regex
		return true
		}
	catch
		return false
	}