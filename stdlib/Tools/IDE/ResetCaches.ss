// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if Sys.Client?()
		ServerEval('ResetCaches')
	LruCache.ResetAll()
	Singleton.ResetAll()
	BookModel.ClearCache()
	MemoizeSingle.ClearAll()
	Memoize.ResetAllExpire()

	for c in Contributions('ResetCaches')
		c()
	}
