// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Xstretch: 	0
	Name:		"EtchedLine"
	New(before = 2, after = 2)
		{
		.CreateElement('div')
		line = CreateElement('div', .El)
		line.SetStyle('border-top', 'solid 1px grey')
		.SetStyles(Object(
			'box-sizing': 'border-box',
			'padding-top': before $ 'px',
			'padding-bottom': after $ 'px'))
		.Ymin = 1 + before + after
		.SetMinSize()
		}
	GetReadOnly()			// read-only not applicable to etchedline
		{ return true }
	}
