// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
ChooseListControl
	{
	New(@args)
		{
		super(@.args(args))
		}
	args(args)
		{
		args.Add(BookTables(), at: 0)
		if not args.Member?("allowOther") or args.allowOther is false
			args.selectFirst = true
		return args
		}
	}