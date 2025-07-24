// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
function (file1, file2, chunkSize = 65536/*= 64Kb */)
	{
	File(file1, "r")
		{ |f1|
		File(file2, "r")
			{ |f2|
			forever
				{
				chunk1 = f1.Read(chunkSize)
				chunk2 = f2.Read(chunkSize)
				if chunk1 is false and chunk2 is false
					return true
				if chunk1 isnt chunk2
					return false
				}
			}
		}
	}