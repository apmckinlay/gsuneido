// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
function (t, block, update = false)
	{
	if t is false
		Transaction(:update, :block)
	else
		block(t)
	}