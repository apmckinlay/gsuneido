// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name:	'RadioButtons'
	CustomizableOptions: #(readonly, tabover)

	New(@args)
		{
		super(.controls(args))
		.Ystretch = args.GetDefault('ystretch', 0)
		.Send("Data")
		.buttons = args.Member?('label') ? .GroupBox.Horz[.dir] : this[.dir]
		if .dir is 'Horz'
			.setHorzXmin(args)

		.mandatory = args.GetDefault('mandatory', false)
		.noInitialValue = args.GetDefault('noInitialValue', false)
		if .mandatory is false and .noInitialValue is true
			ProgrammerError('RadioButton has noInitialValue but is not mandatory')

		// send initial setting to record control if there is no value
		value = .Send("GetField", .Name)
		if not .noInitialValue and (value is '' or value is 0)
			.Picked(.values[0])
		.Top = .buttons.Top
		}
	setHorzXmin(args)
		{
		lastb = .values.Size() - 1
		xmin = .buttons.Xmin + (args.Member?('label') ? 16 /*= label offset*/ : 0)
		if xmin < .Xmin	// if space left, divide between items
			{
			if lastb > 0
				{
				xtra = ((.Xmin - xmin) / lastb).Int()
				for (i = 0; i < lastb; i++)
					.buttons.Group_ctrls[i*2 + 1].Xmin += xtra
				}
			}
		else
			.Xmin = xmin
		}
	controls(args)
		{
		.dir = args.GetDefault('horz', false) ? 'Horz' : 'Vert'
		buttons = Object(.dir)
		skip = args.GetDefault('skip', 6) /*= default skip size*/
		.values = .buildValues(args)
		.addButtons(args, buttons, skip)
		.value = .values[0]
		if args.Member?('label')
			buttons = Object('GroupBox', args.label, Object('Horz' #(Skip 4) buttons))
		return buttons
		}
	buildValues(args)
		{
		evalOptions = args.GetDefault('evalOptions', false)
		return evalOptions is true
			? args.Values(list:).Map(Global)
			: args.Values(list:)
		}
	addButtons(args, buttons, skip)
		{
		for i in args.Members(list:)
			{
			if args[i] is ""
				throw "Empty RadioButton labels not allowed"
			if .dir is 'Horz' and i isnt 0	// have to use Skip, setting Xmin doesn't
				buttons.Add(Object('Skip' skip))	// work with text on the left side
			if .dir is 'Vert' and args.Member?('ymin') and i isnt 0
				buttons.Add(Object('Fill'))
			.addButton(buttons, args, i, str: .values[i])
			}
		}
	addButton(buttons, args, i, str)
		{
		buttons.Add(Object('RadioButton', str, name: "Rb" $ i,
			left: args.GetDefault('left', false)))
		}
	Picked(value)
		{
		if .GetReadOnly()
			return
		.Set(value)
		.Send("NewValue", value)
		}
	Get()
		{ return .value }
	Set(value)
		{
		.value = value
		if value is '' and not .noInitialValue
			return .Picked(.values[0])
		for i in .values.Members()
			.buttons['Rb' $ i].Set(.values[i] is value)
		}
	readonly: false
	SetReadOnly(.readonly)
		{
		super.SetReadOnly(readonly)
		if .readonly is false and .value is "" and not .noInitialValue
			.Picked(.values[0])
		}
	GetReadOnly()
		{
		return .readonly
		}
	SetEnabledChild(i, state)
		{
		.buttons['Rb' $ i].SetEnabled(state)
		}
	Valid?()
		{
		return .valid?(.value, .values, .mandatory)
		}
	valid?(value, values, mandatory)
		{
		return value isnt ''
			? values.Has?(value)
			: mandatory isnt true
		}
	ValidData?(@args)
		{
		values = .buildValues(args)
		value = values.Extract(0)
		return .valid?(value, values, args.GetDefault('mandatory', false))
		}
	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}
