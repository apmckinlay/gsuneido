// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// maintains data for controls
PassthruController
	{
	Name: "Data"
	New(.ctrl = false, disableFieldProtectRules = false, custom = false,
		.accessControl = false)
		{
		super(.setup(ctrl, disableFieldProtectRules, custom))
		.dirty? = false
		}
	setup(ctrl, disableFieldProtectRules, custom)
		{
		// need to do setup before super New
		.ctrls = Object().Set_default(Object())
		.newctrls = Object()
		.observers = Object()
		.set_observers = Object()
		.data = Record()
		.data.Observer(.Change)
		.disableFieldProtectRules = disableFieldProtectRules
		.Custom = custom
		return ctrl
		}
	GetControlLayout()
		{
		return .ctrl
		}
	HasControl?(name)
		{
		return .ctrls.Member?(name)
		}
	GetControl(name)
		{
		return .ctrls.Member?(name) ? .ctrls[name][0] : false
		}
	AddObserver(fn)
		{
		.observers.Add(fn)
		.data.Observer(fn)
		}
	AddObserversToCurrentRecord()
		{
		.data.Observer(.Change)
		for observer in .observers
			.data.Observer(observer)
		}
	RemoveObserver(fn)
		{
		.observers.Remove(fn)
		.data.RemoveObserver(fn)
		}
	RemoveObserversFromCurrentRecord()
		{
		.data.RemoveObserver(.Change)
		for observer in .observers
			.data.RemoveObserver(observer)
		}
	AddSetObserver(fn)
		{
		.set_observers.Add(fn)
		}
	notifySetObservers()
		{
		for (set_observer in .set_observers)
			(set_observer)()
		}
	RemoveSetObserver(fn)
		{
		.set_observers.Remove(fn)
		}

	// Observer - update controls when data changes
	Change(member)
		{
		if .Destroyed?()
			return

		if member is ""
			ProgrammerError('RecordControl: record has blank member name')

		if member isnt .protectField and not member.Suffix?('__protect')
			{
			if not .ignoreDirty? and .readonly? and .accessControl
				if not Suneido.Member?("RecordNotInEditMode")
					{
					Suneido.RecordNotInEditMode = [:member,
						callstack: FormatCallStack(GetCallStack(limit: 20), levels: 20)]
					SujsAdapter.CallOnRenderBackend(#DumpStatus,
						'Record not in edit mode - ' $ member)
					}
			.dirty? = true
			}

		// check protect rule
		if member is .protectField
			.handle_protects()

		.syncChange(member)
		}

	syncChange(member)
		{
		// connect rules to controls (before ignore)
		// e.g. field__valid => ctrl.valid(.data[field])
		if member.Has?('__') and .ctrls.Member?(m = member.BeforeLast('__'))
			.connectRuleChangeToControl(member, m)
		else if .ignore is member and .ignore_value is .data[member]
			.ignore = ""
		else if .ctrls.Member?(member)
			for ctrl in .ctrls[member]
				ctrl.Set(.data[member])
		}

	connectRuleChangeToControl(member, m)
		{
		method = member.AfterLast('__').Capitalize()
		if (method is "Protect")
			.handle_protects()
		else
			for ctrl in .ctrls[m]
				ctrl[method](.data[member])
		}

	/* messages from controls */

	TabsControl_SelectTab(source)
		{
		dirty? = .dirty?
		for fld in .newctrls
			for ctrl in .ctrls[fld]
				{
				ctrl.Set(.data[fld])
				if ctrl.Method?("SetValid")
					ctrl.SetValid(true)
				}
		.newctrls = Object()
		.dirty? = dirty?
		.handle_protects()
		.notifySetObservers()

		.Send('TabsControl_SelectTab', :source)
		}

	// sent by controls to register
	Data(source)
		{
		fld = source.Name
		.ctrls[fld].Add(source)

		// source.Set(.data[fld])
			// should be able to do this instead of SelectTab
			// but some controls aren't completely set up
			// when they Send("Data")
		.newctrls.Add(fld)
		}
	ControlInRecord?(source)
		{
		if '' is fld = source.Name
			return false
		return .ctrls[fld].FindIf({ Same?(it, source) }) isnt false
		}
	// sent by controls after a new value is entered
	ignore: ""
	ignore_value: ""
	NewValue(value, source)
		{
		newval = false
		if .data[source.Name] isnt value
			newval = true
		// Needed to trigger validation when NumberControl gets invalid data and returns
		// "" from .Get() which looks like the value is not changed (see Suggestion 25234)
		else if source.Valid?() is false
			.dirty? = true
		.ignore_value = value
		.data[.ignore = source.Name] = value
		if .ctrls.Member?(source.Name)
			for ctrl in .ctrls[source.Name]
				{
				ctrl.Dirty?(false)
				if newval and ctrl isnt source
					ctrl.Set(value)
				}
		// needed a message that only gets sent when the value changes.  This will
		// only work when the control changes the value, it will not handle calls to
		// RecordControl.SetField or modifying the record directly.
		if newval
			.Send("Record_NewValue", source.Name, value)
		}
	// sent by controls when they destroy
	NoData(source)
		{
		.ctrls[source.Name].Remove(source)
		if .ctrls[source.Name].Empty?()
			.ctrls.Delete(source.Name)
		}

	/* methods for applications */

	protectField: false
	SetProtectField(field)
		{
		.protectField = field
		}

	SetField(name, value)
		{
		.data[name] = value
		}

	// set fields without affecting the dirty? flag
	SetFieldDefault(name, value)
		{
		.DoWithoutDirty({ .SetField(name, value) })
		}

	ignoreDirty?: false
	DoWithoutDirty(block)
		{
		dirty? = .dirty?
		.ignoreDirty? = true
		try
			block()
		catch (e)
			{
			.ignoreDirty? = false
			.dirty? = dirty?
			throw e
			}
		.ignoreDirty? = false
		.dirty? = dirty?
		}

	GetField(name)
		{
		return .data[name]
		}

	Set(newdata)
		{
		if newdata.Member?("")
			ProgrammerError('RecordControl: Set record has blank member name')

		.RemoveObserversFromCurrentRecord()
		.data = newdata
		.AddObserversToCurrentRecord()

		// set dirty? before setting each control's value because some
		// controls Send NewValue in Set (ex. RadioButtonsControl)
		.dirty? = false
		for name in .ctrls.Members().Copy()
			for ctrl in .ctrls[name]
				{
				ctrl.Set(newdata[name])
				if ctrl.Member?("SetValid")
					ctrl.SetValid(true)
				}
		.handle_protects()
		.newctrls = Object()

		.notifySetObservers()
		}
	Get(excludeHandleFocus = false)
		{
		if not excludeHandleFocus
			.HandleFocus()
		return .data
		}
	GetControlData()
		{
		data = Record()
		for (ctrl in .ctrls)
			data[ctrl[0].Name] = ctrl[0].Get()
		return data
		}
	dirty?: false
	Dirty?(dirty = '')
		{
		Assert(dirty is true or dirty is false or dirty is '')
		.HandleFocus()
		if dirty isnt ''
			.dirty? = dirty
		if dirty is false
			Suneido.Delete(#RecordNotInEditMode)
		return .dirty?
		}
	RecordDirty?(dirty = '')
		{ .Dirty?(dirty) }
	InvalidateFields(fields)
		{
		if Object?(fields)
			for f in fields.UniqueValues()
				.data.Invalidate(f)
		else
			.data.Invalidate(fields)
		}

	handle_protects()
		{
		protect_val = .evalProtectValue()
		if protect_val is true or (String?(protect_val) and protect_val isnt "")
			.handleProtectString()
		else if protect_val is false or protect_val is ""
			.handleProtectByRule()
		else if Object?(protect_val)
			.handleProtectObject(protect_val)
		else
			throw "invalid return type from protect rule"
		}

	evalProtectValue()
		{
		return .readonly?
			? true
			: .protectField is false ? false : .data[.protectField]
		}

	handleProtectString()
		{
		for i in .ctrls.Members().Copy()
			for ctrl in .ctrls[i]
				.setControlReadOnly(ctrl, true)
		}

	handleProtectByRule()
		{
		for field in .ctrls.Members().Copy()
			{
			protect = .fieldProtectedByRule?(field)
			for ctrl in .ctrls[field]
				.setControlReadOnly(ctrl, protect)
			}
		}

	handleProtectObject(protect_val)
		{
		allbut? = protect_val.GetDefault(0, '') is 'allbut'
		for field in .ctrls.Members().Copy()
			{
			protect1 = protect_val.Member?(field) isnt allbut?
			protect2 = .fieldProtectedByRule?(field)
			for ctrl in .ctrls[field]
				.setControlReadOnly(ctrl, protect1 or protect2)
			}
		}

	setControlReadOnly(ctrl, readOnly)
		{
		if ctrl.GetReadOnly() isnt readOnly
			ctrl.SetReadOnly(readOnly)
		}

	fieldProtectedByRule?(field)
		{
		return not .disableFieldProtectRules and .data[field $ "__protect"] is true
		}

	HandleFocus()
		{
		try
			if _inHandleFocus is true
				return // avoid infinite loop
		_inHandleFocus = true
		hwnd = GetFocus()
		ctrl = .Window.HwndMap.GetDefault(hwnd, false)
		if .isControlled?(ctrl) and ctrl.Dirty?()
			{
			// can't use .ClearFocus here because that causes dialogs
			// to select the first control, and if that is a field
			// then the entire contents of the field gets selected
			// this is a problem for things like Refactor_Extract_Method
			// that are modifying fields after every character typed
			// so you end up typing over the first character
			SetFocus(NULL)
			SetFocus(hwnd)
			}
		}
	isControlled?(ctrl)
		{
		while ctrl isnt false
			{
			for c in .ctrls
				if c.Any?({ Same?(it, ctrl) })
					return true
			ctrl = ctrl.GetDefault(#Parent, false)
			}
		return false
		}

	SetAllValid()
		{
		for ob in .ctrls
			for ctrl in ob
				if ctrl.Member?("SetValid")
					ctrl.SetValid(true)
		}
	Valid(forceCheck = false)
		{
		if not .dirty? and not forceCheck
			return true
		invalid_list = Object()
		for m in .ctrls.Members().Copy()
			for ctrl in .ctrls[m]
				.validControl(ctrl, forceCheck, invalid_list)
		return invalid_list.Size() is 0 ? true : "Invalid: " $ invalid_list.Join(", ")
		}

	validControl(ctrl, forceCheck, invalid_list)
		{
		control_valid? = ctrl.Valid?(:forceCheck) or
			(ctrl.GetReadOnly() and .data[ctrl.Name] is '' and
				.isCustomizedMandatory(ctrl))
		if not control_valid?
			invalid_list.Add(PromptOrHeading(ctrl.Name))
		if ctrl.Method?(#SetValid)
			ctrl.SetValid(control_valid?)
		}

	isCustomizedMandatory(ctrl)
		{
		customized = ctrl.GetDefault('Custom', false) // Custom could be false
		if not Object?(customized)
			return false
		customizedOptions = customized.GetDefault(ctrl.Name, false)
		if not Object?(customizedOptions)
			return false
		return customizedOptions.GetDefault('mandatory', false)
		}

	readonly?: false
	SetReadOnly(readonly, hwnd = false)
		{
		.readonly? = readonly is true
		super.SetReadOnly(readonly)
		if readonly isnt true
			.handle_protects()
		else
			{
			if hwnd is false
				hwnd = .Window.Hwnd
			// to clean up white border left around fields
			InvalidateRect(hwnd, NULL, false)
			}
		}
	GetReadOnly()
		{
		return .readonly?
		}
	GetRecordControl()
		{
		// if inside a RecordControl
		// Send('GetRecordControl') will return the record control
		// otherwise, when nothing responds, you'll get 0
		return this
		}

	// used by RepeatControl, requires ctrls to be the same in other
	MoveStateTo(other)
		{
		.RemoveObserversFromCurrentRecord()
		other.RemoveObserversFromCurrentRecord()
		other.RecordControl_data = .data
		other.AddObserversToCurrentRecord()
		other_ctrls = other.RecordControl_ctrls
		for m in .ctrls.Members().Copy()
			{
			list = .ctrls[m]
			other_list = other_ctrls[m]
			for i in list.Members()
				{
				c = list[i]
				oc = other_list[i]
				oc.Set(.data[m])
				oc.Dirty?(c.Dirty?())
				if oc.Method?(#SetValid)
					oc.SetValid(c.Valid?())
				}
			}

		other.RecordControl_dirty? = .dirty?
		// need to re-kick in the protection after the move.
		if .protectField isnt false
			other.Change(.protectField)
		// don't want this record control to still point to the moved data
		// else next Set would remove observers added by other record control
		.data = []
		// NOTE: the calling code must call Set after moves
		}
	}
