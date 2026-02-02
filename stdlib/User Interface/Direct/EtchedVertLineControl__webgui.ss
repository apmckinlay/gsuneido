// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Ystretch: 		1
	Name:			"EtchedVertLine"
	ComponentName: 	"EtchedVertLine"
	New(before = 2, after = 2)
		{
		.ComponentArgs = Object(before, after)
		}
	GetReadOnly()			// read-only not applicable to etchedline
		{ return true }
	}