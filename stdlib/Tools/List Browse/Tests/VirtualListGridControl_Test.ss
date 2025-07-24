// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
VirtualListModelTests
	{
	Test_SetModel_clears_selects()
		{
		model = VirtualListModel(.VirtualList_Table, linked?: true)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(10)

		grid = Mock(VirtualListGridControl)
		grid.VirtualListGridControl_focusedRow = false
		grid.When.clearSelects().CallThrough()
		grid.When.repaintRow([anyArgs:]).Do({})
		grid.Controller = FakeObject(Send: 0)
		grid.VirtualListGridControl_painter = FakeObject(SetModel: true)

		grid.Eval(VirtualListGridControl.SetModel, model)
		Assert(model.Selection.GetSelectedRecords() is: #())
		Assert(grid.VirtualListGridControl_focusedRow is: false)

		grid.VirtualListGridControl_focusedRow = 1
		model.Selection.SelectRows(false, false, 0)
		Assert(model.Selection.GetSelectedRecords()[0].num is: 0)

		grid.Eval(VirtualListGridControl.SetModel, model)
		Assert(model.Selection.GetSelectedRecords() is: #())
		Assert(grid.VirtualListGridControl_focusedRow is: false)
		grid.Verify.repaintRow([anyArgs:])
		}

	Test_recycle_clears_selects()
		{
		_grid = grid = Mock(VirtualListGridControl)
		grid.Controller = FakeObject(Send: 0)
		grid.VirtualListGridControl_painter = FakeObject(SetModel: true)
		grid.VirtualListGridControl_rowHeight = 0
		grid.When.GetClientRect().Return(Rect(0, 0, 0, 0))
		grid.When.scroll([anyArgs:])
		grid.When.Send([anyArgs:])
		grid.When.SetFocus([anyArgs:])

		modelClass = VirtualListModel
			{
			VirtualListModel_limit: 10
			VirtualListModel_segment: 3
			UpdateOffset(offset, saveAndCollapse = false)
				{
				return super.UpdateOffset(offset,
					{ |rec, row_num| _grid.Eval(saveAndCollapse, rec, row_num) })
				}
			}
		model = modelClass(.VirtualList_Table, enableMultiSelect:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(10)

		grid.Eval(VirtualListGridControl.SetModel, model)
		model.Selection.SelectRows(false, false, 0) // select 0
		model.Selection.SelectRows(false, true, 4) 	// select 0 - 4
		Assert(model.Selection.GetSelectedRecords().Size() is: 5)

		// scroll down to 2, no recycle
		grid.Eval(VirtualListGridControl.VirtualListGridControl_vertScroll, 2)
		Assert(model.Offset is: 2)
		Assert(model.Selection.GetSelectedRecords().Size() is: 5)

		// scroll down to 4, recycled 0 - 2
		grid.Eval(VirtualListGridControl.VirtualListGridControl_vertScroll, 2)
		Assert(model.Offset is: 4)
		Assert(model.Selection.GetSelectedRecords().Size() is: 0)

		// shift select 5,
		// since VirtualListGridSelection.shiftStart is cleared, only 5 is selected
		model.Selection.SelectRows(false, true, 5)
		Assert(model.Selection.GetSelectedRecords().Size() is: 1)
		Assert(model.Selection.GetSelectedRecords()[0].num is: 5)
		}
	}
