// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	DynamicTypes: false
	New(args)
		{
		if args.Member?("dynamicTypes")
			.DynamicTypes = args.dynamicTypes
		.stickyFields = args.GetDefault("stickyFields", Object())
		.sticky_values = .DynamicTypes isnt false
			? Object().Set_default(Object())
			: Object()
		}

	SetStickyValues(record)
		{
		for f in .stickyFields
			{
			if .DynamicTypes isnt false
				.sticky_values[.current_type][f] = record[f]
			else
				.sticky_values[f] = record[f]
			}
		}

	ClearStickyFieldValues()
		{
		.sticky_values.Delete(all:)
		}

	DynamicTypeList(omitTypes = false)
		{
		if .DynamicTypes is false
			return false
		types = .DynamicTypes.Members()
		types.Remove("typefield")
		types.Remove("default")

		if omitTypes and .DynamicTypes.Member?("omitTypes")
			for type in .DynamicTypes.omitTypes
				types.Remove(type)
		return types.Remove("omitTypes")
		}

	GetTypeName()
		{
		if .current_type is false
			return false

		for member in .DynamicTypes.Members()
			{
			if (Object?(.DynamicTypes[member]) and
				.DynamicTypes[member].Member?("value") and
				.DynamicTypes[member].value is .current_type)
				return member
			}
		throw "missing type in Access dynamic types (" $ .current_type $ ")"
		}

	current_type: false
	SetCurrentType(args)
		{
		old_type = .current_type
		if (.DynamicTypes isnt false)
			{
			if (args.Member?(0))
				.current_type = .DynamicTypes[.find_typename(args[0])].value
			else
				{
				// empty table so On_New called with no argument
				// use default if supplied, else arbitrarily pick a type
				if (.DynamicTypes.Member?("default"))
					m = .DynamicTypes.default
				else
					{
					members = .DynamicTypes.Members()
					m = members[0] is "typefield" ? members[1] : members[0]
					}
				.current_type = .DynamicTypes[m].value
				}
			}
		return old_type is .current_type
		}
	find_typename(button)
		{
		// pre: dynamicTypes is an object
		if not .DynamicTypes.Member?(button)
			throw "missing type in Access dynamic types"
		return button
		}

	last_type: false
	ApplyStickyValues(record)
		{
		if (.current_type isnt false)
			record[.DynamicTypes.typefield] = .current_type
		record.Merge(.DynamicTypes isnt false
			? .sticky_values[.current_type]
			: .sticky_values)
		}

	DetectTypeChange(newrec, x, changeTypeFn)
		{
		if .DynamicTypes isnt false
			{
			if newrec is false
				.current_type = x[.DynamicTypes.typefield]
			if .current_type isnt .last_type
				{
				typename = .GetTypeName()
				changeTypeFn(typename, .DynamicTypes[typename].control)
				}
			.last_type = .current_type
			}
		}
	}