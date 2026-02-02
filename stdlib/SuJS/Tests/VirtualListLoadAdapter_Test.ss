// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ModelIsReset?()
		{
		mock = Mock(VirtualListLoadAdapter)
		mock.When.ModelIsReset?().CallThrough()
		mock.VirtualListLoadAdapter_model = model = Mock()

		mock.VirtualListLoadAdapter_loadedBottom =
			mock.VirtualListLoadAdapter_loadedTop = false
		Assert(mock.ModelIsReset?() is: false)

		mock.VirtualListLoadAdapter_loadedTop = 1
		Assert(mock.ModelIsReset?() is: false)

		mock.VirtualListLoadAdapter_loadedBottom = 9
		model.When.GetLoadedData().Return(#())
		Assert(mock.ModelIsReset?())

		model.When.GetLoadedData().Return(#(5: [], 1: [], 9: [], 8: []))
		Assert(mock.ModelIsReset?() is: false)

		mock.VirtualListLoadAdapter_loadedBottom = 8
		Assert(mock.ModelIsReset?())

		mock.VirtualListLoadAdapter_loadedBottom = 9
		mock.VirtualListLoadAdapter_loadedTop = 0
		Assert(mock.ModelIsReset?())

		mock.VirtualListLoadAdapter_loadedTop = 1
		Assert(mock.ModelIsReset?() is: false)
		}
	}