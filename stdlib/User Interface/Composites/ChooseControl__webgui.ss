// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
PassthruController
	{
	ComponentName: "Choose"
	New(field, buttonBefore = false, .allowReadOnlyDropDown = false)
		{
		hidden = Object?(field) and field.GetDefault('hidden', false)
		.Button = .Construct(DropDownButtonControl, :hidden, :allowReadOnlyDropDown)
		.Field = .Construct(field)

		SuRenderBackend().AddOverrideProc(.Field.UniqueId, .fieldproc)
		.children =  Object(.Field, .Button)
		.Send("Data")

		.ComponentArgs = Object(.Field.GetLayout(), .Button.GetLayout(), :buttonBefore,
			:hidden)
		}

	fieldProcOverrideCheck: false
	SetFieldProcOverrideCheck(fn)
		{
		.fieldProcOverrideCheck = fn
		}
	allowOverride?()
		{
		if .fieldProcOverrideCheck is false
			return true
		return (.fieldProcOverrideCheck)()
		}
	fieldproc(uniqueId/*unused*/, event, args)
		{
		if not .allowOverride?()
			return false

		if event is 'KEYDOWN'
			{
			if args[0] in (VK.UP, VK.DOWN)
				{
				.On_DropDown()
				return false
				}
			if args[0] is VK.RETURN
				{
				.FieldReturn()
				return false
				}
			if args[0] is VK.ESCAPE
				{
				.Send('FieldEscape')
				return false
				}
			}

		.handleFocus(event)
		return true
		}

	handleFocus(event)
		{
		if event.Has?('SETFOCUS')
			.FieldSetFocus()
		if event.Has?('KILLFOCUS')
			.FieldKillFocus()
		}

	On_DropDown()
		{
		throw "must be defined by derived class"
		}

	FieldSetFocus()
		{
		.Send('Field_SetFocus')
		}
	FieldKillFocus()
		{
		}
	FieldReturn()
		{
		}

	children: false
	GetChildren()
		{
		return .children isnt false ? .children : Object()
		}

	Set(value)
		{
		.Field.Set(value)
		}

	Get()
		{
		return .Field.Get()
		}

	Dirty?(dirty = "")
		{
		return .Field.Dirty?(dirty)
		}

	NewValue(value /*unused*/)
		{
		// resend newvalue so that this controller becomes the source
		// and thus responsible for Get method.
		.Field.Dirty?(true)
		.Send("NewValue", .Get())
		}

	SetValid(valid)
		{
		.Field.SetValid(valid)
		}

	SetFocus()
		{
		.Field.SetFocus()
		}

	InitDropDown()
		{
		if .Destroyed?() or .dropDownReadOnly()  // already destroyed or readonly
			return false
		SetFocus(.Field.Hwnd)
		// Focus change could cause this field to be destroyed or protected
		if .Destroyed?() or .dropDownReadOnly()
			return false
		return .Field.Hwnd
		}
	dropDownReadOnly()
		{
		if .allowReadOnlyDropDown is true
			return false
		return .GetReadOnly()
		}

	SetFont(font, size)
		{
		.Field.SetFont(font, size)
		}
	SetStatus(status)
		{
		.Field.SetStatus(status)
		}
	SetBgndColor(color)
		{
		.Field.SetBgndColor(color)
		}
	SetTextColor(color)
		{
		.Field.SetTextColor(color)
		}
	SetReadOnly(readOnly)
		{
		.Field.SetReadOnly(readOnly)
		if .Button isnt false
			.Button.SetReadOnly(readOnly)
		}
	GetReadOnly()
		{
		if .Field.GetReadOnly() is false
			return false
		return .Button is false or .Button.GetReadOnly()
		}

	Data()
		{
		// Block the Data message from the Field control
		}

	NoData()
		{
		}

	EditHwnd()
		{
		return .Field.Hwnd
		}

	Edit_ParentValid?()
		{
		return .Valid?()
		}

	RemoveButton() // used by KeyControl
		{
		.children.Remove(.Button)
		.Button.Destroy()
		.Button = false
		}

	Destroy()
		{
		.Send("NoData")
		SuRenderBackend().RemoveOverrideProc(.Field.UniqueId, .fieldproc)
		super.Destroy()
		}
	}
