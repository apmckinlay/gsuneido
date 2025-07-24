// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// test [i ::], [i :: n], [i ..], and [.. n]
function (@args)
	{
	s = args[0]
	pos = Number(args[1])
	expected = args.Last()
	ob = s.Divide(1)
	expected_ob = expected.Divide(1)
	switch args.Size(list:)
		{
	case 3:
		if s[pos ::] isnt expected or s[pos ..] isnt expected or
			ob[pos ::] isnt expected_ob	or ob[pos ..] isnt expected_ob
			return false
	case 4:
		len = Number(args[2])
		if s[pos :: len] isnt expected or
			(pos is 0 and s[pos .. len] isnt expected) or
			ob[pos :: len] isnt expected_ob or
			(pos is 0 and ob[pos .. len] isnt expected_ob)
			return false
		}
	return true
	}