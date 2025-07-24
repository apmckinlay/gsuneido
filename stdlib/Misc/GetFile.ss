// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (filename, limit = false)
	{
	max = 10_000_000
	Assert(limit is false or (Number?(limit) and limit < max))
	try
		File(filename)
			{|f|
			s = f.Read(limit is false ? max : limit)
			}
	catch
		return false
	if s isnt false and limit is false and s.Size() is max
		throw "GetFile bigger than 10mb from " $ filename
	return s is false ? "" : s
	}
