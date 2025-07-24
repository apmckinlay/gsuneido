// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: 'Address'
	New(prefix = '', suffix = '', title = '', extra_field1 = '', extra_field2 = '',
		noAddressButton = false)
		{
		super(.Layout(prefix, suffix, title, extra_field1, extra_field2, noAddressButton))
		.Left = title is '' ? .Form.Left : .GroupBox.Form.Left
		.address1 = .FindControl(prefix $ 'address1' $ suffix)
		.address2 = .FindControl(prefix $ 'address2' $ suffix)
		.city = .FindControl(prefix $ 'city' $ suffix)
		.state_prov = .FindControl(prefix $ 'state_prov' $ suffix)
		.zip_postal = .FindControl(prefix $ 'zip_postal' $ suffix)
		.country = .FindControl(prefix $ 'country' $ suffix)
		.addContextToAddress1()
		}

	Layout(prefix = '', suffix = '', title = '', extra_field1 = '', extra_field2 = '',
		noAddressButton = false)
		{
		map = noAddressButton ? #(Skip 0) : Object('MapButton', name: 'mapbutton')
		ctrl = Object('Form',
			Object(prefix $ 'address1' $ suffix, group: 0), map,
				'Skip', extra_field1, 'nl',
			Object(prefix $ 'address2' $ suffix, group: 0), extra_field2, 'nl',
			Object(prefix $ 'city' $ suffix, group: 0),
				Object(prefix $ 'state_prov' $ suffix),
				Object('StaticText', '       ', name: prefix $ 'country' $ suffix),
				Object('StaticText', name: prefix $ 'region' $ suffix),
				Object(prefix $ 'zip_postal' $ suffix), 'nl',
			xstretch: 0, ystretch: 0)
		if title isnt ''
			ctrl = Object('GroupBox', title, ctrl)
		return ctrl
		}

	Map_GetAddress()
		{
		return Object(address1: .address1.Get(),
			address2: .address2.Get(),
			city: .city.Get(),
			state_prov: .state_prov.Get(),
			zip_postal: .zip_postal.Get()
			country: .country.Get().Trim())
		}

	addContextToAddress1()
		{
		if not .address1.Method?('AddContextMenuItem')
			return
		.address1.AddContextMenuItem('', '')
		.address1.AddContextMenuItem('Map', .On_Map)
		}

	On_Map()
		{
		if false isnt ctrl = .FindControl('mapbutton')
			ctrl.On_Map()
		}

	MakeSummary() // used by Expand and Accordion
		{
		cityJur = Join(' ', .city.Get(), .state_prov.Get())
		fullAddress = Join(', ', .address1.Get(), .address2.Get(), cityJur)
		return Join('  ', fullAddress, .zip_postal.Get())
		}
	}
