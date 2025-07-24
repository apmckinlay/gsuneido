// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function (block)
	{
	SpyManager().RemoveAll()
	Finally(block, { SpyManager().RemoveAll() })
	}