// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Addon
	{
	For: VirtualListViewControl
	filterByWhere: ''
	Filterby: false
	ensureFilterValidFn: false
	ExtraLayout()
		{
		filterParams = Object('Record', name: 'VirtualListParams')
		if .Filterby isnt false
			{
			.ensureFilterValidFn = .Filterby.Extract('ensureValidFn', false)
			filterForm = Object('Form')
			for field in .Filterby
				{
				excludeOps = #()
				if Object?(field)
					{
					excludeOps = field.Extract('excludeOps', #())
					field = field[0]
					}
				filterForm.Add(Object('ParamsSelect', field, :excludeOps, group: 0))
				filterForm.Add('nl')
				}
			filterParams.Add(filterForm)
			.Filterby = .Filterby.Flatten()
			}
		return filterParams
		}

	ResetWhere()
		{
		.filterByWhere = .getWhere()
		}

	ExtraWhere()
		{
		return .filterByWhere
		}

	SetFilter(filters = false)
		{
		if filters is false
			return ''
		for field in filters.Members()
			.Parent.Vert.VirtualListParams.SetField(field, filters[field])
		.filterByWhere = .getWhere()
		return .filterByWhere
		}

	getWhere()
		{
		paramsArgs = .Filterby.Copy()
		paramsArgs.data = .Parent.Vert.VirtualListParams.Get()
		if .ensureFilterValidFn isnt false
			(.ensureFilterValidFn)(paramsArgs.data)
		return GetParamsWhere(@paramsArgs)
		}
	}