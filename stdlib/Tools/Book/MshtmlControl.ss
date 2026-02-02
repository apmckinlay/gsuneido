// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Container
	{
	Name:		'Mshtml'
	Xstretch:	1
	Ystretch:	1

	New(.text = '', .allowReadOnly = false, .style = false)
		{
		.ctrl = .Construct(Object('WebBrowser', .convertValue(text)))
		.setStyle()
		.Send('Data')
		}

	toRemove: false
	convertValue(text)
		{
		.cleanInMemory()

		if text.Prefix?('http://') or text.Prefix?('https://')
			return text

		if text.Size() < 60.Kb() /*=64 Kb is the heap size allocated for c args.
									Leave 4 Kb for other args*/
			return "MSHTML:" $ text

		return .toRemove = InMemory.Add(text)
		}

	cleanInMemory()
		{
		if .toRemove isnt false
			{
			InMemory.Remove(.toRemove)
			.toRemove = false
			}
		}

	setStyle()
		{
		if .style isnt false
			.ctrl.SetCssStyle(.style)
		}
	x: false
	Set(.text)
		{
		// frequent setting could cause html not displaying
		.Defer({
			if not .Destroyed?()
				{
				.ctrl.Load(.convertValue(text))
				.setStyle()
				}
			}, uniqueID: .Name)
		}
	Get()
		{
		return .text
		}
	SetReadOnly(readonly)
		{
		if .allowReadOnly is true
			super.SetReadOnly(readonly)
		}
	Dirty?(dirty/*unused*/ = "")
		{
		return false
		}
	Valid?(@unused)
		{
		return true
		}
	Resize(.x, .y, .w, .h)
		{
		.ctrl.Resize(x, y, w, h)
		}
	GetChildren()
		{
		return Object(.ctrl)
		}
	// need to RedirAccels for Ctrl + F
	On_Find()
		{
		.ctrl.DoFind()
		}

	On_Copy()
		{
		.ctrl.DoCopy()
		}

	Default(@args)
		{
		return (.ctrl[args[0]])(@+1args)
		}

	Destroy()
		{
		.Send('NoData')
		.cleanInMemory()
		super.Destroy()
		}
	}
