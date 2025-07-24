// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
SplitterComponent
	{
	SplitName: "Splitter"
	styles: `
		.su-spliter-handle-buton-vert {
			margin: 0 auto;
			width: 34px;
			height: 100%;
			cursor: pointer;
			text-align: center;
		}
		.su-spliter-handle-buton-horz {
			margin: auto 0;
			width: 100%;
			height: 34px;
			cursor: pointer;
			position: absolute;
			top: calc(50% - 17px)
		}
		.su-spliter-handle-buton-vert:hover,
		.su-spliter-handle-buton-horz:hover {
			background-color: lightblue;
		}
		.su-spliter-handle-buton-caret-vert {
			width: 0;
			height: 0;
			display: inline-block;
			vertical-align: middle;
			border-left: 4px solid transparent;
			border-right: 4px solid transparent;
		}
		.su-spliter-handle-buton-caret-horz {
			width: 0;
			height: 0;
			display: inline-block;
			position: absolute;
			top: calc(50% - 4px);
			border-top: 4px solid transparent;
			border-bottom: 4px solid transparent;
		}`
	New()
		{
		LoadCssStyles('handle-splitter-control.css', .styles)
		.associate = .Parent.Associate
		.dir = .associate in (#north, #south) ? 'vert' : 'horz'
		.handle = CreateElement('div', .El, className: 'su-spliter-handle-buton-' $ .dir)
		.button = CreateElement('span', .handle,
			className: 'su-spliter-handle-buton-caret-' $ .dir)
		.handle.AddEventListener('click', .click)
		}

	click()
		{
		if .Parent.Open?
			.Parent.Close()
		else
			.Parent.Open()
		}

	currentStatus: ''
	buttonStyles: (
		close:	(east: right, north: top, south: bottom, west: left)
		open:	(east: left, north: bottom, south: top, west: right))

	UpdateButton()
		{
		if .Parent.Open? is .currentStatus
			return
		.currentStatus = .Parent.Open?
		.button.SetStyle(
			'border-' $ .buttonStyles[.Parent.Open? ? 'open' : 'close'][.associate],
			'6px dashed')
		.button.SetStyle(
			'border-' $ .buttonStyles[.Parent.Open? ? 'close' : 'open'][.associate],
			'')
		}
	}
