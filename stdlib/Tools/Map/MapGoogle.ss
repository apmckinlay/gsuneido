// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
MapWeb
	{
	BuildUrl(address1, address2, city, state_prov, lat_long = '', zip_postal = '')
		{
		if not lat_long.Blank?()
			loc = lat_long.Trim()
		else if address1.Blank?() and address2.Blank?() and city.Blank?() and
			state_prov.Blank?()
			loc = zip_postal
		else
			loc = address1 $ ' ' $ address2 $ ',' $ city $ ', ' $ state_prov

		return 'http://maps.google.com/maps?q=' $ loc $ '&hl=en'
		}
	BuildMultiLocationsUrl(locationOb)
		{
		multiLocs = .buildMultiLocOb(locationOb)
		if not Object?(multiLocs)
			return multiLocs

		url = 'http://maps.google.com/maps?f=d&source=s_d'
		if multiLocs.Size() is 1
			return url $ '&saddr=' $ multiLocs[0]
		url $= '&saddr=' $ multiLocs[0] $ '&daddr=' $ multiLocs[1]
		for (i = 2; i < multiLocs.Size(); i++)
			url $= ' to:' $ multiLocs[i]
		return url $ '&hl=en'
		}
	buildMultiLocOb(locationOb)
		{
		if locationOb.Empty?()
			return ''

		apiLocationsLimit = 10
		if locationOb.locations.Size() > apiLocationsLimit
			return 'ERROR: Cannot Map more than ' $ apiLocationsLimit $ ' locations'

		multiLocs = Object()
		for loc in locationOb.locations
			{
			if loc.Has?('\t')
				{
				ob = loc.Split("\t")
				zipIndex = 3
				zipPostal = ob.Member?(zipIndex) ? ob[zipIndex].Tr(" \t") : ""
				locStr = ob[1] $ "," $ ob[2] $ Opt(",", zipPostal)
				}
			else // lat long
				locStr = loc
			multiLocs.Add(locStr)
			}

		if multiLocs.Size() is 0
			return ''
		return multiLocs
		}
	}
