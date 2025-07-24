// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// Priority: 1: number/string 2: date/timestamp 3: composite key
	CallClass(keys, rec, formatFunc = false)
		{
		if false is formatFunc
			formatFunc = .formatRecAndKey
		dateKeyDisplay = ''
		compKeyDisplay = ''
		for key in keys
			{
			if .hasPriorityKey?(rec[key])
				return formatFunc(rec, key)
			if Date?(rec[key])
				dateKeyDisplay = formatFunc(rec, key)
			if key.Has?(',')
				compKeyDisplay = key.Split(',').Map({ formatFunc(rec, it) }).Join(',')
			}
		return dateKeyDisplay isnt ""
			? dateKeyDisplay
			: compKeyDisplay // will return '' if no key in table
		}

	hasPriorityKey?(value)
		{
		// checking "" because this would be value for a composite key
		return Number?(value) or (String?(value) and value isnt '')
		}

	formatRecAndKey(rec, field)
		{
		val = Date?(rec[field])
			? Display(rec[field])
			: rec[field]
		return field $ ' ' $ val
		}
	}