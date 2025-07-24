// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(query_columns, option)
		{
		.Afterfield_plugins = Object().Set_default(#())
		.Observer_plugins = Object().Set_default(#())
		Plugins().ForeachContribution('MultiView', false)
			{|x|
			.Collect(x, query_columns, option)
			}
		}

	Collect(x, query_columns, option, pluginOb = false)
		{
		if x[2] isnt option
			return
		x = x.Copy() // so plugins can store stuff in it
		if String?(x.func)
			x.func = x.func.Compile()
		fields = Object?(x.fields) ? x.fields : (x.fields)()

		if pluginOb is false
			pluginOb = x[1] is 'AfterField'
				? .Afterfield_plugins : x[1] is 'Observer' ? .Observer_plugins : false

		Assert(pluginOb isnt false)
		.addPlugins(fields, x, query_columns, pluginOb)
		}

	addPlugins(fields, x, query_columns, pluginOb)
		{
		fieldsToUse = Object?(fields)
			? fields.Copy().Delete('dupCheck')
			: fields
		for f in fieldsToUse
			{
			if f isnt 'setdata' and not query_columns.Has?(f)
				{
				// Possible for Multiview to have valid duplicate checking on access
				// but not on list due to excludes.
				if fields.GetDefault('dupCheck', false)
					continue
				throw "invalid MultiView plugin field: " $ f
				}
			pluginOb[f].Add(x)
			}
		}

	Observers(@args)
		{
		.ExecuteObservers(@args)
		}

	ExecuteObservers(@args)
		{
		member = args.GetDefault('member', "")
		collectionOb = args.GetDefault('collectionOb', .Observer_plugins)
		for plugin in collectionOb[member]
			{
			args.plugin = plugin
			(plugin.func)(@args)
			}
		}

	AfterField(@args)
		{
		field = args.GetDefault('field', "")
		if .Afterfield_plugins.Member?(field)
			for x in .Afterfield_plugins[field]
				{
				(x.func)(@args)
				}
		}
	}