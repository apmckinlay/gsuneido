// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function(key, replace = "[0-9]*$", length = 50)
	{
	return key.Replace(replace, "")[:: length]
	}
