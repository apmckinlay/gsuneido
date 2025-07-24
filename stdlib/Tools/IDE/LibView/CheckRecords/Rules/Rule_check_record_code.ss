// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	CheckCode(.text, .name, .table, results = Object())
	return results.Map!({ it.msg }).Join('\r\n')
	}