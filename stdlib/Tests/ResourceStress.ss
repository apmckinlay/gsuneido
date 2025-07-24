// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	stressTestLimit: 500
	Controls()
		{
		ob = Object('Vert')
		for .. .stressTestLimit
			ob.Add(#(Static hello))
		return ob
		}
	}