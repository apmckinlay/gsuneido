// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Startup(view)
		{
		view.Defer()
			{
			model = view.GetModel()
			if view.IsLinked?() is false and model isnt false
				model.SetFirstSelection()
			if view.FiltersOnTop?()
				view.Addons.Send('LoadSavedFilters', view)
			}
		}

	RowNumWithOffset(rowNum)
		{
		return rowNum
		}

	On_Context_Edit_Field(addon)
		{
		// delay to avoid ListEditWindow from being closed
		// by the SetFocus triggered by closing context menu
		_forceOnBrowser = true
		addon.Defer(addon.ContextEditField)
		}

	On_Context_New(addon)
		{
		// delay to avoid ListEditWindow from being closed
		// by the SetFocus triggered by closing context menu
		_forceOnBrowser = true
		addon.Defer(addon.ContextNew)
		}

	VirtualListGrid_AfterExpand(rec, ctrl, view)
		{
		view.RepaintExpandBar()
		view.Send('VirtualList_AfterExpand', :rec, :ctrl)
		}
	}