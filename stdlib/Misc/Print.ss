// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	// Suneido.Print is set by WorkSpace
	if not Suneido.Member?('Print')
		return // fallback is to do nothing
	s = ""
	for x in args.Values(list:)
		s $= String(x) $ ' '
	for m in args.Members(named:)
		s $= String(m) $ ": " $ String(args[m]) $ ' '
	(Suneido.Print)(s[..-1] $ '\r\n')
	return // no return value
	}
