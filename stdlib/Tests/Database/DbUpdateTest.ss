// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_gsuneido_bug()
		{
		tbl = .MakeTable("(k, _TS) key(k)")
		r = [k: 1]
		r.Observer({ throw "observer should not be called" })
		QueryOutput(tbl, r)
		}
	}