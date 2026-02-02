// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(table)
		{
		.table = table
		.Create(table)
		// BookEditControl initializes with table: '' when opened during persistent load
		if table isnt ''
			.svcTable = SvcTable(table)
		}

	Create(table)
		{
		if table is ""
			return false
		BookModel.Create(table)
		return true
		}

	NewItem(x)
		{
		Transaction(update:)
			{ |t|
			y = t.Query1(.table, num: x.parent)
			if y is false
				x.path = ""
			else
				x.path = y.path $ "/" $ y.name
			.svcTable.Output(x, :t, committed: x.lib_committed isnt '')
			}
		.ClearCache()
		return true
		}

	TreeSort(lParam1, lParam2)
		{
		item1 = .Get(lParam1)
		item2 = .Get(lParam2)
		return item1.order is item2.order
			? item1.name < item2.name ? -1 : 1
			: item1.order < item2.order ? -1 : 1
		}

	ClearCache()
		{
		ClearBookImageCache(.table, clearQuery1Cache?:)
		}

	Save(x)
		{
		.Update(x)
		}

	Update(x, t = false)
		{
		if x.name is 'HtmlPrefix' or x.name is 'HtmlSuffix'
			Query1CacheReset()
		DoWithTran(t, update:)
			{ |t|
			t.QueryApply1(.svcTable.Table(), num: x.num)
				{
				.formatData(x, it)
				.svcTable.Update(it, t, x.text)
				}
			}
		.ClearCache()
		return true
		}

	formatData(modified, saved)
		{
		modified.order = modified.GetDefault('order', saved.order)
		modified.path = modified.GetDefault('path', saved.path)
		// Appends order to text in order for version comparisons
		.svcTable.GetData(modified)
		.svcTable.GetData(saved)
		saved.order = modified.order
		}

	Rename(rec, newName)
		{
		.svcTable.Rename(rec, newName)
		.ClearCache()
		}

	// NOTE: Calling code should verify that the newParent is not a child of the moved
	// record (Reference ExplorerMultiTreeControl.Move for an example)
	Move(rec, newParent)
		{
		.svcTable.Move(rec, newParent)
		.ClearCache()
		}

	Modified?(rec)
		{
		// root folder is defined by -1 and 0, depends where the call is coming from
		if rec.num in (-1, 0)
			return false
		return rec.lib_committed is '' and SvcSettings.Set?() or rec.lib_modified isnt ''
		}

	Synced?(rec, savedRec)
		{
		return rec.lib_modified is savedRec.lib_modified and
			rec.lib_committed is savedRec.lib_committed and
			Adler32(rec.text) is Adler32(savedRec.text)
		}

	DeleteItem(name = false, path = false, num = false)
		{
		.deleteItems(:name, :path, :num)
		.ClearCache()
		}

	deleteItems(name = false, path = false, num = false)
		{
		for child in .Children(num)
			{
			childPath = Paths.Combine(.table, child.path, child.name)
			.deleteItems(name: child.name, path: childPath, num: child.num)
			}
		path = Opt('/', path.RemovePrefix(.table $ '/').BeforeLast('/' $ name))
		.svcTable.StageDelete(.svcTable.MakeName([:path, :name]))
		}

	Children(parent)
		{
		if parent is 0
			return Object(
				Object(group:, name: .table, parent: 0, num: -1, text: ""))
		else if parent is -1
			query = .table $ " where path=''"
		else
			{
			x = Query1(.table, num: parent)
			if x is false
				return Object()
			query = .table $ " where path=" $ Display(x.path $ "/" $ x.name)
			}
		children = Object()
		QueryApply(query $ ' sort order,name')
			{ |x|
			x.group = true
			children.Add(x)
			}
		return children
		}

	Children?(num)
		{
		if num is -1
			return not QueryEmpty?(.table, path: '')

		x = Query1(.table, :num)
		if x is false
			return false

		return not QueryEmpty?(.table, path: x.path $ "/" $ x.name)
		}

	Container?(item /*unused*/)
		{
		return true
		}

	Editable?(item)
		{
		if .Static?(item)
			return false

		rec = .Get(item)
		if rec.plugin is true
			return false

		name = SvcBook.MakeName([name: rec.name, path: rec.path])
		return not BookResource?(name, readOnly?:)
		}

	Static?(item)
		{
		return item is -1
		}

	Get(num)
		{
		if num is -1
			x = [num: 0, name: .table]
		else if false is x = Query1(.table, :num)
			return false
		x.group = true
		x.table = .table
		return x
		}

	EnsureUnique(x)
		{
		// Return a version of the record 'x' whose fields do not violate table keys
		origName = x.name
		if x.parent is -1
			x.path = ""
		else
			{
			y = Query1(.table, num: x.parent)
			if y is false
				return false
			x.path = y.path $ "/" $ y.name
			}
		while not QueryEmpty?(.table, name: x.name, path: x.path)
			x.name = .copyName(x.name)
		if x.name isnt origName
			x.lib_committed = ''
		return x
		}

	copyName(name)
		{
		copyNum = 1
		if false isnt copy = name.Extract(' Copy\s\d+$')
			{
			name = name.RemoveSuffix(copy)
			copyNum = Number(copy.AfterFirst(' Copy ')) + 1
			}
		return name $ ' Copy ' $ copyNum
		}

	GetTable()
		{
		return .table
		}

	TableName(unused)
		{
		return .GetTable()
		}
	}
