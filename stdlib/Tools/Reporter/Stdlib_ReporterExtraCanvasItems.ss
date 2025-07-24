// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
#((
	name: 'DAF',
	needData?:,
	button: (command: #DAF image: 'D', name: 'daf', tip: 'Insert Data Field'),
	handler: function (canvas, palette)
		{
		palette.SetButtons('daf')
		canvas.SetTracker(DrawClickTracker, DrawDAFAdapter)
		}))