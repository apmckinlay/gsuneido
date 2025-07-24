// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function (path = "*.*", files = false, details = false)
	{
	list = Dir(path, :files, :details)
	Assert(list.Size() lessThanOrEqualTo: 1, msg: "Dir1 got more than one record")
	return list.GetDefault(0, false)
	}