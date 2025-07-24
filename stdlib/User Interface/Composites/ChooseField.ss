// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
ChooseControl
	{
	New(field = 'Field', mandatory = false, buttonBefore = false, width = 20,
		font = "", size = "", weight = "", status = "", tabover = false, hidden = false,
		.readonly = false, allowReadOnlyDropDown = false)
		{
		super(.field_control(field, Object(:mandatory, :width, :font, :size, :weight,
			:status, :tabover, :hidden, :readonly)), buttonBefore, allowReadOnlyDropDown)
		}
	field_control(field, options)
		{
		field = String?(field) ? Object(field) : field.Copy()
		field.name = 'Value'
		field.style = ES.MULTILINE |
			WS.CLIPSIBLINGS | (field.Member?('style') ? field.style : 0)
		// do not want to override field's options
		for option in options.Members()
			if not field.Member?(option)
				field[option] = options[option]
		return field
		}

	dialog?: false
	Dialog?()
		{
		return .dialog?
		}
	Title: ''
	On_DropDown()
		{
		if .readonly or .NoData?() or false is posRect = .InitDropDown()
			return
		.dialog? = true
		ctrl = .DialogControl
		border = ctrl.GetDefault(#border, 2)
		// use MainHwnd as parent in case we are in a ListEditWindow so that the parent is
		// NOT the ListEditWindow, otherwise we can run into a variety of problems

		// I would prefer to use OkCancel here instead of OkCancelWrapper,
		// but we need border, posRect and keep_size. (They are what handle ensuring
		// that the pop-up appears under the Choose Field.
		// I would also like to only use OkCancelWrapper here, but there are controls
		// that do not use ok/cancel buttons which we want to leave "as is"
		displayCtrl = ctrl.GetDefault('closeButton?', false) is true
			? ctrl
			: Object(OkCancelWrapper, ctrl)

		result = ToolDialog(.Window.MainHwnd(), displayCtrl, :border, :posRect,
			keep_size: .GetDropDownKeepSizeName(), title: .Title,
			closeButton?: ctrl.GetDefault('closeButton?', false))

		// in Browse, in some cases the focus can be set back to the field
		// and allow it to be destroyed before the dialog is finished.
		// This is because of .Window.MainHwnd being used
		// The following check for Empty can be removed when MainHwnd is removed
		if .Empty?()
			return
		if result isnt false
			.ProcessResults(result)
		else
			.ReprocessValue()

		.dialog? = false
		SetFocus(.Field.Hwnd)
		// don't know why we need this, or why it has to be .Window.Hwnd
		// if you don't have it, sometimes the field ends up underneath
		// this is also in KeyControl
		SetWindowPos(.Window.Hwnd, HWND.TOP, 0, 0, 0, 0, SWP.NOMOVE | SWP.NOSIZE)
		}

	GetDropDownKeepSizeName()
		{
		return .Name
		}

	// to give extra checking in inheriting classes
	NoData?()
		{
		return false
		}

	ProcessResults(result)
		{
		.Set(result)
		.NewValue(.Get())
		}

	ReprocessValue()
		{
		}

	FieldReturn()
		{
		dirty? = .Dirty?()
		.Field.KillFocus()
		if (dirty?)
			.NewValue(.Get())
		}

	Valid?(forceCheck = false)
		{
		return .Field.Valid?(:forceCheck)
		}
	}
