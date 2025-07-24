// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: 'DataSource'
	New()
		{
		super(.layout())
		.src = .FindControl('Source')
		}

	layout()
		{
		.sources = .Contributions(.all_sources = Object())
		return LastContribution('ReporterDataSource').Controls(.sources)
		}

	Contributions(all_sources = false)
		{
		sources = Object()
		.ForeachQuery()
			{ |c|
			if Object?(all_sources)
				all_sources.Add(c.name)
			if Suneido.User isnt 'default'
				{
				name = Reporter_table.Authorization(c)
				if not name.Has?('/')
					{
					if false is x = QueryFirst(FindCurrentBook() $
						' where name = ' $ Display(name) $
						' and not path.Has?("Reporter Reports")' $
						' and not path.Has?("Reporter Forms")' $
						' project path, name sort path')
						continue
					name = x.path $ '/' $ x.name
					}
				if AccessPermissions(name) is false
					continue
				}
			if sources.Member?(c.name)
				Alert('duplicate Reporter contribution: ' $ c.name)
			sources[c.name] = c.Set_default('')
			}
		return sources.MergeNew(Customizable.GetPermissableDataSources())
		}

	ForeachQuery(block)
		{
		Plugins().ForeachContribution('Reporter', 'queries')
			{|c|
			c = c.Copy()
			if c.query.BeforeFirst(".").GlobalName?()
				c.query = Global(c.query)()
			block(c)
			}
		}

	Source()
		{
		return LastContribution('ReporterDataSource').Source(.src, .sources)
		}

	SourceName()
		{
		return .src.Get()
		}

	sources: false
	all_sources: false
	Authorized?(source_name = false)
		{
		if source_name is false
			source_name = .SourceName()
		all_sources = .all_sources is false ? Object() : .all_sources
		sources = .sources is false
			? .Contributions((all_sources.Empty?() ? all_sources : false))
			: .sources
		return sources.Member?(source_name) or
			not all_sources.Has?(source_name) // invalid are "authorized"
		}
	}
