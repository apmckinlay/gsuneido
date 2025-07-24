// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// part of AutoChooseControl
Window
	{
	Xmin: 10
	Ymin: 10
	New(.parent, list)
		{
		super(Object('AutoListBox' list, xmin: 1000, ymin: 1000),
			style: WS.POPUP,
			exStyle: WS_EX.TOOLWINDOW | WS_EX.TOPMOST,
			show: SW.SHOWNA)
		.move()
		.SelectFirst()
		}
	SelectFirst()
		{
		.ListBox.SetCurSel(0)
		}
	SelectLast()
		{
		.ListBox.SetCurSel(.ListBox.Count() - 1)
		}
	Down()
		{
		.ListBox.SetCurSel((.ListBox.GetCurSel() + 1) % .ListBox.Count())
		}
	Up()
		{
		if -1 is cursel = .ListBox.GetCurSel()
			cursel = 0
		count = .ListBox.Count()
		.ListBox.SetCurSel((cursel + count - 1) % count)
		}
	Get()
		{
		return .ListBox.Get()
		}
	move()
		{
		w = .ListBox.GetHorizontalExtent() + 40
		nrows = Min(10, .ListBox.Count())
		h = nrows * .ListBox.GetItemHeight() + 4

		r = .parent.GetListPos()
		x = r.left
		y = r.bottom
		wr = GetWorkArea(r)
		if y + h > wr.bottom
			y -= h + (r.bottom - r.top)
		if x + w > wr.right
			x = wr.right - w

		SetWindowPos(.Hwnd, NULL, x, y, w, h,
			SWP.NOZORDER | SWP.NOACTIVATE)
		}
	MOUSEACTIVATE()
		{
		return MA.NOACTIVATE
		}
	AutoListBox_Click(item)
		{
		.parent.Picked(.ListBox.GetText(item))
		}
	DESTROY()
		{
		.parent.ListClosed()
		super.DESTROY()
		}
	}