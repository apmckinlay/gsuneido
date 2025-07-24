// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	try
		for cont in Contributions('CustomFontSaveName').Reverse!()
			if '' isnt saveName = (cont)()
				return saveName
	catch (e)
		SuneidoLog('ERROR: CustomFontSaveName ' $ e)
	return GetComputerName() $ '_logfont'
	}