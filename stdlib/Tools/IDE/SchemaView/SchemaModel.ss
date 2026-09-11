// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New()
		{
		.tables = Object()
		QueryApply(#views)
			{|x|
			.tables.Add([name: x.view_name, group: false])
			}
		QueryApply(#tables)
			{|x|
			.tables.Add([name: x.table, group:])
			}
		.tables.Sort!(By(#name))
		i = 0
		for x in .tables
			x.num = ++i
		}

	colstart: 10_000
	Children(parent)
		{
		if parent is 0
			return .tables
		if parent >= .colstart
			return []
		x = .tables[parent-1]
		if not x.group
			return []
		cols = QueryAll("columns where table = " $ Display(x.name) $ " sort column")
		return .children(cols, parent * .colstart)
		}

	children(cols, mangle)
		{
		children = Object()
		prevNum = 0 // Sudo "num" for uppercase rule columns as they have no true "num"
		for col in cols
			{
			num = col.field is -1 // Uppercase rules always have a value of: -1
				? ++prevNum
				: col.field
			children.Add([name: col.column, num: mangle + num, group: false])
			prevNum = num
			}
		return children
		}

	Get(i, name)
		{
		return i >= .colstart
			? .getColumn(i, name)
			: .tables[--i].group
				? Object(:name, group:, table: name, type: #table)
				: Object(:name, group: false, table: name, type: #view)
		}

	getColumn(i, name)
		{
		table = .tables[(i / .colstart).Int() - 1].name
		return Object(group: false, :name, :table, type: #column)
		}

	Children?(parent/*unused*/)
		{
		return true
		}

	Container?(i)
		{
		return i < .colstart and .tables[Max(i - 1, 0)].group
		}

	Static?(x/*unused*/)
		{
		return true
		}

	Nextnum() { }

	NewItem(x/*unused*/) { }

	Update(x/*unused*/) { }

	DeleteItem(num/*unused*/) { }
	}
