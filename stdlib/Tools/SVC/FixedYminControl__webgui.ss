// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	ComponentName: 'FixedYmin'
	New(ymin, control)
		{
		super(control)
		.ComponentArgs.Add(ymin, at: 0)
		}
	}