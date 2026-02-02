// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
WindowBase
	{
	ComponentName: 'Window'
	CallClass(@args)
		{
		reservation = SuRenderBackend().ReserveAction()
		window = new this(@args)
		layout = window.GetLayout()
		args[0] = layout
		args.uniqueId = window.UniqueId
		args.title = window.GetWindowTitle()
		args.menubar = window.GetMenubar()
		args.Delete(#onDestroy)
		excludeModalWindow = args.GetDefault(#excludeModalWindow, false)
		args.Delete(#excludeModalWindow)
		SuRenderBackend().RecordAction(false,
			window.GetWindowComponentName(excludeModalWindow), args, reservation.at)
		return window
		}

	keep_placement: false
	onDestroy: false
	New(control, title = false, .keep_placement = false, .onDestroy = false,
		useDefaultSize = false, skipStartup? = false)
		{
		.DoActivate()
		.SetupCommands(control)
		if control isnt false
			.Ctrl = .Construct(control)

		.menutips = Object().Set_default("")
		.dynmenu = Object()
		.curdynmenu = false
		.dynid = Object().Set_default(Object())
		.createmenu()

		.title = .FindTitle(title, control)

		if .keep_placement is false or .restoreWindowPlacement() is false
			.showWindow(useDefaultSize)

		if not skipStartup?
			DoStartup(.Ctrl)

		if .Ctrl.DefaultButton isnt "" and
			false isnt (ctrl = .FindControl(ToIdentifier(.Ctrl.DefaultButton))) and
			ctrl.Base?(ButtonControl)
			.Act(#SetDefaultButton, ctrl.UniqueId)
		}

	menubar: #()
	menubarMap: #()
	createmenu()
		{
		if (not .Ctrl.Member?("Menu"))
			return
		menus = .Ctrl.Val_or_func(#Menu).Copy()
		.addPluginMenus(menus)

		if menus.Empty?()
			return

		.menubar = Object()
		.menubarMap = Object()
		for menu in menus
			{
			name = menu[0].Tr('&')
			.menubar.Add(name)
			.menubarMap[name] = Object()
			.makesubmenu(menu, .menubarMap[name])
			}
		}
	GetMenubar()
		{
		return .menubar
		}

	prevMenuFocus: false
	SyncPrevMenuFocus(.prevMenuFocus) { }

	MenuBar(menu, x, y)
		{
		id = ContextMenu(.menubarMap[menu]).Show(false, x, y, parent: this)
		if .prevMenuFocus isnt false
			// Need to force setting the focus to the browser side because
			// the context menu doesn't change the focus in SuRenderBackend().Status()
			SetFocus(.prevMenuFocus, force:)
		if 0 isnt id
			{
			.COMMAND(id)
			}
		.curdynmenu = false
		}

	LoadDynamic(id)
		{
		if (not .dynmenu.Member?(id) or
			0 is menu = .Send('Menu_' $ .dynmenu[id]))
			return #()
		newMenu = Object()
		for (m in menu)
			{
			if false isnt entry = .makemenuentry(m, newMenu)
				{
				.dynid[id][entry.id] = Object('On_' $ .dynmenu[id], m)
				}
			}
		return newMenu
		}

	UpdateCurDynamicMenu(.curdynmenu) { }

	addPluginMenus(menus)
		{
		attached = Object().Set_default(false)
		Plugins().ForeachContribution("UI", 'attach')
			{|c|
			if c.to is .Ctrl.Title
				attached[c.menu] = true
			}
		map = Object()
		Plugins().ForeachContribution("UI", 'action')
			{|c|
			if attached[c.menu] is false
				continue
			if not map.Member?(c.menu)
				menus.Add(map[c.menu] = Object(c.menu))
			map[c.menu].Add(c.name)
			if c.Member?(#target)
				.menu_redir(c.name, c.target)
			if c.Member?(#accel)
				.AddAccel(.menuify(c.name), c.accel)
			}
		.SetAccels()
		order = Object()
		Plugins().ForeachContribution("UI", 'menu')
			{|c|
			order[c.menu] = c.order
			}
		defaultOrder = 5
		menus.Sort!({|x, y|
			order.GetDefault(x[0], defaultOrder) < order.GetDefault(y[0], defaultOrder) })
		}

	menu_redir(name, target)
		{
		if String?(target)
			target = Global(target)
		if Object?(name)
			.Ctrl.Redir("Menu_" $ .menuify(name = name[0]), target)
		.Ctrl.Redir("On_" $ .menuify(name), target)
		}

	menuify(name)
		{
		return ToIdentifier(name)
		}

	makesubmenu(menu, pulldown)
		{
		for (i = 1; i < menu.Size(); ++i)
			{
			m = menu[i]
			if String?(m)
				.makemenuentry(m, pulldown)
			else if Object?(m)
				{
				p = Object()
				entry = .makemenuentry(m[0], pulldown)
				if (m.Size() is 1)
					{
					.dynmenu[entry.id] = .menuify(m[0])
					entry.dynamic = true
					}
				.makesubmenu(m, p)
				pulldown.Add(p)
				}
			}
		}

	makemenuentry(m, pulldown)
		{
		if (m is '')
			{
			pulldown.Add('')
			return false
			}

		cmd = .menuify(m)
		c = .Commands().GetDefault(cmd, { Object(id: .Mapcmd(cmd), accel: "", help: "") })
		m = TranslateLanguage(m)
		accel = TranslateLanguage(c.accel)
		if m is '&Copy' and accel is ''
			accel = 'Ctrl+C'
		entry = Object(name: m $ '\t' $ accel, id: c.id)
		pulldown.Add(entry)
		.menutips['On_' $ cmd] = TranslateLanguage(c.help)
		return entry
		}

	COMMAND(id)
		{
		if (.dynid.Member?(.curdynmenu) and .dynid[.curdynmenu].Member?(id))
			{
			.Send(@.dynid[.curdynmenu][id])
			return 0
			}
		else
			return super.COMMAND(id)
		}

	focus: 0
	ACTIVATE(active)
		{
		if active is false
			{
			.focus = GetFocus()
			.Send("Inactivate")
			PubSub.Publish(#WindowInactivated)
			}
		else
			{
			if .focus not in (false, 0)
				SetFocus(.focus)

			.Send("Activate")
			PubSub.Publish(#WindowActivated)
			}
		return 0
		}

	GetWindowTitle()
		{
		return .title
		}

	SetTitle(text)
		{
		.title = text
		SuRenderBackend().WindowManager.UpdateTaskbar()
		super.SetTitle(text)
		}

	place: false
	SYNCWINDOWPLACEMENT(rect, maximized, minimized)
		{
		super.WINDOWPOSCHANGING()
		.place = Object(showCmd: maximized
			? SW.SHOWMAXIMIZED
			: minimized
				? SW.SHOWMINIMIZED
				: SW.SHOWNORMAL,
			rcNormalPosition: rect)
		}

	IsMinimized?()
		{
		return .place isnt false and .place.showCmd is SW.SHOWMINIMIZED
		}

	// called by GetWindowPlacement
	GetWindowPlacement(wpPlace)
		{
		if .place isnt false
			wpPlace.Merge(.place)
		else
			{
			wpPlace.showCmd = SW.SHOWNORMAL
			wpPlace.rcNormalPosition = Object(top: 0, left: 0,
				width: .DefaultSize().w, height: .DefaultSize().h,
				right: .DefaultSize().w, bottom: .DefaultSize().h)
			}
		}

	// called by SetWindowPlacement
	SetWindowPlacement(wpPlace)
		{
		.place = wpPlace.Project(#showCmd, #rcNormalPosition)
		if not .place.rcNormalPosition.Member?(#width)
			.place.rcNormalPosition.width =
				.place.rcNormalPosition.right - .place.rcNormalPosition.left
		if not .place.rcNormalPosition.Member?(#height)
			.place.rcNormalPosition.height =
				.place.rcNormalPosition.bottom - .place.rcNormalPosition.top
		.Act(#SetWindowPlacement, .place)
		}

	getter_windowPlacementSaveName()
		{
		saveName = 'wp: '
		if String?(.keep_placement) and .keep_placement isnt ''
			saveName $= .keep_placement
		else
			saveName $= .title
		return saveName
		}

	restoreWindowPlacement()
		{
		if false is info = KeyListViewInfo.Get(
			.windowPlacementSaveName $ .keyListViewInfoSuffix)
			return false
		if info.window_info.showCmd is SW.SHOWMINIMIZED
			info.window_info.showCmd = SW.SHOWNORMAL
		place = info.window_info
		if not place.rcNormalPosition.Member?(#width)
			place.rcNormalPosition.width =
				place.rcNormalPosition.right - place.rcNormalPosition.left
		if not place.rcNormalPosition.Member?(#height)
			place.rcNormalPosition.height =
				place.rcNormalPosition.bottom - place.rcNormalPosition.top
		.Act(#SetWindowPlacement, info.window_info)
		return true
		}

	showWindow(useDefaultSize)
		{
		if useDefaultSize
			{
			r = Object(left: 0, top: 0,
				width: .DefaultSize().w,
				height: .DefaultSize().h)
			.Act(#SetWindowPlacement, Object(showCmd: SW.SHOWNORMAL, rcNormalPosition: r))
			}
		}

	// not share position with exe client, which may have multiple monitors
	keyListViewInfoSuffix: ' - web'
	saveWindowPlacement()
		{
		// program errors while creating print previews causes title to be uninitialized,
		// could also be other cases
		if not .Member?("Window_title") or .place is false
			return
		KeyListViewInfo.Save(.windowPlacementSaveName $ .keyListViewInfoSuffix, .place)
		}

	Show(unused)
		{
		}

	ResetStyle()
		{
		.Ctrl.Send('ResetAddons')
		.Ctrl.TopDown('ResetAddons')
		.Ctrl.Send('ResetTheme')
		.Ctrl.TopDown('ResetTheme')
		}

	HeaderContextMenu(status, extra, x, y)
		{
		menu = Object()
		menu.Add(Object(name: 'Restore', state: .getState(status.restore)))

		menu.Add(Object(name: 'Maximize', state: .getState(status.maximize)))
		menu.Add('Set Size',
			#('640 x 480',
			'800 x 600',
			'1024 x 768',
			'1280 x 1024',
			'1600 x 900',
			'1920 x 1080',
			'2560 x 1440'))
		menu.Add(@extra)
		.windowMenuOptions(menu)

		menu.Add('', Object(name: 'Close', state: .getState(status.close)))
		if menu.Empty?()
			return
		ContextMenu(menu).ShowCallCascade(this, x, y)
		}

	windowMenuOptions(menu)
		{
		if .windowMenus is false
			return
		menu.Add('')
		menu.Add(@.windowMenus)
		}

	windowMenus: false
	AddWindowMenuOptions(menus)
		{
		.windowMenus = Object()
		for menu in menus
			{
			prefix = 'On_Context_' $ menu.root
			.windowMenus.Add(menu.root)
			if menu.GetDefault('options', false) is false
				.addWindowMenuCmd(prefix, menu.cmd)
			else
				{
				subMenu = Object()
				for option, callable in menu.options
					{
					.addWindowMenuCmd(prefix $ '_' $ option, callable)
					subMenu.Add(option)
					}
				.windowMenus.Add(subMenu)
				}
			}
		}

	addWindowMenuCmd(name, callable)
		{
		.windowMenuCommands[ToIdentifier(name)] = callable
		}

	getter_windowMenuCommands()
		{
		return .windowMenuCommands = Object()
		}

	On_Context_Maximize()
		{
		.Act('MAXIMIZE')
		}

	On_Context_Restore()
		{
		.Act('MAXIMIZE', forceRestore?:)
		}

	On_Context_Set_Size(item)
		{
		size = item[0]
		width = Number(size.BeforeFirst(' '))
		height = Number(size.AfterLast(' '))
		top = .place is false ? 0 : .place.rcNormalPosition.top
		left = .place is false ? 0 : .place.rcNormalPosition.left
		place = Object(
			showCmd: SW.SHOWNORMAL,
			rcNormalPosition: Object(:top, :left, :width, :height
				right: left + width, bottom: top + height))
		.WINDOWRESIZE(place.rcNormalPosition)
		.SetWindowPlacement(place)
		}

	On_Context_Close()
		{
		.CLOSE()
		}

	EscapeCancel()
		{
		if IsWindowEnabled(.Hwnd)
			.Send('On_Cancel')
		}

	Default(@args)
		{
		if not args[0].Prefix?('On_Context')
			return
		if .windowMenuCommands.Member?(args[0])
			(.windowMenuCommands[args[0]])()
		else
			.Act('ExtraContextCall', args)
		}

	getState(enabled?)
		{
		return enabled? ? 0 : 3
		}

	AfterMinimized()
		{
		SuRenderBackend().WindowManager.ActivateNextWindow(this)
		}

	DESTROY(_windowResult = false)
		{
		if .keep_placement isnt false
			.saveWindowPlacement()
		if .onDestroy isnt false
			(.onDestroy)(:windowResult)
		super.DESTROY()
		}
	}
