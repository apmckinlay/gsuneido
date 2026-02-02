// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
HorzControl
	{
	Name: 'Flow'
	ComponentName: 'Flow'

	New(controls, skip = 8)
		{
		super(@.convert(controls, skip))
		}

	convert(controls, skip)
		{
		first = true
		converted = Object()
		for control in controls
			{
			if first is true
				first = false
			else
				converted.Add(Object('Skip', skip))
			converted.Add(control)
			}
		return converted
		}
	}