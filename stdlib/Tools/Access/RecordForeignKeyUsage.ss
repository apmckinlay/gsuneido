// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function (query, record)
	{
	func = false
	try func = Global((QueryGetTable(query) $ "_show_fk_usage").Capitalize())
	if not Function?(func)
		return ""
	return record.Eval(func)
	}