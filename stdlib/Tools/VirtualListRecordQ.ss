// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function (rec)
	{
	return rec.Members().Any?({ it.Prefix?('vl_') })
	}