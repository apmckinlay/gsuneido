// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
FormControl
	{
	Name: 'Grid'
	ComponentName: 'Grid'
	New(controls, setLeft/*unused*/ = false)
		{
		super(@.convertControls(controls))
		}

	convertControls(controls)
		{
		converted = Object()
		for row in controls
			{
			group = 0
			for ctrl in row
				{
				inc = .span?(ctrl) ? ctrl.span : 1
				converted.Add(not Object?(ctrl)
					? Object(ctrl, :group)
					: ctrl.Copy().Add(group, at: 'group'))
				group += inc
				}
			converted.Add('nl')
			}
		return converted
		}

	span?(ctrl)
		{
		return Object?(ctrl) and ctrl.Member?('span')
		}
	}