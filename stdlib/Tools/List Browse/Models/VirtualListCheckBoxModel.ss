// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
/*
READ ME FIRST
To correctly process the checkmarked records in the list, you need to check for the three
different selection states:

1. selected: you can loop through checked records in memory using
	VirtualListControl.GetCheckedRecords method (calls GetSelectedInfo in this class)
2. all: you will need to read through the query but can assume all records are checked
3. allbut: you will need to read through the query and pass each record to IsSelected
		method to determine if it is checked or not

IMPORTANT: make sure to use the same query the List is using (handling any filters etc)
*/
class
	{
	New(.CheckBoxColumn, .keyField, .checkBoxAmountField)
		{
		Assert(.keyField not in (false, ''),
			msg: 'VirtualList - keyField is required for checkBoxColumn option')
		.selected = Object(state: 'selected', list: Object())
		}

	CheckRecord(rec)
		{
		result = false
		if .IsSelected(rec)
			.UnselectItem(rec)
		else
			{
			.SelectItem(rec)
			result = true
			}
		if rec isnt false
			rec[.CheckBoxColumn] = result
		}

	SelectItem(rec)
		{
		Assert(.selected.state isnt 'all')
		if rec is false
			return

		item = rec[.keyField]
		rec[.CheckBoxColumn] = true
		if .selected.state is 'allbut'
			.selected.list.Delete(item)
		else if .selected.state is 'selected'
			.selected.list[item] = rec
		}

	UnselectItem(rec)
		{
		if rec is false
			return

		item = rec[.keyField]
		if .selected.state is 'all'
			{
			.selected.state = 'allbut'
			.selected.list[item] = rec
			}
		else if .selected.state is 'allbut'
			.selected.list[item] = rec
		else
			.selected.list.Delete(item)

		rec[.CheckBoxColumn] = false
		}

	GetSelectedInfo()
		{
		state = .selected.state
		list = .selected.list.Values()
		return Object(:state, :list)
		}

	IsSelected(rec)
		{
		item = rec[.keyField]
		if .selected.state is 'all'
			return true
		else if .selected.state is 'allbut'
			return not .selected.list.Member?(item)
		else
			return .selected.list.Member?(item)
		}

	SelectAll()
		{
		.selected.state = 'all'
		.selected.list = Object()
		}

	UnselectAll()
		{
		.selected.state = 'selected'
		.selected.list = Object()
		}

	ReloadRecord(oldRec, newRec)
		{
		.selected.list.Replace(oldRec, newRec)
		}

	total: false
	GetSelectedTotal(recalc = false)
		{
		if not recalc and .total isnt false
			return .total
		Assert(.selected.state isnt 'all')
		.total = 0
		for rec in .selected.list
			.total += rec[.checkBoxAmountField]
		return .total
		}

	AutoSelectByAmount(col, data, rec)
		{
		if col isnt .checkBoxAmountField
			return

		if Number(data) is 0 and rec[.CheckBoxColumn] is true
			.UnselectItem(rec)
		if Number(data) isnt 0 and rec[.CheckBoxColumn] isnt true
			.SelectItem(rec)
		}
	}
