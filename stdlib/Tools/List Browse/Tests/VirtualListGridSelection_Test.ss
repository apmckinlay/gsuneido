// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Setup()
		{
		super.Setup()
		.table = .MakeTable('(num, field1) key(num)')
		.emptyTable = .MakeTable('(num, field1) key(num)')
		for(i = 0; i < 50; i++)
			QueryOutput(.table, [num: i])

		.teardownModels = Object()
		}

	Test_PageKey_down()
		{
		_model = VirtualListModel(.table)
		_model.UpdateVisibleRows(10)
		.teardownModels.Add(_model)
		selection = VirtualListGridSelection(_model)

		.assertPageKey(selection, #(focusedRow: 0, shift:, up?: false), assertRow: 9)
		.assertPageKey(selection, #(focusedRow: 5, shift:, up?: false), assertRow: 9)
		.assertPageKey(selection, #(focusedRow: 15, shift:, up?: false), assertRow: 24)
		.assertPageKey(selection, #(focusedRow: 48, shift:, up?: false), assertRow: 49)
		.assertPageKey(selection, #(focusedRow: 49, shift:, up?: false), assertRow: 49)
		}

	Test_PageKey_down_middle()
		{
		_model = VirtualListModel(.table)
		_model.UpdateVisibleRows(10)
		.teardownModels.Add(_model)
		selection = VirtualListGridSelection(_model)

		_model.UpdateOffset(20)
		Assert(_model.Offset is: 20)
		.assertPageKey(selection, #(focusedRow: 15, shift:, up?: false), assertRow: 24)
		.assertPageKey(selection, #(focusedRow: 20, shift:, up?: false), assertRow: 29)
		.assertPageKey(selection, #(focusedRow: 25, shift:, up?: false), assertRow: 29)
		.assertPageKey(selection, #(focusedRow: 22, shift:, up?: false), assertRow: 29)
		.assertPageKey(selection, #(focusedRow: 29, shift:, up?: false), assertRow: 38)
		.assertPageKey(selection, #(focusedRow: 35, shift:, up?: false), assertRow: 44)
		.assertPageKey(selection, #(focusedRow: 48, shift:, up?: false), assertRow: 49)
		.assertPageKey(selection, #(focusedRow: 49, shift:, up?: false), assertRow: 49)
		}

	Test_PageKey_down_startLast()
		{
		_model = VirtualListModel(.table, startLast:)
		_model.UpdateVisibleRows(10)
		Assert(_model.Offset is: -9)
		.teardownModels.Add(_model)
		selection = VirtualListGridSelection(_model)

		.assertPageKey(selection, #(focusedRow: -1, shift:, up?: false), assertRow: -1)
		.assertPageKey(selection, #(focusedRow: -5, shift:, up?: false), assertRow: -1)
		.assertPageKey(selection, #(focusedRow: -15, shift:, up?: false), assertRow: -6)
		.assertPageKey(selection, #(focusedRow: -48, shift:, up?: false), assertRow: -39)
		.assertPageKey(selection, #(focusedRow: -49, shift:, up?: false), assertRow: -40)
		.assertPageKey(selection, #(focusedRow: -50, shift:, up?: false), assertRow: -41)
		}

	Test_PageKey_up()
		{
		_model = VirtualListModel(.table)
		_model.UpdateVisibleRows(10)
		.teardownModels.Add(_model)
		selection = VirtualListGridSelection(_model)

		Assert(_model.Offset is: 0)
		.assertPageKey(selection, #(focusedRow: 0, shift:, up?:), assertRow: 0)
		.assertPageKey(selection, #(focusedRow: 5, shift:, up?:), assertRow: 0)
		.assertPageKey(selection, #(focusedRow: 15, shift:, up?:), assertRow: 6)
		.assertPageKey(selection, #(focusedRow: 48, shift:, up?:), assertRow: 39)
		.assertPageKey(selection, #(focusedRow: 49, shift:, up?:), assertRow: 40)
		}

	Test_PageKey_up_startLast()
		{
		_model = VirtualListModel(.table, startLast:)
		_model.UpdateVisibleRows(10)
		Assert(_model.Offset is: -9)
		.teardownModels.Add(_model)
		selection = VirtualListGridSelection(_model)

		.assertPageKey(selection, #(focusedRow: -1, shift:, up?:), assertRow: -9)
		.assertPageKey(selection, #(focusedRow: -5, shift:, up?:), assertRow: -9)
		.assertPageKey(selection, #(focusedRow: -15, shift:, up?:), assertRow: -24)
		.assertPageKey(selection, #(focusedRow: -48, shift:, up?:), assertRow: -50)
		.assertPageKey(selection, #(focusedRow: -49, shift:, up?:), assertRow: -50)
		.assertPageKey(selection, #(focusedRow: -50, shift:, up?:), assertRow: -50)
		}

	Test_selection_when_list_empty()
		{
		model = VirtualListModel(.emptyTable)
		model.UpdateVisibleRows(10)
		.teardownModels.Add(model)
		selection = VirtualListGridSelection(model)
		selection.SelectRows(false, false, 0)
		Assert(selection.GetSelectedRecords() is: #())
		}

	Test_PageKey_down_expand()
		{
		_model = VirtualListModel(.table)
		_model.UpdateVisibleRows(10)
		.teardownModels.Add(_model)
		selection = VirtualListGridSelection(_model)

		_model.SetRecordExpanded(7, 3)

		.assertPageKey(selection, #(focusedRow: 0, shift:, up?: false), assertRow: 9)
		.assertPageKey(selection, #(focusedRow: 7, shift:, up?: false), assertRow: 16)
		}

	assertPageKey(selection, args, assertRow)
		{
		_called = Object(false, assertRow, startLast: args.focusedRow < 0)
		args = args.Copy()
		args.selectRowFn = function(row, shift) {
			if _called.startLast and row >= 0
				{
				_model.UpdateOffset(50)
				return false
				}
			if _called.startLast and row < -50
				{
				_model.UpdateOffset(-50)
				return false
				}
			if not _called.startLast and row < 0
				return false
			if row >= 50
				{
				_model.UpdateOffset(50)
				return false
				}
			_called[0] = true
			Assert(row is: _called[1])
			Assert(shift)
			return true
		}
		selection.PageKey(@args)
		Assert(_called[0])
		}

	Test_UpdateShiftStart()
		{
		update = VirtualListGridSelection.UpdateShiftStart
		mock = Mock(VirtualListGridSelection)
		mock.VirtualListGridSelection_model = Object()
		mock.VirtualListGridSelection_model.Offset = -10
		mock.VirtualListGridSelection_shiftStart = -9
		mock.Eval(update, -8, 2)
		Assert(mock.VirtualListGridSelection_shiftStart is: -11)

		mock.VirtualListGridSelection_model.Offset = -10
		mock.VirtualListGridSelection_shiftStart = -2
		mock.Eval(update, -9, 2)
		// no change below the focused row
		Assert(mock.VirtualListGridSelection_shiftStart is: -2)

		mock.VirtualListGridSelection_model.Offset = 10
		mock.VirtualListGridSelection_shiftStart = 12
		mock.Eval(update, 14, 2)
		// no change above the focused row
		Assert(mock.VirtualListGridSelection_shiftStart is: 12)

		mock.VirtualListGridSelection_model.Offset = 10
		mock.VirtualListGridSelection_shiftStart = 16
		mock.Eval(update, 10, 2)
		Assert(mock.VirtualListGridSelection_shiftStart is: 18)
		}

	Test_selections_and_reload()
		{
		selection = VirtualListGridSelection(false)

		oldRec = [test: 1]
		selection.ReloadRecord(oldRec, newRec = [test: 2])
		Assert(selection.HasSelectedRow?(oldRec) is: false)
		Assert(selection.HasSelectedRow?(newRec) is: false)

		Assert(selection.NotEmpty?() is: false)
		Assert(selection.GetSelectedRecords() is: #())
		Assert(selection.HasSelectedRow?(oldRec) is: false)

		model = VirtualListModel(.table, enableMultiSelect:)
		model.UpdateVisibleRows(10)
		.teardownModels.Add(model)
		selection = model.Selection

		selection.SelectRows(false, false, 0)

		oldRec = [num: 0]
		newRec = [test: 200]
		Assert(selection.NotEmpty?())
		Assert(selection.GetSelectedRecords() is: Object(oldRec))
		Assert(selection.HasSelectedRow?(oldRec))
		Assert(selection.HasSelectedRow?(newRec) is: false)

		selection.ReloadRecord(oldRec, newRec)
		Assert(selection.HasSelectedRow?(oldRec) is: false)
		Assert(selection.HasSelectedRow?(newRec))

		Assert(selection.NotEmpty?())
		Assert(selection.GetSelectedRecords() is: Object(newRec))
		Assert(selection.HasSelectedRow?(oldRec) is: false)
		Assert(selection.HasSelectedRow?(newRec))

		selection.SelectRows(true, false, 1)
		selection.SelectRows(true, false, 2)
		selection.ReloadRecord(newRec, newRec2 = [test: 2000])
		Assert(selection.HasSelectedRow?(newRec) is: false)
		Assert(selection.HasSelectedRow?(newRec2))

		Assert(selection.NotEmpty?())
		Assert(selection.GetSelectedRecords() is: Object([num: 1], [num: 2], newRec2))
		Assert(selection.HasSelectedRow?(newRec) is: false)
		Assert(selection.HasSelectedRow?(newRec2))

		selection.ClearSelect()
		Assert(selection.NotEmpty?() is: false)
		Assert(selection.GetSelectedRecords() is: Object())
		Assert(selection.HasSelectedRow?(newRec) is: false)
		Assert(selection.HasSelectedRow?(newRec2) is: false)

		// using ctrl + shift makes no sense so we disable selection change in that case
		selection.ClearSelect()
		selection.SelectRows(true, true, 1)
		selection.SelectRows(true, true, 2)
		Assert(selection.NotEmpty?() is: false)
		}

	Test_AdjustFocusedRow()
		{
		modelMock = Mock(VirtualListModel)
		modelMock.Offset = 0
		modelMock.VirtualListModel_data = Object()
		modelMock.When.ValidateRow([anyArgs:]).CallThrough()

		mockGridSelection = Mock(VirtualListGridSelection)
		mockGridSelection.When.AdjustFocusedRow([anyArgs:]).CallThrough()
		mockGridSelection.VirtualListGridSelection_model = modelMock

		// focusedRow is currently false
		Assert(mockGridSelection.AdjustFocusedRow(false, 0) is: false)

		// .data is empty
		Assert(mockGridSelection.AdjustFocusedRow(0, 0) is: false)

		// .data has values
		modelMock.VirtualListModel_data = data = Object(
			[row: 0],
			[row: 1],
			[row: 2],
			[row: 3],
			[row: 4],
			[vl_rows: 3],
			[vl_rows: 3],
			[vl_rows: 3],
			[row: 5])
		modelMock.VirtualListModel_curTop = Object(Pos: 0)
		modelMock.VirtualListModel_curBottom = Object(Pos: data.Size())
		modelMock.When.GetRecord([anyArgs:]).CallThrough()
		modelMock.When.ValidateRow([anyArgs:]).CallThrough()
		// Focused row is before curTop.Pos, returns curTop.Pos
		Assert(mockGridSelection.AdjustFocusedRow(-2, -2) is: 0)
		Assert(mockGridSelection.AdjustFocusedRow(-1, -1) is: 0)
		// Focused row is unchanged
		Assert(mockGridSelection.AdjustFocusedRow(0, 0) is: 0)
		Assert(mockGridSelection.AdjustFocusedRow(1, 1) is: 1)
		Assert(mockGridSelection.AdjustFocusedRow(2, 2) is: 2)
		Assert(mockGridSelection.AdjustFocusedRow(3, 3) is: 3)
		// Focused row is an expanded row, returns the last valid row
		Assert(mockGridSelection.AdjustFocusedRow(4, 4) is: 4)
		Assert(mockGridSelection.AdjustFocusedRow(5, 5) is: 4)
		Assert(mockGridSelection.AdjustFocusedRow(6, 6) is: 4)
		Assert(mockGridSelection.AdjustFocusedRow(7, 7) is: 4)
		// Selects the row after the expanded rows
		Assert(mockGridSelection.AdjustFocusedRow(8, 8) is: 8)
		// FocusedRow is after curBottom.Pos, returns curBottom.Pos
		Assert(mockGridSelection.AdjustFocusedRow(9, 9) is: 8)
		Assert(mockGridSelection.AdjustFocusedRow(10, 10) is: 8)

		// .data only has expanded rows (while not plausible, testing outliers handling)
		modelMock.VirtualListModel_data = data = Object(
			[vl_rows: 5],
			[vl_rows: 5],
			[vl_rows: 5],
			[vl_rows: 5],
			[vl_rows: 5])
		modelMock.VirtualListModel_curTop = Object(Pos: 0)
		modelMock.VirtualListModel_curBottom = Object(Pos: data.Size())
		Assert(mockGridSelection.AdjustFocusedRow(0, 0) is: false)
		Assert(mockGridSelection.AdjustFocusedRow(1, 1) is: false)
		Assert(mockGridSelection.AdjustFocusedRow(2, 2) is: false)
		Assert(mockGridSelection.AdjustFocusedRow(3, 3) is: false)
		Assert(mockGridSelection.AdjustFocusedRow(4, 4) is: false)
		Assert(mockGridSelection.AdjustFocusedRow(5, 5) is: false)
		}

	Test_focusRow()
		{
		mock = Mock(VirtualListGridSelection)
		mock.VirtualListGridSelection_model = Object(Offset: 0)
		mock.When.focusRow([anyArgs:]).CallThrough()
		Assert(mock.focusRow(0, 0) is: 0)
		Assert(mock.focusRow(0, 3) is: 0)
		Assert(mock.focusRow(3, 0) is: 2)

		mock.VirtualListGridSelection_model = Object(Offset: -5)
		Assert(mock.focusRow(-3, 3) is: -2)
		Assert(mock.focusRow(-2, 3) is: -1)
		Assert(mock.focusRow(-1, 3) is: -1)
		// Note: 0 is invalid, need to be fixed by VirtualListModel.ValidateRow
		Assert(mock.focusRow(-1, 4) is: 0)
		}

	Teardown()
		{
		for m in .teardownModels
			m.Destroy()
		super.Teardown()
		}
	}
