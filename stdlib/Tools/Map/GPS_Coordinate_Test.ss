// Copyright (C) 2014 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_valid()
		{
		// garbage data cases
		Assert(GPS_Coordinate(#()).Valid?() is: false)
		Assert(GPS_Coordinate(false).Valid?() is: false)
		Assert(GPS_Coordinate(100).Valid?() is: false)
		Assert(GPS_Coordinate('').Valid?() is: false)
		Assert(GPS_Coordinate('abc').Valid?() is: false)
		Assert(GPS_Coordinate('abc,efg').Valid?() is: false)
		Assert(GPS_Coordinate('Saskatoon, SK').Valid?() is: false)
		Assert(GPS_Coordinate('S7N 2X8 Saskatoon, SK').Valid?() is: false)
		Assert(GPS_Coordinate('abc,efg').Valid?() is: false)
		Assert(GPS_Coordinate('N,W').Valid?() is: false)

		// invalid lat number value
		point = GPS_Coordinate('90.454629,-112.784675')
		Assert(point.Valid?() is: false)
		Assert(point.Latitude() is: false)
		Assert(point.Longitude() is: -112.784675)

		// invalid lat number value
		point = GPS_Coordinate('-90.454629,-112.784675')
		Assert(point.Valid?() is: false)
		Assert(point.Latitude() is: false)
		Assert(point.Longitude() is: -112.784675)

		// invalid lat number value with direction char
		point = GPS_Coordinate('90.454629S,112.784675W')
		Assert(point.Valid?() is: false)
		Assert(point.Latitude() is: false)
		Assert(point.Longitude() is: -112.784675)

		// invalid lat number value with direction char
		point = GPS_Coordinate('90.454629N,112.784675W')
		Assert(point.Valid?() is: false)
		Assert(point.Latitude() is: false)
		Assert(point.Longitude() is: -112.784675)

		// invalid long number value
		point = GPS_Coordinate('51.454629,-180.784675')
		Assert(point.Valid?() is: false)
		Assert(point.Latitude() is: 51.454629)
		Assert(point.Longitude() is: false)

		// invalid long number value
		point = GPS_Coordinate('51.454629,180.784675')
		Assert(point.Valid?() is: false)
		Assert(point.Latitude() is: 51.454629)
		Assert(point.Longitude() is: false)

		// invalid long number value with direction char
		point = GPS_Coordinate('51.454629N,180.784675W')
		Assert(point.Valid?() is: false)
		Assert(point.Latitude() is: 51.454629)
		Assert(point.Longitude() is: false)

		// invalid long number value with direction char
		point = GPS_Coordinate('51.454629N,180.784675W')
		Assert(point.Valid?() is: false)
		Assert(point.Latitude() is: 51.454629)
		Assert(point.Longitude() is: false)

		// valid cases
		point = GPS_Coordinate('49.454629,-112.784675')
		Assert(point.Valid?())
		Assert(point.Latitude() is: 49.454629)
		Assert(point.Longitude() is: -112.784675)

		point = GPS_Coordinate('49.454629, -112.784675')
		Assert(point.Valid?())
		Assert(point.Latitude() is: 49.454629)
		Assert(point.Longitude() is: -112.784675)

		point = GPS_Coordinate('-49.454629,112.784675')
		Assert(point.Valid?())
		Assert(point.Latitude() is: -49.454629)
		Assert(point.Longitude() is: 112.784675)

		point = GPS_Coordinate('49.454629N,112.784675W')
		Assert(point.Valid?())
		Assert(point.Latitude() is: 49.454629)
		Assert(point.Longitude() is: -112.784675)

		point = GPS_Coordinate('49.454629N, 112.784675W')
		Assert(point.Valid?())
		Assert(point.Latitude() is: 49.454629)
		Assert(point.Longitude() is: -112.784675)

		point = GPS_Coordinate('49.454629S,112.784675E')
		Assert(point.Valid?())
		Assert(point.Latitude() is: -49.454629)
		Assert(point.Longitude() is: 112.784675)
		}

	Test_DistanceFrom()
		{
		a = GPS_Coordinate('52.104481N,106.568403W')  // S7V 1G8
		b = GPS_Coordinate('50.452384N,104.624119W')  // Mosaic Stadium, Regina
		Assert(a.DistanceFrom(b).Round(3) is: 228.097)

		a = GPS_Coordinate('52.170934N,106.700972W')  // Saskatoon
		b = GPS_Coordinate('43.617692N,79.375939W')	  // Toronto
		Assert(a.DistanceFrom(b).Round(3) is: 2230.568)

		a = GPS_Coordinate('61.384467N,150.060611W')  // Anchorage, AK, USA
		b = GPS_Coordinate('31.986815S,115.720640E')  // Perth, AU
		Assert(a.DistanceFrom(b).Round(3) is: 13305.957)
		}

	Test_ToString()
		{
		a = GPS_Coordinate('52.104481N,106W')
		Assert(a.ToString() is: '52.104481,-106.000000')
		}

	Test_DMS()
		{
		a = GPS_Coordinate('53.818621N,105.389653W')
		Assert(a.DMS_Latitude(round: 1) is:
			#(minutes: "49", degrees: "53", seconds: "07"))
		Assert(a.DMS_Longitude(round: 1) is:
			#(minutes: "23", degrees: "-105", seconds: "22.8"))

		a = GPS_Coordinate('35.832486N,115.431899W')
		Assert(a.DMS_Latitude(round: 3) is:
			#(minutes: "49", degrees: "35", seconds: "56.95"))
		Assert(a.DMS_Longitude(round: 3) is:
			#(minutes: "25", degrees: "-115", seconds: "54.836"))

		a = GPS_Coordinate('44.391878N,68.204084W')
		Assert(a.DMS_Latitude() is:
			#(minutes: "23", degrees: "44", seconds: "31"))
		Assert(a.DMS_Longitude() is:
			#(minutes: "12", degrees: "-68", seconds: "15"))
		}
	}