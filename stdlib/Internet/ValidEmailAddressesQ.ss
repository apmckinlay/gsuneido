// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
function (addrs)
	{
	return not addrs.Blank?() and
		addrs.Tr(';', ',').Split(',').Map!(#Trim).Every?(ValidEmailAddress?)
	}