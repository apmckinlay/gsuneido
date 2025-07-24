// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// all this does is default keep_size to true, maybe that should be in Dialog ?
class
	{
	CallClass(parentHwnd, control, title = false, closeButton? = true,
		keep_size = true, border = 5, posRect = false, useDefaultSize = false)
		{
		return Dialog(parentHwnd, control,
			:title, :keep_size, :border, :closeButton?, :posRect, :useDefaultSize)
		}
	}
