// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 'DropDownButton'
	ComponentName: 'DropDownButton'
	isReadOnly: false
	New(hidden = false, .allowReadOnlyDropDown = false)
		{
		.SuSetHidden(hidden)
		.ComponentArgs = Object(hidden)
		}
	CLICK(x, y, rcExclude, rect)
		{
		_posInfo = Object(:x, :y, :rcExclude, :rect)
		.Send('On_DropDown')
		return 0
		}

	SetReadOnly(readOnly)
		{
		.isReadOnly = readOnly
		if .allowReadOnlyDropDown is true
			return
		super.SetReadOnly(readOnly)
		}

	GetReadOnly()
		{
		if .allowReadOnlyDropDown is true
			return .isReadOnly
		return super.GetReadOnly()
		}
	}