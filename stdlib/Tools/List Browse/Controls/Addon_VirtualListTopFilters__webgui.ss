// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
_Addon_VirtualListTopFilters
	{
	OpenFilters(view)
		{
		filtersWrapper = view.FindControl('select')  // possibly call once ???
		filtersWrapper.Remove(0)
		UserSettings.Put(view.Option $ ' - Split Open', true)
		colModel = view.GetModel().ColModel
		filtersWrapper.Append(Object('SelectRepeat',
			view.GetSelectFields(), view.Select_vals, colModel.GetSelectMgr().Name(),
			option: view.Option, title: view.GetTitle(), fromFilter:,
			selChanged: view.GetDefault('SelectChanged?', false),
			noUserDefaultSelects?: not colModel.UserDefaultSelectEnabled?()))
		split = view.FindControl('VertSplit')
		split.UpdateSplitter()
		if false is split.SetSplitSaveName(view.Option) // no default
			split.MaximizeSecond()
		if .Model.CheckAboveSortLimit?()
			{
			filtersWrapper.FindControl('buttons').
				Insert(0, Object('VirtualListSortWarningButton', .Parent))
			if false isnt warnButton = .FindControl('VirtualListSortWarningButton')
				.setWarnButtonState(warnButton)
			}
		view.AfterTopFilter("open")
		}
	}