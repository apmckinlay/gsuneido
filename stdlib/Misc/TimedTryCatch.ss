// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function (block)
	{
	t = Timer()
		{
		result = ''
		try
			result = block()
		catch(err)
			result $= '\tERROR (CAUGHT) ' $ err
		}
	return result $ ' in ' $ t $ ' seconds\r\n\r\n'
	}