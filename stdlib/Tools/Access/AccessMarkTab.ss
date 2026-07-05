// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
// adds an icon to a tab when one of a set of fields isnt empty
// normally used in an Access,Observer plugin e.g.
class
	{
	CallClass(@args)
		{
		.checkArgs(args)
		access = args.access
		plugin = args.plugin
		.init(access, plugin)
		if plugin.tabs isnt false and not plugin.tabs.Destroyed?()
			plugin.tabs.SetImage(plugin.tab_i, .tabHasData?(access, plugin) ? 0 : -1)
		}

	checkArgs(args)
		{
		for field in #('access', 'plugin')
			if not args.Member?(field)
				throw 'missing argument'
		}

	tabHasData?(access, plugin)
		{
		hasdata = false
		data = access.GetData()
		for field in plugin.fields
			if data[field] isnt "" and data[field] isnt false
				hasdata = true
		return hasdata
		}

	TabAlwaysExists?: true
	init(access, plugin)
		{
		if false is (tabs = access.FindControl('Tabs'))
			return
		if tabs is plugin.GetDefault(#tabs, false)
			return
		plugin.tabs = tabs
		plugin.tabs.SetImageList(.GetInitImageList())
		tab = plugin.tabs.Tab
		for (i = tab.Count() - 1; i >= 0; --i)
			if tab.GetText(i) is .Tab
				break
		// some tabs may not exist for legit reason (ie. dynamic types)
		if i < 0 and not .TabAlwaysExists?
			{
			plugin.tabs = false
			return
			}
		Assert(i >= 0)
		plugin.tab_i = i
		}

	GetInitImageList()
		{
		return not Sys.SuneidoJs?()
			? Suneido.GetInit(#DocumentImagelistTab, .images)
			: Suneido.GetInit(#DocumentImagelistTabJs, .imagesJs)
		}

	images()
		{
		return Object(
			Object(ImageResource(#document), CLR.BLACK),
			Object(ImageResource(#document), CLR.red))
		}

	imagesJs()
		{
		codeOb = IconFont().MapToCharCode(#document)
		return Object(
			Object(char: codeOb.char, font: codeOb.font, color: #black),
			Object(char: codeOb.char, font: codeOb.font, color: #red))
		}
	}
