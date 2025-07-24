// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (hDrop)
	{
	list = Object()
	for i in .. DragQueryFileCount(hDrop)
		list.Add(DragQueryFile(hDrop, i))
	return list
	}