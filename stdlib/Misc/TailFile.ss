// Copyright (C) 2025 Axon Development Corporation All rights reserved worldwide.
function (filename, size)
	{
	maxSize = 10_000_000
	Assert(Number?(size) and size < maxSize)
	try
		File(filename)
			{|f|
			if f.Size() > size
				f.Seek(-size, 'end')
			s = f.Read(size)
			return s is false ? "" : s
			}
	catch
		return false
	}