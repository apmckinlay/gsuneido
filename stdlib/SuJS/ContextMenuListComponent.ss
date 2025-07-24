// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
// TODO: handle hotkey
Component
	{
	Xstretch: 0
	Ystretch: 0
	styles: `
		.su-contextmenu-list {
			background-color: white;
			position: fixed;
			border: 1px solid black;
			border-spacing: 0px;
			border-radius: 0.3em;
			user-select: none;
			box-shadow: 3px 3px 3px grey;
		}
		.su-contextmenu-list hr {
			margin-block-start: 0px;
			margin-block-end: 0px;
		}
		.su-contextmenu-item .su-contextmenu-label{
			min-width: 10em;
		}
		.su-contextmenu-item .su-contextmenu-shortcut{
			text-align: right;
		}
		.su-contextmenu-item.su-selected {
			background-color: lightblue;
		}
		.su-contextmenu-item td {
			padding: 5px;
		}
		.su-contextmenu-item.su-disabled {
			color: darkgray;
		}
		`
	New(.menu, left = false, top = false, rcExclude = 0, buttonRect = #())
		{
		LoadCssStyles('context-menu.css', .styles)

		.outstandingExpands = Object()
		.outstandingExpandId = 0

		.CreateElement('table', className: 'su-contextmenu-list')
		.buildMenu(.menu, .El, left, top, rcExclude, buttonRect)
		.menus = Object(Object(submenu: menu, selected: false))
		.keydownCB = SuUI.GetCurrentDocument().AddEventListener('keydown', .keydown)
		}

	buildMenu(menu, container, left = false, top = false, rcExclude = 0, buttonRect = #())
		{
		for item in menu
			{
			if item is ''
				{
				hr = CreateElement('tr', container)
				hr.innerHTML = '<td></td><td><hr></td><td></td><td></td>'
				continue
				}

			item.el = CreateElement('tr', container, className: 'su-contextmenu-item')

			statusEl = CreateElement('td', item.el)
			nameEl = CreateElement('td', item.el, className: 'su-contextmenu-label')
			nameEl.innerText = item.name.BeforeFirst('\t').Tr('&')

			shortcutEl = CreateElement('td', item.el,
				className: 'su-contextmenu-shortcut')
			shortcutEl.innerText = item.name.AfterFirst('\t')

			submenuIndicatorEl= CreateElement('td', item.el)
			if item.Member?(#type)
				statusEl.innerHTML = '&#8226;'
			if item.Member?(#submenu)
				submenuIndicatorEl.innerHTML = '&gt'
			if item.GetDefault(#disable, false) is true
				item.el.classList.Add('su-disabled')
			item.el.AddEventListener('mouseenter', .listenerFactory(.mouseenter, item))
			item.el.AddEventListener('mouseleave', .listenerFactory(.mouseleave, item))
			item.el.AddEventListener('click', .listenerFactory(.mouseclick, item))
			}

		PlaceElement(container, left, top, rcExclude, buttonRect)
		}

	mouseenter(item)
		{
		.removeMenusUntil(item)
		.addSelect(item)
		.expandMenu(item)
		}

	removeMenusUntil(item)
		{
		for (i = .menus.Size() - 1; i >= 0; i--)
			{
			.removeSelect(.menus[i])
			if .menus[i].submenu.Has?(item)
				break
			.menus[i].submenuEl.SetStyle('display', 'none')
			.menus.Delete(i)
			}
		}

	outstandingExpands: #()
	expandMenu(item, selectFirst? = false)
		{
		if item is false or
			item.GetDefault(#disable, false) is true or
			not item.Member?(#submenu)
			return

		if item.Member?(#submenuEl)
			item.submenuEl.SetStyle('display', '')
		else
			{
			if .loadingDynamicMenu(item, selectFirst?)
				return
			r = SuRender.GetClientRect(item.el)
			item.submenuEl = CreateElement('table', .ParentEl,
				className: 'su-contextmenu-list')
			.buildMenu(item.submenu, item.submenuEl, r.right, r.top,
				Object(top: -9999, bottom: 9999, left: r.left, right: r.right), r)
			}

		if item.GetDefault(#dynamic, false) is true
			.Event('ContextMenuList_UpdateCurDynamicMenu', item.id)

		.menus.Add(item)
		if selectFirst? is true
			.move(1)
		}

	loadingDynamicMenu(item, selectFirst?)
		{
		if item.submenu is #() and item.GetDefault(#dynamic, false) is true
			{
			if not .outstandingExpands.Any?({ it.item is item })
				{
				.outstandingExpands[.outstandingExpandId] = Object(:item, :selectFirst?)
				.EventWithFreeze('ContextMenuList_Load', item.id, .outstandingExpandId++)
				}
			return true
			}
		return false
		}

	LoadDynamicMenu(expandId, menu)
		{
		if false is task = .outstandingExpands.GetDefault(expandId, false)
			return
		task.item.submenu = menu
		if .menus.Last().GetInit(#selected, false) is task.item
			.expandMenu(@task)
		.outstandingExpands.Delete(expandId)
		}

	mouseleave(item)
		{
		if .Empty?()
			return
		if .menus.Last().GetInit(#selected, false) is item
			.removeSelect(.menus.Last())
		}

	removeSelect(menu)
		{
		if menu.GetInit(#selected, false) is false
			return
		menu.selected.el.classList.Remove('su-selected')
		menu.selected = false
		}

	addSelect(item)
		{
		item.el.classList.Add('su-selected')
		.menus.Last().selected = item
		}

	mouseclick(item)
		{
		if item.GetDefault(#disable, false) is true or not .menus.Last().Has?(item)
			return
		if item.GetDefault(#submenu, #()).NotEmpty?()
			.mouseenter(item)
		else
			.EventWithOverlay(#CLICKED, item.id)
		}

	keydown(event)
		{
		switch (event.key)
			{
		case 'ArrowUp':
			.move(-1)
		case 'ArrowDown':
			.move(1)
		case 'ArrowRight':
			.expandMenu(.menus.Last().GetInit(#selected, false), selectFirst?:)
		case 'ArrowLeft':
			.fold()
		case 'Escape':
			.cancel()
		case 'Enter':
			.enter()
		default:
			}
		event.PreventDefault()
		event.StopPropagation()
		}

	move(offset)
		{
		menu = .menus.Last()
		next = offset > 0 ? 0 : menu.submenu.Size() - 1
		if menu.GetInit(#selected, false) isnt false
			{
			cur = menu.submenu.Find(menu.selected)
			.removeSelect(menu)
			next = (cur + offset + menu.submenu.Size()) % menu.submenu.Size()
			}
		while (menu.submenu[next] is '')
			{
			next = (next + offset + menu.submenu.Size()) % menu.submenu.Size()
			}
		.addSelect(menu.submenu[next])
		}

	fold()
		{
		if .menus.Size() < 2
			return

		prevSelect = .menus[.menus.Size() - 2].selected
		.removeMenusUntil(prevSelect)
		.addSelect(prevSelect)
		}

	cancel()
		{
		.Event(#On_Cancel)
		}

	enter()
		{
		item = .menus.Last().GetInit(#selected, false)
		if item is false
			.cancel()
		else if item.Member?(#submenu)
			.expandMenu(item, selectFirst?:)
		else
			.mouseclick(.menus.Last().selected)
		}

	listenerFactory(@args)
		{
		fn = args[0]
		return { fn(@+1args) }
		}

	Destroy()
		{
		if .keydownCB isnt false
			SuUI.GetCurrentDocument().RemoveEventListener('keydown', .keydownCB)
		super.Destroy()
		}
	}
