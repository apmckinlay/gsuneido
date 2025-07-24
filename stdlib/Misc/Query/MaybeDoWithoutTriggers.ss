// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
function (tables, update, block)
	{
	if update is true
		DoWithoutTriggers(tables, block)
	else
		block()
	}