// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
HtmlDivComponent
	{
	New(@args)
		{
		super(@.buildArgs(args))
		}

	buildArgs(args)
		{
		.ymin = args[0]
		args.Delete(0)
		return args
		}

	Recalc(@args)
		{
		super.Recalc(@args)
		.Ymin = .ymin
		.SetMinSize()
		}
	}