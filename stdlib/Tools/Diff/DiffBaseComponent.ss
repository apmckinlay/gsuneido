// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
HtmlDivComponent
	{
	Name: "Diff"
	New(@args)
		{
		super(@args)
		.ovbar = .FindControl(#OverviewBar)
		}

	SetProcs(@listIds)
		{
		.lists = listIds.Map({ SuRender().GetRegisteredComponent(it) })
		for list in .lists
			{
			list.AddEventListenerToCM('scroll', .factory(#OnScroll, list))
			list.AddEventListenerToCM('mousedown', .factory(#OnMouseDown, list))
			list.AddEventListenerToCM('keydown', .factory(#OnKeyDown, list))
			}
		}

	factory(method, list)
		{
		return { |@args| (this[method])(list, :args) }
		}

	OnScroll(source)
		{
		scrollInfo = source.CM.GetScrollInfo()
		for list in .lists
			{
			if source.UniqueId is list.UniqueId
				continue
			list.CM.ScrollTo(scrollInfo.left, scrollInfo.top)
			}
		//TODO .ovbar
		}

	delay: false
	OnMouseDown(source)
		{
		if .delay isnt false
			{
			.delay.Kill()
			.delay = false
			}
		.delay = SuDelayed(0)
			{
			.Event('SyncSelectedLineByClick',
				.lists.FindIf({ it.UniqueId is source.UniqueId }))
			.delay = false
			}
		}

	OnKeyDown(source, args)
		{
		event = args[1]
		if event.key not in (#ArrowUp, #ArrowDown, #ArrowLeft, #ArrowRight)
			return
		if .delay isnt false
			{
			.delay.Kill()
			.delay = false
			}
		.delay = SuDelayed(0)
			{
			.Event('SyncSelectedLineByKeyDown',
				.lists.FindIf({ it.UniqueId is source.UniqueId }))
			.delay = false
			}
		}

	Destroy()
		{
		if .delay isnt false
			{
			.delay.Kill()
			.delay = false
			}
		super.Destroy()
		}
	}