// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	ComponentName: 'ContextMenuList'
	New(menu, left = false, top = false, rcExclude = 0, buttonRect = #(),
		.parent = false)
		{
		.ComponentArgs = Object(menu, left, top, rcExclude, buttonRect)
		}

	CLICKED(id)
		{
		.Window.Result(id)
		}

	ContextMenuList_Load(id, expandId)
		{
		menu = .parent is false
			? #()
			: .parent.LoadDynamic(id)
		.Act('LoadDynamicMenu', expandId, menu)
		}

	ContextMenuList_UpdateCurDynamicMenu(id)
		{
		if .parent is false
			return
		.parent.UpdateCurDynamicMenu(id)
		}
	}