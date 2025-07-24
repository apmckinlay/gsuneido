// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// test [from .. to]
function (s, from, to, expected)
	{
	from = Number(from)
	to = Number(to)
	ob = s.Divide(1)
	expected_ob = expected.Divide(1)
	if s[from .. to] is expected and ob[from .. to] is expected_ob
		return true
	Print("\t", "got", s[from..to], "should be", expected)
	return false
	}