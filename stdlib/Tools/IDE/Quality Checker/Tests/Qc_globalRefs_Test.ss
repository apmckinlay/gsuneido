// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		code = Query1(#stdlib, name: #Qc_globalRefs_Test).text
		Assert(Qc_globalRefs(code, true) is: #(Assert: 1, Test: 1, Qc_globalRefs: 1))
		}
	}