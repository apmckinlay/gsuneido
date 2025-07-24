// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Test_SetFirstSelection()
		{
		// start first, all read
		model = VirtualListModel(.VirtualList_Table)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(35)
		Assert(model.SetFirstSelection() is: 0)
		Assert(model.Selection.GetSelectedRecords() isSize: 1)
		Assert(model.Selection.GetSelectedRecords()[0].num is: 0)

		// start first, not all read
		model = VirtualListModel(.VirtualList_Table)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(10)
		Assert(model.SetFirstSelection() is: 0)
		Assert(model.Selection.GetSelectedRecords() isSize: 1)
		Assert(model.Selection.GetSelectedRecords()[0].num is: 0)

		// start last, not all read
		model = VirtualListModel(.VirtualList_Table, startLast:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(10)
		Assert(model.SetFirstSelection() is: -1)
		Assert(model.Selection.GetSelectedRecords() isSize: 1)
		Assert(model.Selection.GetSelectedRecords()[0].num is: 29)

		// start last, all read
		model = VirtualListModel(.VirtualList_Table, startLast:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(35)
		Assert(model.SetFirstSelection() is: -1)
		Assert(model.Selection.GetSelectedRecords() isSize: 1)
		Assert(model.Selection.GetSelectedRecords()[0].num is: 29)

		// load all if small
		model = VirtualListModel(.VirtualList_Table $ ' sort num', startLast:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(35)
		Assert(model.SetFirstSelection() is: 29)
		Assert(model.Selection.GetSelectedRecords() isSize: 1)
		Assert(model.Selection.GetSelectedRecords()[0].num is: 29)
		}
	}