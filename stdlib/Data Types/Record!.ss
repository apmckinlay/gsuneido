// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
function(@args)
	{
	if args.Size(list:) is 1 and not Function?(args[0])
		x = args[0]
	else
		x = (args[0])(@+1args)
	if not Record?(x)
		throw "expected record got " $ Type(x) $ " from: " $ Display(args)
	return x
	}
