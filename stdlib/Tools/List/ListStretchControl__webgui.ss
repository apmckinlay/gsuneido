// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
_ListStretchControl
	{
	New(@args)
		{
		super(@args)
		.ComponentArgs.stretch = args.GetDefault('stretchColumn', true)
		}

	Adjust(@unused)
		{
		}
	}
