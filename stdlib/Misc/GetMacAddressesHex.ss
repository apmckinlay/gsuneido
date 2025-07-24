// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	return GetMacAddresses().Map(#ToHex)
	}