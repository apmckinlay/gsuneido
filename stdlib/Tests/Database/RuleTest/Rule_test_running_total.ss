// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (setup = false)
	{
	if setup
		return
	.DoWithTran()
		{ |t|
		q = t.Query("testrule_tmp")
		s = ""
		while false isnt x = q.Next()
			s $= x.name $ " "
		}
	return s
	}
