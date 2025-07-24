// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (args)
	{
	s = ""
	for x in args.Values(list:)
		s $= (String?(x) ? x : Display(x)) $ ' '
	for m in args.Members(named:)
		s $= m $ ': ' $ Display(args[m]) $ ' '
	return s[..-1]
	}