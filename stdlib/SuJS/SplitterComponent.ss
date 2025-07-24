// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Xmin: 0
	Xstretch: 0
	Ymin: 0
	Ystretch: 0

	styles: `
		.su-spliter-handle {
			position: relative;
			user-select: none;
			background-color: lightgrey;
		}`

	New()
		{
		LoadCssStyles('splitter-control.css', .styles)

		.CreateElement('div', className: 'su-spliter-handle')

		if _parent.Dir is 'vert'
			{
			.dir = 'vert'
			.Ymin = 6
			.Xstretch = 1
			.SetStyles([cursor: 'row-resize'])
			}
		else
			{
			.dir = 'horz'
			.Xmin = 6
			.Ystretch = 1
			.SetStyles([cursor: 'col-resize'])
			}
		.SetMinSize()

		.El.AddEventListener('mousedown', .mousedown)
		}

	mousedown(event)
		{
		if event.target isnt .El or event.button isnt 0
			return
		.Parent.Splitter_mousedown(.dir is 'vert' ? event.y : event.x)
		.StartMouseTracking(.mouseup, .mousemove)
		}

	mousemove(event)
		{
		if .dir is 'vert'
			.Parent.Splitter_mousemove(event.y)
		else
			.Parent.Splitter_mousemove(event.x)
		}

	mouseup(event/*unused*/)
		{
		.Parent.Splitter_mouseup()
		.StopMouseTracking()
		}
	}
