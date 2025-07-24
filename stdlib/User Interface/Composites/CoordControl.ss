// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'Coord'
	New(readonly = false, .xPrompt = "Left", .yPrompt = "Top", .tabover = false,
		.mandatory = false)
		{
		super(.layout())
		.xValue = .FindControl('xValue')
		.yValue = .FindControl('yValue')
		.Left = .xValue.Left
		.Top = .Horz.Top
		.SetReadOnly(readonly)
		.readonly = readonly
		.Send('Data')
		}
	layout()
		{
		xctrl = Object('Pair'
			Object('Static', .xPrompt)
			Object('Number' mask: '##.####', rangefrom: 0, rangeto: 10, name: 'xValue'
				mandatory: .mandatory))
		xctrl.tabover = .tabover
		yctrl = Object('Pair'
			Object('Static', .yPrompt)
			Object('Number' mask: '##.####', rangefrom: 0, rangeto: 10, name: 'yValue'
				mandatory: .mandatory))
		yctrl.tabover = .tabover
		return Object('Horz', xctrl, 'Skip', yctrl, overlap:)
		}
	NewValue(value/*unused*/)
		{
		.Send('NewValue', .Get())
		}
	Set(x)
		{
		ob = .SplitCoord(x)

		.xValue.Set(ob.x)
		.yValue.Set(ob.y)
		}
	Get()
		{
		xVal = .xValue.GetUnvalidated()
		yVal = .yValue.GetUnvalidated()
		return String(xVal) $ ',' $ String(yVal)
		}
	SplitCoord(value)
		{
		Assert(String?(value))
		if value is ''
			return Object(x: '', y: '')

		Assert(value has: ',')
		// set default value since '2,'.Split(',') => #(2)
		coord = value.Split(',').Set_default('').Map!({ it.Number?() ? Number(it) : it })
		return Object(x: coord[0], y: coord[1])
		}
	Valid?()
		{
		if not .xValue.Valid?()
			return false
		if not .yValue.Valid?()
			return false

		x = .xValue.Get()
		y = .yValue.Get()
		return .valid?(x, y, .mandatory)
		}
	valid?(x, y, mandatory = false)
		{
		bothNumber? = Number?(x) and Number?(y)
		if mandatory
			return bothNumber?
		return bothNumber? or (x is '' and y is '')
		}
	ValidData?(@args)
		{
		value = args[0]
		ob = .SplitCoord(value)
		return .valid?(ob.x, ob.y, args.GetDefault('mandatory', false))
		}

	Dirty?(dirty = "")
		{
		return .xValue.Dirty?(dirty) or .yValue.Dirty?(dirty)
		}

	readonly: false
	SetReadOnly(on = true)
		{
		if (.readonly)
			return

		.xValue.SetReadOnly(on)
		.yValue.SetReadOnly(on)
		}

	HandleTab()
		{
		if .xValue.HasFocus?()
			{
			.yValue.SetFocus()
			return true
			}
		return false
		}

	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}
