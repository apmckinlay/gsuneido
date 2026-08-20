// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: LibLocate

	New()
		{
		super([LocateAutoChooseControl, LibLocateList.GetMatches, width: 15, xstretch: 1,
			cue: 'Locate by name (Ctrl+L)', allowOther:])
		.sub = PubSub.Subscribe(#LibraryTreeChange, LibLocateList.ForceRun)
		LibLocateList.Start()
		}
	NewValue(value)
		{
		if value !~ `^` $ GlobalRegExForGoTo $ ` - \(?[[:alpha:]]\w*\)?$`
			return
		// convert e.g. 'Name - lib' to 'lib:Name' (for GotoLibView)
		.value = value.AfterFirst(' - ').Tr('()') $ ':' $ value.BeforeFirst(' ')
		.Send('Locate', .value)
		}
	value: ''
	Get()
		{
		return .value
		}
	Set(value)
		{
		.AutoChoose.Set(value)
		.AutoChoose.Escape() // close the list
		}
	SetFocus() // called by LibView
		{
		.AutoChoose.SetFocus()
		}
	FieldEscape()
		{
		.Send('LocateEscape')
		}

	Destroy()
		{
		.sub.Unsubscribe()
		super.Destroy()
		}
	}