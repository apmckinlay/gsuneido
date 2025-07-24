// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// Ensure that a table contains a given record
// If the record exists but is different it will be deleted
// If the record doesn't exist (or was different) it will be output
// Block runs if the record did not already exist
function (table, rec, block = function(unused){})
	{
	keyfields = ShortestKey(table).Split(',')
	args = rec.Project(keyfields).Add(table)
	RetryTransaction()
		{|t|
		x = t.Query1(@args)
		if x isnt rec
			{
			if x isnt false
				x.Delete()
			t.QueryOutput(table, rec)
			block(t)
			}
		}
	}