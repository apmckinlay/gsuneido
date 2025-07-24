// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
function (state_prov)
	{
	if Province?(state_prov) or state_prov is 'PQ' or state_prov is 'NF'
		return 'CA'
	if State?(state_prov)
		return 'US'
	if MexicanState?(state_prov)
		return 'MX'
	return ''
	}