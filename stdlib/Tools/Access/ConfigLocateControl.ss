// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: ConfigLocate

	New()
		{
		super([LocateAutoChooseControl, .getMatches, width: 8, xstretch: 0,
			cue: 'Find in Options', allowOther:, noTruncateValue:])
		}

	getMatches(prefix)
		{
		fields = .Send('GetConfigFields')
		fields.Each()
			{ |f|
			f.path = f.Members(list:).Sort!().Map({ f[it].section }).Join(' > ')
			f.lower = f.path.Lower()
			}
		prefix = prefix.Lower()
		found =  fields.Filter({ it.lower.Has?(prefix) }).Map({ it.path })
		if found.Size() is 0
			return #('No matches')
		return found
		}

	NewValue(.value)
		{
		.Send('LocateConfigField', value)
		}
	value: ''
	Get()
		{
		return .value
		}
	SetFocus()
		{
		.AutoChoose.SetFocus()
		}
	FieldEscape()
		{
		.Send('LocateEscape')
		}
	}
