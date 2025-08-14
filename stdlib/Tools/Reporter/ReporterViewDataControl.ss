// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'View Data'
	CallClass(query, excludeCols, sf, hwnd = 0)
		{
		cols = QuerySelectColumns(query).Difference(excludeCols).RemoveIf(Internal?)
		fields_ob = sf.Fields
		suffix = .suffix()
		renamed = Object()
		.makeDatadicts(cols, fields_ob, suffix, renamed, sf.GetConverted())
		sort = .getSort(query, renamed)
		baseQuery = QueryStripSort(query)

		query = baseQuery $
			renamed.Map2(
				{ |oldCol, newCol| ' rename ' $ oldCol $ ' to ' $ newCol }).Join(' ') $
			Opt(' sort ', sort)
		cols.SortWith!(SelectPrompt)

		.showData(hwnd, query, cols, renamed)

		return Object(:query, :cols, :renamed)
		}

	showData(hwnd, query, cols, renamed)
		{
		ToolDialog(hwnd, Object(this, query, cols, renamed), title: .Title, border: 0)
		}

	New(query, cols, .renamed)
		{
		super(.layout(query, cols))
		.list = .FindControl('VirtualList')
		}

	layout(query, cols)
		{
		return Object('VirtualList', query, cols, headerSelectPrompt:,
			disableSelectFilter:, preventCustomExpand?:)
		}

	// extracted for testing
	suffix()
		{
		return '_' $ String(Timestamp()).Tr('#.')
		}

	getSort(query, renamed)
		{
		sort = QueryGetSort(query)
		reverseStr = ''
		if sort isnt ''
			{
			if sort.Prefix?('reverse ')
				sort = sort.RemovePrefix(reverseStr = 'reverse ')
			newSort = Object()
			for field in sort.Split(',')
				newSort.Add(renamed.GetDefault(field, field))
			sort = reverseStr $ newSort.Join(', ')
			}
		return sort
		}

	makeDatadicts(cols, fields_ob, suffix, renamed, converted)
		{
		numsToRemove = Object()
		for fn in GetContributions('SelectFieldsReporterRemoveNum')
			fn(cols)
		cols.Map!()
			{ |col|
			prompt = fields_ob.Find(col)
			selectPrompt = SelectPrompt(col)
			if ((selectPrompt is col or converted.Member?(col)) and
				false isnt prompt)
				{
				if col =~ '_(name|abbrev)'
					numsToRemove.Add(col.Replace('(name|abbrev)', 'num'))
				col = .makeColDD(col, suffix, prompt, renamed)
				}
			else if selectPrompt isnt col and col =~ '_(name|abbrev)'
				numsToRemove.Add(col.Replace('(name|abbrev)', 'num'))
			col
			}
		cols.RemoveIf({ numsToRemove.Has?(it) })
		}

	makeColDD(col, suffix, prompt, renamed)
		{
		newCol = col.RemoveSuffix("?") $ suffix
		ReporterModel.MakeDD(newCol, prompt, 'string')
		renamed[col] = newCol
		return newCol
		}

	VirtualList_AddGlobalMenu?()
		{
		return false
		}

	VirtualList_BuildContextMenu(rec /*unused*/)
		{
		return #('Summarize...', 'CrossTable...', 'Export...')
		}

	On_Context_Summarize()
		{
		RecordMenuManager.Summarize(.list)
		}

	On_Context_CrossTable()
		{
		RecordMenuManager.CrossTable(.list)
		}

	On_Context_Export()
		{
		RecordMenuManager.Export(.list)
		}

	Destroy()
		{
		for fld in .renamed.Values()
			Reporter.DeleteDD('Field_' $ fld)
		super.Destroy()
		}
	}
