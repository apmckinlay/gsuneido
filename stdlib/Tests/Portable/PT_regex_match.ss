// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	if args.Size(list:) is 2 or args[2] is true
		return args[0] =~ args[1]
	else if args[2] is 'false'
		return args[0] !~ args[1]
	if false is m = args[0].Match(args[1])
		return false
	for (i = 2; i < args.Size(list:); ++i)
		if args[i] isnt args[0][m[i - 2][0] :: m[i-2][1]]
			return false
	return true
	}