// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if Suneido.Member?('IM_Available?')
		return Suneido.IM_Available?
	Suneido.IM_Available? = true
	Plugins().ForeachContribution('IM_Available', 'im_available')
		{ |c|
		if false is (c.func)()
			Suneido.IM_Available? = false
		}
	return Suneido.IM_Available?
	}