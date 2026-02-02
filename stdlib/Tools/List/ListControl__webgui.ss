// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: "List"
	ComponentName: "List"
	lcf:			#(SELECTED: 0x0100, HIGHLIGHTED:0x0200, FOCUSED: 0x0400,
						HOVERED: 0x0800)
	columns: 		false
	data: 			false
	markCol: 		false
	contextOnly: 	false
	readOnly:		false
	allowTab:		true
	editor:			false
	focused:		false
	forceDelete:	false
	sortCol:		false
	shiftAnchor: 	false
	maxLines: 		3000

	multiSelect:	false
	New(columns = false, data = false, .defWidth = 100, noShading = false,
		noDragDrop = false, highlightColor = 0x00FF9957,
		noHeaderButtons = false,
		.headerSelectPrompt = false, booleansAsBox = false, fontSize = "",
		.resetColumns = false, .customizeColumns = false,
		.alwaysHighlightSelected = false, indicateHovered = false,
		.columnsSaveName = false, .checkBoxColumn = false, .sortSaveHandler = false,
		.trackValid = false, .limitHandler = false)
		{
		Assert((Integer?(defWidth) and (defWidth > -1)) or (defWidth is false))
		.origColumns = columns isnt false ? columns : Object()
		.header = SuJsListHeader(:noDragDrop, :noHeaderButtons, :headerSelectPrompt)
		.formatting = ListFormatting(fontSize isnt "", booleansAsBox)
		.SetColumns(.origColumns)
		.Set(data isnt false ? data : Object())
		.highlight_colors = Object()
		.AddHighlight(false, highlightColor)	// set default at index 0

		.ComponentArgs = Object(noShading, indicateHovered, :noDragDrop, :noHeaderButtons)
		if .columnsSaveName isnt false
			UserColumns.Load(.GetColumns(), .columnsSaveName, this, deletecol: false,
				initialized?:, load_visible?:)
		if .trackValid
			.Window.AddValidationItem(this)
		}

	SetColumns(columns, reset = false)
		{
		Assert(Object?(columns) and not columns.HasNamed?())
		if (.columns is columns and not reset)
			return
		contentChanged? = .isColumnContentChanged?(.columns, columns)
		.columns = columns.Copy()
		.widths = Object()
		.header.Clear()
		.formatting.SetFormats(.columns) // create formats for displaying columns
		.processColumns()

		if (reset)
			.header_changed? = true
		.Act(#UpdateHead, .header.Get(), .markCol,
			data: contentChanged? ? .generateDisplayData(.data) : false,
			showSortIndicator: .showSortIndicator)
		}

	isColumnContentChanged?(oldCols, newCols)
		{
		// this is when initializing the list , no need to pass data
		if oldCols is false
			return false

		return newCols.Size() isnt oldCols.Size() or
			newCols.Copy().Sort!() isnt oldCols.Copy().Sort!()
		}

	processColumns()
		{
		colno = 0
		for (col in .columns.Members())
			if (col is 0 and .columns[0] is "listrow_deleted")
				{
				.markCol = true
				.widths.Add('1em')
				}
			else
				.addHeader(.columns[col], colno++)
		}

	AppendColumns(newColumns)
		{
		.formatting.SetFormats(.columns.Copy().Append(newColumns))
		offset = .columns.Has?("listrow_deleted") ? -1 : 0
		for column in newColumns
			{
			.addHeader(column, .columns.Size() + offset)
			.columns.Add(column)
			}
		.Act(#UpdateHead, .header.Get(), .markCol, data: .generateDisplayData(.data),
			showSortIndicator: .showSortIndicator)
		}

	addHeader(column, colno)
		{
		.header.AddItem(column, .defWidth)
		.header.SetItemFormat(colno, .formatting.GetHeaderAlign(column))
		.widths.Add(.header.GetItemWidth(colno))
		}

	GetColumns()
		{
		return .columns.Copy()
		}

	GetCheckBoxField()
		{
		return .checkBoxColumn
		}

	GetVisibleColumns()
		{
		cols = Object()
		for (i = 0; i < .columns.Size(); i++)
			if .GetColWidth(i) > 0
				cols.Add(.columns[i])
		return cols
		}

	GetNumFullyVisibleRows()
		{
		return 9999/*=big number*/
		}

	Get()
		{
		return .data
		}

	Set(data, continueWhenLimitReached? = false)
		{
		if .data is data
			return
		.data = Object()
		.focused = false
		.Act(#ClearData)
		.AddRows(data, :continueWhenLimitReached?)
		.focused = false
		.select(0, true, false, false)
		}

	Clear()
		{
		.Set(Object())
		}

	GetField(field)
		{
		if 1 isnt size = .getSelection().Size()
			throw "List.GetField requires a single row selection; size is " $
			Display(size)
		return .GetCurrentRecord()[field]
		}

	SetField(field, value, idx = false)
		{
		if idx is false
			{
			selected = .getSelection()
			if selected.Size() isnt 1
				throw "List.SetField requires a single row selection"
			idx = selected[0]
			}
		record = .GetRow(idx)
		record[field] = value
		}

	GetCurrentRecord()
		{
		selected = .GetSelection()
		return selected.Size() is 1 ? .GetRow(selected[0]) : false
		}
	GetSelection()
		{
		return .getFocusSelection()
		}
	getFocusSelection()
		{
		selected = .getSelection()
		if selected.Size() is 0 and .focused isnt false
			{
			.SetSelection(.focused)
			selected = Object(.focused)
			}
		return selected
		}
	getSelection()
		{
		selection = Object()
		for (row in .data.Members())
			if .RowSelected?(row)
				selection.Add(row)
		return selection
		}
	SetSelection(row)
		{
		Assert(.data.Member?(row))
		.selectFocus(row)
		}

	ClearSelectFocus()
		{
		.ClearSelect()
		.focused = false
		}

	ClearSelect()
		{
		.Send("List_BeforeClearSelect")
		for row in .data.Members().Filter(.RowSelected?).Copy()
			.DeSelectRow(row)
		.Send("List_AfterClearSelect")
		}

	AddRows(data, continueWhenLimitReached? = false)
		{
		Assert(Object?(data) and not data.HasNamed?(), "data must be a list")
		runInfo = Object(limitReached?: false, :continueWhenLimitReached?)
		batch = .generateDisplayData(data, { |row| .addRow(row, runInfo) })
		if .hideContent isnt true
			{
			.Act(#AddBatch, batch)
			.data.Members().Each()
				{
				if .RowHighlighted?(it)
					.Act('AddHighlightRow', it,
						.highlight_colors[HIWORD(.data[it].listrow_flags)])
				}
			}
		.setSortIndicator(.sortCol = false)
		.repaintHeader()
		}

	addRow(row, runInfo)
		{
		if not runInfo.limitReached? and .reachLimit?()
			{
			runInfo.limitReached? = true
			if runInfo.continueWhenLimitReached? isnt true
				return false
			}
		.data.Add(row.Set_default(""))
		return true
		}

	reachLimit?()
		{
		if .limitHandler isnt false and .data.Size() >= .maxLines
			{
			(.limitHandler)('Load limit of ' $ .maxLines $ ' reached.')
			return true
			}
		return false
		}

	FindRowIdx(field, value)
		{
		return .data.FindIf({ it[field] is value })
		}

	generateDisplayData(data, block = false)
		{
		batch = Object()
		for i in data.Members()
			{
			.checkRow(data[i])
			if block isnt false
				{
				if false is block(data[i])
					break
				}
			batch[i] = .formatting.PaintRow(data[i])
			}
		return batch
		}

	checkRow(dataRow)
		{
		Assert(Object?(dataRow), "row must be an object")
		Assert(not dataRow.Readonly?(), "row cannot be read-only object")
		}

	GetInvalidFieldData(rec, field)
		{
		if rec.Member?('List_InvalidData') and Object?(rec.List_InvalidData)
			return rec.List_InvalidData.GetDefault(field, '')
		return ''
		}

	GetColWidth(col)
		{
		return .widths[col]
		}

	SetCol(col /*unused*/, text /*unused*/)
		{
		throw "not implemented"
		}
	GetCol(col)
		{
		// returns the name of column at index col
		return .columns.GetDefault(col, false)
		}
	GetNumCols()
		{
		return .columns.Size()
		}
	GetRow(row)
		{
		return .data[row]
		}
	SetRow(row, value)
		{
		Assert(.data.Member?(row))
		.checkRow(value)
		.data[row] = value.Set_default("")
		.sortCol = false
		.RepaintRow(row)
		}
	AddRow(value)
		{
		.InsertRow(.GetNumRows(), value)
		}
	InsertRow(row, value)
		{
		Assert(.data.Member?(row) or row is .GetNumRows())
		.checkRow(value)
		.data.Add(value.Set_default(""), at: row)
		.sortCol = false
		if .hideContent isnt true
			.Act(#InsertData, row, .formatting.PaintRow(.data[row]))
		}
	CheckAndInsertRow(row, newRecord, useDefaultsIfEmpty? = false)
		{
		if false is .wantNewRow(prevRow: .focused,
			record: newRecord, :useDefaultsIfEmpty?)
			return false;			// not allowed by parent
		.InsertRow(row, newRecord)
		.Send("List_NewRowAdded", row, record: newRecord)
		return .data[row]
		}

	wantNewRow(prevRow, record, useDefaultsIfEmpty? = false)
		{
		if .reachLimit?()
			return false
		return .Send("List_WantNewRow", :prevRow, :record, :useDefaultsIfEmpty?)
		}

	GetNumRows()
		{
		return .data.Size()
		}

	SetMultiSelect(multi)
		{
		Assert(Boolean?(multi))
		if .multiSelect isnt multi and false is (.multiSelect = multi)
			.select(0, false, false, true)
		}
	GetMultiSelect()
		{
		return .multiSelect
		}

	AllowTab(allow)
		{
		Assert(Boolean?(allow))
		.allowTab = allow
		}

	LBUTTONDOWN(row, col, shift, control, mouseEventId = false)
		{
		if .contextOnly
			return 0

		SetFocus(.Hwnd)

		// when another ctrl loses focus, it could trigger list destroy (dynamic layouts)
		if .Destroyed?()
			return 0
		if 0 isnt .Send("List_SingleClick", .data.Member?(row) ? row : false, col)
			return 0

		return .updateSelection(row, shift, control, mouseEventId)
		}

	updateSelection(row, shift, control, mouseEventId)
		{
		if row >= .GetNumRows()
			{
			if not shift and not control
				.ClearSelect()
			return 0
			}
		if not shift or .multiSelect is false
			.selectFocus(row, control, false, true)
		else
			.select(row - .focused, false, true, true)
		if .dragging?(shift, control)
			{
			.dragging = true
			.Act('List_AllowDragging', .focused, :mouseEventId)
			}
		return 0
		}

	dragging?(shift, control)
		{
		return not .readOnly and not shift and
			not control and true is .Send("List_AllowMove")
		}

	ListMoveRow(focused, newRow)
		{
		Assert(focused is: .focused)
		for (inc = newRow > .focused ? 1 : -1; .focused isnt newRow; .focused += inc)
			{
			.data.Swap(.focused, .focused + inc)
			.Send("List_Move", .focused, .focused + inc)
			}
		.sortCol = false
		}

	LBUTTONUP(row, col)
		{
		if .contextOnly
			return 0

//		.EndDrag()

		// get row/col from coordinates in lParam
		.Send("List_LButtonUp", .data.Member?(row) ? row : false, col)
		}

	LBUTTONDBLCLK(row, col)
		{
		if .contextOnly
			return 0

		if not .HasFocus?()
			return 0

		datarow = .data.Member?(row) ? row : false
		col = Min(col, .GetNumCols() - 1)
		if datarow is false
			.ClearSelect()
		else
			.selectFocus(row)
		dbl_click_result = .Send("List_DoubleClick", datarow, col)
		// added zoom here to make sure the logic still follows the same sequence
		if .Send("List_AllowZoom", datarow, col) is true
			.zoomOnField(datarow, col, .GetCol(col))
		.editFromDblClick(dbl_click_result, datarow, row, col)
		return 0
		}

	zoomOnField(row, col, field)
		{
		if row isnt false and col isnt false and
			String?(value = .GetField(field)) and value isnt "" and
			false isnt ctrl = GetControlClass.FromField(field)
			ctrl.ZoomReadonly(value)
		}

	editFromDblClick(dbl_click_result, datarow, row, col)
		{
		if .edtiable?(dbl_click_result, datarow, row)
			{
			col = (col is 0 and .markCol is true) ? 1 : col
			if not .Edit(col, row, 1, canMoveRows: false) and col > 0 // begin edit
				.Edit(col - 1, row, -1, canMoveRows: false)
			}
		}

	edtiable?(dbl_click_result, datarow, row)
		{
		return ((dbl_click_result is 0 and not .readOnly) and
			(datarow is false or .data[row].listrow_deleted isnt true))
		}

	ContextNew()
		{
		// delay to avoid ListEditWindow from being closed
		// by the SetFocus triggered by closing context menu
		_forceOnBrowser = true
		.Defer({ .KEYDOWN(VK.INSERT, 0) })
		}

	ContextEdit(pt)
		{
		Assert(.tempPos isnt: false)
		Assert(.tempPos.x is: pt.x)
		Assert(.tempPos.y is: pt.y)
		tempPos = .tempPos
		// delay to avoid ListEditWindow from being closed
		// by the SetFocus triggered by closing context menu
		_forceOnBrowser = true
		if .data[tempPos.row].listrow_deleted isnt true
			.Defer({ .Edit(Min(tempPos.col, .GetNumCols() - 1), tempPos.row, 1) })
		}

	Edit(col, row, amt, canMoveRows = false)
		{
		// attempts to edit cell at (col, row)
		// if this is not possible, attempts to edit next cell amt cells away
		// if canMoveRows is false,
		// it will not attempt to edit a cell in another row
		while (.edit(col, row) is false)
			{
			nextCell = .getNextCell(col, row, amt)	// try next cell
			if not .nextCellAvailable?(nextCell, col, row, amt, canMoveRows)
				return false
			col = nextCell.col
			row = nextCell.row
			}
		return true
		}
	nextCellAvailable?(nextCell, col, row, amt, canMoveRows)
		{
		if nextCell.col is col and nextCell.row is row
			return false	// no available next cell
		if nextCell.row isnt row and not canMoveRows
			return false	// can't move (or create) rows
		// only check next cell available when tabbing to next row
		if nextCell.row isnt row and amt is 1 and
			true is .Send('List_NextCellNotAvailable?', row)
			return false
		return true
		}

	customFields: false
	SetCustomFields(customFields)
		{
		.customFields = customFields
		}
	edit(col, row)
		// attempts to edit a cell
		// if the row specified by row does not exist, a new row is added
		// returns true if
		// 		editing begun successfully OR readOnly OR
		// 		no cols OR callback disallowed row adding
		// false otherwise
		{
		if (row >= .GetNumRows())	// add a row?
			{
			row = .GetNumRows()
			newRow = Record()
			if (false is .wantNewRow(prevRow: row - 1, record: newRow))
				return true
			.AddRow(newRow)
			.Send("List_NewRowAdded", row, record: newRow)
			}
		if (false is .Send("List_AllowCellEdit", col, row) or
			false is control = .Send("List_WantEditField",
			:col, :row, data: .data[row][.columns[col]]))
			return false						// no editing allowed by parent

		// The following line MUST ALWAYS be done (ie. when inserting rows)
		.selectFocus(row)
		.ScrollColToView(col)

		if .getSelection().Size() isnt 1 // should be assert but don't want to annoy users
			SuneidoLog('ERROR: ListControl edit selection size isnt 1. It is ' $
				Display(.getSelection().Size()), calls:)
		custom = .customFields isnt false
			? .customFields.GetDefault(.GetCol(col), false)
			: false
		.editor = new ListEditWindow(control,
			.Send('List_EditFieldReadonly', col, row) is true,
			col, row, this, :custom, customFields: .customFields)
		return true;							// field editing begun!
		}

	selectFocus(row, ctrl = false, shift = false, select = true)
		{
		row = Min(row, .GetNumRows() - 1)
		oldFocus = .focused
		.focused = row
		if oldFocus is false
			.select(0, ctrl, shift, select)
		else
			.select(shift ? .focused - oldFocus : 0, ctrl, shift, select)
		}

	select(amt, ctrl = false, shift = false, select = false)
		{
		if not .selectable?()
			return
		if not shift
			.shiftAnchor = false
		newFocus = Max(0, Min(.focused + amt, .GetNumRows() - 1))
		oldSelection = .getSelection()
		newSelection = oldSelection.Copy()
		if ctrl
			newSelection = .ctrlSelect(select, newFocus, newSelection)
		else if shift
			newSelection = .shiftSelect(newFocus)
		else if select
			newSelection = Object(newFocus)
		.RepaintRow(newFocus)
		.RepaintRow(.focused)
		.focused = newFocus
		.handleSelectionChanged(newSelection, oldSelection)
		.ScrollRowToView(.focused)		// ensure focus is in view
		}

	selectable?()
		{
		if .GetNumRows() < 1
			{
			.focused = false
			.Send("List_Selection", selection: false)
			return false
			}
		if .focused is false
			.focused = 0
		return true
		}

	handleSelectionChanged(newSelection, oldSelection)
		{
		if newSelection.Sort!() isnt oldSelection	// if selection changed
			{
			for (row in oldSelection.Difference(newSelection))
				.DeSelectRow(row)
			for (row in newSelection.Difference(oldSelection))
				.SelectRow(row)
			.Send("List_Selection",
				selection: newSelection.Empty?() ? false : newSelection)
			}
		}

	ctrlSelect(select, newFocus, newSelection)
		{
		if select
			{
			if .RowSelected?(newFocus)
				newSelection.Remove(newFocus)
			else
				newSelection = .multiSelect
					? newSelection.Add(newFocus) : Object(newFocus)
			}
		return newSelection
		}

	shiftSelect(newFocus)
		{
		if .shiftAnchor is false
			.shiftAnchor = .focused
		inc = .shiftAnchor < newFocus ? 1: -1
		newSelection = Object(newFocus)
		for (row = .shiftAnchor; row isnt newFocus; row += inc)
			newSelection.Add(row)
		return newSelection
		}

	RowSelected?(row)
		{
		return (.data[row].listrow_flags & .lcf.SELECTED) isnt 0
		}

	RowHighlighted?(row)
		{
		return (.data[row].listrow_flags & .lcf.HIGHLIGHTED) isnt 0
		}

	SelectRow(row)
		{
		if .hideContent is true or not .data.Member?(row) or .RowSelected?(row)
			return

		.setRowFlags(.data[row], .data[row].listrow_flags | .lcf.SELECTED)
		.Act(#SelectRow, row)
		}

	DeSelectRow(row)
		{
		if .hideContent is true or not .data.Member?(row) or not .RowSelected?(row)
			return

		.FinishEdit()
		.setRowFlags(.data[row], .data[row].listrow_flags ^ .lcf.SELECTED)
		.Act(#DeSelectRow, row)
		}

	setRowFlags(rowOb, flags)
		{
		if Record?(rowOb)
			rowOb.PreSet("listrow_flags", flags)
		else
			rowOb.listrow_flags = flags
		}

	FinishEdit()
		{
		if .editor isnt false
			.editor.Return()

		// ensure any drag process being done by the user is ended. This helps prevent
		// situations like AccessControl leaving edit mode while the user is still
		// dragging a record, which can lead to Access not in edit mode when saving
//		.EndDrag()
		}

	ListEditWindow_Commit(col, row, dir, data, valid?, unvalidated_val= '')
		{
		.editor = false
		if not .commitable?()
			return

		if unvalidated_val isnt '' and unvalidated_val is
			.GetInvalidFieldData(.data[row], .columns[col])
			return
		_committing = .columns[col]
		.commit(col, row, data, valid?, unvalidated_val)

		if .allowTab and dir isnt 0
			{
			nextCell = .getNextCell(col, row, dir)
			.Edit(nextCell.col, nextCell.row, dir, canMoveRows:)
			}
		}

	commitable?()
		{
		if .data is false or .data.Empty?()
			return false
		if .Destroyed?() or .readOnly is true
			{
			SuneidoLog("ERRATIC: ListEditCommit on destroyed or readonly list",
				params: Record(destroyed: .Destroyed?(), readonly: .readOnly))
			return false
			}
		return true
		}

	commit(col, row, data, valid?, unvalidated_val)
		{
		// from invalid to valid
		corrected? = valid? is true and
			'' isnt .GetInvalidFieldData(.data[row], .columns[col])
		if .SetInvalidFieldData(.data[row], .columns[col], unvalidated_val)
			.Send('List_InvalidDataChanged', .data[row])
		valueChanged? = data isnt .data[row][.columns[col]]
		.updateCellValid(col, row, valid?, unvalidated_val)
		if valueChanged? and false isnt .Send("List_CellEdit", :col, :row, :data, :valid?)
			{
			.data[row][.columns[col]] = data
			.Act(#UpdateDataCell, :row, :col,
				newCell: .formatting.PaintCell(.columns[col], 0, 0, 0, 0, .data[row]))
			.Send("List_CellValueChanged", :col, :row, :data)
			.Send("List_AfterEdit", :col, :row, :data, :valid?)
			}
		.Send('List_AfterEditWindowCommit', :col, :row, :data, :valid?,
			valueChanged?: valueChanged? or unvalidated_val isnt '' or corrected?)
		}

	updateCellValid(col, row, valid?, unvalidated_val)
		{
		if not .trackValid
			return

		if not valid? and unvalidated_val isnt ''
			.AddInvalidCell(col, row)
		else
			.RemoveInvalidCell(col, row)
		}

	SetInvalidFieldData(rec, field, val)
		{
		if not rec.Member?('List_InvalidData')
			rec.List_InvalidData = Object().Set_default('')
		if rec.List_InvalidData[field] isnt val
			{
			rec.List_InvalidData[field] = val
			return true
			}
		return false
		}

	AllowContextOnly(allow)
		{
		Assert(Boolean?(allow))
		.contextOnly = allow
		.Act(#AllowContextOnly, allow)
		}

	getNextCell(col, row, amt)
		{
		if row is 0 and col is 0 and amt < 0 //no cells left to get
			return Object(col: 0, row: 0)

		numCols = .GetNumCols()
		numRows = .GetNumRows()
		forever
			{
			col += amt
			row = Max(0, row + (col / numCols).Floor())
			col %= numCols
			if (col < 0)
				col += numCols
			if (row >= numRows)		// will give non-existent row at end for edit
				return Object(col: 0, row: numRows)

			if .skip?(row, col)
				continue

			return Object(:col, :row)
			}
		}

	skip?(row, col)
		{
		// skip deleted rows
		if .data[row].GetDefault("listrow_deleted", false) is true
			return true
		if .hiddenInvalidCol?(col, row)
			return true
		if true is .Send('List_Tabover?', .GetCol(col))
			return true
		return false
		}

	hiddenInvalidCol?(col, row)
		{
		minWidth = 5
		return (.GetColWidth(col) < minWidth and
			not (.data[row].Member?("list_invalid_cells") and
				.data[row].list_invalid_cells.Member?(.columns[col])) and
			.widths[1..].Max() >= minWidth)
		}

	AddInvalidCell(col, row)
		{
		rec = .data[row]
		Assert(.columns.Member?(col))
		field = .columns[col]
		if not rec.Member?("list_invalid_cells")
			rec.list_invalid_cells = Object()
		rec.list_invalid_cells[field] = rec[field]
		.RepaintRow(row)
		}
	HasInvalidCell?(record, member)
		{
		return record.Member?("list_invalid_cells") and
			record.list_invalid_cells.Member?(member)
		}
	RowHasInvalidCell?(record)
		{
		return record.Member?("list_invalid_cells") and
			not record.list_invalid_cells.Empty?()
		}
	InvalidCellValue(record, member)
		{
		return record.list_invalid_cells[member]
		}
	RemoveInvalidCell(col, row)
		{
		rec = .data[row]
		Assert(.columns.Member?(col))
		if (not rec.Member?("list_invalid_cells"))
			return
		rec.list_invalid_cells.Delete(.columns[col])
		.RepaintRow(row)
		}

	AddHighlight(row, color = false, rowData = false)
		{
		// color is added to the colors/brushes if not already present
		// the order in which colors are used is preserved for sorting/grouping
		// the default highlightcolor is at index 0 !
		// returns the index number of the highlightcolor/brush
		if color is false
			cidx = 0
		else if false is cidx = .highlight_colors.Find(color)
			{
			cidx = .highlight_colors.Size()
			.highlight_colors.Add(color)
			}
		// if row is false, just add color to force order for sorting/grouping
		if row isnt false
			{
			.addHeightLightFlag(.data[row], cidx)
			.Act(#AddHighlightRow, row, .highlight_colors[cidx])
			}
		// rowData is used to set a highlight prior to row construction
		else if rowData isnt false
			.addHeightLightFlag(rowData, cidx)
		return cidx
		}
	addHeightLightFlag(rec, cidx)
		{
		.setRowFlags(rec,
			LOWORD(rec.listrow_flags) | .lcf.HIGHLIGHTED + (cidx << 16)) /*= add first hex
				as color index*/
		}
	GetHighlighted()
		{
		highlight = Object()
		for (row in .data.Members())
			if .RowHighlighted?(row)
				highlight.Add(row)
		return highlight
		}
	HighlightValues(member, values, color = false, sortHighlight = false,
		group/*unused*/ = false)
		{	// allways add new colors, so these are eventually set for sorting/grouping
		cidx = color is false ? 0 : .AddHighlight(false, color)
		for (row in .data.Members())
			if values.Has?(.data[row][member])
				{
				.addHeightLightFlag(.data[row], cidx)
				.Act(#AddHighlightRow, row, .highlight_colors[cidx])
				}
		Assert(sortHighlight is: false, msg: '.SortHighlight is not implemented')
		}
	ClearHighlight(row = false)
		{
		if row isnt false
			{
			Assert(.data.Member?(row))
			if .RowHighlighted?(row)
				{
				.setRowFlags(.data[row],
					.data[row].listrow_flags & (0xffff - .lcf.HIGHLIGHTED))
				.Act(#RemoveHighlightRow, row)
				}
			}
		else
			for (row in .data.Members())
				.ClearHighlight(row)
		}

	RepaintRow(row)
		{
		if .hideContent is true
			return
		Assert(.data.Member?(row))
		.rowsToRepaint.AddUnique(row)
		.Defer(.repaintRows, uniqueID: 'RepaintRows')
		}

	getter_rowsToRepaint()
		{
		return .rowsToRepaint = Object()
		}

	repaintRows()
		{
		.rowsToRepaint.Each(.repaintRow)
		.rowsToRepaint = Object()
		}

	repaintRow(row)
		{
		if .hideContent is true or not .data.Member?(row)
			return

		.Act(#UpdateData, row, .formatting.PaintRow(.data[row]))
		.Act(.RowSelected?(row) ? #SelectRow : #DeSelectRow, row)
		}

	Repaint()
		{
		.Defer(.repaintAll, uniqueID: 'Repaint')
		}

	repaintAll()
		{
		if .hideContent is true
			return

		batch = Object()
		for row in .data.Members()
			batch[row] = Object(.formatting.PaintRow(.data[row]), .RowSelected?(row))
		.Act(#UpdateBatch, batch)
		}

	DeleteSelection()
		{
		// selected rows are marked for deletion / deleted
		// returns actual number of rows deleted
		return .DeleteRows(@ .getSelection())
		}
	DeleteAll() // WARNING: bad on large lists
		{
		return .DeleteRows(@ .data.Members())
		}
	SetForceDelete()
		{
		// next DeleteRows() DELETES rows
		// wo confimation or notification of controller !
		.forceDelete = true
		}
	DeleteRows(@args)
		{
		// args is a list of row indices to delete
		// rows are marked for deletion / deleted
		// returns actual number of rows deleted
		rowsToDelete = .getDeletedRow(args)
		if not rowsToDelete.Empty?()
			{
			dataRowsDeleted = Object()
			for (row in rowsToDelete)
				{
				dataRowsDeleted.Add(.data[row])
				.data.Delete(row)
				if .focused isnt false and row < .focused
					.focused--
				}
			.checkAndSetFocusedRow()
			.Act(#DeleteRows, rowsToDelete)
			// notify controller of deletions
			if .forceDelete is false
				.Send("List_Deletions", deletions: dataRowsDeleted.Reverse!())
			}
		.forceDelete = false
		return rowsToDelete.Size()
		}
	getDeletedRow(args)
		{
		rowsToDelete = Object()					// build valid list of rows to delete
		for (row in args.Sort!().Reverse!())			// traverse in decending order
			if row >= 0 and row < .GetNumRows()	// filter valid rownumbers
				if not rowsToDelete.Has?(row)		// ignore duplicates
					if true is .forceDelete or
						false isnt .Send("List_DeleteRecord", .data[row])
						rowsToDelete.Add(row)
		return rowsToDelete
		}
	checkAndSetFocusedRow()
		{
		if .focused isnt false
			if .GetNumRows() is 0
				.focused = false
			else if .focused >= .GetNumRows()
				.focused = .GetNumRows() - 1
		}

	header_changed?: false
	HeaderChanged?()
		{
		return .header_changed?
		}
	SetHeaderChanged(status)
		{
		.header_changed? = status
		}
	HeaderResize(col, width)
		{
		// sent by header control
		Assert(.columns.Member?(col))
		minWidth = .Send("Header_TrackMinWidth", .markCol isnt true ? col : col - 1)
		width = Max(width, minWidth)
		.SetColWidth(col, width)
		.header_changed? = true
		}

	HeaderDividerDoubleClick(col, width)
		{
		.HeaderResize(col, Max(width, 50/*=min width*/))
		}

	HeaderClick(col)
		{
		Assert(.columns.Member?(col))
		if .GetNumRows() < 1
			return
		scol = col + 1
		scol = .sortCol is false or .sortCol.Abs() isnt scol ? scol : -.sortCol
		.SortListData(scol)
		.Send("List_AfterSort")
		}

	SortListData(col = false)
		{
		if col isnt false
			.sortCol = col
		if .sortCol is false
			return
		compareFunc = .Send("List_GetCompareFunc", col: .columns[.sortCol.Abs() - 1])
		if compareFunc is false
			return // parent doesn't allow sorting on this column
		if compareFunc is 0			// use default comparison if no callback
			compareFunc = .defCompareFunc	// or parent has no dedicated function
		if .focused isnt false		// mark focused row
			.setRowFlags(.data[.focused], .data[.focused].listrow_flags | .lcf.FOCUSED)
		cmpfn = .sortCol < 0 ? {|x, y| compareFunc(y, x) } : compareFunc
		.sort(cmpfn)
		}

	defCompareFunc(x, y)
		{
		return .formatting.CompareRows(.columns[.sortCol.Abs() - 1], x , y)
		}

	sort(cmpfn)
		{
		for i in .data.Members()
			Record?(.data[i])
				? .data[i].PreSet('listrow_sort', i)
				: .data[i].listrow_sort = i

		.data.Sort!(cmpfn)
		.restoreFocused()
		.setSortIndicator(.sortCol)

		newOrders = Object().AddMany!(0, .data.Size())
		for m, v in .data
			newOrders[v.listrow_sort] = m

		.Act(#ReorderList, newOrders)
		.repaintHeader()
		}

	restoreFocused()
		{
		// try to restore screenposition as much as possible
		if .focused is false
			return
		for (row in .data.Members()) // eventually changed after sorting
			if 0 isnt (.data[row].listrow_flags & .lcf.FOCUSED)
				{
				.setRowFlags(.data[row], .data[row].listrow_flags ^ .lcf.FOCUSED)
				.focused = row
				.Send("List_SelectedRowPositionChanged", selection: Object(.focused))
				.ScrollRowToView(.focused)
				return
				}
		}

	showSortIndicator: true
	setSortIndicator(sortCol)
		{
		if not .showSortIndicator
			return

		sortDown? = sortCol < 0
		if sortCol isnt false
			sortCol = sortCol.Abs() - 1
		offset = .markCol is true ? 1 : 0
		for (idx = 0; idx < .header.GetItemCount(); idx++)
			{
			.header.SetItemSort(idx,
				(idx + offset) isnt sortCol ? false : sortDown? ? -1 : 1)
			}
		}

	signedSort(savedSort)
		{
		if savedSort.Has?(',')
			savedSort = savedSort.BeforeFirst(',')
		sortDown? = false
		if savedSort.Has?('reverse')
			{
			savedSort = savedSort.AfterFirst(' ')
			sortDown? = true
			}
		if false is colNum = .columns.Find(savedSort)
			return false
		sortCol = (colNum + 1) * (sortDown? ? -1 : 1)
		.SetSortCol(sortCol)
		}

	SavedSortIndicator(savedSort)
		{
		if savedSort is false
			return
		.signedSort(savedSort)
		}

	GetSortCol()
		{
		return .sortCol
		}
	GetSort()
		{
		if not Number?(sortCol = .GetSortCol())
			return ''

		if false is fieldname = .GetCol(sortCol.Abs() - 1)
			return ''

		return sortCol.Sign() is -1 ? 'reverse ' $ fieldname : fieldname
		}
	SetSortCol(col)
		{
		.sortCol = col
		.SortListData()
		}

	Header_AllowDrag(col)
		{
		return .Send("List_AllowHeaderReorder", col)
		}
	HeaderReorder(col, newIdx)
		{
		newIdx = Max(newIdx, .markCol is true ? 1 : 0)
		// sent by header control
		Assert(.columns.Member?(col))
		Assert(.columns.Member?(newIdx))
		if col is newIdx or false is .Header_AllowDrag(col)
			return
		oldColumns = .columns.Copy()
		.columns.Delete(col).Add(oldColumns[col], at: newIdx)
		org = .widths[col]
		.widths.Delete(col).Add(org, at: newIdx)
		if .markCol is true
			{
			--col
			--newIdx
			}
		.header.Reorder(col, newIdx)
		if .sortCol isnt false		// adjust index sortcolumn
			{
			sortCol = .columns.Find(oldColumns[.sortCol.Abs() - 1]) + 1
			.sortCol = .sortCol > 0 ? sortCol : -sortCol
			.setSortIndicator(.sortCol)
			}
		.header_changed? = true
		.Act(#UpdateHead, .header.Get(), .markCol, showSortIndicator: .showSortIndicator)
		}


	SetColWidth(col, width)
		{
		Assert(.columns.Member?(col))
		if width is false
			width = .header.GetDefaultColumnWidth(.columns[col])
		if ((col is 0 and .markCol is true) or width - .widths[col] is 0)
			return
		.widths[col] = width
		.header.SetItemWidth(.markCol isnt true ? col : col - 1, width)
		.Act("SetColWidth", col, width)
		}

	SetReadOnly(readOnly, grayOut = true)
		{
		Assert(Boolean?(readOnly) and Boolean?(grayOut))
		if readOnly
			.FinishEdit()
		.readOnly = readOnly
		.Act("SetReadOnly", .readOnly, grayOut)
		}

	GetReadOnly()
		{
		return .readOnly
		}

	ScrollRowToView(row)
		{
		Assert(.data.Member?(row))
		if .hideContent isnt true
			.Act('ScrollRowToView', row)
		}

	ScrollToBottom()
		{
		rows = .GetNumRows()
		if rows > 1
			.ScrollRowToView(rows - 1)
		}

	ScrollColToView(col)
		{
		Assert(.columns.Member?(col))
		if .hideContent isnt true
			.Act('ScrollColToView', col)
		}

	DoWithCurrentVScrollPos(block)
		{
		.Act('SaveVScrollPos')
		block()
		.Act('RestoreVscrollPos')
		}

	CONTEXTMENU(x, y, row, col)
		{
		if .Destroyed?() is true or .contextOnly
			return

		.doWithTempPos(x, y, row, col)
			{
			.SetListFocus()
			.contextMenuFromRow(x, y, row)
			}
		}

	CONTEXTMENU_HEADER(x, y, col)
		{
		if .Destroyed?() is true or .contextOnly
			return

		.doWithTempPos(x, y, -1, col)
			{
			result = .Send("List_HeaderContextMenu", x, y)
			if result is 0
				.buildHeaderContextMenu(x, y)
			}
		}

	buildHeaderContextMenu(x, y)
		{
		contextItems = Object()
		if .resetColumns
			contextItems.Add('Reset Columns')
		if .customizeColumns
			contextItems.Add('Customize Columns...')
		if .sortSaveHandler isnt false
			contextItems.Append(.buildSortMenu())
		if contextItems.Empty?() is false
			ContextMenu(contextItems).ShowCallCascade(this, x, y)
		}

	buildSortMenu()
		{
		sortMenu = Object('Reset Sort to System Default', '')
		if .sortCol is false
			sortMenu.Add(Object(name: 'Set as Default Sort for Current User',
				state: MFS.DISABLED))
		else
			{
			field = .columns[.sortCol.Abs() - 1]
			header = .header.GetHeaderText(field)
			sortMenu.Add(Object(
				name: 'Set (' $ header $ ') as Default Sort for Current User',
				cmd: 'Set as Default Sort'))
			}
		return sortMenu
		}

	On_Context_Reset_Columns()
		{
		if not .resetColumns
			return
		reverse = .sortCol < 0
		fieldName = Number?(.sortCol)
			? .columns[.sortCol.Abs() - 1]
			: false

		if .customizeColumns is true
			UserColumns.Reset(
				this, .columnsSaveName, .origColumns, deletecol: false, load_visible?:)
		else
			.SetColumns(columns: .origColumns, reset:)
		if false is colNum = .columns.Find(fieldName)
			return
		.sortCol = (colNum + 1) * (reverse ? -1 : 1)
		.setSortIndicator(.sortCol)
		.repaintHeader()
		}

	repaintHeader()
		{
		.Act(#UpdateHead, .header.Get(), .markCol, data: 'skip',
			showSortIndicator: .showSortIndicator)
		}

	On_Context_Customize_Columns()
		{
		if 0 is mandatory = .Send('List_MandatoryColumns')
			mandatory = #()
		CustomizeColumnsDialog(.Hwnd, this, .origColumns, .columnsSaveName, mandatory,
			headerSelectPrompt: .headerSelectPrompt)
		}

	On_Context_Reset_Sort_to_System_Default()
		{
		if .sortSaveHandler is false or false is .sortCol
			return

		defaultSort = (.sortSaveHandler)(reset:)
		.signedSort(defaultSort)
		}

	On_Context_Set_as_Default_Sort()
		{
		if .sortSaveHandler is false or false is .sortCol
			return

		fieldName = .columns[.sortCol.Abs() - 1]
		save = ((.sortCol < 0) ? 'reverse ' : '') $ fieldName
		(.sortSaveHandler)(:save)
		}

	GetColFromX(x)
		{
		Assert(.tempPos isnt: false)
		Assert(x is: .tempPos.x)
		return .tempPos.col
		}

	tempPos: false
	doWithTempPos(x, y, row, col, block)
		{
		Assert(.tempPos is: false)
		.tempPos = Object(:x, :y, :row, :col)
		Finally(block, { .tempPos = false })
		}

	contextMenuFromRow(x, y, row)
		{
		sel = .getSelection()
		if not .contextOnly or (sel isnt #() and sel[0] is row)
			{
			if row >= .GetNumRows()
				.ClearSelectFocus()
			else if not .RowSelected?(row)
				.SetSelection(row)	// new selection if not selected
			.Send("List_ContextMenu", x, y)
			}
		}

	SETFOCUS()
		{
		if .editor isnt false
			{
			.editor.Return()						// or end it
			.lastEdit = Date()
			}
		}

	LIST_KILLFOCUS()
		{
		}

	SetListFocus()
		{
		.FinishEdit()		// tell editor window to close
		SetFocus(.Hwnd)
		}

	KEYDOWN(wParam, lParam, ctrl = false, shift = false)
		{
		if .contextOnly
			return 0

		.keydown(wParam, shift, ctrl)
		if .focused isnt false	// if actual rows, .focused should be set !
			{
			// two methods because switch on key is too big for one method...
			.keydown_focused_scroll_keys(wParam, shift, ctrl)
			.keydown_focused_selection_keys(wParam, shift, ctrl)
			}
		.Send('List_KeyDown', :wParam, :lParam)
		return 0
		}

	keydown(wParam, shift, ctrl) // independent of actual rows/focused
		{
		(.keydown_fns[wParam])(:shift, :ctrl)
		}
	getter_keydown_fns()
		{
		ob = Object().Set_default(function () { })
		ob[VK.F2] = // edit focused row
			{
			if .readOnly isnt true and .focused isnt false
				.Edit(0, .focused, 1, canMoveRows: false)
			}
		ob[VK.INSERT] = .insertNewRow
		ob[VK.F5] =	// refresh display
			{
			.Repaint()
			}
		ob[VK.F8] =
			{
			if Suneido.User is 'default'
				Inspect(this)
			}
		ob[VK.SPACE] = .listToggle
		return .keydown_fns = ob // once only
		}
	insertNewRow(shift)
		{
		if .readOnly isnt true
			{
			if .getSelection() is #() or .focused is false
				{
				insertAt = .GetNumRows()
				prevRow = insertAt - 1
				}
			else
				{
				insertAt = shift ? .focused + 1 : .focused
				prevRow = .focused
				}
			newRow = Record()
			if false isnt .wantNewRow(:prevRow, record: newRow)
				{ // parent didn't disallow
				.InsertRow(insertAt, newRow)
				.Send("List_NewRowAdded", insertAt, record: newRow)
				.Edit(0, insertAt, 1, canMoveRows: false)
				}
			}
		}
	listToggle(shift = false, ctrl = false)
		{
		if .checkBoxColumn is false or
			shift isnt false or ctrl isnt false
			return

		sel = .getSelection()
		if sel.Empty?()
			return

		for row in sel
			{
			data = .GetRow(row)
			if false isnt .Send('List_AllowToggle', :data, :row)
				{
				data[.checkBoxColumn] = data[.checkBoxColumn] isnt true
				.Send('List_AfterToggle', :data, :row)
				.RepaintRow(row)
				}
			}
		}

	keydown_focused_scroll_keys(wParam, shift, ctrl) // based on focused row
		{
		(.keydown_focused_fns[wParam])(:shift, :ctrl)
		}
	getter_keydown_focused_fns()
		{
		ob = Object().Set_default(function () { })
//		ob[VK.PRIOR] = .vk_prior
//		ob[VK.NEXT] = .vk_next
		ob[VK.HOME] = .vk_home
		ob[VK.END] = .vk_end
		return .keydown_focused_fns = ob // once only
		}
	vk_home(shift, ctrl)
		{
		if not shift or .multiSelect is false
			.selectFocus(0, ctrl, false, not ctrl)
		else
			.select(-.focused, false, true, true)
		}
	vk_end(shift, ctrl)
		{
		if not shift or .multiSelect is false
			.selectFocus(.GetNumRows() - 1, ctrl, false, not ctrl)
		else
			.select(.GetNumRows() - .focused - 1, false, true, true)
		}

	keydown_focused_selection_keys(wParam, shift, ctrl) // based on focused row
		{
		(.keydown_selection_fns[wParam])(:ctrl, :shift)
		}
	getter_keydown_selection_fns()
		{
		ob = Object().Set_default(function () { })
		ob[VK.UP] = {|ctrl, shift| .vk_up_down(ctrl, shift, dir: -1) }
		ob[VK.DOWN] = {|ctrl, shift| .vk_up_down(ctrl, shift, dir: 1) }
		ob[VK.SPACE] =
			{ |ctrl, shift|
			if ctrl
				.select(0, ctrl, shift, true)
			}
		ob[VK.DELETE] =
			{
			if .readOnly isnt true and false isnt .Send("List_DeleteKeyDown")
				.DeleteSelection()
			}
		.keydown_selection_fns = ob // once only
		}
	vk_up_down(ctrl, shift, dir)
		{
		if ctrl
			.move(dir)
		else
			.select(dir, ctrl, shift and .multiSelect isnt false, not ctrl)
		}
	move(direction)
		{
		sels = .getSelection()
		if .readOnly or sels.Size() isnt 1 or true isnt .Send("List_AllowMove")
			return
		fromRow = sels[0]
		toRow = fromRow + direction
		if toRow < 0 or toRow >= .GetNumRows()
			return
		.data.Swap(fromRow, toRow)
		.Act(#SwapRows, fromRow, toRow)
		.focused = toRow
		.RepaintRow(fromRow)
		.RepaintRow(toRow)
		.ScrollRowToView(toRow)
		.sortCol = false
		.Send("List_Move", fromRow, toRow)
		}

	GetRowHeight()
		{
		return 1
		}

	hideContent: false
	HideContent(hide)
		{
		if hide is true
			.Act(#ClearData)
		.hideContent = hide
		}

	trackValid: false
	Valid?()
		{
		if .trackValid
			{
			// loop through the lines and check if any invalid data
			for rec in .data
				if .RowHasInvalidCell?(rec)
					return false
			}
		return true
		}

	ConfirmDestroy()
		{
		return not .trackValid ? true : .Valid?()
		}

	Destroy()
		{
		if .columnsSaveName isnt false
			UserColumns.Save(.columnsSaveName, this, .origColumns)
		if .trackValid
			.Window.RemoveValidationItem(this)
		super.Destroy()
		}
	}
