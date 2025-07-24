// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (name)
	{
	return name.Has?('_') and
		Libraries().Map!(#Capitalize).Has?(name.BeforeFirst('_'))
	}