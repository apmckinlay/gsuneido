// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(address1, address2, city, state_prov, lat_long = "", zip_postal = "")
		{
		url = .BuildUrl(.stripNonEncodeableChars(address1),
			.stripNonEncodeableChars(address2), city, state_prov, :lat_long, :zip_postal)
		.openUrl(url)
		}
	stripNonEncodeableChars(text)
		{
		// having '#' in front of the unit in an address does not work
		// stripping it out before passing the address on to the encode
		// if other chars prove problematic, add them here.
		return text.Tr("#")
		}
	BuildUrl(@unused)
		{ throw 'must be defined by inheriting class' }
	HandleMultiLocations(locations)
		{
		url = .BuildMultiLocationsUrl(locations)
		if url is ""
			return ''
		if url.Prefix?('ERROR:')
			return url[7 ..] /*= ERROR: prefix size */
		.openUrl(url)
		return ''
		}
	BuildMultiLocationsUrl(locations /*unused*/)
		{
		return 'ERROR: This mapping option does not support Multiple Locations'
		}
	openUrl(url)
		{
		SetFocus(NULL) // prevent "function call overflow" if user lean on space key
		ShellExecute(0, 'open', Url.Encode(url))
		}
	}