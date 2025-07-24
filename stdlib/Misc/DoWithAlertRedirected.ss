// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function(alertFn, block)
	{
	oldAlert = Suneido.GetDefault('Alert', false)
	Suneido.Alert = alertFn
	Finally(block,
		{
		if oldAlert is false
			Suneido.Delete("Alert")
		else
			Suneido.Alert = oldAlert
		})
	}