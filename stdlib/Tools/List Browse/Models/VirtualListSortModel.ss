// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	sortOb: false
	customSort?: false
	New(query, .sf, .sortSaveName = false, .loadAll? = false,
		.disableCheckSortLimit? = false)
		{
		sort = QueryGetSort(query)
		.origSort = sort
		.initSortOb(sort)

		.loadSavedSort(query)
		}

	initSortOb(sort)
		{
		sortDir = sort.Prefix?('reverse ') ? -1 : 1
		.sortOb = sort.Replace('^reverse ', '').Split(',').Map!({
			col = it.Trim()
			colDir = col.Prefix?('reverse__') ? -1 : 1
			.buildSortOb(col, colDir * sortDir)
			})
		}

	buildSortOb(col, dir, displayCol = false)
		{
		id? = String?(col) and col.Has?('_num') and
			.sf.NameAbbrev?(col.Replace('_num', '_name'))
		retVal = [:col, :dir, :id?]
		if displayCol isnt false
			retVal.displayCol = displayCol
		return retVal
		}

	loadSavedSort(query)
		{
		if not .SaveSort?()
			return
		if false is val = UserSettings.Get(.sortSaveName)
			return

		.customSort? = true
		availableCols = QueryColumns(query)
		cols = Object?(val)
			? val
			: val.Split(',').Map({
				Object(dir: it.Prefix?('-') ? -1 : 1, col: it.Replace('^-', '')) })
		.sortOb = Object()
		for c in cols
			{
			dir = c.dir
			col = c.col
			if availableCols.Has?(col)
				.sortOb.Add(.buildSortOb(col, dir,
					c.GetDefault('displayCol', false)))
			}
		if .loadAll? isnt true
			.sortOb = .sortOb[..1]
		}

	BuildQuery(query, where = '', sortCol = false)
		{
		// needed if sort or where has a join so we can find the base table
		initialQuery = query
		query = .StripSort(query) $ Opt(' ', where)
		sortOb = sortCol is false ? false : .getSortOb(.sortOb.DeepCopy(), sortCol)
		// Add /* tableHint: */ when sort on foreign name to prevent Complex Queries
		// also need to handle when the where clause gets a leftjoin added
		// (e.g. from GetForeignNumsFromNameAbbrevFilter)
		sortStr = .getSortStr(query, sortOb)
		if ((.join? is true or where.Has?('leftjoin by')) and
			not query.Has?('/* tableHint: '))
			query = `/* tableHint: ` $ QueryGetTable(initialQuery) $ ` */ ` $ query
		return query $ sortStr
		}

	StripSort(query)
		{
		sups = QueryGetSuppressions(query)
		query = QueryStripSort(query)
		if .join? is true and query.Has?('leftjoin by')
			query = query.BeforeLast('leftjoin by')
		if query.Has?(' extend reverse__')
			query = query.BeforeFirst(' extend reverse__')
		keep = ''
		if sups isnt QueryGetSuppressions(query) // sups were removed by BeforeFirst/Last
			keep = sups.Empty?() ? '' : ' ' $ sups.Join(' ')
		return query.Trim() $ keep
		}

	getSortStr(query = '', sortOb = false) // query is used for check duplicate leftjoin
		{
		if sortOb is false
			sortOb = .sortOb
		spec = .getSortSpec(query, sortOb)
		return spec.sortExtend $ ' ' $ spec.sortStr
		}

	getSortSpec(query, sortOb)
		{
		sortSpec = Object(sortExtend: '', sortStr: '')
		if sortOb is false or sortOb.Empty?()
			return sortSpec
		if .loadAll? isnt true and .customSort?
			return .sortSingleColOnly(sortSpec, query, sortOb)

		// REFACTOR: remove unnecessary extend for in-memory sorting
		if sortOb.Every?({ it.dir is 1 })
			{
			sortSpec.sortStr = 'sort ' $ sortOb.Map({ it.col }).Join(',')
			return sortSpec
			}

		if sortOb.Every?({ it.dir is -1 })
			{
			sortSpec.sortStr = 'sort reverse ' $ sortOb.Map({ it.col }).Join(',')
			return sortSpec
			}

		col1 = sortOb[0].col
		col2 = sortOb[1].col
		sortSpec.sortStr = (sortOb[0].dir is 1 ? 'sort ' : 'sort reverse ') $ col1 $
			',reverse__' $ col2
		sortSpec.sortExtend = ' extend reverse__' $ col2 $ ' = true'
		return sortSpec
		}

	join?: false
	sortSingleColOnly(sortSpec, query, sortOb)
		{
		sortOb.Delete(1)
		sort = sortOb[0]
		sortCol = sort.col
		if sort.id? is true
			sortCol = sortCol.Replace("_num", "_name")
		leftJoin = .sf.Joins(sortCol)
		.join? = false
		if not query.Has?(leftJoin.Trim()) and not QueryColumns(query).Has?(sortCol)
			{
			sortSpec.sortExtend = .sf.Joins(sortCol)
			.join? = true
			}
		sortSpec.sortStr = 'sort ' $ (sortOb[0].dir is 1 ? '' : 'reverse ') $ sortCol
		return sortSpec
		}

	UsingDefaultSort?(query)
		{
		return .origSort is QueryGetSort(query)
		}

	SetSort(displayCol, dataCol = false)
		{
		// displayCol is the column shows to the user what was sorted on.
		// (i.e. what they clicked on)
		// However, there are some scenarios where that is not the column we want to
		// sort on, in that case dataCol is what we actually sorted on (behind the scenes)

		.customSort? = true

		// if dataCol does not exist, then what we show the user IS what we sorted on
		//		In this case we just save displayColumn as the sortColumn
		// if dataCol does exist, then we are sorting on it and not what we show the user.
		//		In that case we need to save the dataColumn as the sortColumn
		//		but we also still need to save the displayColumn.
		if dataCol is false
			.sortOb = .getSortOb(.sortOb, displayCol)
		else
			.sortOb = .getSortOb(.sortOb, dataCol, displayCol)
		}

	getSortOb(sortOb, col, displayCol = false)
		{
		// just let buildSortOb handle/worry about when displayCol is false
		if sortOb is false or sortOb.Empty?()
			{
			sortOb = Object(.buildSortOb(col, 1, displayCol))
			return sortOb
			}

		if sortOb[0].col is col
			{
			sortOb[0].dir *= -1
			return sortOb
			}

		if .loadAll? is true
			sortOb[1] = sortOb[0]
		sortOb[0] = .buildSortOb(col, 1, displayCol)
		return sortOb
		}

	GetPrimarySort()
		{
		return .sortOb.GetDefault(0, .buildSortOb(false, 1))
		}

	CheckSortable(query, col)
		{
		if UnsortableField?(col) or not QueryColumns(query).Has?(col)
			{
			InfoWindowControl(SelectPrompt(col) $ ' is not sortable.', titleSize: 0)
			return false
			}
		return true
		}

	SortInMemory(data, sort)
		{
		reverseAll? = sort.Prefix?('reverse ')
		sortFields = sort.RemovePrefix('reverse ').Split(',').Map!(#Trim)
		compare = {|a, b, reverse?| reverse? ? a > b : a < b  }
		data.Sort!({|x, y|
			blockReturn = false
			for f in sortFields
				{
				reverse? = f.Prefix?('reverse__')
				f = f.RemovePrefix('reverse__')
				if compare(x[f], y[f], reverse?)
					{
					blockReturn = true
					break
					}
				else if x[f] is y[f]
					continue
				else
					{
					blockReturn = false
					break
					}
				}
			blockReturn
			})

		if reverseAll?
			data.Reverse!()
		}

	SaveSort?()
		{
		return .sortSaveName isnt false
		}

	ResetSort()
		{
		.customSort? = false
		.initSortOb(.origSort)

		if .SaveSort?()
			UserSettings.Remove(.sortSaveName)
		}

	SetDefaultSort()
		{
		if not .SaveSort?()
			return

		if .sortOb is false or .sortOb.Empty?()
			return

		if .loadAll? isnt true
			.sortOb = .sortOb[..1]
		UserSettings.Put(.sortSaveName, .sortOb)
		}

	sortLimit: 100000
	QueryAboveSortLimit?(query, where)
		{
		if not .CheckAboveSortLimit?()
			return false
		// Note: where should only include where on indexes. Other fields will cause
		// inaccurate estimate
		estimated = QueryCost(.BuildQuery(query, where, false)).nrecs
		return estimated > .sortLimit
		}

	overSortLimit?: ''
	SetOverSortLimit(overSortLimit?)
		{
		if not .CheckAboveSortLimit?()
			return

		if overSortLimit?
			{
			.customSort? = false
			.initSortOb(.origSort)
			}
		.overSortLimit? = overSortLimit?
		}

	OverSortLimit?()
		{
		return .CheckAboveSortLimit?() and .overSortLimit? is true
		}

	SortLimitChecked?()
		{
		return .overSortLimit? isnt ""
		}

	CheckAboveSortLimit?()
		{
		return not .disableCheckSortLimit?
		}

	Destroy()
		{
		}
	}
