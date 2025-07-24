// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
Controller
	{
	Xmin: 600
	Ymin: 500
	CallClass(data, sf, title)
		{
		return OkCancel(Object("DateCodeChooseFilter", data, sf), title)
		}

	New(data, .sf)
		{
		.filters = .FindControl("Filters")
		.filters.Set(data)
		}

	Controls: (Scroll (ChooseFilters))
	FieldPrompt_GetSelectFields()
		{
		return .sf
		}

	ChooseFilters_AliasParamSelectFields?()
		{
		return false
		}

	DateControl_ConvertDateCodes()
		{
		return false
		}

	OK()
		{
		if .filters.Get().Empty?()
			return #()
		if true isnt .filters.Valid?()
			return false
		return .filters.Get().Each({ it.check = true })
		}

	Cancel()
		{
		return false
		}
	}
