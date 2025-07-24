// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
MultiViewPlugins
	{
	New(query_columns, option)
		{
		super(query_columns, option)
		.accessObserver_plugins = Object().Set_default(#())
		Plugins().ForeachContribution('Access', false)
			{|x|
			.Collect(x, query_columns, option, .accessObserver_plugins)
			}
		}

	AccessObservers(@args)
		{
		args.collectionOb = .accessObserver_plugins
		.ExecuteObservers(@args)
		}
	}

