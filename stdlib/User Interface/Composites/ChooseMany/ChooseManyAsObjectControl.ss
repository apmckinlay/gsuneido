// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	Unsortable: true
	listarg: false
	New(.idField, .displayField, cols, mandatory = false, buttonBefore = false,
		list = false, .listField = false, .delimiter = ', ',
		.saveColName = 'ChooseManyAsObject', width = 20, height = 2,
		.editableListColumns = false, tabover = false, .allowOtherField = "")
		{
		// the control itself isn't read only, we just don't want users typing in the
		// field. However, readonly automatically comes with tabover which means
		// the user cannot tab into the control. Using style WS.TABSTOP at this level
		// to ensure that the user can still tab into this control.
		super(Object('Editor', readonly:,
			style: tabover is false ? WS.TABSTOP : 0, :height),
				mandatory, buttonBefore, :width)
		if .Button.Method?('AddBorder')
			.Button.AddBorder()
		.columns = cols.Copy().Add('choosemany_select', at: 0)
		.ob = Object()
		// if list field is used we want to look up the list everytime the drop down
		// is selected because the results could change depending on other options on
		// the screen, but if a list is passed in then we only want to use that list
		// and don't need to keep re-looking up the list field.
		.listarg = list
		}

	List()
		{
		return .listarg is false ? .GetList() : .listarg
		}

	GetList()
		{
		list = .Send("GetField", .listField)
		return list
		}

	FieldSetFocus()
		{
		super.FieldSetFocus()
		if .listField is false
			return
		.Send('DoWithoutDirty')
			{
			.Send("InvalidateFields", Object(.listField))
			}
		}

	allowOtherField: ""
	GetOtherList()
		{
		if 0 is (allowOther = .Send('GetField', .allowOtherField)) or
			allowOther is ""
			return #()
		return allowOther
		}

	NoData?()
		{
		// allow the drop down if field has a value even if list is empty so users
		// can clean up invalid data
		return .List().Empty?() and .Get().Empty?()
		}

	Getter_DialogControl()
		{
		selected = .Get() is "" ? Object() : .Get()
		list = .GetListItems(.List().DeepCopy(), selected)
		return Object(ChooseManyAsObjectList, list, .columns, .saveColName,
			.editableListColumns)
		}

	ProcessResults(result)
		{
		.SetListResult(result)
		.NewValue(.Get())
		}

	GetListItems(list, selected)
		{
		for item in list
			item.choosemany_select = selected.Has?(item[.idField])
		return list
		}

	SetListResult(result)
		{
		.Set(.setListOb(result, .idField))
		}

	setListOb(result, idField)
		{
		ob = Object()
		for item in result
			if item.choosemany_select is true
				ob.AddUnique(item[idField])
		return ob
		}

	ob: #()
	Get()
		{
		return .ob.Copy()
		}
	Set(val)
		{
		.ob = .setField(.Field, val)
		}

	setField(field, val)
		{
		list = .List()
		ob = Object()
		if not Object?(val)
			{
			field.Set("")
			return ob
			}

		ob = val.Copy()

		field.Set(.ellipseIfNeeded?(.validate(ob, list)))
		return ob
		}

	validate(ob, list)
		{
		setval = ''
		if Object?(ob)
			{
			li = list.Copy().Append(.GetOtherList())
			descs = Object()
			for item in ob.Copy()
				{
				i = li.FindIf(){|x| x[.idField] is item }
				if i isnt false
					{
					desc = li[i]
					descs.Add(.BuildDesc(desc))
					}
				else
					ob.Remove(item)
				}
			setval = descs.Join(.delimiter)
			}
		return setval
		}

	BuildDesc(desc)
		{
		return desc[.displayField]
		}

	// Max display size 1024, left room for ellipse
	displayLimit: 1021
	ellipseIfNeeded?(setval)
		{
		if setval.Size() < .displayLimit
			return setval

		start = setval[.. .displayLimit]
		return start.BeforeLast(.delimiter) $ .delimiter $ ' ...'
		}
	}
