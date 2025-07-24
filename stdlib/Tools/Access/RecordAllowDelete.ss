// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function (query, record)
	{
	// keep the allow delete function separate from the reason protected rule
	// since the allow_delete will in some cases be doing queries
	// and we need to limit how often the rule is used which is difficult
	// to do with the protect rules
	func = false
	try func = Global((QueryGetTable(query) $ "_allow_delete").Capitalize())
	if not Function?(func)
		return ""
	return record.Eval(func)
	}
