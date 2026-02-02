// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
WindowBase
	{
	parent: 		false
	row:			false
	col:			false
	ComponentName: 'ListEditWindow'

	New(control, readonly, .col, .row, .parent, custom = false,
		.customFields = false)
		{
		reservation = SuRenderBackend().ReserveAction()
		.createControl(control, readonly, custom)
		args = Object(control: this.GetLayout(), :readonly, :col, :row,
			uniqueId: .UniqueId, parent: parent.UniqueId)
		.DoActivate()
		SuRenderBackend().RecordAction(false, .GetWindowComponentName(), args,
			reservation.at)
		}

	createControl(control, readonly, custom)
		{
		name = .parent.GetCol(.col)
		ctrl = .ctrlProperties(control, name, custom, readonly)

		// needed so the pop up control will inherit .Custom from this window
		// - needed for KeyControl using CustomizableMap
		.Custom = .customFields
		.Ctrl = .Construct(ctrl)
		if (ctrl.Member?('readonly') and
			ctrl.readonly is true and not .Ctrl.GetReadOnly())
			.Ctrl.SetReadOnly(true)

		.Send("Setup")				// call Controller's Setup Method
		dataRow = .parent.GetRow(.row)
		data = dataRow[name]	// get the data
		if not .Ctrl.GetReadOnly() and
			'' isnt invalidVal = ListControl.GetInvalidFieldData(dataRow, name)
			data = invalidVal
		.Ctrl.Set(data)
		}

	ctrlProperties(control, name, custom, readonly)
		{
		field = false
		if 0 is control						// no control set
			{
			field = Datadict(name)
			control = field.Control
			}
		ctrl = Object?(control) ? control.Copy() : Object(control)
		if field isnt false and custom isnt false
			ctrl.Merge(custom)
		ctrl.name = name
		if (readonly)
			ctrl.readonly = true
		return ctrl
		}

	sending?: false
	sendToParent(dir)
		{
		if .sending? is true
			return
		.sending? = true // .Destroy will clear the flag
		if .parent isnt false
			{
			.ClearFocus()			// so control gets KillFocus and validates
			.listCommit(dir)
			}
		.Destroy()
		}
	// triggered by brwoser side keyboard event
	ListEditWindow_SendToParent(dir)
		{
		.sendToParent(dir)
		}

	dirty?: false
	Msg(args)
		{
		// TODO: should maybe only pass on specific messages
		// (like Status and GetField)
		// instead of only stopping certain messages
		msg = args[0]
		if msg isnt 'Data' and msg isnt 'NewValue' and .parent.Member?('Controller')
			{
			if ((msg is "GetField" or msg is "SetField") and
				not .parent.Controller.Base?(BrowseControl))
				return .parent[msg](@+1 args)
			if .parent.Controller.Base?(VirtualListEdit)
				return .parent.Controller[msg](@+1 args)
			if .parent.Controller.Method?(msg)
				return .parent.Controller[msg](@+1 args)
			}
		if msg is 'NewValue'
			.dirty? = true
		return 0
		}
	Return()
		{
		.listCommit()
		}
	listCommit(dir = 0)
		{
		if not this.Member?('Ctrl') or .Ctrl.Empty?() // got destroyed somehow already
			return
		val = .Ctrl.Get()
		valid? = .Ctrl.Valid?()
		readonly = .Ctrl.GetReadOnly()
		dirty? = .dirty?
		unvalidated_val = not valid? and .Ctrl.Method?('GetUnvalidated')
			? .Ctrl.GetUnvalidated()
			: ""
		parent = .parent
		col = .col
		row = .row
		.Destroy()
		parent.ListEditWindow_Commit(col, row, dir, val, valid?,
			:unvalidated_val, :readonly, :dirty?)
		if dir is 0
			parent.SetFocus()
		}

	focus: 0
	ACTIVATE(active)
		{
		if active is false
			{
			.focus = GetFocus()
			.Send("Inactivate")
			}
		else
			{
			SetFocus(.focus not in (false, 0) ? .focus : .UniqueId)
			.Send("Activate")
			}
		}

	On_Cancel(@unused)
		{
		if .Ctrl.GetReadOnly() is true
			.Destroy()
		}

	Destroy()
		{
		if .Destroyed?()
			return
		super.DESTROY()
		}
	}