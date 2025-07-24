// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_recycleTab()
		{
		explorerMock = Mock(ExplorerMultiControl)
		explorerMock.When.CloseActiveTab().Do({ })

		libViewMock = Mock()
		libViewMock.Editor = false
		libViewMock.Explorer = explorerMock
		libViewMock.When.CurrentTable().Return(#lib0)
		libViewMock.When.CurrentName().Return(#name0)
		libViewMock.When.CloseTab().Do({ })

		findMock = Mock(FindInLibrariesControl)
		findMock.When.recycleTab([anyArgs:]).CallThrough()
		findMock.When.libView().Return(libViewMock)
		findMock.When.modifiedRec?(#lib0, #name0).Return(false) // will be closed
		findMock.When.modifiedRec?(#lib1, #name1).Return(false) // will be closed
		findMock.When.modifiedRec?(#lib2, #name2).Return(true) // will not be closed

		// First use of Next/Prev
		Assert(findMock.FindInLibrariesControl_nextPrevTab is: false)
		findMock.recycleTab(#lib0, #name0)
		Assert(findMock.FindInLibrariesControl_nextPrevTab is: [lib: #lib0, name: #name0])
		findMock.Verify.Never().libView()

		// Second use of Next/Prev, same library and record
		findMock.recycleTab(#lib0, #name0)
		Assert(findMock.FindInLibrariesControl_nextPrevTab is: [lib: #lib0, name: #name0])
		findMock.Verify.Never().libView()

		// Third use of Next/Prev, new record is selected, Editor is false (no tabs open)
		findMock.recycleTab(#lib1, #name1)
		Assert(findMock.FindInLibrariesControl_nextPrevTab is: [lib: #lib0, name: #name0])
		findMock.Verify.libView()
		libViewMock.Verify.Never().CurrentTable()
		libViewMock.Verify.Never().CurrentName()
		explorerMock.Verify.Never().CloseActiveTab()

		// Fourth use of Next/Prev, new record is selected, Editor has tabs,
		// nextPrevTab is closed and updated
		libViewMock.Editor = true
		findMock.recycleTab(#lib1, #name1)
		Assert(findMock.FindInLibrariesControl_nextPrevTab is: [lib: #lib1, name: #name1])
		findMock.Verify.Times(2).libView()
		libViewMock.Verify.CurrentTable()
		libViewMock.Verify.CurrentName()
		explorerMock.Verify.CloseActiveTab()

		// Fifth use of Next/Prev, new record is selected, Editor has tabs,
		// nextPrevTab is not closed as the active tab is NOT our nextPrevTab.
		// nextPrevTab is updated to the newly opened tab
		libViewMock.Editor = true
		findMock.recycleTab(#lib2, #name2)
		Assert(findMock.FindInLibrariesControl_nextPrevTab is: [lib: #lib2, name: #name2])
		findMock.Verify.Times(3).libView()
		libViewMock.Verify.Times(2).CurrentTable()
		libViewMock.Verify.Times(2).CurrentName()
		explorerMock.Verify.CloseActiveTab()

		// Sixth use of Next/Prev, new record is selected, Editor has tabs,
		// nextPrevTab is not closed as the record is modified.
		// nextPrevTab is updated to the newly opened tab.
		libViewMock.Editor = true
		findMock.recycleTab(#lib0, #name0)
		Assert(findMock.FindInLibrariesControl_nextPrevTab is: [lib: #lib0, name: #name0])
		findMock.Verify.Times(4).libView()
		libViewMock.Verify.Times(3).CurrentTable()
		libViewMock.Verify.Times(3).CurrentName()
		explorerMock.Verify.CloseActiveTab()
		}
	}