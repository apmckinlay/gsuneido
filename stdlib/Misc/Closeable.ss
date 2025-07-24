// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// Base class for classes that hold a resource which needs to be "closed" to
// prevent a memory leak
class
	{
	open()
		{
		if .open?()
			throw "already open"
		.open_ = true
		return
		}
	open?()
		{
		.GetDefault("Closeable_open_", false)
		}
	Open?()
		{
		.open?()
		}
	Close()
		{
		.Delete("Closeable_open_")
		}
	}
