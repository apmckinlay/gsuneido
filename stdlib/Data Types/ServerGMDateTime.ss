// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if Sys.Client?()
		return ServerEval(#ServerGMDateTime)

	return Date().GMTime()
	}