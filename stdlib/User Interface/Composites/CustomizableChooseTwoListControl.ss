// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
/* example:
	list = #(
		(Prompt: 'Test Column1', Width: 20, Font: 'Arial')
		(Prompt: 'Test Column2', Width: 15, Font: 'Helvetica')
		(Prompt: 'Test Column3', Width: 30, Font: 'Courier')
		)
	CustomizableChooseTwoListControl('Prompt', list)
*/
ChooseTwoListControl
	{
	FieldControl: 'CustomizableChooseTwoListField'
	New(.displayMemberName, .list, width = 50)
		{
		super(.list, :width)
		.availableList = .getAvailableList()
		.selected = Object()
		}

	ProcessFieldValues()
		{
		list = Object()
		for item in .Field.Get().Split(',').Map!(#Trim)
			{
			if false is val = .findOne(.selected, item)
				val = .findOne(.availableList, item)
			if val isnt false
				list.AddUnique(val)
			}
		.Set(list)
		}

	findOne(list, item)
		{
		return list.FindOne({ it[.displayMemberName] is item })
		}

	Get()
		{
		return .selected
		}

	Set(data)
		{
		.selected = data
		super.Set(data.Map({ it[.displayMemberName] }).Join(','))
		}

	Valid?() // called by field control validation
		{
		textList = .Field.Get().Split(',').Map!(#Trim)
		if textList.Size() isnt textList.UniqueValues().Size()
			return false
		for choice in textList
			if not .availableList.HasIf?({ it[.displayMemberName] is choice })
				return false
		return true
		}

	Getter_DialogControl()
		{
		.ProcessFieldValues()
		return Object('TwoListDlg', .availableList, initial_list: .selected, title: '',
			control: 'CustomizableTwoList',	displayMemberName: .displayMemberName)
		}

	getAvailableList()
		{
		return .list
		}
	}