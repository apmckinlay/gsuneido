// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Xstretch: 		0
	Name:			"EtchedLine"
	ComponentName:	"EtchedLine"
	New(before = 2, after = 2)
		{
		.ComponentArgs = Object(before, after)
		}
	GetReadOnly()			// read-only not applicable to etchedline
		{ return true }
	}