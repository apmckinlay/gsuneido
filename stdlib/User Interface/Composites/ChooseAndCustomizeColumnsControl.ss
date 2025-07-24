// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// data example:  #(field_1, field_2, #(field_3, width: 16))
Controller
	{
	Unsortable: true
	New(.width = 50, .customFieldTable = false, .defaultCols = #(), .additionalCols = #(),
		.logNoPromptAsError = true)
		{
		super(.layout())
		.twoList = .FindControl('twoList')
		.Send('Data')
		}

	layout()
		{
		return Object('Horz',
			Object('CustomizableChooseTwoList' displayMemberName: 'column',
				list: .buildList(), width: .width, name: 'twoList')
			'Skip'
			#(ResetCols))
		}

	buildList()
		{
		.chooseColumns = ChooseColumns(.allColumns(), .logNoPromptAsError)
		return .chooseColumns.AvailableList()
		}

	NewValue(unused = false)
		{
		.Send('NewValue', .Get())
		}

	Get()
		{
		return .chooseColumns.GetSaveData(.twoList.Get())
		}

	Set(data)
		{
		.twoList.Set(.chooseColumns.SetSaveList(data))
		}

	Dirty?(dirty = "")
		{
		return .twoList.Dirty?(dirty)
		}

	Valid?()
		{
		return .twoList.Valid?()
		}

	TwoListValid(data)
		{
		.chooseColumns.ValidList(data)
		}

	On_Reset_Columns()
		{
		if false isnt .Send('SetEditMode')
			{
			.Set(.DefaultColumns())
			.NewValue()
			}
		}

	// Override this function in field definition to specify the extra columns
	AdditionalColumns()
		{
		return .additionalCols
		}

	// Override this function in field definition to add reset columns button and
	// specify the default columns
	DefaultColumns()
		{
		return .defaultCols
		}

	allColumns()
		{
		cols = .DefaultColumns().Copy().MergeUnion(.AdditionalColumns())
		if .customFieldTable isnt false
			cols.MergeUnion(Customizable.GetPermissableFields(.customFieldTable))
		return cols
		}

	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}