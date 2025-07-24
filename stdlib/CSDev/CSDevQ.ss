// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	try
		return ServerSuneido.Get('CSDev?') is true
	catch(unused, 'not authorized') // if not authorized assume we are a real client
		return false
	}

