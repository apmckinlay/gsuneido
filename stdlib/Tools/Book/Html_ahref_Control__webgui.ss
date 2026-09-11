// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Control
	{
	ComponentName: #Html_ahref_
	New(text, .href, hrefOnBrowser = false)
		{
		.ComponentArgs = [text, :hrefOnBrowser]
		}

	LBUTTONUP()
		{
		.Send(#Goto, .href)
		return 0
		}
	}
