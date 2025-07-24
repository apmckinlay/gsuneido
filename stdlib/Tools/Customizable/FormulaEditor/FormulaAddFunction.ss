// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(fnName, selectFields)
		{
		return OkCancel(Object(this, fnName, selectFields), "Add " $ fnName)
		}

	New(.fnName, .selectFields)
		{
		super(.layout())
		.data = .FindControl('Data')
		}

	layout()
		{
		.fnSpec = GetContributions('FormulaReserved')[.fnName]
		vert = Object('Vert')
		vert.Add(Object('StaticWrap', .fnSpec.desc, xmin: 600))
		vert.Add('Skip')
		prompts = .selectFields.Prompts().Sort!()
		.convertArgs(.fnSpec.args, vert, prompts)
		return Object('Record', vert)
		}

	convertArgs(args, ob, prompts)
		{
		args.Each()
			{
			if it[0] isnt false
				ob.Add(Object('Pair', Object('Static', it[0]),
					.mapArgToControl(it[1], prompts)))
			else
				ob.Add(.mapArgToControl(it[1], prompts))
			}
		}

	mapArgToControl(arg, prompts)
		{
		if arg.Prefix?('field')
			return Object('ChooseList', prompts, name: arg)
		if arg.Prefix?('text')
			return Object(.textDisplay, width: 12, name: arg)
		if arg.Prefix?('unit')
			return Object(.unitChooser, Uom_Conversions.Members(), allowOther:, name: arg)
		if arg.Prefix?('number')
			return Object('Number', mask: false, width: 20, name: arg)
		if arg.Prefix?('checkbox')
			return Object('CheckBox', name: arg)
		if arg.Prefix?('date')
			return Object('ChooseDate', name: arg)
		if arg.Prefix?('fmt')
			return Object('ChooseList', width: 14, listField: 'shortdate_fmt_list',
				set: Settings.Get('ShortDateFormat'), mandatory:, name: arg)
		if arg.Prefix?('if')
			return Object('FormulaIf', prompts, name: arg)
		return Object('Field', name: arg)
		}

	unitChooser: ChooseListControl
		{
		Get()
			{
			return Display(super.Get())
			}
		}

	textDisplay: FieldControl
		{
		Get()
			{
			return Display(super.Get())
			}
		}

	OK()
		{
		if true isnt .data.Valid()
			return false
		value = .data.Get()
		return .formFunctionString(value)
		}

	formFunctionString(values)
		{
		argsList = Object()
		.extractArgs(.fnSpec.args, values, argsList)
		return .fnName $ '( ' $ argsList.Join(', ') $ ' )'
		}

	extractArgs(args, values, list)
		{
		args.Each()
			{
			val = values.GetDefault(it[1], it.GetDefault(#default, ''))
			if it.Member?(#format)
				val = (it.format)(val, record: values)
			list.Add(val)
			}
		}

	DateControl_ConvertDateCodes()
		{
		return false
		}
	}
