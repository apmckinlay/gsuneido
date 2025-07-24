// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	New(.ymin, control)
		{
		super(control)
		}
	Recalc()
		{
		super.Recalc()
		.Ymin = ScaleWithDpiFactor(.ymin)
		}
	}