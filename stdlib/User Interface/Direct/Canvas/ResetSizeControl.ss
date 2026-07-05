// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Please enter the new coordinates'

	CallClass(hwnd, fields)
		{
		return OkCancel(Object(this, fields), .Title, hwnd)
		}

	New(fields)
		{
		super(.controls(fields))
		}

	controls(fields)
		{
		ctrl = Object('Vert')
		for field in fields.Members()
			ctrl.Add(Object('Pair'
				Object('Static' field)
				Object('Number', name: field, set: fields[field],
					mask: '#####.##', mandatory:, xstretch: 1)))
		return Object('Record', ctrl)
		}

	OK()
		{
		return .Data.Valid(forceCheck:) is true
			? .Data.Get()
			: false
		}
	}
