// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.args, dynamicTypes)
		{
		if .model_name?(args[0]) or Object?(args[0]) or Class?(args[0])
			.innerModel = Construct(args[0])
		else // old style
			{
			.innerModel = new class { Before_Delete(unused) { } }
			.innerModel.Query = args[0]
			}
		.base_query = .innerModel.Query
		.SetQuery(.innerModel.Query)
		.fields = .c.Columns().Set_readonly()
		.keys = .c.Keys()
		.shortest_key = ShortestKey(.keys)

		locateArg = args.GetDefault('locate', #())
		locateKeyArgs = locateArg.GetDefault("keys", .keys)
		.locateKeys = .buildLocateKeys(locateKeyArgs)
		.mainLocatekey = locateArg.GetDefault("mainkey", false)
		.excludeSelectFields = args.GetDefault('excludeSelectFields', #())
		.option = .args.GetDefault('option', 'Access')
		.observers = Object()
		headerSelectPrompt = args.Size(list:) is 1 and dynamicTypes is false
			? 'no_prompts'
			: false
		.sf = SelectFields(.GetFields(), .GetExcludeSelectFields(), :headerSelectPrompt,
			includeMasterNum:)
		}

	model_name?(name)
		{
		return String?(name) and name.GlobalName?() and not TableExists?(name)
		}

	buildLocateKeys(locateKeyArgs)
		{
		locateKeys = Object()
		for k in locateKeyArgs
			{
			// composite keys do not work in locate
			if k.Has?(",") or k.Suffix?('_num')
				continue
			if k is "" or locateKeys.Member?(prompt = SelectPrompt(k))
				locateKeys[k] = k
			else
				locateKeys[prompt] = k
			}
		return locateKeys
		}

	GetQuery()
		{
		return .query
		}

	GetBaseQuery()
		{
		return .base_query
		}

	Before_Delete(@args)
		{
		.innerModel.Before_Delete(@args)
		}

	SetQuery(.query)
		{
		.setCursor(query)
		}

	GetCursor()
		{
		return .c
		}

	c: false
	setCursor(query)
		{
		if .c isnt false
			.c.Close()
		.c = SeekCursor(query)
		}

	GetFields()
		{
		return .fields
		}

	GetSelectFields()
		{
		return .sf
		}

	shortest_key: false
	GetKeyField()
		{
		// keep the order if we can
		if false is (field = .c.Order()) and
			false is (field = .queryKeyOrder())
			field = .shortest_key
		return field
		}

	queryKeyOrder()
		{
		sort = QueryGetSort(.query).RemovePrefix("reverse ")
		return .GetKeys().Has?(sort) ? sort : false
		}

	GetKeys()
		{
		return .keys
		}

	GetLockKey(record)
		{
		return .shortest_key.Split(',').Map({ record[it] }).Join('\x01')
		}

	LookupRecord(field, value)
		{
		return QueryFirst(QueryAddWhere(.base_query,
			" where " $ field $ " is " $ Display(value)))
		}

	SetKeyQuery(record)
		{
		if record is false
			throw "unable to generate key where clause"
		.keyquery = QueryStripSort(.base_query) $ "\n" $ .getKeyWhere(record)
		}

	getKeyWhere(record)
		{
		where = ""
		for keyfield in .shortest_key.Split(',')
			where $= " where " $ keyfield $ " is " $ Display(record[keyfield])
		return where
		}

	keyquery: false
	GetKeyQuery()
		{
		return .keyquery
		}

	AddWhere(where)
		{
		query = QueryAddWhere(.base_query, where)
		if where.Has?('leftjoin by') and not .base_query.Has?('/* tableHint: ')
			query = `/* tableHint: ` $ QueryGetTable(.base_query) $ ` */ ` $ query
		try
			if QueryEmpty?(query)
				return 'No Records Found'
		catch (err, '*regex')
			return "Invalid matcher - " $
				err.AfterFirst('regex: ').BeforeFirst('(from server)')
		.SetQuery(query)
		return true
		}

	GetLocateKeys()
		{
		return .locateKeys
		}

	GetLocateKey(by)
		{
		return .mainLocatekey isnt false ? .mainLocatekey : .locateKeys[by]
		}

	GetLocateLayout()
		{
		locateArg = .args.GetDefault('locate', #())
		locatequery = locateArg.GetDefault("query", .query)
		columns = locateArg.GetDefault('columns', false)
		sortFields = locateArg.GetDefault('sortFields', #())
		optionalRestrictions = locateArg.GetDefault('optionalRestrictions', #())
		startLast = locateArg.
			GetDefault('startLast', .args.GetDefault('startLast', false))
		customizeQueryCols = locateArg.GetDefault('customizeQueryCols', false)

		option = .option is 'Access' or .option is '' ? false : .option
		return Object('Locate', locatequery, .locateKeys, :columns, :sortFields,
			:option, excludeSelect: .excludeSelectFields, :optionalRestrictions,
			:customizeQueryCols, :startLast)
		}

	FindGotoField(field)
		{
		if .fields.Has?(field)
			return field

		// could be renamed on the Access to fill in
		// values on new records, but the field coming
		// from a KeyListView may not be renamed
		if .fields.Has?(field $ "_new")
			return field $ "_new"
		else
			return false
		}

	GetExcludeSelectFields()
		{
		return .excludeSelectFields.Copy().MergeUnion(
			Customizable.GetNonPermissableFields(.query))
		}

	plugins: false
	Plugins_Init()
		{
		.plugins = AccessPlugins(.fields, .option)
		}

	Plugins_Execute(@args)
		{
		if .plugins is false or not args.Member?('pluginType')
			return

		// pluginType = 'AccessObservers', 'Observers', 'AfterField'
		.plugins[args.pluginType](@args)
		}

	AddObserver(fn, at = false)
		{
		if at isnt false
			.observers.Add(fn, :at)
		else
			.observers.Add(fn)
		}

	RemoveObserver(fn)
		{
		.observers.Remove(fn)
		}

	NotifyObservers(@args)
		{
		ok = true
		// have to copy because some times obervers get removed in the process
		for observer in .observers.Copy()
			ok = ok and observer(@args)
		return ok
		}

	TableNotEmpty?()
		{
		table = QueryGetTable(.query, nothrow:)
		return .keyquery isnt false and table isnt "" and
			not QueryEmpty?(table) and QueryEmpty?(.query)
		}

	GetTableName()
		{
		return QueryGetTable(.base_query)
		}

	ForeignKeyUsage(record)
		{
		func = false
		try func = Global((QueryGetTable(.query) $ "_show_fk_usage").Capitalize())
		if not Function?(func)
			return ""
		return record.Eval(func)
		}
	}
