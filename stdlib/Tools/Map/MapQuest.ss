// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
MapWeb
	{
	BuildUrl(address1, address2, city, state_prov, lat_long /*unused*/, zip_postal = '')
		{
		loc = 'address=' $ address1 $ ' ' $ address2 $
			'&city=' $ city $ '&state=' $ state_prov $
			'&country=' $ CountryFromStateProv(state_prov) $
			'&zipcode='
		if address1.Blank?() and address2.Blank?() and city.Blank?() and
			state_prov.Blank?()
			loc = 'address= &city=&state=&country=&zipcode=' $ zip_postal

		return 'http://mapquest.com/maps/map.adp?' $ loc $ '&cid=lfmaplink'
		}
	}