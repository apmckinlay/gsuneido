// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	s = args.Join(" ")
	if (not Suneido.Member?("log"))
		Suneido.log = Object()
	Suneido.log.Add(s)
	}