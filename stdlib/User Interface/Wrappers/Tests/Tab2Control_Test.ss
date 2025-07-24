// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_closeTab?()
		{
		// Check is required to handle the empty Tab2Control override in sujswebgui
		if not Tab2Control.Method?(#Tab2Control_closeTab?)
			return
		calcsMock = Mock()
		calcsMock.When.ImageRect([anyArgs:]).
			Return([left: 5, right: 25, top: 5, bottom: 25])
		tabMock = Mock(Tab2Control)
		tabMock.Tab2Control_closeImage = 'closeImage'
		tabMock.When.closeTab?([anyArgs:]).CallThrough()
		tabMock.Tab2Control_calcClass = calcsMock

		tabMock.Tab2Control_tabItems = #(#(tabName: 'abc', image: false))
		Assert(tabMock.closeTab?(0, 0, 0) is: false)
		calcsMock.Verify.Never().ImageRect([anyArgs:])

		tabMock.Tab2Control_tabItems = #(#(tabName: 'abc', image: 'closeImage'))
		Assert(tabMock.closeTab?(0, 0, 0) is: false)
		calcsMock.Verify.ImageRect([anyArgs:])

		// Test position logic
		Assert(tabMock.closeTab?(0, 27, 0) is: false)
		Assert(tabMock.closeTab?(0, 18, 30) is: false)
		Assert(tabMock.closeTab?(0, 5, 4) is: false)
		Assert(tabMock.closeTab?(0, 15, 17))
		Assert(tabMock.closeTab?(0, 25, 25))
		Assert(tabMock.closeTab?(0, 5, 5))
		Assert(tabMock.closeTab?(0, 25, 5))
		Assert(tabMock.closeTab?(0, 5, 25))
		Assert(tabMock.closeTab?(0, 5, 26) is: false)
		Assert(tabMock.closeTab?(0, 25, 4) is: false)
		}

	Test_baseImage()
		{
		// Check is required to handle the empty Tab2Control override in sujswebgui
		if not Tab2Control.Method?(#Tab2Control_baseImage)
			return
		mock = Mock(Tab2Control)
		mock.When.baseImage([anyArgs:]).CallThrough()
		mock.Tab2Control_closeImage = false
		mock.Tab2Control_staticTabs = #()
		Assert(mock.baseImage('tabName', false) is: false)

		mock.Tab2Control_imageList = false
		Assert(mock.baseImage('tabName', -1) is: false)
		Assert(mock.baseImage('tabName', 0) is: false)

		mock.Tab2Control_imageList = #('image')
		Assert(mock.baseImage('tabName', -1) is: false)
		Assert(mock.baseImage('tabName', 0) is: 'image')

		mock.Tab2Control_closeImage = 'closeImage'
		mock.Tab2Control_staticTabs = #(staticTab)
		Assert(mock.baseImage('tabName', false) is: 'closeImage')
		Assert(mock.baseImage('tabName', 0) is: 'image')
		Assert(mock.baseImage('staticTab', false) is: false)
		}

	Test_hideTab?()
		{
		// Check is required to handle the empty Tab2Control override in sujswebgui
		if not Tab2Control.Method?(#Tab2Control_hideTab?)
			return

		m = Tab2Control.Tab2Control_hideTab?
		i = firstTab = left = availableSpace = 0
		imageWidth = 0

		// Relatively empty state, tab is not hidden
		Assert(m(i, firstTab, left, imageWidth, availableSpace) is: false)

		// Tab left point is past the available space, tab is hidden
		left = 10
		Assert(m(i, firstTab, left, imageWidth, availableSpace))

		// Tab left point is within the available space, tab is not hidden
		availableSpace = 11
		Assert(m(i, firstTab, left, imageWidth, availableSpace) is: false)

		// Tab is before the first visible tab, tab is hidden
		firstTab = 1
		Assert(m(i, firstTab, left, imageWidth, availableSpace))

		// Tab is after the first visible tab, tab is not hidden
		i = 2
		Assert(m(i, firstTab, left, imageWidth, availableSpace) is: false)

		// Tab image width makes tab to large to for the availableSpace, tab is hidden
		imageWidth = 3
		Assert(m(i, firstTab, left, imageWidth, availableSpace))
		}

	Test_tabFullyVisible?()
		{
		// Check is required to handle the empty Tab2Control override in sujswebgui
		if not Tab2Control.Method?(#Tab2Control_tabFullyVisible?)
			return
		fn = Tab2Control.Tab2Control_tabFullyVisible?
		Assert(fn([hide?: false, renderWidth: 100, width: 100]))
		Assert(fn([hide?:, renderWidth: 100, width: 100]) is: false)
		Assert(fn([hide?: false, renderWidth: 10, width: 100]) is: false)
		Assert(fn([hide?: false, renderWidth: 100, width: 10]) is: false)
		}

	Test_attemptDrag_previous()
		{
		// Check is required to handle the empty Tab2Control override in sujswebgui
		if not Tab2Control.Method?(#Tab2Control_attemptDrag)
			return

		mock = .tabDragMock(.tabs(), draggedIdx = 4)

		// Dragging tab within / boarder of the dragged tabs start/end
		Assert(mock.attemptDrag(4, 90) is: false)
		mock.Verify.Never().dragTab?([anyArgs:])
		Assert(mock.attemptDrag(3, 80) is: false)
		mock.Verify.Never().dragTab?([anyArgs:])

		// Dragging over the remaining visible tabs
		draggedIdx = .assertDrag(mock, draggedIdx, 3, 75, true)
		draggedIdx = .assertDrag(mock, draggedIdx, 2, 50, true)
		draggedIdx = .assertDrag(mock, draggedIdx, 1, 15, true)

		// Dragging past the end of the tab bar range, hidden tab is scrolled into view
		.assertDrag(mock, draggedIdx, false, -1, true, 0)
		}

	tabs()
		{
		return Object(
			[tabName: 'Tab 0', hide?:, renderRect: [start: 0, end: 0]],
			[tabName: 'Tab 1', hide?: false, renderRect: [start: 0,   end: 30]],
			[tabName: 'Tab 2', hide?: false, renderRect: [start: 30,  end: 65]],
			[tabName: 'Tab 3', hide?: false, renderRect: [start: 65,  end: 80]],
			[tabName: 'Tab 4', hide?: false, renderRect: [start: 80,  end: 105]],
			[tabName: 'Tab 5', hide?:, renderRect: [start: 0, end: 0]]).Each()
			{
			it.width = it.renderRect.start + it.renderRect.end
			it.renderWidth = it.hide? ? 0 : it.width
			}
		}

	tabDragMock(tabs, draggedIdx)
		{
		mock = Mock(Tab2Control)
		mock.Tab2Control_tabItems = tabs
		mock.Tab2Control_draggedIdx = draggedIdx
		mock.When.Send([anyArgs:]).Do({ })
		mock.When.lastTabIdx().Return(5)
		mock.When.attemptDrag([anyArgs:]).CallThrough()
		mock.When.visibleTabRange().Return([first: 1, last: 4])
		mock.When.On_Previous().Do({ })
		mock.When.On_Next().Do({ })
		return mock
		}

	assertDrag(mock, draggedIdx, i, check, previous?, finalI = false)
		{
		if finalI is false
			finalI = i
		Assert(mock.Tab2Control_draggedIdx is: draggedIdx)
		Assert(mock.attemptDrag(i, check))
		Assert(mock.Tab2Control_draggedIdx is: draggedIdx += previous? ? -1 : 1)
		mock.Verify.dragTab(finalI, previous?)
		return draggedIdx
		}

	Test_attemptDrag_next()
		{
		// Check is required to handle the empty Tab2Control override in sujswebgui
		if not Tab2Control.Method?(#Tab2Control_attemptDrag)
			return

		mock = .tabDragMock(tabs = .tabs(), draggedIdx = 1)

		// Dragging tab within / boarder of the dragged tabs start/end
		Assert(mock.attemptDrag(1, 10) is: false)
		mock.Verify.Never().dragTab?([anyArgs:])
		Assert(mock.attemptDrag(2, 30) is: false)
		mock.Verify.Never().dragTab?([anyArgs:])

		// Dragging tab over the next tab but before it crosses the drag threshold
		Assert(mock.attemptDrag(2, 35) is: false)
		mock.Verify.dragTab?(35, false, tabs[2], tabs[draggedIdx])

		// Dragging the first non-hidden tab
		// ensureVisibleIdx is set to avoid adjusting tab positions
		Assert(mock.Tab2Control_ensureVisibleIdx is: 0)
		draggedIdx = .assertDrag(mock, draggedIdx, 2, 36, false)
		Assert(mock.Tab2Control_ensureVisibleIdx is: 4)

		// Dragging over the remaining visible tabs
		draggedIdx = .assertDrag(mock, draggedIdx, 3, 70, false)
		draggedIdx = .assertDrag(mock, draggedIdx, 4, 90, false)

		// Dragging past the end of the tab bar range, hidden tab is scrolled into view
		.assertDrag(mock, draggedIdx, false, 110, false, 5)
		}
	}
