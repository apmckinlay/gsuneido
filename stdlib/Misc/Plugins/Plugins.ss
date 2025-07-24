// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Singleton
	{
	New()
		{
		LibUnload.AddObserver(#Plugins, Plugins.Unload)
		.extenpts = Object()
		.contribs = Object().Set_default(Object())
		.errors = Object()
		.get()
		.check()
		}
	Reset()
		{
		LibUnload.RemoveObserver(#Plugins)
		super.Reset()
		}
	Unload(name)
		{
		if name.Prefix?('Plugin_')
			.Reset()
		}
	get()
		{
		.foreachPluginLibraryRecord()
			{ |x|
			name = x.name.Replace("Plugin_", "")
			try
				{
				plugin = x.text.Compile()
				if plugin.Member?('ExtensionPoints')
					{
					if .extenpts.Member?(name)
						throw "multiple definitions of ExtensionPoints"
					extenpts = plugin.ExtensionPoints
					.check_extens(extenpts)
					.extenpts[name] = extenpts
					}
				if plugin.Member?('Contributions')
					{
					contribs = plugin.Contributions
					.check_contribs(contribs)
					for c in contribs
						.contribs[c[0]].
							Add(c.Copy().Add(x.lib $ ':' $ x.name, at: 'from'))
					}
				}
			catch (err)
				.log_error(x.name $ ': ' $ err)
			}
		}
	foreachPluginLibraryRecord(block)
		{
		for lib in .libraries()
			{
			query = lib $ ' where group is -1 and name > "Plugin_" and name < "Plugin_~"
				extend lib = ' $ Display(lib)
			QueryApply(query, block)
			}
		}
	libraries()
		{
		return Libraries()
		}
	check_extens(extenpts)
		{
		if not Object?(extenpts)
			throw "ExtensionPoints: should be a list"
		for e in extenpts
			if not Object?(e) or not e.Member?(0) or not String?(e[0])
				throw "invalid extension: " $ Display(e)
		}
	check_contribs(contribs)
		{
		if not Object?(contribs)
			throw "Contributions: should be a list"
		for c in contribs
			if not Object?(c) or
				not c.Member?(0) or not String?(c[0]) or
				not c.Member?(1) or not String?(c[1])
				throw "invalid contribution: " $ Display(c)
		}

	check()
		{
		for contribs in .contribs
			for contrib in contribs
				{
				plugin = contrib[0]
				extenpt = contrib[1]
				if not .extenpts.Member?(plugin) or
					not .extenpts[plugin].Any?({ it[0] is extenpt})
					.log_error("invalid contribution: " $ Display(contrib))
				}
		}

	log_error(msg)
		{
		SuneidoLog('ERROR: Plugins - ' $ msg)
		.errors.Add(msg)
		}
	Errors()
		{
		return .errors
		}
	ShowErrors()
		{
		Alert(.errors.Empty?() ? "No Errors" : .errors.Join('\n'),
			"Plugin Errors")
		}
	Contributions(plugin, extenpt = false)
		{
		.check_plugin(plugin, extenpt)
		return .contribs[plugin].
			Filter({ extenpt is false or it[1] is extenpt })
		}
	check_plugin(plugin, extenpt = false)
		{
		if not .extenpts.Member?(plugin)
			{
			err = "Plugins: nonexistent plugin name: " $ plugin
			.throw_error(err)
			}
		if extenpt isnt false and not .Extenpts(plugin).Has?(Object(extenpt))
			{
			err = "Plugins: nonexistent extension point: " $ extenpt $
				" in " $ plugin $ " plugin"
			.throw_error(err)
			}
		}
	throw_error(err)
		{
		throw err
		}
	Extenpts(plugin)
		{
		.check_plugin(plugin)
		return .extenpts.GetDefault(plugin, Object())
		}
	ForeachContribution(plugin, extenpt, block, showErrors = false, sort = false)
		{
		contributions = .Contributions(plugin, extenpt)
		if sort
			contributions.Sort!(By(#seq))

		for x in contributions
			try
				block(x)
			catch (ex/*, 'block:'*/)
				{
				if ex is "block:break"
					break
				else if ex is "block:continue"
					continue
				else
					{
					f = showErrors ? .throw_error : .log_error
					f("error in ForeachContribution(" $
						plugin $ ", " $ extenpt $ ", " $ x.from $ ")\n " $ ex)
					}
				}
		}

	// used by ShowPluginsControl
	ForeachContrib(block)
		{
		for cs in .contribs
			for c in cs
				block(c)
		}
	}