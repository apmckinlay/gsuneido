// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
MapWeb
	{
	BuildUrl(address1, address2, city, state_prov, lat_long, zip_postal = '')
		{
		if not lat_long.Blank?()
			loc = lat_long.Trim()
		else if address1.Blank?() and address2.Blank?() and city.Blank?() and
			state_prov.Blank?()
			loc = zip_postal
		else
			loc = address1 $ ' ' $ address2 $ ', ' $ city $ ', ' $ state_prov

		return 'http://www.bing.com/maps/?v=2&where1=' $ loc $ '&encType=1'
		}
	}