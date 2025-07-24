// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
VirtualListThumbImageButtonControl
	{
	New(.width = false)
		{
		super(@.layout())
		tip = .Controller.Send('VirtualList_EditButtonHasAccelerator?') isnt false
			? "Edit (Alt + E)"
			: "Edit"
		.ToolTip(tip)
		}

	layout()
		{
		.height = .width
		return Object(command: 'Edit', image: 'edit.emf',
			buttonWidth: .width, buttonHeight: .width,
			mouseEffect:, imagePadding: 0.15)
		}

	GetRecord()
		{
		winRec = GetWindowRect(.Hwnd)
		pRec = GetWindowRect(.Parent.Hwnd)
		return .Parent.GetRecordFromY(winRec.top - pRec.top)
		}

	MoveTo(row_num, rowHeight, headerYmin)
		{
		.Resize(0, row_num * rowHeight + 1 + headerYmin, .width, .height)
		}
	}
