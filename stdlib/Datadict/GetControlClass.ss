// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	FromField(fieldName)
		{
		dd = Datadict(fieldName)
		return .FromControl(dd.Control)
		}

	FromControl(controlOb)
		{
		// Based on the implementaion of Control.Construct method, a valid Datadict Control
		// has to be an object with its first element having a string or class type
		if not Object?(controlOb) or not controlOb.Member?(0)
			return false

		control = controlOb[0]

		if Class?(control)
			return control

		if not String?(control)
			return false

		if not control.Suffix?('Control')
			control $= 'Control'
		try
			ctrl = Global(control)
		catch (err, "can't find")
			{
			SuneidoLog("ERROR: (CAUGHT) " $ err)
			return false
			}
		return ctrl
		}
	}