// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: "NameAbbrev"
	ComponentName: "NameAbbrev"
	New(prefix = '', suffix = '')
		{
		super(.Layout(prefix, suffix))
		firstChild = .Horz.GetChildren()[0]
		.Left = firstChild.Left
		.Top = firstChild.Top
		}

	Layout(prefix = '', suffix = '')
		{
		return Object('Horz'
			Object('MainField1', prefix $ 'name' $ suffix),
			'Skip'
			Object('MainField1', prefix $ 'abbrev' $ suffix),
			#(Skip 4)
			#(Static '(optional)')
			'Skip'
			)
		}
	}