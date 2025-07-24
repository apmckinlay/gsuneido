// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Test_resize_and_scroll_down()
		{
		model = VirtualListModel(.VirtualList_Table, startLast:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(10)
		Assert(model.Offset is: -9)
		Assert(model.VisibleRows is: 10)

		model.UpdateVisibleRows(20)
		Assert(model.Offset is: -20)
		Assert(model.VisibleRows is: 20)

		model.UpdateOffset(5)
		Assert(model.Offset is: -19) // show empty row

		model.UpdateOffset(5)
		Assert(model.Offset is: -19) // no move

		model.UpdateOffset(-5)
		Assert(model.Offset is: -24)
		}
	}
