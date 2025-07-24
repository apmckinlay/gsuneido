// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
// TODO: ChooseListControl.SelectItem may not work properly becase it calls .Get()
ChooseListControl
	{
	allowOther?: false
	New(@args)
		{
		super(@.getList(args))
		}

	getList(args)
		{
		.tuples = args[0]
		.mandatory = args.GetDefault(#mandatory, false)
		.keyField = args.GetDefault(#keyField, false)
		.listField = args.GetDefault(#listField, false)
		list = .tuples.Map({ it[0] })
		args[0] = list
		args.allowOther = false
		args.listSeparator = ''
		return args
		}

//	To use ChooseListTuple with dynamic lists (rule listField). Expected format:
//	Object(dynamicType1: #(#('shown value1', 'hidden value1'),
//		#('shown value2', 'hidden value2')),
//		dynamicType2: #('shown value3', 'hidden value3'))
	dynamicTuple: false
	UpdateTuple(list = false)
		{
		.dynamicTuple = true
		if list is false
			list = .GetList()
		if not Object?(list[0])
			return
		.tuples = list
		if .tuples isnt #(#())
			.list = .tuples.Map({ it[0] })
		else
			.list = #()
		.SetList(.list)
		}

	ListGet(ctrl, listfield, listarg, splitValue = ",")
		{
		if .dynamicTuple
			return listarg

		return super.ListGet(ctrl, listfield, listarg, splitValue)
		}

	Set(value)
		{
		if value is ''
			{
			super.Set(value)
			return
			}
		if false is val = .tuples.FindOne({ value is it[0] or .equal?(value, it[1]) })
			return
		super.Set(val[0])
		}

	equal?(val1, val2)
		{
		if .keyField is false
			return val1 is val2
		if not Object?(val1) or not val1.Member?(.keyField)
			return false
		if not Object?(val2) or not val2.Member?(.keyField)
			return false
		return val1[.keyField] is val2[.keyField]
		}

	Get()
		{
		listVal = super.Get()
		if listVal is '' or false is val = .tuples.FindOne({ it[0] is listVal })
			return ''
		return val[1]
		}

	Valid?()
		{
		val = .Field.Get()
		if not .mandatory and val is ''
			return true

		return .tuples.FindOne({ it[0] is val }) isnt false
		}

	FindValueIndex(value)
		{
		return .tuples.FindIf({ it[1] is value})
		}
	}