// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'MapZipPostal'
	New(mandatory = false, readonly = false, hidden = false, tabover = false)
		{
		super(.controls(mandatory, readonly, hidden, tabover))
		.zipPostalField = .FindControl('zip_postal_field')
		.zipPostalField.AddContextMenuItem("", "")
		.zipPostalField.AddContextMenuItem("Map", .On_Map)
		.Send('Data')
		}
	controls(mandatory, readonly, hidden, tabover)
		{
		ctrls = [#Horz,
			[ZipPostalControl, :mandatory, :readonly, :hidden, :tabover,
				name: 'zip_postal_field']]
		if not hidden
			ctrls.Add(#(Skip, small:), #(MapButton, name: 'mapbutton'))
		return ctrls
		}

	NewValue(value/*unused*/)
		{
		.Send('NewValue', .Get())
		}
	Set(x)
		{
		.zipPostalField.Set(x)
		}
	Get()
		{
		return .zipPostalField.Get()
		}
	Dirty?(dirty = "")
		{
		return .zipPostalField.Dirty?(dirty)
		}

	Valid?()
		{
		return .GetReadOnly() or .zipPostalField.Valid?()
		}
	ValidData?(@args)
		{
		return ZipPostalControl.ValidData?(@args)
		}

	On_Map()
		{
		if false isnt ctrl = .FindControl("mapbutton")
			ctrl.On_Map()
		}
	Map_GetAddress()
		{
		return Object(address1: '', address2: '', city: '',
			state_prov: '', zip_postal: .Get())
		}
	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}