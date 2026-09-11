// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
MapWeb
	{
	SearchUrl:    "https://www.google.com/maps/search/?api=1"
	DirectionUrl: "https://www.google.com/maps/dir/?api=1"

	BuildUrl(address1, address2, city, state_prov, lat_long = "", zip_postal = "")
		{
		if not lat_long.Blank?()
			loc = lat_long.Trim()
		else if address1.Blank?() and address2.Blank?() and city.Blank?() and
			state_prov.Blank?()
			loc = zip_postal
		else
			loc = address1 $ ' ' $ address2 $ ',' $ city $ ", " $ state_prov

		return .SearchUrl $ "&query=" $ loc.Trim()
		}

	BuildMultiLocationsUrl(locationOb)
		{
		multiLocs = .buildMultiLocOb(locationOb)
		if not Object?(multiLocs)
			return multiLocs

		if multiLocs.Size() is 1
			return .SearchUrl $ "&query=" $ multiLocs[0]

		url =
			.DirectionUrl $ "&origin=" $ multiLocs[0] $ "&destination=" $ multiLocs.Last()

		if multiLocs.Size() > 2
			{
			waypoints = multiLocs[1..-1].Join('|')
			url $= "&waypoints=" $ waypoints
			}

		return url
		}

	buildMultiLocOb(locationOb)
		{
		if locationOb.Empty?()
			return ""

		apiLocationsLimit = 10
		if locationOb.locations.Size() > apiLocationsLimit
			return "ERROR: Cannot Map more than " $ apiLocationsLimit $ " locations"

		multiLocs = Object()
		for loc in locationOb.locations
			{
			if loc.Has?('\t')
				{
				ob = loc.Split('\t')
				zipIndex = 3
				zipPostal = ob.Member?(zipIndex) ? ob[zipIndex].Tr(" \t") : ""
				locStr = ob[1] $ ',' $ ob[2] $ Opt(',', zipPostal)
				}
			else // lat long
				locStr = loc
			multiLocs.Add(locStr)
			}

		if multiLocs.Size() is 0
			return ""
		return multiLocs
		}
	}
