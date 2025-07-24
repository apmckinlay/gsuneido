// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// an ExplorerMultiControl model
// access a tree stored in a database table
// the table must have the following columns (and presumably additional ones for data)
//		num:int, parent:int, group:int, name:string
// each record is assigned a unique number (starting at 1)
// in the table: group=-1 for leaf items, group=parent for groups
// in the objects: group = true or false
// name+group key ensures unique item names overall, and unique group names within groups
// .original_table is the table passed in with no restrictions
// .table is the table passed in and any restrictions applied using .ChangeQuery
// root items have a parent of 0 (shouldn't be an item with a num of 0)
class
	{
	New(table)
		{
		.original_table = .table = table
		.nextnum = table $ " summarize max num"
		.Create(table)
		}
	Create(table)
		{
		Database("ensure " $ table $
			" ( num, parent, group, name ) " $
				"key (num) index (parent, name)")
		}
	Nextnum()
		{
		max_record = Query1(.nextnum)
		return max_record is false ? 1 : max_record.max_num + 1
		}
	NewItem(x)
		{
		x.num = .Nextnum()
		x.group = x.group ? x.parent : -1
		Plugins().ForeachContribution('TreeModel', 'newrecords',
			{|plugin| (plugin.func)(x, table: .table) })
		QueryOutput(.table, x)
		return true
		}
	Update(x)
		{
		Transaction(update:)
			{|t|
			y = .save_old(t, x.num)
			if y is false
				return false
			for i in x.Members()
				y[i] = x[i]
			if x.Member?("parent") or x.Member?("group")
				y.group = x.group ? x.parent : -1
			y.Update()
			}
		return true
		}
	DeleteItem(num)
		{
		Transaction(update:)
			{ |t|
			.delete(t, num)
			}
		}
	delete(t, num) // recursive
		{
		// itself
		if false is y = .save_old(t, num)
			return
		y.Delete()
		// children (if any)
		q = t.Query(.original_table $ " where parent=" $ Display(num))
		while false isnt (x = q.Next())
			.delete(t, x.num)
		}
	save_old(t, num)
		{
		if false is y = t.Query(.original_table $ " where num = " $ Display(num)).Next()
			Alert("can't get " $ num, title: 'Error', flags: MB.ICONERROR)
		return y
		}
	Children(parent)
		{
		// get the items using the restriction (.table)
		items = QueryAll(.table, :parent, group: -1)
		items.Each({ it.group = false })
		items.SortWith!(.sort)

		// get the folders without using the restriction (.original_table)
		folders = QueryAll(.original_table $ ' where group > -1', :parent)
		folders.Each({ it.group = true })
		folders.SortWith!(.sort)

		return folders.Add(@items)
		}
	sort(item)
		{
		return item.name.Tr('_', ' ')
		}
	Children?(parent)
		{
		// Return true if item has children
		return not QueryEmpty?(.original_table $ ' where group >= -1', :parent)
		}
	Container?(item)
		{
		// Return true if item has children
		if false is x = Query1(.original_table, num: item)
			return false
		return x.group > -1
		}
	Static?(item/*unused*/)
		{ return false }
	Get(num)
		{
		if false isnt x = Query1(.original_table, :num)
			{
			x.group = x.group is x.parent
			x.table = .original_table
			}
		return x
		}
	EnsureUnique(x, table = false, t = false)
		{
		if table is false
			table = .original_table
		// Return a version of record 'x' whose name follows uniqueness rules of table
		DoWithTran(t)
			{ |t|
			while false isnt t.Query(.uniqueQuery(x, table)).Next()
				x.name = .copyName(x.name)
			}
		return x
		}

	uniqueQuery(x, table)
		{
		where = Boolean?(x.group) and x.group or x.group > -1
			? 'parent is ' $ Display(x.parent) $ Opt(' and num isnt ', x.num) $
				' and group > -1'
			: 'group is -1'
		return table $ ' where name is ' $ Display(x.name) $ ' and ' $ where
		}

	copyName(name)
		{
		copyNum = 1
		if false isnt copy = name.Extract('_Copy\d+$')
			{
			name = name.RemoveSuffix(copy)
			copyNum = Number(copy.AfterFirst('_Copy')) + 1
			}
		return name $ '_Copy' $ copyNum
		}

	ChangeQuery(query)
		{
		.table = query
		}
	}
