// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(@args)
		{
		super(@(.setup(args)))
		}
	setup(args)
		{
		// convert data passed to TextFormat to the actual prompt
		for m in #(0 data)
			if args.Member?(m) and String?(args[m])
				args[m] = SelectPrompt(args[m])
		return args
		}
	// needed for browse control in Gear Option
	Print(x, y, w, h, data = "")
		{
		if .Data isnt false
			data = .Data
		data = SelectPrompt(data)
		super.Print(x, y, w, h, data)
		}
	}