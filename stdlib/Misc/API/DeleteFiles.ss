// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// NOTE: does not delete sub-folders if they are included in the pattern
function (pattern, failedToDelete = false)
	{
	Assert(pattern.Has?('*') or pattern.Has?('?'))

	if failedToDelete is false
		failedToDelete = Object()
	dir = Paths.ParentOf(pattern) $ '/'
	for file in Dir(pattern, files:)
		if true isnt DeleteFile(path = dir $ file)
			failedToDelete.Add(path)
	return failedToDelete
	}