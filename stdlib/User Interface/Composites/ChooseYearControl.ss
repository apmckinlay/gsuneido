// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
ChooseListControl
	{
	Name: ChooseYear
	New(@args)
		{
		super(@.args(args))
		}
	args(args)
		{
		args.Add(.yearlist(), at: 0)
		if not args.Member?('width')
			args.width = 4
		if not args.Member?('status')
			args.status = "a four digit year e.g. 2003"
		return args
		}
	yearlist()
		{
		startYear = 2000
		futureYears = 7
		return Seq(startYear, Date().Year() + futureYears)
		}
	}