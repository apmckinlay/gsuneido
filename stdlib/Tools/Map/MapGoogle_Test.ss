// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_BuildUrl()
		{
		Assert(MapGoogle.BuildUrl('100 Broadway', '', 'San Diego', 'CA', '', '90210')
			is: "http://maps.google.com/maps?q=100 Broadway ,San Diego, CA&hl=en")

		Assert(MapGoogle.BuildUrl('100 Broadway', '', 'San Diego',
			'CA', '42.217777778N,83.278888889W')
			is: "http://maps.google.com/maps?q=42.217777778N,83.278888889W&hl=en")

		Assert(MapGoogle.BuildUrl('', '', '', '', zip_postal: '90210'),
			is: "http://maps.google.com/maps?q=90210&hl=en")
		}
	Test_MultiLocations()
		{
		Assert(MapGoogle.BuildMultiLocationsUrl(#()) is: '')
		Assert(MapGoogle.BuildMultiLocationsUrl(#(locations: #())) is: '')
		Assert(MapGoogle.BuildMultiLocationsUrl(
			#(locations: #('	Boulder City	NV')))
			is: 'http://maps.google.com/maps?f=d&source=s_d&saddr=Boulder City,NV')
		Assert(MapGoogle.BuildMultiLocationsUrl(
			#(locations: #('	Boulder City	NV', '	Goodsprings	NV')))
			is: 'http://maps.google.com/maps?f=d&source=s_d&saddr=Boulder City,NV&' $
				'daddr=Goodsprings,NV&hl=en')
		Assert(MapGoogle.BuildMultiLocationsUrl(
			#(locations: #('	Boulder City	NV', '	Goodsprings	NV',
				'	Las Vegas	NV')))
			is: 'http://maps.google.com/maps?f=d&source=s_d&saddr=Boulder City,NV&' $
				'daddr=Goodsprings,NV to:Las Vegas,NV&hl=en')
		Assert(MapGoogle.BuildMultiLocationsUrl(
			#(locations: #('	Boulder City	NV', '	Goodsprings	NV',
				'	Las Vegas	NV	89101')))
			is: 'http://maps.google.com/maps?f=d&source=s_d&saddr=Boulder City,NV&' $
				'daddr=Goodsprings,NV to:Las Vegas,NV,89101&hl=en')
		}

	Test_buildMultiLocOb()
		{
		b = MapGoogle.MapGoogle_buildMultiLocOb
		Assert(b(#()) is: "")
		Assert(b(.buildLocations(15)) is: 'ERROR: Cannot Map more than 10 locations')
		Assert(b(.buildLocations(3))
			is: #('city0,state0,zip0', 'city1,state1,zip1', 'city2,state2,zip2'))
		}
	buildLocations(count)
		{
		ob = Object()
		locations = Object()
		for (i = 0; i < count; i++)
			locations.Add(' \tcity' $ i $ '\tstate' $ i $ '\tzip' $ i)
		ob.locations = locations
		return ob
		}
	}
