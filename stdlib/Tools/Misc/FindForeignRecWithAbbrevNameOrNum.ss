// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// when useNum? the rec does not have the num field required
	// but will build it based on the name/abbrev field it does have
	CallClass(rec, field, useNum? = false, filter = "")
		{
		if false is (field.Has?('name') or field.Has?('abbrev'))
			return false

		if false is table = .getNumTable(field.BeforeFirst("_"))
			return false

		keys = .queryKeys(table)
		indexes = .uniqueIndexes(table)

		return useNum?
			? .getForeignRecFromNum(rec, field, keys, indexes, table)
			: .getForeignRecFromNameOrAbbrev(rec, field, keys, indexes, table, filter)
		}

	getForeignRecFromNum(rec, field, keys, indexes, table)
		{
		numfield = field.Replace('_name|_abbrev', '_num')
		baseNum = numfield.Has?('_num_')
			? numfield.BeforeLast('_')
			: numfield

		if .keysIndexHasField?(keys, indexes, baseNum) is false
			return false

		if rec[numfield] is ''
			{
			if false is numPos = rec.Members().FindIf({ it.Prefix?(baseNum) })
				return false
			numfield = rec.Members()[numPos]
			}
		return Query1(table $ ' where ' $ baseNum $ ' is ' $ Display(rec[numfield]))
		}

	getForeignRecFromNameOrAbbrev(rec, field, keys, indexes, table, filter = "")
		{
		basename = field.BeforeFirst("_")
		lookupField = field.Has?("_name") ? basename $ "_name" : basename $ "_abbrev"

		if .keysIndexHasField?(keys, indexes, lookupField) is false
			return false

		lookupVal = rec[field]
		if keys.Has?(lookupField $ '_lower!')
			{
			lookupField $= '_lower!'
			lookupVal = lookupVal.Lower()
			}
		return Query1(table $ ' where ' $ lookupField $ ' is ' $ Display(lookupVal) $
			" " $ Opt("and ", filter))
		}

	keysIndexHasField?(keys, indexes, field)
		{
		return keys.Has?(field) or indexes.Has?(field)
		}

	// extracted for tests
	getNumTable(baseName)
		{
		return GetNumTable(baseName)
		}
	queryKeys(table)
		{
		return QueryKeys(table)
		}
	uniqueIndexes(table)
		{
		return UniqueIndexes(table)
		}
	}
