// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(reporter = false, filterBy = false, _customizable = false)
		{
		types = .getTypes()

		which = reporter is true ? #reporter : #customize
		types = types.Filter({ it.GetDefault(which, false) })

		if filterBy isnt false
			{
			filterByOb = types.Filter({ it.base is filterBy })
			if false is compatible = .GetCompatible(filterByOb.GetDefault(0, #()))
				return filterByOb

			types = types.Filter({ .compatible?(it, compatible) })
			}
		if customizable isnt false
			types = customizable.FilterTable(types)
		return types
		}

	// extracted for testing
	getTypes()
		{
		return Plugins().Contributions("FieldTypes", "type")
		}

	compatible?(type, compatible)
		{
		curCompatible = .GetCompatible(type)
		return curCompatible isnt false and curCompatible is compatible
		}

	GetCompatible(type)
		{
		return type.Member?(#compatible) ? type.Val_or_func(#compatible) : false
		}

	GetFormat(c)
		{
		return Datadict(c.base).Format
		}

	GetControl(c)
		{
		return Datadict(c.base).Control
		}
	}
