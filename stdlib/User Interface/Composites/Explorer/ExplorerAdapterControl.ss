// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// allow single controls to be used as Explorer model
// by converting Get/Set object to Get/Set value
// transfers GetState() and SetState() calls to its child control
PassthruController
	{
	Name: "ExplorerAdapter"
	New(control, field)
		{
		super(control)
		.field = field
		.ctrl = .GetChild()
		.curob = Object()
		}
	Get()
		{
		ob = Object()
		ob[.field] = .ctrl.Get()
		return ob
		}
	GetAll()
		{
		.curob[.field] = .ctrl.Get()
		return .curob
		}
	Get_Ctrl()
		{ return .ctrl }
	Set(ob)
		{
		.curob = ob.Copy()
		.ctrl.Set(ob.GetDefault(.field, ""))
		}
	Dirty?(dirty = '')
		{ return .ctrl.Dirty?(dirty) }
	GetState()
		{ return .ctrl.GetState() }
	SetState(stateobject)
		{ .ctrl.SetState(stateobject) }
	SetReadOnly(readonly = true)
		{
		if (.ctrl.Method?('SetReadOnly'))
			.ctrl.SetReadOnly(readonly)
		}
	}
