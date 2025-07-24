// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_CloseTabs()
		{
		// Close all tabs / default
		mock = .mockCloseTabs()
		mock.CloseTabs()
		mock.Verify.closeTab(0, skipCollapse?:)
		mock.Verify.closeTab(1, skipCollapse?:)
		mock.Verify.closeTab(2, skipCollapse?:)
		mock.Verify.closeTab(3, skipCollapse?:)
		mock.Verify.closeTab(4, skipCollapse?:)

		// Ignore block dictates that all tabs should be ignored
		mock = .mockCloseTabs()
		mock.CloseTabs({|i| i < 5})
		mock.Verify.Never().closeTab(0, skipCollapse?:)
		mock.Verify.Never().closeTab(1, skipCollapse?:)
		mock.Verify.Never().closeTab(2, skipCollapse?:)
		mock.Verify.Never().closeTab(3, skipCollapse?:)
		mock.Verify.Never().closeTab(4, skipCollapse?:)

		// Ignore block dictates that all even tabs should be ignored
		mock = .mockCloseTabs()
		mock.CloseTabs({|i| i % 2 isnt 0})
		mock.Verify.closeTab(0, skipCollapse?:)
		mock.Verify.Never().closeTab(1, skipCollapse?:)
		mock.Verify.closeTab(2, skipCollapse?:)
		mock.Verify.Never().closeTab(3, skipCollapse?:)
		mock.Verify.closeTab(4, skipCollapse?:)

		// Ignore block dictates that all even tabs should be ignored
		mock = .mockCloseTabs()
		mock.CloseTabs({|i| i % 2 is 0})
		mock.Verify.Never().closeTab(0, skipCollapse?:)
		mock.Verify.closeTab(1, skipCollapse?:)
		mock.Verify.Never().closeTab(2, skipCollapse?:)
		mock.Verify.closeTab(3, skipCollapse?:)
		mock.Verify.Never().closeTab(4, skipCollapse?:)
		}

	mockCloseTabs()
		{
		mock = Mock(ExplorerMultiControl)
		mock.When.CloseTabs([anyArgs:]).CallThrough()
		mock.When.closeTab([anyArgs:]).Return(true)
		mock.When.collapsenode([anyArgs:]).Return(true)
		mock.ExplorerMultiControl_tabsCtrl = class {
			GetTabCount() { return 5 }
			NoCtrls?() { return true } // avoids calling .tree methods
			}
		return mock
		}
	}
