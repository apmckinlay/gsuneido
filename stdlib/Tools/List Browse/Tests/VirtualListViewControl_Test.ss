// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Test_collapseAllWhenSetRecsTop()
		{
		fn = VirtualListViewControl.VirtualListViewControl_collapseAllWhenSetRecsTop

		model = VirtualListModel(.VirtualList_Table)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(35)
		vlCtrl = Object()
		vlCtrl.VirtualListViewControl_model = model
		vlCtrl.Eval(fn, #(1, 2), 'num')

		model = VirtualListModel(.VirtualList_Table, useExpandModel?:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(35)
		vlCtrl = Mock(VirtualListViewControl)
		vlCtrl.VirtualListViewControl_model = model
		vlCtrl.When.GetLoadedData().CallThrough()
		vlCtrl.Eval(fn, #(40, 88), 'num') // nothing to set records to top
		vlCtrl.Verify.Never().CollapseAll([anyArgs:])

		vlCtrl.Eval(fn, #(15, 20), 'num') // set records to top but nothing expanded
		vlCtrl.Verify.Never().CollapseAll([anyArgs:])

		model.SetRecordExpanded(0, 3) // "0"
		model.SetRecordExpanded(4, 5) // "1"
		vlCtrl.Eval(fn, #(1, 2), 'num') // set records to top and has rows to collapse
		vlCtrl.Verify.CollapseAll([anyArgs:])
		}

	Test_GetTitle()
		{
		mock = Mock(VirtualListViewControl)
		mock.When.GetTitle().CallThrough()
		mock.When.GetAccessCustomKey().Return(#CustomKey, false)
		mock.Option = #BookOption

		mock.VirtualListViewControl_title = #Title
		Assert(mock.GetTitle() is: #Title)

		mock.VirtualListViewControl_title = false
		Assert(mock.GetTitle() is: #CustomKey)

		Assert(mock.GetTitle() is: #BookOption)
		}

	Test_applySelect()
		{
		mock = Mock
			{
			Name: 'VirtualListViewControl'
			}(VirtualListViewControl)
		mock.Name = 'VirtualListViewControl'
		addons = Object(Addon_VirtualListTopFilters: )
		mock.Addons = AddonManager(mock, addons)
		mock.When.Send([anyArgs:]).Return(true)
		mock.When.UpdateTopFilters([anyArgs:]).Do({ })
		mock.VirtualListViewControl_model = model =
			VirtualListModel(
				.VirtualList_Table, stickyFields: #(field1), observerList: mock)
		.AddTeardownModel(model)
		mock.When.GetModel().Return(model)
		mock.When.GetColumns().Return(#(num))
		mock.Parent = [Window: [Parent: 0]]
		mock.When.FindControl('SelectRepeat').Return(filters = Mock())
		mock.When.SetSelectVals([anyArgs:]).Do({ })
		mock.When.SetWhere([anyArgs:]).Do({ })

		model.UpdateStickyField([field1: 'sticky'], 'field1')
		rec = model.InsertNewRecord()
		Assert(rec.field1 is: 'sticky')

		mock.VirtualListViewControl_grid = grid = Mock()
		filters.When.SelectChanged?().Return(false)
		mock.Select_vals = Object()
		mock.When.Recv([anyArgs:]).CallThrough()
		mock.Recv('Addon_VirtualListTopFilters_applySelect', #(), #())
		grid.Verify.SetFocusedRow([anyArgs:])

		rec = model.InsertNewRecord()
		Assert(rec.field1 is: 'sticky')

		filters.When.SelectChanged?().Return(true)
		mock.Recv('Addon_VirtualListTopFilters_applySelect', #(), #())

		rec = model.InsertNewRecord()
		Assert(rec.field1 is: '') // sticky value should be cleared
		}
	}
