// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xstretch: 1
	Ystretch: 1
	Xmin: 200
	Ymin: 100
	Name: "TwoList"
	New(list = #(), initial_list = #(), mandatory_list = #(), .noSort = false,
		.delimiter = ',')
		{
		.list1 = .Horz.Vert1.list1
		.list2 = .Horz.Vert2.list2
		// note: items in initial list should be in list
		initial_list.MergeUnion(mandatory_list)
		for item in initial_list
			.list2.AddItem(item)
		for item in list
			if not initial_list.Has?(item)
				.list1.AddItem(item)
		.newlist = initial_list.Copy()
		.mandatory_list = mandatory_list
		.Send("Data")
		}

	Controls()
		{
		return Object('Horz'
			Object('Vert'
			#(Static " Available")
			#(Skip 2)
			Object('ListBox' name: "list1", sort: not .noSort) name: 'Vert1')
			#Skip,
			.AddRemoveButtons()
			#Skip
			#(Vert
			(Static " Selected")
			(Skip 2)
			(ListBox name: "list2") name: 'Vert2')
			#Skip,
			.MoveButtons()
			#Skip
		)
		}

	MoveButtons()
		{
		return #(Vert
			Fill
			(Button, "Move Up" xstretch: 0,
				tip: "Move selected item up in the list")
			Fill
			(Button "Move Down" xstretch: 0,
				tip: "Move selected item down in the list")
			Fill)
		}

	AddRemoveButtons()
		{
		return #(Vert
			Fill (Button ">" "Move" xstretch: 0,
				tip: "Add selected item")
			Fill (Button "<" "MoveBack" xstretch: 0,
				tip: "Remove selected item")
			Fill (Button ">>" "All" xstretch: 0,
				tip: "Add all items")
			Fill (Button "<<" "AllBack" xstretch: 0,
				tip: "Remove all items")
			Fill)
		}

	ListBoxDoubleClick(item, source)
		{
		if item is -1
			return

		target = .list1 is source ? .list2 : .list1
		data = source.GetText(item)
		if target is .list1 and .mandatory_list.Has?(data)
			{
			AlertError("Cannot move mandatory item.", .Window.Hwnd)
			return
			}

		.addItem(target, data)
		if target is .list2
			.newlist.Add(data)
		else
			.newlist.Delete(item)
		source.DeleteItem(item)
		.Send("NewValue", .Get())
		}

	On_All()
		{
		for (i = 0; (item = .list1.GetText(i)) isnt ''; ++i)
			{
			.list2.AddItem(item)
			.newlist.Add(item)
			}
		.list1.DeleteAll()
		.Send("NewValue", .Get())
		}

	On_Move()
		{
		if ((sel = .list1.GetCurSel()) is -1)
			return
		data = .list1.GetText(sel)
		.addItem(.list2, data)
		.newlist.Add(data)
		.list1.DeleteItem(sel)
		.Send("NewValue", .Get())
		}

	On_MoveBack()
		{
		if ((sel = .list2.GetCurSel()) is -1)
			return
		data = .list2.GetText(sel)
		if .mandatory_list.Has?(data)
			{
			AlertError("Cannot move mandatory item " $ data, .Window.Hwnd)
			return
			}
		.list1.AddItem(data)
		.newlist.Delete(sel)
		.list2.DeleteItem(sel)
		.Send("NewValue", .Get())
		}

	addItem(list, item)
		{
		list.AddItem(item)
		if list is .list2
			list.SetCurSel(list.Count() - 1)
		}

	On_AllBack()
		{
		.AllBack()
		.Send("NewValue", .Get())
		}

	AllBack()
		{
		for (i = 0; (item = .list2.GetText(i)) isnt ''; ++i)
			if not .mandatory_list.Has?(item)
				.list1.AddItem(item)
		.newlist.Delete(all:)
		.list2.DeleteAll()
		if not .mandatory_list.Empty?()
			{
			.newlist = .mandatory_list.Copy()
			for item in .mandatory_list
				.list2.AddItem(item)
			AlertError("Cannot move mandatory items.", .Window.Hwnd)
			}
		}

	On_Move_Up()
		{
		if ((sel = .list2.GetCurSel()) is -1 or sel is 0)
			return

		.newlist.Swap(sel, sel - 1)

		// swap items in listbox
		.list2.DeleteItem(sel)
		.list2.DeleteItem(sel - 1)
		.list2.InsertItem(.newlist[sel - 1], sel - 1)
		.list2.InsertItem(.newlist[sel], sel)
		.list2.SetCurSel(sel - 1)
		}

	On_Move_Down()
		{
		if ((sel = .list2.GetCurSel()) is -1 or sel is .newlist.Size() - 1)
			return

		.newlist.Swap(sel, sel + 1)

		// swap items in listbox
		.list2.DeleteItem(sel + 1)
		.list2.DeleteItem(sel)
		.list2.InsertItem(.newlist[sel], sel)
		.list2.InsertItem(.newlist[sel + 1], sel + 1)
		.list2.SetCurSel(sel + 1)
		}

	Get()
		{
		return .newlist.Join(.delimiter)
		}

	Set(list_str = '')
		{
		list = list_str.Split(.delimiter)
		.newlist = list
		.list2.DeleteAll()
		for item in list
			{
			.list2.AddItem(item)
			idx = .list1.FindString(item)
			if idx >= 0
				.list1.DeleteItem(idx)
			}
		}

	SetList(list_str = '')
		{
		for item in list_str.Split(.delimiter)
			.list1.AddItem(item)
		.newlist = Object()
		.list2.DeleteAll()
		}

	GetNewList()
		{
		return .newlist
		}

	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}
