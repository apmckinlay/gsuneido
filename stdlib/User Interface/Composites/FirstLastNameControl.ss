// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'FirstLastName'
	ComponentName: 'FirstLastName'
	New(.readonly = false, width = 20, bothMandatory = false, size = '+2',
		weight = 400, heading = false)
		{
		super(.controls(width, bothMandatory, size, heading, weight))
		.first = .Horz.first.CapitalizeWords
		.last = .Horz.last.CapitalizeWords
		.Left = .Horz.first.Left
		.Top = .Horz.Top
		.SetReadOnly(readonly)
		.Send('Data')
		}
	controls(width, bothMandatory, size, heading, weight)
		{
		ctrl = heading is false ? 'Static' : 'Heading'
		return Object('Horz'
			Object('Pair' Object(ctrl, 'First Name')
				Object('CapitalizeWords', :width, :size, :weight,
					mandatory: bothMandatory, readonly: .readonly),
				name: "first")
			'Skip'
			Object('Pair' Object(ctrl, 'Last Name')
				Object('CapitalizeWords', :width, :size, :weight,
					mandatory: bothMandatory, readonly: .readonly),
				name: "last")
			)
		}

	NewValue(value/*unused*/)
		{
		.setNames(.Get())
		.first.KillFocus()
		.last.KillFocus()
		// send AFTER formatting
		.Send('NewValue', .Get())
		}

	Set(x)
		{
		.setNames(x)
		}

	setNames(name)
		{
		// formatting - remove ","
		ob = NameSplit(name, split_on: ',')
		.first.Set(ob.first)
		.last.Set(ob.last)
		}

	Get()
		{
		if .last.Destroyed?() or .first.Destroyed?()
			return ""

		last = .last.Get().Tr(',', ' ').Trim()
		first = .first.Get().Tr(',', ' ').Trim()
		return Opt(last, ', ') $ first
		}

	Dirty?(dirty = "")
		{
		return .first.Dirty?(dirty) or .last.Dirty?(dirty)
		}

	Valid?()
		{
		return .first.Valid?() and .last.Valid?() and
			FieldControl.ValidTextLength?(.Get(), FieldControl.MaxCharacters)
		}

	readonly: false
	SetReadOnly(on = true)
		{
		if .readonly
			return
		.first.SetReadOnly(on)
		.last.SetReadOnly(on)
		}

	HandleTab()
		{
		if .first.HasFocus?()
			{
			.last.SetFocus()
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
