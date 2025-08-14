// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Ymin: 150
	Name: 'ChooseManyList'
	New(list, list_desc = #(), value = "", nobuttons = false,
		text = "", name = "", .additionalButtons = #())
		{
		super(.controls(nobuttons is false, text))
		.setLocalValues(list, list_desc, value, name)
		.setColumns(list_desc)
		if list_desc.NotEmpty?()
			.listview.SetMaxWidth('value')
		}

	On_Default(@args)
		{
		name = args.source.Name.Replace('_', ' ')
		if false isnt button = .additionalButtons.FindOne({ name is it.text	})
			{
			.listview.CheckAll(false)
			if button.Member?(#command) and Object?(items = (button.command)())
				.SetChecked(items)
			}
		}

	controls(buttons?, text = '')
		{
		controls = Object('Vert',
			Object('ListView', menu: false, ystretch: 1, stretch:,
				noHeader: /*suneido.js*/))
		if text isnt ''
			controls.Add(#(Skip 3), Object('Static' text))
		if buttons?
			{
			buttonLayout = Object(#HorzEqual)
			.additionalButtons.Each(
				{
				buttonLayout.Add(
					Object(#Button text: it.text, command: #Default), #Skip)
				})
			buttonLayout.Add(#AllNoneOkCancel)
			controls.Add(#(Skip 5), buttonLayout)
			}
		return controls
		}

	setLocalValues(list, .list_desc, value, .name)
		{
		.list = list.Copy()
		.value = value.Has?('(All)') ? .list.Join(',') : value
		.listview = .Vert.List
		.listview.SetStyle(LVS.REPORT | LVS.NOCOLUMNHEADER | LVS.SINGLESEL |
			LVS.SHOWSELALWAYS)
		.listview.SetExtendedStyle(LVS_EX.FULLROWSELECT | LVS_EX.CHECKBOXES)
		}

	setColumns(list_desc)
		{
		.listview.AddColumn('value')
		if not list_desc.Empty?()
			.listview.AddColumn('desc')

		values = .value.Split(',')
		for i in values.Members()
			values[i] = values[i].Trim()
		for m in .list.Members()
			{
			item = .list[m]
			if Number?(item)
				item = Display(item)
			else if String?(item)
				.list[.list.Find(item)] = item = item.Trim()
			if not list_desc.Empty?()
				i = .listview.Addrow(Object(value: item,
					desc: list_desc.Member?(m) ? list_desc[m] : ''))
			else
				i = .listview.AddItem(item)
			.listview.SetCheckState(i, values.Has?(String(item)))
			}
		}

	AddItem(s)
		{
		.listview.AddItem(s)
		.list.Add(s)
		}
	DeleteValue(value)
		{
		i = .list.Find(value)
		.listview.DeleteItem(i)
		.list.Delete(i)
		}
	DeleteAll()
		{
		.listview.DeleteAll()
		.list = Object()
		}
	GetSelected()
		{ return .listview.GetSelected() }
	SelectItem(i)
		{
		.listview.SelectItem(i)
		}
	UnSelect()
		{
		.listview.UnSelect()
		}
	On_All()
		{
		.listview.CheckAll()
		}
	On_None()
		{
		.listview.CheckAll(false)
		}
	On_Cancel()
		{
		.Window.Result(.value)
		}
	SetMenu(menu)
		{
		.listview.SetMenu(menu)
		}
	GetChecked()
		{
		checked = Object()
		for i in .list.Members()
			if .listview.GetCheckState(i)
				checked.Add(.list[i])
		return checked
		}
	SetChecked(items)
		{
		for name in items
			if -1 isnt i = .list.Find(name)
				.listview.SetCheckState(i, true)
		}
	GetList()
		{
		return .list
		}
	OK()
		{
		return .GetChecked().Join(',')
		}
	SelectChanged(olditem, newitem) // pass on from ListView
		{
		.Send('SelectChanged', olditem, newitem)
		}
	}
