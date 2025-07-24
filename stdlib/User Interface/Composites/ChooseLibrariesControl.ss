// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	Name: "ChooseLibraries"
	Title: 'Choose Libraries'

	New(.all = false)
		{
		super(.layout())
		.Set(.all ? '(All)' : "Contrib")
		.Send('NewValue', .Get())
		}

	layout()
		{
		.list = Libraries().MergeUnion(LibraryTables())
		if .all
			.list.Add('(All)', '(In Use)')
		return Object("AutoFillField", candidates: .list,
			font: '@mono', size: '+1', width: 30,)
		}

	lastSet: ''
	Getter_DialogControl()
		{
		.Set(.Get())
		additionalButtons = Object(Object(text: 'In Use', command: { Libraries() }))
		return Object(ChooseManyListControl, .list, :additionalButtons, value: .lastSet)
		}

	processSelection(libs)
		{
		if .all and (libs is '' or libs.Split(',').Size() is .list.Size())
			libs = '(All)'
		return .lastSet = libs
		}

	Get()
		{
		return .processSelection(super.Get())
		}

	Set(libs)
		{
		super.Set(.processSelection(libs))
		}
	}
