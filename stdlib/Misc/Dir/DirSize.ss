// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
function(dir = "")
	{
	if dir is ""
		dir = "."
	size = 0
	Dir(dir $ '/*.*', details:) // used block to handle too many files (>10000)
		{
		size += it.name.Suffix?('/') ? DirSize(dir $ '/' $ it.name) : it.size
		}
	return size
	}