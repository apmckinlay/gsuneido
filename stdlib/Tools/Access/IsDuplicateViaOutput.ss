// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// SEE ALSO: IsDuplicate
function (query, field, value)
	{
	try
		Transaction(update:)
			{ |t|
			rec = Record()
			rec[field] = value
			t.QueryOutput(query, rec)
			t.Rollback()
			}
	catch (unused, "*duplicate key")
		return true
	return false
	}
