// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// NOTE: the use of the datadict format may not always be appropriate
HorzFormat
	{
	New(@args)
		{
		super(@.format(args))
		}
	format(args)
		{
		fmt = Object()
		for (i = 0; args.Member?(i); ++i)
			{
			if i > 0
				fmt.Add("Hskip")
			arg = args[i]
			field = Object?(arg) ? arg.field : arg
			fmt.Add(Object("Text", Opt(Heading(field), ": ")))
			fmt.Add(arg)
			}
		if args.Member?("data")
			fmt.data = args.data
		if args.Member?("font")
			fmt.font = args.font
		return fmt
		}
	}