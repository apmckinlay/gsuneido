// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
function (table, rec)
	{
	if Client?()
		return ServerEval('QueryOutputIfNew', table, rec)
	try
		QueryOutput(table, rec)
	catch (unused, "duplicate key")
		return false
	return true
	}
