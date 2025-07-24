// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	New(list, .listField, .saveAll, .saveNone, .mandatory, .allowOther = false,
		.allowOtherField = false, width = 20, style = 0, status = "", tabover = false)
		{
		super(:width, :style, :status, :tabover)
		.listarg = list
		}
	getter_list()
		{
		return ChooseListControl.ListGet(this, .listField, .listarg)
		}

	SetList(list)
		{
		.listarg = list
		}

	Valid?()
		{
		value = .Get()
		if value is ''
			return not .mandatory
		.match_prefix(value)
		return .valid?
		}
	KillFocus()
		{
		val = .match_prefix(.Get())
		// Set marks field as not dirty, so we have to restore the dirty flag
		dirty = .Dirty?()
		.Set(val)
		.Dirty?(dirty)
		}

	match_prefix(val)
		{
		.valid? = true

		if false isnt allOrNone = .allOrNone(val)
			return allOrNone

		displayValues = Object()
		otherValidFields = .allowOtherField isnt false
			? ChooseListControl.ListGet(this, .allowOtherField, false)
			: #()
		validListItems = .list.Copy().MergeUnion(otherValidFields)
		for chosenItem in val.Split(',').Map(#Trim)
			displayValues.AddUnique(.getItem(validListItems, chosenItem))
		return displayValues.Join(',')
		}

	allOrNone(val)
		{
		if val is ""
			return .saveNone is true ? "None" : ""
		if val is "None" and .saveNone is true
			return "None"
		if val is "(All)" and .saveAll is true
			return "(All)"
		return false
		}

	getItem(validListItems, chosenItem)
		{
		if validListItems.Has?(chosenItem)
			return chosenItem

		// only want to autocomplete items that are available in the list
		prefixMatches = .list.Filter({ it.Trim().Lower().Prefix?(chosenItem.Lower()) })

		if prefixMatches.Size() is 1
			return prefixMatches[0]

		.setValid(prefixMatches.Size() is 0 ? .allowOther : false)
		return chosenItem
		}

	setValid(valid)
		{
		.valid? = valid
		}
	}
