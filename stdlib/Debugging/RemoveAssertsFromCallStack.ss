// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
function (calls)
	{
	while calls.Size() > 0 and Display(calls[0].fn).Prefix?('Assert')
		calls.Delete(0)
	}