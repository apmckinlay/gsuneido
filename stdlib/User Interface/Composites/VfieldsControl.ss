// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
FormControl
	{
	New(@controls)
		{
		super(@.makecontrols(controls))
		}
	makecontrols(controls)
		{
		if controls.Size() is 1 and controls.Member?(0)
			controls = controls[0]
		form = Object()
		j = 0
		for i in controls.Members()
			{
			c = controls[i]
			if String?(c)
				c = Object(c, group: 1)
			else
				{
				c = c.Copy()
				c.group = 1
				}
			form[j++] = c
			form[j++] = 'nl'
			}
		return form
		}
	}