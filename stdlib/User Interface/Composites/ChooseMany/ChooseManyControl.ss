// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// TODO: Field should be a custom field that handles validation using the list passed in
ChooseField
	{
	Name: 'ChooseMany'
	Unsortable:,
	New(list = #(), .listField = false, listDesc = #(), .listDescField = false,
		.saveAll = false, .saveNone = false, mandatory = false, allowOther = false,
		allowOtherField = false, width = 20, status = "",
		.text = '', tabover = false, .additionalButtons = #())
		{
		super(Object('ChooseManyField', list, listField, saveAll, saveNone,
			:status, :mandatory, :allowOther, :allowOtherField, :width, :tabover))
		.listarg = list
		.listdescarg = listDesc
		}

	SetList(list)
		{
		.listarg = list
		.Field.SetList(list)
		}

	// TODO don't invalidate every time you get the lists
	// maybe do it in setfocus like ChooseListControl
	// and maybe make it an option - not all list rules need to invalidate

	getter_list()
		{
		if (.listField isnt false)
			.Send("InvalidateFields", Object(.listField))
		return ChooseListControl.ListGet(this, .listField, .listarg)
		}

	getter_list_desc()
		{
		if (.listDescField isnt false)
			.Send("InvalidateFields", Object(.listDescField))
		return ChooseListControl.ListGet(this, .listDescField, .listdescarg)
		}

	Getter_DialogControl()
		{
		return Object('ChooseManyList', .list, .list_desc, .Get(),
			text: .text, name: .Name, additionalButtons: .additionalButtons)
		}

	Set(val)
		{
		if val.Blank?() and .saveNone
			val = "None"
		if .saveAll
			{
			all = ChooseListControl.ListGet(this, .listDescField, .listdescarg).Join(',')
			if val.Trim() is all
				val = "(All)"
			}
		super.Set(val)
		}

	Get()
		{
		val = super.Get()
		return (val.Blank?() and .saveNone) ? "None" : val
		}

	ValidData?(@args)
		{
		value = args[0]

		saveNone = args.GetDefault('saveNone', false)
		if value is ""
			return not saveNone and not args.Member?('mandatory')

		if saveNone and value is 'None'
			return true

		// ChooseMany does not use a listSeparator (list and desc are separate)
		// we don't want the default " - " from ChooseListControl to be used
		// because users can enter it in the list data
		args.listSeparator = "***NO LIST SEPARATOR USED***"
		listOb = ChooseListControl.GetValidDataList(args)
		return .allValsAreInList?(value.Split(','), listOb)
		}

	allValsAreInList?(vals, list)
		{
		return vals.Every?({ list.Has?(it.Trim()) })
		}
	}
