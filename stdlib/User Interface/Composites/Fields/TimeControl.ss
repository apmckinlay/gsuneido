// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	Name: "time"
	New(width = 4, readonly = false, .mandatory = false,
		status = "A time e.g. 930 or 1500", hidden = false, tabover = false)
		{
		super(:width, :status, :readonly, justify: 'RIGHT', :hidden, :tabover)
		}
	Valid?()
		{
		t = super.Get()
		return .valid?(t, .mandatory)
		}
	valid?(value, mandatory)
		{
		return value is "" ? not mandatory : Time?(value)
		}
	ValidData?(@args)
		{
		return .valid?(args[0], args.GetDefault('mandatory', false))
		}
	KillFocus()
		{
		t = super.Get()
		if not Time?(t)
			return
		if t.Size() is 4 and t[0] is "0"
			t = t[1..]
		t = t.LeftFill(3, '0')
		super.Set(t)
		}
	Set(value)
		{
		if Number?(value)
			value = String(value).LeftFill(3, '0')
		super.Set(value)
		}
	Get()
		{
		if '' is time = super.Get()
			return ''
		return .Valid?() ? Number(time) : time
		}
	}
