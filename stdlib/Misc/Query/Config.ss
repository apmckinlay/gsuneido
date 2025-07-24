// Copyright (C) 2002 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(query)
		{
		Retry()
			{
			.config(query)
			}
		}
	config(query)
		{
		config = .getConfig()
		result = config.GetDefault(query, false)
		if not Object?(result) or
			not Date?(result.asof) or
			Date().MinusSeconds(result.asof) > 10 /*= expiry seconds */
			{
			if false isnt from_db = Query1(query) // query outside Synchronized
				from_db.asof = Date()
			.Synchronized()
				{
				if result isnt config.GetDefault(query, false)
					throw "concurrent modification of Config " $ query
				result = config[query] =
					from_db is false ? Record(asof: Date()) : from_db
				}
			}
		return result
		}

	// tests should not use Invalidate, use Override method so values are
	// overridden on the client and the server
	Invalidate(query, alsoServer = false)
		{
		.Synchronized({ .getConfig()[query] = Random() })
		if alsoServer
			ServerEval('Config.Invalidate', query)
		}

	// Only for tests
	// WARNING - NOT multi-user safe
	Override(query, values)
		{
		if Sys.Client?()
			ServerEval("Config.Override", query, values)
		config = .getConfig()
		orig = 'config_original_' $ query

		if not Object?(config.GetDefault(query, false))
			Config(query)
		if not config.Member?(orig)
			config[orig] = config[query].Copy()
		config[query].Merge(values)
		asof = Date().Plus(minutes: 1)
		config[query].asof = config[query].asof_override = asof
		return asof
		}

	// Only for tests
	Restore(query)
		{
		if Sys.Client?()
			ServerEval("Config.Restore", query)
		config = .getConfig()
		orig = 'config_original_' $ query

		if config.Member?(orig)
			{
			expired? = (Object?(config.GetDefault(query, false)) and
				(not config[query].Member?('asof_override') or
				(Date?(config[query].asof_override) and
					config[query].asof_override < Date())))
			config[query] = config[orig]
			config.Delete(orig)
			if expired?
				throw "Config.Restore called on an expired override: " $ query
			}
		}

	GetCachedConfig(query)
		{
		x = .getConfig().GetDefault(query, false)
		return Object?(x) ? x : Config(query)
		}

	// used by tests and Company_Date_Protected
	ConfigCached?(query)
		{
		return Object?(.getConfig().GetDefault(query, false))
		}

	getConfig()
		{
		return Suneido.GetInit("Config", { Object() })
		}
	}
