// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: 'VirtualListExpandButtons'
	New(.switchToForm)
		{
		super(Object(#WndPane, .layout()))
		.expandBtn = .FindControl('VirtualListExpandButton')
		.switchBtn = .FindControl('VirtualListExpandSwitchToFormButton')
		}

	layout()
		{
		ob = Object(#Horz #VirtualListExpandButton)
		if .switchToForm
			ob.Add(#VirtualListExpandSwitchToFormButton)
		return ob
		}

	expandBtn: false
	EnsureExpandButton(visible = false)
		{
		if visible is false
			{
			if .expandBtn isnt false
				{
				.WndPane.Horz.Remove(0)
				.expandBtn = false
				}
			return
			}

		if .expandBtn is false
			{
			.WndPane.Horz.Insert(0, 'VirtualListExpandButton')
			.expandBtn = .FindControl('VirtualListExpandButton')
			}
		}

	row_num: false
	MoveTo(.row_num, rowHeight, headerYmin, minus, curLeft = false)
		{
		width = 0
		if .expandBtn isnt false
			{
			.expandBtn.SetSize(rowHeight, rowHeight)
			.expandBtn.UpdateButton(row_num, minus)
			width += rowHeight
			}
		if .switchBtn isnt false
			{
			.switchBtn.SetSize(rowHeight, rowHeight)
			if curLeft < ScaleWithDpiFactor(
				VirtualListGridBodyComponent.DistanceToShowSwitchBtn)
				width += rowHeight
			}
		.Resize(0, row_num * rowHeight + headerYmin, width, rowHeight)
		.SetVisible(true)
		}

	MouseOver?()
		{
		return (.expandBtn isnt false and .expandBtn.MouseOver?()) or
			(.switchBtn isnt false and .switchBtn.MouseOver?())
		}

	VirtualListExpand_SwitchToForm()
		{
		if .row_num isnt false
			.Send('VirtualListExpand_SwitchToForm', .row_num)
		}

	Reset()
		{
		.row_num = false
		if .expandBtn isnt false
			.expandBtn.Reset()
		}
	}