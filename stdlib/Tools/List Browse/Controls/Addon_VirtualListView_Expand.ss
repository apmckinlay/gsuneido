// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Addon_VirtualListViewBase
	{
	VirtualListExpandBar_MouseMove(y)
		{
		if y < .Header.Ymin
			return
		row_num = .Grid.GetRows(y - .Header.Ymin)
		.Grid.ShowExpandButton(row_num)
		}

	VirtualListExpandBarButton_MouseMove(lParam, source)
		{
		rect = GetWindowRect(source.Hwnd)
		pRect = GetWindowRect(.ExpandBar.Hwnd)
		y = (rect.top - pRect.top) + HISWORD(lParam)
		.VirtualListExpandBar_MouseMove(y)
		}

	VirtualListExpand_SwitchToForm(rowNum = false)
		{
		.Grid.SelectRow(.rowNumWithOffset(rowNum))
		.ExpandBar.HideButton()
		.Send('VirtualList_SwitchToForm')
		}

	rowNumWithOffset(rowNum)
		{
		c = OptContribution('VirtualListView', class { })
		if c.Member?('RowNumWithOffset')
			return c.RowNumWithOffset(rowNum)
		return rowNum + .Model.Offset
		}

	GetExpandCtrlAndRecord(source, curFocus = false)
		{
		if false in (.Model, source)
			return false
		if .Model.ExpandModel is false
			return false
		if source.Base?(Window) // redirected from shortcut
			{
			.SelectRecordByFocus(curFocus)
			if false is rec = .GetSelectedRecord()
				return false
			if rec.vl_expanded_rows is ''
				return false
			ctrl = .Model.ExpandModel.GetExpandedControl(rec).ctrl
			}
		else
			{
			if false is ctrl = .Model.ExpandModel.GetControl(source)
				return false
			rec = ctrl.GetControl().Get()
			}
		return Object(:rec, :ctrl)
		}

	ExpandByField(values, field, keepPos? = false, collapse? = false)
		{
		if not values.Empty?()
			.expand(collapse?, keepPos?, {|rec| values.Has?(rec[field]) })
		}

	CollapseAll(keepPos? = false)
		{
		.expand(collapse?:, :keepPos?, block: { true })
		}

	expand(collapse? = false, keepPos? = false, block = false)
		{
		data = .GetLoadedData()
		for (i = data.Size() -1; i >= 0; i--)
			if block(rec: data[i])
				.VirtualListThumb_Expand(i-.Model.Offset,expand: not collapse?, :keepPos?)
		}

	GetExpandedControl(rec)
		{
		if false is layout = .Model.ExpandModel.GetExpandedControl(rec)
			return false
		return layout.ctrl
		}
	}