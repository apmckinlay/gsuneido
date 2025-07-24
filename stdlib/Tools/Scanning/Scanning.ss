// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		// To keep 1 instance of Clscan, So the getter on Scanners won't fetch every time
		return Suneido.GetInit(#Scanning, { new Clscan })
		}
	}
