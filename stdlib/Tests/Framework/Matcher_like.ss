// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// trims leading and trailing whitespace from each line
// removes \r to standardize line endings
class
	{
	ellipsis: 100
	Match(value, args, _debugTest = false)
		{
		if not String?(value)
			return false
		if .Tr(args) is .Tr(value)
			return true
		if debugTest
			Diff2Control('Assert FAILED', .Tr(value), .Tr(args), 'Actual', 'Expected')
		return false
		}
	Tr(s)
		{
		return s.Trim().Tr('\r').Tr(' \t', ' ').
			Replace('^ +', '').Replace(' +$', '')
		}
	Expected(args)
		{
		return "a string like:\n" $ Display(.Tr(args)).Ellipsis(.ellipsis) $ '\n'
		}
	Actual(value, args/*unused*/)
		{
		return String?(value)
			? "was:\n" $ Display(.Tr(value)).Ellipsis(.ellipsis)
			: "was a " $ Type(value)
		}
	}