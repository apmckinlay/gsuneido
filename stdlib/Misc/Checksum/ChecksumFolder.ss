// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(dir)
		{
		cksum = Adler32()
		.Update(dir, cksum)
		return cksum.Value()
		}
	Update(dir, cksum)
		{
		for file in Dir(dir $ '/*.*').Sort!()
			{
			if file.Suffix?('/')
				.Update(dir $ '/' $ file, cksum)
			else
				ChecksumFile.Update(dir $ '/' $ file, cksum)
			}
		}
	}