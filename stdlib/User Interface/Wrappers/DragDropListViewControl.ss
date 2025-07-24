// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
ListViewControl
	{
	New(@args)
		{
		super(@.args(args))
		}
	args(args)
		{
		args.allowDrag = true
		args.style = args.GetDefault(#style, 0) |
			LVS.REPORT | LVS.SHOWSELALWAYS
		args.exStyle = args.GetDefault(#exStyle, 0) |
			LVS_EX.GRIDLINES | LVS_EX.FULLROWSELECT | LVS_EX.HEADERDRAGDROP
		.noSort = args.GetDefault(#noSort, false)
		return args
		}
	GetColumns()
		{
		return .GetHeaderFields()
		}
	GetHeaderFields()
		{
		cols = Object()
		hdr = .GetHeader()
		n = SendMessage(hdr, HDM.GETITEMCOUNT, 0, 0)
		Assert(n < 256) /*= max columns in LV_ORDERCOLUMN */
		colOrder = Object()
		SendMessageListColumnOrder(hdr, HDM.GETORDERARRAY, n, colOrder)
		.ColumnIndexes = Object()
		for (i = 0; i < n; i++)
			{
			hdi = Object(
				mask: 			HDI.TEXT | HDI.WIDTH | HDI.FORMAT
				cchTextMax:		200)
			SendMessageHditem(hdr, HDM.GETITEM, colOrder.order[i], hdi)
			field = .FieldMap[hdi.pszText]
			cols.Add(field)
			.ColumnIndexes.Add(colOrder.order[i], at: field)
			}
		return cols
		}
	GetColWidth(i)
		{
		cols = .GetColumns()
		return SendMessage(.Hwnd, LVM.GETCOLUMNWIDTH, .ColumnIndexes[cols[i]], 0)
		}

	LVN_COLUMNCLICK(lParam)
		{
		if .noSort
			return super.LVN_COLUMNCLICK(lParam)
		nm = NMLISTVIEW(lParam)
		.SortList(nm.iSubItem)
		return super.LVN_COLUMNCLICK(lParam)
		}

	SortList(colIdx)
		{
		.columns = .GetColumns() // need to call this to set the ColumnIndexes
		if false is col = .ColumnIndexes.Find(colIdx)
			return

		model = .GetModel()
		if model is false
			.sortListObjectData(col)
		else
			.sortListQuery(col)
		}
	sortCol: false
	sortListObjectData(sortCol)
		{
		data = .buildData()
		.DeleteAll()

		.formatting = new ListFormatting(true, true)
		.formatting.SetFormats(.columns)

		prevSortCol = .sortCol
		.sortCol = sortCol
		.setSortVal(data)

		compareFunc = .defCompareFunc
		cmpfn = .sortCol is false ? {|x, y| compareFunc(y, x) } : compareFunc
		if .sortCol is prevSortCol
			{
			for rec in data.Sort!(cmpfn).Reverse!()
				.addrecord(rec)
			.sortCol = false
			}
		else
			for rec in data.Sort!(cmpfn)
				.addrecord(rec)
		}
	buildData()
		{
		data = Object()
		for (i = 0; i < .GetRowCount(); i++)
			{
			row = .Getrow(i)
			row.checkstate = .GetCheckState(i)
			data.Add(row)
			}
		return data
		}
	setSortVal(data)
		{
		type = DatadictType(.sortCol)
		if type is 'date'
			data.Map!({|x| x.Add(Date(x[.sortCol]), at: 'sortcol') })
		else if type is 'number'
			data.Map!({|x| x.Add(Number(x[.sortCol]), at: 'sortcol') })
		else
			data.Map!({|x| x.Add(x[.sortCol], at: 'sortcol') })
		}
	addrecord(rec)
		{
		.Addrow(rec)
		.SetCheckState(.GetRowCount() - 1, rec.checkstate)
		}

	defCompareFunc(x, y)
		{ return .formatting.CompareRows('sortcol', x , y) }

	sortListQuery(col)
		{
		sortCol = .Send('DragDropListView_GetSortByColumn', col)
		if sortCol is 0 or sortCol is false
			return
		query = QueryStripSort(.GetQuery())
		sort = ' sort'
		if sortCol is .sortCol
			{
			sort $= ' reverse'
			.sortCol = false
			}
		else
			.sortCol = sortCol
		sort $= ' ' $ sortCol
		.ResetQueryWithoutDestroy(query $ sort)
		}
	SetReadOnly(readOnly/*unused*/)
		{
		// override default method so control doesn't get disabled
		}
	}