// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(file)
		{
		cksum = Adler32()
		.Update(file, cksum)
		return cksum.Value()
		}
	Update(file, cksum)
		{
		File(file, mode: "r")
			{ |src|
			while (false isnt s = src.Read(65536 /*= 64k at a time */))
				cksum.Update(s)
			}
		}
	}