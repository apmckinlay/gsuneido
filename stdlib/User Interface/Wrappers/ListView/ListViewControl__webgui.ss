// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 		"List"
	ComponentName: "ListView"
	ComponentArgs: #()
	Xmin: 		100
	Ymin: 		50
	Xstretch:	1
	Ystretch:	1

	New(@args)
		{
		.ComponentArgs = Object(
			stretch: args.GetDefault('stretch', false),
			noHeader: args.GetDefault('noHeader', false))
		.cols = Object()
		.rows = Object()
		.ncols = .nrows = 0

		.updateHead()
		}

	SetStyle(style/*unused*/) {}
	SetExtendedStyle(ex_style/*unused*/) {}

	maxFieldWidth: 300
	AddColumn(name, header = false)
		{
		if header is false
			header = PromptOrHeading(name)
		.cols[.ncols] = [:name, :header]
		++.ncols
		.updateHead()
		}

	updateHead()
		{
		headCols = .cols.Map({
			Object(text: it.header, field: it.name, width: .maxFieldWidth,
				tip: false, sort: false) })
		headCols.Add(#(text: '', field: #CheckBox, width: 'calc(1em + 8px)', tip: false,
			sort: false),
			at: 0)
		.CancelAct(#UpdateHead)
		.Act(#UpdateHead, headCols)
		}

	Addrow(ob, unused = 0)
		{
		ob.CheckBox = false
		.rows.Add(ob)
		rec = []
		for (i = 0; i < .ncols; i++)
			{
			rec[.cols[i].name] = [type: 'text', data: ob[.cols[i].name], font: false,
				justify: 'left', color: false, ellipsis?:, html: false, bkColor: '']
			}
		rec[#CheckBox] = [type: 'text', data: "&#9744;", font: false,
				justify: 'center', color: false, html:, bkColor: '',
				extra: #('font-size': '150%', 'font-style': 'normal',
					'font-weight': 'normal')]
		.Act(#InsertData, .nrows, rec)
		return .nrows++
		}

	AddItem(label, image /*unused*/= 0, lParam/*unused*/ = 0)
		{
		return .Addrow(Object().Add(label, at: .cols[0].name))
		}

	SetMenu(menu)
		{
		.menu = menu
		.use_menu = true
		}

	GetMenu()
		{
		return .use_menu ? .menu : false
		}

	SetCheckState(i, state)
		{
		.rows[i].CheckBox = state is true
		.Act(#UpdateDataCell, i, 0, [type: 'text',
			data: state is true ? "&#09745;" : "&#9744;",
			font: false, justify: 'center', color: false, html:, bkColor: '',
			extra: #('font-size': '150%', 'font-style': 'normal',
					'font-weight': 'normal')])
		}

	GetCheckState(i)
		{
		return .rows[i].CheckBox
		}

	CheckAll(checked = true)
		{
		for (i = 0; i < .nrows; ++i)
			.SetCheckState(i, checked)
		}

	selected: false
	LBUTTONDOWN(row, col)
		{
		if row >= .rows.Size()
			return
		if col is 0
			.toggleRow(row)
		if .selected isnt false
			.Act(#DeSelectRow, .selected)
		.Act(#SelectRow, .selected = row)
		}

	toggleRow(row = false)
		{
		if row is false
			row = .selected
		.SetCheckState(row, not .rows[row].CheckBox)
		}

	KEYDOWN(wParam, ctrl = false, shift = false)
		{
		if .selected isnt false and .keydown_fns.Member?(wParam)
			(.keydown_fns[wParam])(:ctrl, :shift)
		return 0
		}

	getter_keydown_fns()
		{
		fns = Object()
		fns[VK.SPACE] =	.toggleRow
		fns[VK.UP] = 	.selectPrev
		fns[VK.DOWN] =	.selectNext
		return .keydown_fns = fns
		}

	selectPrev()
		{
		.selectRow(Max(.selected - 1, 0))
		}

	selectRow(row)
		{
		.Act(#DeSelectRow, .selected)
		.Act(#SelectRow, row)
		.Act('ScrollRowToView', .selected = row)
		}

	selectNext()
		{
		.selectRow(Min(.selected + 1, .rows.Size() - 1))
		}

	SetMaxWidth(column)
		{
		if String?(column)
			column = .columnNameToIndex(column) + 1 /*checkbox is the first column*/
		.Act(#SetMaxWidth, column)
		}

	SetColWidth(col, width)
		{
		if width is false
			SuneidoLog('ERROR: (CAUGHT) width: false is not supported in Suneido.js')
		.Act(#SetColWidth, col, width)
		}

	SetColumnWidth(column, width)
		{
		if String?(column)
			column = .columnNameToIndex(column)
		if (column is false)
			return
		.SetColWidth(column + 1/*checkbox is the first column*/, width)
		}

	columnNameToIndex(column)
		{
		return .cols.FindIf({ it.name is column })
		}
	}
