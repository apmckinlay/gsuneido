// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'LatLong'
	New(readonly = false, tabover = false, hidden = false)
		{
		super(.controls(readonly, tabover, hidden))
		.latLongField = .FindControl('lat_long_field')
		.latLongField.AddContextMenuItem("", "")
		.latLongField.AddContextMenuItem("Map", .On_Map)
		.Send('Data')
		}
	tip: "e.g. 52.140658,-106.625616 " $
		`Right click the point on Google Maps and click "What's Here?" to obtain`
	controls(readonly, tabover, hidden)
		{
		ctrls = [#Horz,
			[.latLongFieldControl, :readonly, :tabover, :hidden, name: 'lat_long_field',
				width: 12, status: .tip]]
		if not hidden
			ctrls.Add(#(Skip, small:), #(MapButton onlyLatLong?:, name: 'mapbutton'))
		return ctrls
		}

	latLongFieldControl: FieldControl
		{
		KillFocus()
			{
			latLongOb = GPS_Coordinate(.Get().Replace('\s',''))
			if latLongOb.Valid?()
				SetWindowText(.Hwnd, latLongOb.ToString())
			}
		Valid?()
			{
			if '' is s = .Get()
				return true
			return GPS_Coordinate(s).Valid?()
			}
		Set(x)
			{
			latLongOb = GPS_Coordinate(x)
			if latLongOb.Valid?()
				x = latLongOb.ToString()
			super.Set(x)
			}
		}

	NewValue(value/*unused*/)
		{
		.Send('NewValue', .Get())
		}
	Set(x)
		{
		.latLongField.Set(x)
		}
	Get()
		{
		return .latLongField.Get()
		}
	Dirty?(dirty = "")
		{
		return .latLongField.Dirty?(dirty)
		}

	Valid?()
		{
		return .GetReadOnly() or .latLongField.Valid?()
		}

	ValidData?(@args)
		{
		val = args[0]
		return val is "" or GPS_Coordinate(val).Valid?()
		}

	On_Map()
		{
		if false isnt ctrl = .FindControl("mapbutton")
			ctrl.On_Map()
		}
	Map_GetAddress()
		{
		return Object(address1: '', address2: '', city: '',
			state_prov: '', zip_postal: '', lat_long: .Get())
		}
	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}
