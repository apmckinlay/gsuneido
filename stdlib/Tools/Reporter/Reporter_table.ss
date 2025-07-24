// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Authorization(source)
		{
		return (false is value = .getValue(source, 'auth', 'Permission'))
			? source.name
			: value
		}

	BookLocation(source)
		{
		return (false is value = .getValue(source, 'bookLocation', 'BookLocation'))
			? false
			: value
		}

	getTable(source)
		{
		table = source.GetDefault('tables', #())
		if table.Empty?()
			return false
		return Tables.GetTable(table[0])
		}

	getSource(source, member)
		{
		source.GetDefault(member, false)
		}

	getValue(source, sourceMember, tableMember)
		{
		if false isnt value = .getSource(source, sourceMember)
			return value

		if false is table = .getTable(source)
			return false

		// In some cases like custom table tabs there may not be a value for things like
		// BookLocation
		return table.Member?(tableMember) ? table[tableMember] : false
		}
	}