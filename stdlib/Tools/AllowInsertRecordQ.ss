// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
function (record, protect_rule)
	{
	if protect_rule is false
		return true
	protect = record[protect_rule]
	return not Object?(protect) or protect.GetDefault("noInsert", false) isnt true
	}