// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Group
	{
	Dir: "vert"
	Name: "Vert"
	New(@controls)
		{
		super(.processControls(controls))
		}

	processControls(controls)
		{
		if Windows11?()
			controls.overlap = false
		return controls
		}
	}