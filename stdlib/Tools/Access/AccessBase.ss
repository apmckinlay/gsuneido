// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
// abstract base class, derived classes must define:
//	ValidField
//	TitleScrollStatus
//	Save() => true/false
CommandParent
	{
	New(@args)
		{
		super(.setup(@args))
		.Window.AddValidationItem(this)
		.Data = .Vert.TitleScroll.Control
		.Data.AddObserver(.RecordChange)
		.Status = .Vert.Status
		}
	setup(args)
		{
		.observers = Object()
		return args
		}
	RecordChange(member/*unused*/)
		{
		if .invalid?
			.CheckValid()
		}
	invalid?: false
	CheckValid(evalRule? = false)
		{
		if not .Data.Dirty?()
			return true

		status = ""
		if true isnt invalid_fields = .Data.Valid()
			status $= ", " $ invalid_fields

		if (.ValidField isnt false and
			"" isnt (s = evalRule? is true
				? .Data.Get().Eval(Global(('Rule_' $ .ValidField)))
				: .Data.GetField(.ValidField)))
			status $= ", " $ s

		if status is ""
			{
			.invalid? = false
			.Status.SetValid()
			.Status.Set('')
			return true
			}
		else
			{
			if not .invalid?
				Beep()
			.invalid? = true
			.Status.SetValid(false)
			.Status.Set(status[2..])
			return false
			}
		}
	Status(status)
		{
		if .Status.GetValid()
			.Status.Set(status)
		}
	AccessObserver(fn, at = false)
		{
		if at isnt false
			.observers.Add(fn, :at)
		else
			.observers.Add(fn)
		}
	RemoveAccessObserver(fn)
		{
		.observers.Remove(fn)
		}

	NotifyObservers(@args)
		{
		ok = true
		// have to copy because some times obervers get removed in the process
		for observer in .observers.Copy()
			ok = ok and observer(@args)
		return ok
		}
	ConfirmDestroy()
		{
		.ClearFocus()
		return .Save()
		}
	Destroy()
		{
		.Window.RemoveValidationItem(this)
		super.Destroy()
		}
	}
