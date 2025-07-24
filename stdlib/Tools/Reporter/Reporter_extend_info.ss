// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(type, infoValues)
		{
		res = Object()
		for infoValue in infoValues
			{
			if type is infoValue.BeforeFirst(':')
				res.Add(infoValue.AfterFirst(':').Trim())
			}
		return res.Join(', ')
		}

	Prefix: 'reporter_info_'
	AddFields(sf)
		{
		infoFields = Object()
		prefixList = Object()
		for fld in sf.Fields
			{
			if fld !~ "_info[1-6]$"
				continue

			prefix = fld.BeforeLast('_')
			if prefixList.Member?(prefix) is false
				{
				if false is prompt = .PrefixPrompt(prefix)
					continue

				prefixList[prefix] = Object(:prompt, deps: Object(fld))
				}
			else
				prefixList[prefix].deps.AddUnique(fld)
			}
		for prefix, ob in prefixList
			for infotype in InfoTypes
				{
				type = infotype[.. -1]
				field = .Prefix $ prefix $ '_' $ type
				displayPrompt = Opt(ob.prompt, ' ') $ type
				sf.AddField(field, displayPrompt)
				infoFields.Add(Object(:field, :prefix, :type,
					prompt: displayPrompt, deps: ob.deps))
				}
		return infoFields
		}

	PrefixPrompt(prefix)
		{
		if prefix $ '_num' is prompt = SelectPrompt(prefix $ '_num')
			prompt = SelectPrompt(prefix $ '_name')

		// if still can't find a datadict at this point - skip
		if prefix $ '_name' is prompt
			return false

		return .infoPrompt(prompt)
		}

	infoPrompt(prompt)
		{
		if prompt is "Name"
			return ""

		if prompt.Has?('Date/Time')
			return prompt.BeforeFirst('Date/Time').Trim()

		return prompt
		}

	Extend(infoFields, fields)
		{
		extends = infoFields.
			Filter({ fields.Has?(it.field) }).
			Map({ it.field $ ' = Reporter_extend_info(' $ Display(it.type) $
				', Object(' $ it.deps.Sort!().Join(', ') $ '))' })
		return Opt('\nextend ', extends.Join(',\n'))
		}
	}