// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function (field, record, protectField = false)
	{
	record.Invalidate(field $ "__protect")
	if record[field $ "__protect"] is true
		return true
	if protectField is false
		return false
	protect_val = record[protectField]
	if Boolean?(protect_val)
		return protect_val
	else if Object?(protect_val)
		{
		allbut? = protect_val.Member?(0) and protect_val[0] is 'allbut'
		return protect_val.Member?(field) isnt allbut?
		}
	else if String?(protect_val)
		return protect_val isnt ''
	else
		throw "invalid return type from protect rule"
	}