// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
CodeViewAddon
	{
	Name: 		HtmlDisplay
	Inject: 	bottomLeft
	InjectControls(container)
		{
		if false isnt control = .control(.Controller.Table, .Controller.RecName)
			container.Insert(0, control)
		}

	control(table, name)
		{
		// The following image extensions require ImageControl to display
		if name =~ `^/res\>` and name =~ "(?i)[.](emf|wmf|ico|cur)$"
			return Object('Image', .image(table, name),
				ystretch: 1, xstretch: 1, name: .Name)
		// Remainder of book records can be displayed via MshtmlControl.
		// However, we do not want to display the text output for the following extensions
		// - Javascript (.js), CSS (.css), FontMetrics (.ttf), large records (.afm/.map)
		else if name !~ `^/res\>` or name !~ '(?i)[.](js|css|ttf|afm|map)$'
			return Object('Scroll',
				Object('Mshtml', .wrapText(table, name), name: .Name) noEdge:)
		return .AddonControl = false
		}

	image(table, name)
		{ return table $ '%' $ name.RemovePrefix('/res/') }

	wrapText(table, name)
		{
		// These dynamic variables are used in e.g. Snippet and SeeAlsoGroup
		// which are referenced in suneidoc book record text directly
		_table = table
		_path = name.BeforeLast('/')
		_name = name.AfterLast('/')
		text = .text(table, name)
		if not text.Prefix?('<')
			try
				{
				x = text.Eval() // needs Eval
				text = String?(x) and x.Prefix?('<') ? x : ''
				}
		return HtmlWrap(text, table)
		}

	text(table, name)
		{
		return BookResource?(name, imagesOnly?:)
			? Xml('img', src: 'suneido:/' $ table $ name)
			: .Get()
		}

	Addon_RedirMethods()
		{ return #(Refresh) }

	IdleAfterChange()
		{ .Refresh() }

	Refresh(force? = false)
		{
		if .AddonControl isnt false and .refresh?(force?)
			.AddonControl.Set(.value(.Controller.Table, .Controller.RecName))
		}

	value(table, name)
		{
		return .AddonControl.Base() is ImageControl
			? .image(table, name)
			: .wrapText(table, name)
		}

	refresh?(force?)
		{ return IDESettings.Get(#ide_book_auto_refresh, defaultVal:) is true or force? }
	}