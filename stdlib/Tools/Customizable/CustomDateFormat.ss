// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
DateFormat
	{
	New(@args)
		{
		super(@.setup(args))
		}

	setup(args)
		{
		// need to be done before constructor for width calculation
		.Method = args.GetDefault('long?', false) is true ? 'LongDate' : 'ShortDate'
		return args
		}
	}