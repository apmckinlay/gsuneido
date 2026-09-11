// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_BuildUrl()
		{
		Assert(MapGoogle.BuildUrl("100 Broadway", "", "San Diego", #CA, "", "90210")
			is: MapGoogle.SearchUrl $ "&query=100 Broadway ,San Diego, CA")

		Assert(MapGoogle.BuildUrl("100 Broadway", "", "San Diego",
				#CA, "42.217777778N,83.278888889W")
			is: MapGoogle.SearchUrl $ "&query=42.217777778N,83.278888889W")

		Assert(MapGoogle.BuildUrl("", "", "", "", zip_postal: "90210"),
			is: MapGoogle.SearchUrl $ "&query=90210")
		}

	Test_MultiLocations()
		{
		Assert(MapGoogle.BuildMultiLocationsUrl(#()) is: "")
		Assert(MapGoogle.BuildMultiLocationsUrl(#(locations: ())) is: "")
		Assert(MapGoogle.BuildMultiLocationsUrl(#(locations: ("\tBoulder City\tNV")))
			is: MapGoogle.SearchUrl $ "&query=Boulder City,NV")
		Assert(MapGoogle.BuildMultiLocationsUrl(
				#(locations: ("	Boulder City	NV", "	Goodsprings	NV")))
			is: MapGoogle.DirectionUrl $
				"&origin=Boulder City,NV&destination=Goodsprings,NV")
		Assert(MapGoogle.BuildMultiLocationsUrl(#(locations: ("	Boulder City	NV",
					"	Goodsprings	NV",
					"	Las Vegas	NV")))
			is: MapGoogle.DirectionUrl $
				"&origin=Boulder City,NV&destination=Las Vegas,NV&" $
				"waypoints=Goodsprings,NV")
		Assert(MapGoogle.BuildMultiLocationsUrl(#(locations: ("	Boulder City	NV",
					"	Goodsprings	NV",
					"	Las Vegas	NV	89101")))
			is: MapGoogle.DirectionUrl $
				"&origin=Boulder City,NV&destination=Las Vegas,NV,89101&" $
				"waypoints=Goodsprings,NV")
		}

	Test_buildMultiLocOb()
		{
		b = MapGoogle.MapGoogle_buildMultiLocOb
		Assert(b(#()) is: "")
		Assert(b(.buildLocations(15)) is: "ERROR: Cannot Map more than 10 locations")
		Assert(b(.buildLocations(3))
			is: #("city0,state0,zip0", "city1,state1,zip1", "city2,state2,zip2"))
		}

	buildLocations(count)
		{
		ob = Object()
		locations = Object()
		for (i = 0; i < count; i += 1)
			locations.Add(" \tcity" $ i $ "\tstate" $ i $ "\tzip" $ i)
		ob.locations = locations
		return ob
		}
	}
