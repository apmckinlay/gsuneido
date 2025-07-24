// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
// Have to run server side to ensure all clients are running with the required sources
MemoizeSingle
	{
	Func()
		{
		Assert(Sys.Client?() is: false)
		sources = Object()
		masterKeys = Object()
		prefix = Customizable.CustomTablePrefix
		QueryApply('customizable where table_name isnt "" and hidden? isnt true')
			{
			table = GetTableName(it.name)
			if '' is table or table is it.name
				continue
			name = prefix $ table $ ' > ' $ it.tab
			sources[name] = Object('Reporter', 'queries', :name,
				auth: Customizable.AuthPath(table, it.tab),
				tables: Object(it.name, it.table_name),
				query: .query(it.table_name, it.name, masterKeys),
				exclude: #(),
				bookLocation: false)
			}
		return sources
		}

	query(custTable, masterTable, masterKeys)
		{
		project = .tableKeys(masterTable, masterKeys).Join(', ')
		fkfield = .foreignKey(custTable)
		return custTable $
			' rename custtable_FK to ' $ fkfield $
			', custtable_num to custtable_num_new' $
			' join by(' $ fkfield $ ') (' $ masterTable $ ' project ' $ project $ ')'
		}

	foreignKey(custTable) // Extracted for tests
		{
		return Query1('indexes where table is ' $ Display(custTable) $
			' and columns is #custtable_FK').fkcolumns
		}

	// tableKeys is hopefully a temporary method. The eventual end goal is to include
	// ALL columns from the masterTable. However, with the current structures in place
	// we cannot do this without duplicating large portions of code
	tableKeys(masterTable, masterKeys)
		{
		if masterKeys.Member?(masterTable)
			return masterKeys[masterTable]
		tableKeys = Object()
		.queryKeys(masterTable).Each()
			{
			keys = it.Tr('()').Split(',').Map('Trim')
			tableKeys.MergeUnion(keys.Filter({ not it.Suffix?('_lower!') }))
			}
		if false isnt numField = tableKeys.FindOne({ it.Suffix?('_num') })
			{
			tableKeys.Remove(numField.Replace('_num', '_name'))
			tableKeys.Remove(numField.Replace('_num', '_abbrev'))
			}

		return masterKeys[masterTable] = tableKeys
		}

	queryKeys(masterTable) // Extracted for tests
		{
		return QueryKeys(masterTable)
		}

	ResetCache()
		{
		if Sys.Client?()
			ServerEval('CustomTableDataSources.ResetCache')
		else
			super.ResetCache()
		}
	}