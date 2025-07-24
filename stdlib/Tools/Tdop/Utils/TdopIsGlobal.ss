// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function(name)
	{
	return String?(name) and name.RemovePrefix('_').GlobalName?()
	}