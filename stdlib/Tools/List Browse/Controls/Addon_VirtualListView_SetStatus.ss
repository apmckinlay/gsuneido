// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Addon
	{
	For: VirtualListViewControl

	SetStatusBar(statusBar, msg, normal = false, warn = false, invalid = false)
		{
		if statusBar is false
			return
		statusBar.Set(msg, :normal, :warn, :invalid)
		if not (normal or warn or invalid)
			statusBar.SetValid(msg is '')
		}

	RefreshValid(view, statusBar, rec)
		{
		if rec is false
			{
			.SetStatusBar(statusBar, '', normal:)
			return
			}

		view.GetGrid().RepaintRecord(rec)
		if 0 isnt extra = .Send('VirtualList_ExtraValidMsg', rec)
			return .SetStatusBar(@extra.Add(statusBar, at: 0))

		result = .getMsg(rec, view, statusBar)
		.SetStatusBar(statusBar, result.msg, warn: result.warn)
		}

	getMsg(rec, view, statusBar)
		{
		warn = false
		model = view.GetModel()
		if '' is msg = model.EditModel.GetInvalidMsg(rec)
			if '' isnt msg = model.EditModel.GetWarningMsg(rec)
				warn = true

		if rec.New?() and not rec.Member?(model.EditModel.ValidField)
			if statusBar isnt false and statusBar.Get() is ''
				{
				msg = ''
				warn = false
				}

		return Object(:msg, :warn)
		}
	}
