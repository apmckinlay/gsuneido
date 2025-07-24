// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(.multiLoc = false, .onlyLatLong? = false)
		{
		.get_map_sources()
		.get_map_choice()
		}

	Controls()
		{
		return Object('EnhancedButton', image: 'location', command: 'Map',
			tip: 'Right click to choose map source', alignTop:,
			size: '-5', imagePadding: .1, imageColor: CLR.Inactive,
			mouseOverImageColor: CLR.Highlight)
		}

	Recalc()
		{
		super.Recalc()
		.Top = .GetChild().Top
		}

	On_Map(option /*unused*/ = "")
		{
		.which = .validateWhich(.funcs)
		BookLog('Map: ' $ .menu[.which])
		if .multiLoc is false
			.mapSingle()
		else
			.mapMulti()
		}

	mapSingle()
		{
		if false is addr_ob = .Send("Map_GetAddress")
			return
		(.funcs[.which])(@addr_ob)
		}

	mapMulti()
		{
		if false is addr_ob = .getLocations()
			return

		if addr_ob.locations.RemoveIf({ String?(it) and it.Blank?() }).Size() is 1 and
			not String?(addr_ob.locations[0])
			{
			if false isnt addr_ob = addr_ob.locations[0]
				(.funcs[.which])(@addr_ob)
			return
			}

		result = (.funcs[.which].HandleMultiLocations)(addr_ob)
		if result isnt ''
			.AlertWarn('Map Locations', result)
		}

	getLocations()
		{
		addr_ob = Object(locations: Object(), options: Object(), hwnd: 0)
		multiOb = .Send("Map_GetMultiLocations")
		if multiOb is 0 or multiOb.Empty?() or multiOb.locations.Empty?()
			return false

		addr_ob.Merge(multiOb)
		return addr_ob
		}

	menu: false
	funcs: false
	which: 0
	orig_which: 0
	On_Map_ContextMenu(x, y, source)
		{
		menu = .menu.Copy()
		.which = .validateWhich(menu)
		menu[.which] = Object(name: menu[.which], state: MF.CHECKED, type: MFT.RADIOCHECK)
		i = ContextMenu(menu).Show(source.Hwnd, x, y) - 1
		if menu.Member?(i)
			.which = i
		}

	PluginName: 'Map'
	get_map_sources()
		{
		.menu = Object()
		.funcs = Object()
		Plugins().ForeachContribution(.PluginName, 'mapSource')
			{ |c|
			if .allowToAddMenu?(c)
				.addToMenu(c)
			}
		}
	allowToAddMenu?(c)
		{
		if .multiLoc
			return c.GetDefault('multiLoc', false)
		if .onlyLatLong?
			return c.GetDefault('allowLatLong', false)
		return true
		}
	addToMenu(c)
		{
		.menu.Add(c.name)
		.funcs.Add(Global(c.func))
		}
	get_map_choice()
		{
		choice = UserSettings.Get(.map_save_settings_key())
		if choice isnt false and .menu.Has?(choice)
			.which = .orig_which = .menu.Find(choice)
		}
	save_map_choice()
		{
		if .which isnt .orig_which
			UserSettings.Put(.map_save_settings_key(), .menu[.which])
		}
	MappingOption()
		{
		if .menu is false
			return 'None'
		return .menu[.validateWhich(.menu)]
		}
	map_save_settings_key()
		{
		return .multiLoc
			? "MultiLocation Map"
			: .onlyLatLong?
				? "Lat Long Map"
				: "AddressControl Map"
		}

	validateWhich(option)
		{
		return option.Member?(.which) ? .which : 0
		}

	Destroy()
		{
		.save_map_choice()
		super.Destroy()
		}
	}
