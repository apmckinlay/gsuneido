// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
PassthruController
	{
	New(@args)
		{
		super(.controls(args))
		.address1 = .FindControl('address1')
		.Send('Data')
		if .Send('ControlInRecord?') is true // only add Map if we have a RecordControl
			{
			.Horz.Append(#(Skip, small:))
			.Horz.Append(#(MapButton, name: 'mapbutton'))
			.address1.AddContextMenuItem("", "")
			.address1.AddContextMenuItem("Map", .On_Map)
			}
		}

	controls(args)
		{
		fieldCtrl = args.Copy().Add('Field', at: 0)
		fieldCtrl.name = 'address1'
		return Object(#Horz, fieldCtrl)
		}

	Data() {} // block data message from address1 field control

	NewValue(value/*unused*/)
		{
		.Send('NewValue', .Get())
		}

	Set(x)
		{
		.address1.Set(x)
		}

	Get()
		{
		return .address1.Get()
		}

	Dirty?(dirty = "")
		{
		return .address1.Dirty?(dirty)
		}

	Valid?()
		{
		return .address1.Valid?()
		}

	ValidData?(@args)
		{
		return EditControl.ValidData?(@args)
		}

	On_Map()
		{
		if false isnt ctrl = .FindControl("mapbutton")
			ctrl.On_Map()
		}

	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}