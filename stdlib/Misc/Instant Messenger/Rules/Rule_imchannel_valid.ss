// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	valid = Object()
	if .imchannel_name.Blank?()
		valid.Add("Name is required")
	else if not .imchannel_name.Tr(' ').AlphaNum?()
		valid.Add("Name must be alphanumeric containing only spaces")

	if .imchannel_abbrev.Blank?() or .imchannel_abbrev.White?()
		valid.Add("Abbreviation is required")
	else if .imchannel_abbrev.Has?(" ")
		valid.Add("Abbreviation cannot have spaces")
	else if not .imchannel_abbrev.Lower?()
		valid.Add("Abbreviation must be all lowercase")
	return valid.Join(', ')
	}
