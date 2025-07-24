// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	Name: 'ParamsChooseList'
	FieldControl: 'Field'
	New(.field, values = false, readonly = true)
		{
		super(Object(.FieldControl, name: 'Value', :readonly), buttonBefore:)
		.values = values is false ? Object() : values
		// SetReadOnly may need done if .FieldControl doesn't handle readonly arg
		.Field.SetReadOnly(true)
		}
	Getter_DialogControl()
		{
		return Object('ParamsChooseListDlg', .field, .values.Copy(), border: 5)
		}
	Get()
		{
		return .values.Copy()
		}

	GetDropDownKeepSizeName()
		{
		return .field
		}

	textLimit: 500
	Set(values)
		{
		.values = Object?(values) ? values.Copy() : Object()
		s = .DisplayValues(.values, .field).Join(', ')
		if s.Size() > .textLimit
			s = s[::(.textLimit - 3/*=ellipsis*/)] $ '...'
		.Field.Set(s)
		}

	DisplayValues(vals, field)
		{
		control = .getControlSpec(field)
		controlClass = Global(control[0].RemoveSuffix('Control') $ 'Control')
		return controlClass.Member?(#DisplayValues)
			? controlClass.DisplayValues(control, vals)
			: vals
		}

	getControlSpec(field)
		{
		if Object?(field) // already a control spec
			return field
		dd = Datadict(field)
		return dd.Member?('SelectControl') ? dd.SelectControl : dd.Control
		}

	// if using "in list", user must specify a list (can't leave it empty)
	Valid?()
		{
		return .Field.Get() isnt ''
		}

	SetReadOnly(readOnly)
		{
		super.SetReadOnly(readOnly)
		if readOnly
			return

		.Field.SetReadOnly(true)
		}
	}
