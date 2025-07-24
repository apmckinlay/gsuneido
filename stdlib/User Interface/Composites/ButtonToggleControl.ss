// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// replacement for CheckBox that shows state as button pressed or not
Controller
	{
	Name: ButtonToggle
	New(text, readonly = false, tip = false, set = '')
		{
		super(['EnhancedButton', text, classic:, :tip, command: 'Click',
			name: 'button', buttonStyle:, mouseEffect:])
		.Send('Data')
		.button = .FindControl('button')
		.SetReadOnly(readonly)
		.readonly = readonly
		if set isnt ''
			.Set(set)
		}
	Get()
		{
		return .button.Pushed?()
		}
	Set(value)
		{
		.button.Pushed?(value)
		}
	dirty?: false
	Dirty?(dirty = '')
		{
		if dirty isnt ''
			.dirty? = dirty is true
		return .dirty?
		}
	readonly: false
	SetReadOnly(readonly = true)
		{
		if not .readonly
			.button.SetEnabled(not readonly)
		}
	GetReadOnly()
		{
		return .button.GetEnabled()
		}
	On_Click()
		{
		.button.Pushed?(value = not .button.Pushed?())
		.Send('NewValue', value)
		}
	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}