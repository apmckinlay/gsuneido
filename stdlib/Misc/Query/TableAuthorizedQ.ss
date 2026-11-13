// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(table)
		{
		return TableExists?(table)
			? .authorized?(table)
			: false
		}

	authorized?(table)
		{
		try
			return .query(table)
		catch (unused, "not authorized")
			return false
		}

	query(table) // extracted for testing
		{
		WithQuery(table)
			{|q|
			q.Columns()
			}
		return true
		}
	}
