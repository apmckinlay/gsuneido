// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Field_internal
	{
	Encode(val)
		{
		if Object?(val)
			return val
		if not String?(val)
			return val
		if val is ""
			return Object()
		try
			ob = val.SafeEval()
		catch
			ob = Object()
		return ob
		}
	}