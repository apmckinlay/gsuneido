// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	if args.Size() is 1
		return args[0]
	helper = class
		{
		UseDeepEquals: true
		New(.args)
			{
			}
		Call(@args2)
			{
			args = .args
			if not args2.Empty?()
				args = args.Copy().Append(args2)
			return (args[0])(@+1 args)
			}
		ToString()
			{
			return "Bind" $ Display(.args)
			}
		}
	return helper(args)
	}