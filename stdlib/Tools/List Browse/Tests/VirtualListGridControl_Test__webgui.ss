// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
_VirtualListGridControl_Test
	{
	Test_SetModel_clears_selects()
		{
		model = VirtualListModel(.VirtualList_Table, linked?: true)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(10)

		grid = Mock(VirtualListGridControl)
		grid.VirtualListGridControl_header = SuJsListHeader()
		grid.VirtualListGridControl_focusedRow = false
		grid.When.clearSelects().CallThrough()
		grid.Controller = FakeObject(Send: 0)
		grid.When.Act([anyArgs:]).Do({})
		grid.When.CancelAct([anyArgs:]).Do({})
		grid.UniqueId = 0

		grid.Eval(VirtualListGridControl.SetModel, model)
		Assert(model.Selection.GetSelectedRecords() is: #())
		Assert(grid.VirtualListGridControl_focusedRow is: false)

		grid.VirtualListGridControl_focusedRow = 1
		model.Selection.SelectRows(false, false, 0)
		Assert(model.Selection.GetSelectedRecords()[0].num is: 0)

		grid.Eval(VirtualListGridControl.SetModel, model)
		Assert(model.Selection.GetSelectedRecords() is: #())
		Assert(grid.VirtualListGridControl_focusedRow is: false)
		grid.Verify.Act(#DeSelectRow, 1)
		}

	Test_recycle_clears_selects()
		{
		gridClass = VirtualListGridControl
			{
			UniqueId: 0
			VirtualListGridControl_virtualVisibleRows: 10
			Send(@unused) { }
			SetFocus(@unused) { }
			VirtualListGridControl_paintRow(rec)
				{
				return rec
				}
			Act(@unused) { }
			CancelAct(@unused) { }
			}
		_ctrlspec = Object()
		_parent = class
			{
			Controller: (Controller:  class { Send(@unused) { return 0 } })
			Window: ()
			}
		grid = new gridClass

		modelClass = VirtualListModel
			{
			VirtualListModel_limit: 10
			VirtualListModel_segment: 3
			}
		model = modelClass(.VirtualList_Table, enableMultiSelect:)
		.AddTeardownModel(model)

		grid.Eval(VirtualListGridControl.SetModel, model)
		model.Selection.SelectRows(false, false, 0) // select 0
		model.Selection.SelectRows(false, true, 4) 	// select 0 - 4
		Assert(model.Selection.GetSelectedRecords().Size() is: 5)

		// scroll down, release 0 and 1
		grid.Eval(VirtualListGridControl.VirtualListGridComponent_Load, 10)
		Assert(model.Offset is: 10)
		Assert(model.Selection.GetSelectedRecords().Size() is: 0)
		Assert(model.VirtualListModel_data isSize: 18)

		// shift select 5,
		// since VirtualListGridSelection.shiftStart is cleared, only 5 is selected
		model.Selection.SelectRows(false, true, 5)
		Assert(model.Selection.GetSelectedRecords().Size() is: 1)
		Assert(model.Selection.GetSelectedRecords()[0].num is: 5)
		}
	}