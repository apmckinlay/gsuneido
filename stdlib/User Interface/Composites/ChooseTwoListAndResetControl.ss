// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(.query = '', .exclude = #(), .extraFieldsPlugin = false, .width = 50,
		.mandatory = false, .listField = false, .defaultCols = false)
		{
		super(.layout())
		.chooseField = .FindControl('chooseField')
		.Send('Data')
		}

	layout()
		{
		ob = Object('Horz',
			Object('ChooseFieldsTwoList' , .query, .exclude, .extraFieldsPlugin,
				.width, .mandatory, .listField, name: 'chooseField'))
		if .defaultCols isnt false
			ob.Add('Skip', #(ResetCols))
		return ob
		}

	NewValue(unused = false)
		{
		.Send('NewValue', .Get())
		}

	Get()
		{
		return .chooseField.Get()
		}

	Set(data)
		{
		return .chooseField.Set(data)
		}

	Valid?()
		{
		return .chooseField.Valid?()
		}

	On_Reset_Columns()
		{
		if false isnt .Send('SetEditMode')
			{
			.chooseField.Set(.DefaultColumns())
			.NewValue()
			.chooseField.SetValid(true)
			}
		}

	GetField(field)
		{
		.Send("GetField", field)
		}

	GetSelectedList()
		{
		return .chooseField.GetSelectedList()
		}

	// Override this function in field definition to add reset columns button and
	// specify the default columns
	DefaultColumns()
		{
		return .defaultCols
		}

	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}
