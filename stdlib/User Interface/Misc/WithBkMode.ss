// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function(hdc, mode, block)
	{
	oldBkMode = SetBkMode(hdc, mode)
	return Finally(block, { SetBkMode(hdc, oldBkMode) })
	}
