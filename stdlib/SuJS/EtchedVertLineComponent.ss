// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Ystretch: 	1
	Name:		"EtchedVertLine"
	New(before = 2, after = 2)
		{
		.CreateElement('div')
		.SetStyles(Object(
			'border-left': 'solid 1px grey',
			'margin-left': before $ 'px',
			'margin-right': after $ 'px'))
		}
	GetReadOnly()			// read-only not applicable to etchedline
		{ return true }
	}
