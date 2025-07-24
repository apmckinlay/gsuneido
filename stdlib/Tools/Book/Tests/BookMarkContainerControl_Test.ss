// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
// TAGS: win32
Test
	{
	Test_getMarkRect?()
		{
		m = BookMarkContainerControl.BookMarkContainerControl_getMarkRect?
		Assert(m(5, 100, 10))
		Assert(m(5, 0, 10) is: false)
		Assert(m(5, -1, 10) is: false)
		Assert(m(11, 100, 10) is: false)
		Assert(m(false, 100, 10) is: false)
		}

	Test_highlight?()
		{
		m = BookMarkContainerControl.BookMarkContainerControl_highlight?
		Assert(m(1, 0, 1))
		Assert(m(1, 1, false))
		Assert(m(0, 1, 1) is: false)
		Assert(m(1, 1, 0) is: false)
		}

	Test_depressed?()
		{
		m = BookMarkContainerControl.BookMarkContainerControl_depressed?
		Assert(m(1, 1, 1))
		Assert(m(0, 1, 1) is: false)
		Assert(m(1, 0, 1) is: false)
		Assert(m(1, 1, 0) is: false)
		}

	Test_skipSwap?()
		{
		m = BookMarkContainerControl.BookMarkContainerControl_skipSwap?
		Assert(m(0, 0))
		Assert(m(0, false))
		Assert(m(false, 0))
		Assert(m(false, false))
		Assert(m(0, 1) is: false)
		}

	Test_swapIndex()
		{
		m = BookMarkContainerControl.BookMarkContainerControl_swapIndex
		index1 = 2
		index2 = 3
		Assert(m(false, index1, index2) is: false)
		Assert(m(1, index1, index2) is: 1)
		Assert(m(2, index1, index2) is: 3)
		Assert(m(3, index1, index2) is: 2)
		Assert(m(4, index1, index2) is: 4)
		}
	}