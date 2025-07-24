// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	New(ctrl = false)
		{
		.container = .FindControl(#Vert)
		if ctrl isnt false
			.Append(ctrl)
		.Send(#Data)
		}

	Controls: #('Vert')

	RemoveAll()
		{
		.container.RemoveAll()
		.valCtrl = false
		}

	value: ''
	ctrlName: 'default_value'
	Append(control, value = '')
		{
		if String?(control)
			control = Object(control)
		control.name = value isnt ''
			? value
			: .ctrlName

		.container.Append(control)
		.valCtrl = .container.GetChildren()[0]
		.Set(.value)
		}

	valCtrl: false
	Data() {}
	NoData() {}

	NewValue(.value, source)
		{
		if source.Name isnt .ctrlName
			return
		.Send(#NewValue, .value)
		}

	Get()
		{
		return .value
		}

	Set(.value)
		{
		if .valCtrl isnt false
			// need "try" because current->restore will set the old value to the current
			// control which may be incompatible (suggestion 26903). The control will be
			// updated later.
			try .valCtrl.Set(.value)
		}

	Dirty?(state = "")
		{
		if .valCtrl is false
			return false
		return .valCtrl.Dirty?(state)
		}

	Valid?()
		{
		return .valCtrl is false ? true : .valCtrl.Valid?()
		}

	Destroy()
		{
		.Send(#NoData)
		super.Destroy()
		}
	}
