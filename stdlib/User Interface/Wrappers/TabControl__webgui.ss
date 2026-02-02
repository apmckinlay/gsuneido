// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 'Tab'
	ComponentName: 'Tab'
	Xstretch: 1

	New(@tabs)
		{
		.tabs = Object()

		// add the tabs
		for tab in tabs.Values(list:)
			.Insert(.Count(), tab)

		if .Count() > 0
			.Select(0)

		.extraControl = false
		if false isnt extra = tabs.GetDefault(#extraControl, false)
			.extraControl = .Construct(extra)

		.ComponentArgs = Object(
			close_button: tabs.GetDefault(#close_button, false),
			orientation: tabs.GetDefault(#orientation, #top),
			staticTabs: tabs.GetDefault(#staticTabs, #()),
			extraControl: .extraControl is false ? false : .extraControl.GetLayout())

		.addTabButton(tabs)
		}

	defaultTip: 'Add Tab'
	addTabButton(args)
		{
		buttonTip = args.GetDefault(#buttonTip, .defaultTip)
		if args.GetDefault(#addTabButton?, false) or buttonTip isnt .defaultTip
			.ComponentArgs.tabButton = buttonTip
		}

	ButtonClicked(cmd, pt)
		{
		cmd = "On_" $ ToIdentifier(cmd)
		if .Method?(cmd)
			(this[cmd])(:pt)
		else
			.Send(cmd, :pt)
		}

	On_Go_to_Tab(pt)
		{
		list = .tabs.Map({ it.text })
		if 0 isnt i = ContextMenu(list).Show(.Window.Hwnd, pt.x, pt.y)
			.goto(i - 1)
		}

	ContextMenu(x, y, hover)
		{
		return hover is false ? 0 : .Send("TabContextMenu", x, y, hover)
		}

	tabId: 0
	nextId()
		{
		return .tabId++
		}

	Insert(i, text, data = #(tooltip: ""), image = -1)
		{
		id = .nextId()
		image = .images.Member?(image) ? .images[image] : -1
		.Act('Insert', i, text, data, image, id)
		.tabs.Add(Object(:id, :text, :image, :data), at: i)
		}

	Count()
		{
		return .tabs.Size()
		}

	Remove(i)
		{
		if i < .selected or .selected is .tabs.Size() - 1
			.selected--
		.Act('Remove', i)
		.tabs.Delete(i)
		return true
		}

	selected: -1
	Select(i)
		{
		.Act('Select', .selected = i)
		}

	Move(i, newPos)
		{
		data = .GetData(i)
		text = .GetText(i)
		.Remove(i)
		.Insert(newPos, text, :data, image: data.image)
		.Select(newPos)
//		if newPos is 0
//			.ensureFirstTabVisible()
		}

	GetSelected()
		{
		return .selected
		}

	SetText(i, text)
		{
		.Act('SetText', i, text)
		.tabs[i].text = text
		}

	GetText(i)
		{
		return .tabs[i].text
		}

	images: #()
	SetImageList(.images) {}

	SetImage(i, img)
		{
		image = .images.Member?(img) ? .images[img] : -1
		.Act('SetImage', i, image)
		.tabs[i].image = image
		}

	GetImage(i)
		{
		return .tabs[i].image
		}

	GetData(i)
		{
		return .tabs[i].data
		}

	SetData(i, data)
		{
		.tabs[i].data = data

		}

	ForEachTab(block)
		{
		i = 0
		.tabs.Each({ block(it.data, idx: i++) })
		}

	Tab_Close(i)
		{
		.Send('Tab_Close', i)
		}

	Click(i)
		{
		.goto(i)
		}

	goto(clicked)
		{
		if clicked is false or true is .Send('TabControl_SelChanging')
			return 0

		if clicked is .selected
			.Send('TabClick', clicked)
		else
			{
			.Act('DoTabChange', .selected = clicked, true)
			.Send('SelectTab', clicked)
			}
		return 0
		}

	Destroy()
		{
		if .extraControl isnt false
			.extraControl.Destroy()
		super.Destroy()
		}
	}
