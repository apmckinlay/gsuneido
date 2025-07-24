// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(from_query, to_query, customFields = false, current_field = false)
		{
		.map = ServerEval('CustomizableMap.GetMap', from_query, to_query,
			customFields, current_field)
		}

	// only public for ServerEval - should not be called from elsewhere
	GetMap(from_query, to_query, customFields = false, current_field = false)
		{
		return (.cached)(from_query, to_query, customFields, current_field)
		}

	cached: Memoize
		{
		Func(from_query, to_query, customFields = false, current_field = false)
			{
			from = Customizable(from_query)
			to = Customizable(to_query)

			fromPrompts = from.CustomFields().Map(Prompt)
			toPrompts = to.CustomFields().Map(Prompt)

			map = Object()
			for prompt in fromPrompts.Intersect(toPrompts)
				{
				from_field = from.PromptToField(prompt)
				to_field = to.PromptToField(prompt)
				toType = DatadictType(to_field)
				fromType = DatadictType(from_field)
				if .allowToFill?(fromType, toType, to_field, customFields, current_field)
					map.AddUnique(Object(:to_field, :from_field,
						trim: .trim?(toType, to_field, from_field)))
				}
			return map
			}

		allowToFill?(fromType, toType, to_field, customFields, current_field)
			{
			if toType isnt fromType or Object(toType, fromType).Has?('image')
				return false
			if customFields is false or current_field is false
				return true
			if not customFields.Member?(to_field)
				return true

			onlyFillinFrom = customFields[to_field].GetDefault(
				'only_fillin_from', '')

			if onlyFillinFrom is ''
				return true
			return onlyFillinFrom is current_field
			}

		trim?(toType, to_field, from_field)
			{
			// if we get here toType and fromType will BOTH be string, only need to check
			// one of them.
			if toType is 'string'
				{
				toDD = Datadict(to_field)
				fromDD = Datadict(from_field)
				if fromDD.Control[0] is 'Editor' and toDD.Control[0] is 'Field'
					return true
				}
			return false
			}

		}

	ResetServerCache()
		{
		if Sys.Client?()
			ServerEval('CustomizableMap.ResetServerCache')
		else
			{
			.cached.ResetCache()
			QueryColumns.ResetCache()
			}
		}

	CopyCustomFields(toRec, fromRec = false)
		{
		if fromRec is false
			fromRec = toRec
		for from_to in .map
			{
			toRec[from_to.to_field] = .TrimValue(fromRec[from_to.from_field],
				from_to.trim)
			}
		}

	// this assumes that if trim is true that x is a string
	TrimValue(x, trim)
		{
		if not trim or x.Size() <= FieldControl.MaxCharacters
			return x
		return x[ .. FieldControl.MaxCharacters - 3 /*=elipsis*/] $ '...'
		}

	ForEachField(block)
		{
		for from_to in .map
			block(from_to.from_field, from_to.to_field, from_to.trim)
		}

	MapEmpty?()
		{
		return .map.Empty?()
		}

	FindMatchedCustomReportColumns(from_fields, from_query, to_query)
		{
		customFields = from_fields.
			Map(ChooseColumns.GetFieldName).
			Filter(Customizable.CustomField?)

		if customFields.Empty?()
			return false

		project = Opt(' project ', customFields.Join(', '))
		return CustomizableMap(from_query $ project, to_query)
		}
	}