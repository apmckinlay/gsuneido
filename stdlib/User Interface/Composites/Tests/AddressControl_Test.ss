// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_layout()
		{
		m = AddressControl.Layout
		expectedLayout = #(Form,
			#(address1, group: 0), #(MapButton, name: mapbutton), Skip, '', nl,
			#(address2, group: 0), '', nl,
			#(city, group: 0), #(state_prov),
				#(StaticText, '       ', name: country),
				#(StaticText, name: region), #(zip_postal), nl,
			ystretch: 0, xstretch: 0)
		Assert(m() is: expectedLayout)
		Assert(m(title: #Test) is: Object(#GroupBox, #Test, expectedLayout))

		extra_field1 = #test_field1
		expectedLayout = #(Form,
			#(address1, group: 0), #(Skip 0), Skip, test_field1, nl,
			#(address2, group: 0), '', nl,
			#(city, group: 0), #(state_prov),
				#(StaticText, '       ', name: country),
				#(StaticText, name: region), #(zip_postal), nl,
			ystretch: 0, xstretch: 0)
		Assert(m(:extra_field1, noAddressButton:) is: expectedLayout)

		prefix = #prefix_
		suffix = #_suffix
		extra_field2 = #test_field2
		expectedLayout = #(Form,
			#(prefix_address1_suffix, group: 0), #(Skip 0), Skip, test_field1, nl,
			#(prefix_address2_suffix, group: 0), test_field2, nl,
			#(prefix_city_suffix, group: 0), #(prefix_state_prov_suffix),
				#(StaticText, '       ', name: prefix_country_suffix),
				#(StaticText, name: prefix_region_suffix),
				#(prefix_zip_postal_suffix), nl,
			ystretch: 0, xstretch: 0)
		layout = m(:prefix, :suffix, :extra_field1, :extra_field2, noAddressButton:)
		Assert(layout is: expectedLayout)
		}
	}