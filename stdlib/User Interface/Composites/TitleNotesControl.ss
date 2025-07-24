// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Xstretch: 0
	New(@args)
		{
		super(.controls(args))
		.static = .Horz.Title
		}
	controls(args)
		{
		title = args[0]
		if title.Has?(' > ')
			title = title.AfterLast(' > ')
		titleLeftCtrl = args.GetDefault('titleLeftCtrl', false)
		if titleLeftCtrl is false
			titleLeftCtrl = #(Skip 25) // to balance button
		args.Delete('titleLeftCtrl')
		controls = ['Horz',
			titleLeftCtrl,
			'Fill',
			['Title', title], // dark blue
			'Skip',
			args.GetDefault('extra', #(Skip 0)),
			'Fill',
			['TitleButtons', args]
			]
		if args.GetDefault('center', false) isnt true
			controls.Delete(2).Delete(1)
		return controls
		}
	Set(title)
		{
		.static.Set(title)
		}
	}