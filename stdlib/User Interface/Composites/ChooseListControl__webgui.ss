// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
ChooseControl
	{
	Name: 'ChooseList'
	New(list = false, width = 10, .mandatory = false, allowOther = false,
		set = false, selectFirst = false, .listField = false, .splitValue = ',',
		.listSeparator = ' - ', status = '', .otherListOptions = #(), font = "",
		trim = true, size = "", bgndcolor = "", textcolor = "", tabover = false,
		hidden = false, weight = '', cue = false, .readonly = false)
		{
		super(Object('Field', name: 'Value',
			:width, :mandatory, :status, :font, :size, :weight, :trim,
			:bgndcolor, :textcolor, :tabover, :hidden, :cue, :readonly))
		.listarg = list
		.allowOther? = allowOther is true

		if selectFirst isnt false and .list.Member?(0)
			{
			.Set(.list[0])
			.Send('NewValue', .Get())
			}

		if set isnt false
			{
			.Set(set)
			.Send('NewValue', .Get())
			}
		}

	GetList()
		{
		return .list
		}
	SetList(list)
		{
		.listarg = list
		}
	getter_list()
		{
		return .ListGet(this, .listField, .listarg, .splitValue)
		}
	ListGet(ctrl, listfield, listarg, splitValue = ",")
		{
		if listfield is false
			return listarg
		list = ctrl.Send('GetField', listfield)
		if list is 0 // nothing handled GetField
			return #()
		return ChooseList_ValidDataListFromRules.SplitListFromString(list, splitValue)
		}

	Dialog?()
		{
		return .dialog isnt false
		}
	dialog: false
	On_DropDown(_posInfo = false)
		{
		Assert(posInfo isnt: false)
		if .readonly or false is .InitDropDown() or .list.Size() is 0
			return
		sel = .matchPrefix(.allowOther?, .list, .listSeparator)

		.dialog = true
		value = Dialog(0,
			Object('ChooseListBox', .list, sel, .listSeparator, fieldHwnd: .UniqueId),
			posRect: .Field.UniqueId, border: 0, style: WS.POPUP, backdropDismiss?:)

		.Result(value)
		}

	Result(value)
		{
		.dialog = false
		if value isnt false
			{
			.Set(value)
			.NewValue(.Get())
			.Field.Valid?()
			}
		if .Member?(#Field) and .Member?(#Window)
			.Field.SetFocus()
		}
	matchPrefix(allowOther?, list, listSeparator)
		{
		value = .Field.Get()
		if value is "" or allowOther?
			return false
		if false isnt i = list.FindIf(
			{|val|
			val = String(val)
			if listSeparator isnt ''
				val = val.BeforeFirst(listSeparator)
			val is value
			})
			return i
		if false isnt i = list.FindIf(
			{|val|
			String(val).Prefix?(value)
			})
			return i
		if false isnt i = list.FindIf(
			{|val|
			String(val).Lower().Prefix?(value.Lower())
			})
			return i
		return false
		}
	itemInList?(value)
		{
		// does not loop through all the items when list contains numbers and
		// strings, thus .list.Copy()
		for item in .list.Copy()
			{
			item = String(item)
			if .listSeparator isnt ''
				item = item.BeforeFirst(.listSeparator)
			if item is value
				return true
			}
		if .otherListOptions.Has?(value)
			return true
		return false
		}

	FieldSetFocus()
		{
		super.FieldSetFocus()
		// refresh list on SETFOCUS without affecting dirty flag
		if .listField isnt false
			{
			.Send('DoWithoutDirty')
				{
				.Send("InvalidateFields", Object(.listField))
				}
			}
		}
	FieldKillFocus()
		{
		if not .otherListOptions.Empty?()
			{
			val = .Field.Get()
			if .Field.Dirty?() is true and .otherListOptions.Has?(val)
				.Field.Set("")
			}
		if not .dialog
			.killfocus()
		}
	FieldReturn()
		{
		.killfocus()
		}
	killfocus()
		{
		if .Dirty?() and
			(false isnt i = .matchPrefix(.allowOther?, .list, .listSeparator))
			{
			item = String(.list[i])
			if .listSeparator isnt ''
				item = item.BeforeFirst(.listSeparator)
			.Field.Set(item)
			.Send('NewValue', .Get()) // use Get in case sub-classes override
			.Send('NotifyKilledFocus')
			}
		}

	Valid?()
		{
		value = .Get()
		if .allowOther? and value isnt ''
			return true
		if not .mandatory and value is ''
			return true
		return .itemInList?(value)
		}

	ValidData?(@args)
		{
		if '' is value = args[0]
			return not args.Member?('mandatory') or args.mandatory is false

		if args.Member?('allowOther') and args.allowOther is true
			return true

		listOb = .GetValidDataList(args)
		otherOptions = args.GetDefault('otherListOptions', #())
		return listOb.Has?(value) or otherOptions.Has?(value)
		}

	GetValidDataList(args)
		{
		listOb = Object()
		if args.Member?('list')
			listOb = args.list
		else if args.Member?('listField') and
			Object?(ob = ChooseList_ValidDataListFromRules(args))
			listOb = ob
		else if args.Member?(1) and Object?(args[1])
			listOb = args[1]
		else
			listOb = .ValidationList()

		if '' isnt sep = args.GetDefault('listSeparator', ' - ')
			listOb = listOb.Map({ it.BeforeFirst(sep) })
		return listOb
		}

	ValidationList()
		{
		return Object()
		}

	SelectItem(i)
		{
		if i >= .list.Size()
			i = 0
		.Set(.list[i])
		.Send('NewValue', .Get())
		}

	SelectedItem()
		{
		return .list.Find(.Get())
		}
	}