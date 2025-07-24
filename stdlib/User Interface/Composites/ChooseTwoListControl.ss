// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	Unsortable: true
	Name: ChooseTwoList
	FieldControl: false
	New(list = #(), .title = "", mandatory = false, width = 20, .listField = false,
		font = "", readonly = false, .mandatory_list = #(), .noSort = false,
		.delimiter = ',', height = 2)
		{
		super(Object(.FieldControl isnt false ? .FieldControl
			: height is 1 ? 'ChooseTwoListField' : 'ChooseTwoListEditor',
			list, :width, :listField, :readonly, :font, :delimiter, :height),
			:mandatory)
		if .Button.Method?('AddBorder')
			.Button.AddBorder()
		.listarg = list
		}

	Getter_DialogControl()
		{
		selected = .GetSelectedList().Map!(#Trim).Remove("")
		return Object('TwoListDlg', .list, initial_list: selected, title: .title,
			mandatory_list: .mandatory_list, noSort: .noSort, delimiter: .delimiter)
		}

	getter_list()
		{
		if .listField isnt false
			.Send("InvalidateFields", Object(.listField))
		return .build_list(.listField, .listarg)
		}

	build_list(field, list_arg)
		{
		if field is false
			return list_arg
		list = .Send("GetField", field)
		if String?(list)
			list = list.Split(.delimiter)
		return list.Map!(#Trim)
		}

	GetSelectedList()
		{
		return .Get().Split(.delimiter)
		}

	SetList(.listarg)
		{
		.Field.SetList(.list)
		}

	Valid?(forceCheck = false)
		{
		if .mandatory_list.Difference(.GetSelectedList()).NotEmpty?()
			return false

		return super.Valid?(:forceCheck)
		}
	}
