// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ChooseManyAsObjectControl
	{
	New(mandatory = false, buttonBefore = false, list = false, listField = false,
		.delimiter = ', ', saveColName = 'ChooseManyAsObject', width = 20,
		.visibleMembers = 'all', .defaultSelectedValue = true, .stringMembers = #())
		{
		super('', '', #(option, value), mandatory, buttonBefore, list,
			listField, .delimiter, saveColName, width, editableListColumns: #(2))
		.options = Object()
		}

	GetList()
		{
		.fullList = super.GetList()
		if .visibleMembers is 'all'
			return .fullList

		.fullList = .fullList.Union(.visibleMembers)
		return .visibleMembers
		}

	GetListItems(list, selected)
		{
		items = Object()
		for item in list
			items.Add(Record(item,
				choosemany_select: selected.Member?(item),
				option: item,
				value: selected.Member?(item) and not Boolean?(selected[item])
					? selected[item] : ''))
		return items
		}

	SetListResult(result)
		{
		selectedItems = Object()
		for item in result
			{
			// if the "option" is set to false, the member will be deleted in Set
			if item.choosemany_select is false
				selectedItems[item.option] = false
			else if item.value isnt ''
				.formatValues(selectedItems, item)
			else
				selectedItems[item.option] = .defaultSelectedValue
			}
		.Set(selectedItems)
		}

	formatValues(selectedItems, item)
		{
		origVal = .origValues.GetDefault(item.option, '')
		if .stringMembers.Has?(item.option)
			{
			selectedItems[item.option] = item.value
			return
			}

		try
			{
			if .number?(origVal, item)
				selectedItems[item.option] = Number(item.value)
			else if .object?(origVal, item)
				{
				selectedItems[item.option] = Object?(item.value)
					? item.value
					: item.value.SafeEval()
				Assert(Object?(selectedItems[item.option]))
				}
			else if Date?(origVal)
				selectedItems[item.option] = Date(item.value)
			else
				selectedItems[item.option] = item.value
			}
		catch (e)
			{
			msg = 'Invalid data type for ' $ item.option
			SuneidoLog('ERROR: (CAUGHT) ' $ e, calls:, caughtMsg: 'user alerted: ' $ msg)
			AlertError(msg)
			}
		}

	number?(origVal, item)
		{
		return Number?(origVal) or (origVal is '' and String?(item.value) and
			item.value.Number?())
		}

	object?(origVal, item)
		{
		return Object?(origVal) or (origVal is '' and String?(item.value) and
			item.value.Prefix?('#('))
		}

	origValues: #()
	Set(selectedItems)
		{
		optionsOb = Object()
		if Object?(selectedItems)
			for option in .List()
				if selectedItems.GetDefault(option, false) isnt false
					optionsOb[option] = selectedItems[option]

		.Field.Set(optionsOb.Members().Join(.delimiter))
		.options = optionsOb
		if .origValues.Empty?()
			{
			.origValues = .options.Copy()
			if .visibleMembers isnt 'all'
				for option in .fullList.Difference(.visibleMembers)
					.origValues[option] = selectedItems[option]
			}
		}

	Get()
		{
		if .visibleMembers is 'all'
			return .options.Copy()
		// only merge in origValues that are not in visibleMembers to allow deleting
		// visible members
		return .options.Copy().MergeNew(.origValues.Copy().DeleteIf(
			{ .visibleMembers.Has?(it) }))
		}
	}