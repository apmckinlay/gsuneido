// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	Name: 'Pane'
	Xmin:		false
	Ymin:		false
	Xstretch:	false
	Ystretch:	false

	Initialize(ctrl)
		{
		super.Initialize(ctrl)
		.SetStyles(#('border-width': '2px',
			'border-style': 'inset',
			'align-self': 'flex-start'))
		}
	}
