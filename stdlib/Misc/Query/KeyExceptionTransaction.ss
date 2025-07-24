// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (t = false, block = false, update = true)
	{
	return KeyException.TryCatch()
		{
		DoWithTran(:t, :block, :update)
		}
	}