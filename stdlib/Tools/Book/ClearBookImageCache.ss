// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function (book, clearQuery1Cache? = false)
	{
	if Sys.Client?()
		ServerEval(#ClearBookImageCache, book, :clearQuery1Cache?)
	if book is false
		Suneido.BookModels = Object()
	else if Suneido.Member?(#BookModels)
		Suneido.BookModels.Delete(book)
	if clearQuery1Cache?
		Query1CacheReset()
	}
