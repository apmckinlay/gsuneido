// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
KeyListViewBase
	{
	New(query, columns, saveInfoName = "", prefix = "", prefixColumn = false,
		keys = false, field = "", value = "", checkBoxColumn = false,
		customizeQueryCols = false, optionalRestrictions = #(), .excludeSelect = #())
		{
		super(query, columns, saveInfoName, prefix, :prefixColumn, :keys, :field, :value,
			enableMultiSelect:, :checkBoxColumn, :customizeQueryCols,
			:optionalRestrictions)
		}

	VirtualList_LeftClick(rec, col)
		{
		.Send('VirtualList_LeftClick', rec, col)
		}

	VirtualList_Space()
		{
		.Send('VirtualList_Space')
		}

	VirtualList_ModelChanged()
		{
		.Send('VirtualList_ModelChanged')
		}

	GetCheckedRecords()
		{
		return .GetList().GetCheckedRecords()
		}

	CheckRecordByKeys(x)
		{
		.GetList().CheckRecordByKeys(x)
		}

	UncheckAll()
		{
		.GetList().UncheckAll()
		}

	GetExcludeSelectFields()
		{
		return .excludeSelect
		}
	}
