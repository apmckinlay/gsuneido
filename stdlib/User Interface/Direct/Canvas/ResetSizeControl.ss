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
	//FIXME: width and height editing are not working on Round Rectangle
	controls(fields)
		{
		ctrl = Object('Vert')
		//QUESTION: sort members instead of random order ?
		for field in fields.Members()
			ctrl.Add(Object('Pair'
				Object('Static' field)
				Object('Field' name: field set: fields[field])))
		return Object('Record', ctrl)
		}

	OK()
		{
		return .Data.Get()
		}
	}
