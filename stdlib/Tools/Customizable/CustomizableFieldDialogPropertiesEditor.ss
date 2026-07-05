// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	// Must be named like this in order to be found by CustomizableFieldDialog
	Name: 'peditor'

	// Expecting an Object (i.e. #('Extra1', 'Extra2')
	// - will match .Name(s) of the contribution
	IgnoreExtra: false

	New()
		{
		super(.layout())
		.controls = Object()
		for item in .c
			{
			for field in item.GetSetFields
				{
				// not everything listed in GetSetFields is guaranteed to have been
				// created on the screen. (e.g. permissions, etc...) Need to check.
				if false isnt x = .FindControl(field.ctrl)
					.controls[field.label] = x
				}
			}
		}

	layout()
		{
		ie = #()
		if Object?(.IgnoreExtra)
			ie = .IgnoreExtra
		ctrls = .GetControls()
		.c = GetContributions('CustomizableFieldDialogProperties')
		for item in .c
			{
			if ie.Has?(item.Name)
				continue
			(item.Controls)(ctrls)
			}
		return ctrls
		}

	GetControls()
		{
		throw "MUST IMPLEMENT"
		}

	Get()
		{
		x = Object()
		for label in .controls.Members()
			x[label] = .controls[label].Get()
		return x
		}

	Set(object)
		{
		for label in .controls.Members()
			.controls[label].Set(object.GetDefault(label, ''))
		}

	Valid?()
		{
		for item in .c
			{
			if item.Member?('Valid') and false is (item.Valid)(.controls)
				return false
			}
		return true
		}
	}