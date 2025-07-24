// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function (block)
	{
	ma = MemoryAlloc()
	r = Timer.Secs(:block, secs: 5)
	ma = MemoryAlloc() - ma
	a = Max(0, (ma / r.reps).Round(0) - 16) /*= overhead is roughly 16 bytes */
	return Timer.Format(r) $ ", " $ ReadableSize(a) $ " allocation"
	}