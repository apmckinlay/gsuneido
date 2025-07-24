// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// TreeModel with Libraries() as roots
// creates a TreeModel for each library and delegates to them
// mangles numbers to identify which library
class
	{
	childrenCacheSize: 50
	New()
		{
		.libs = .Libs()
		.lts = .Lts(.libs)
		.cache = .initCache()
		.codeState = new CodeState(this)
		}

	Libs()
		{
		libs = Libraries()
		otherlibs = LibraryTables().Difference(libs)
		libs.Add(@otherlibs.Map!({ '(' $ it $ ')' }))
		return libs
		}

	Lts(libs)
		{
		lts = Object()
		for lib in libs
			lts.Add(new TreeModel(lib.Tr('()')))
		return lts
		}

	initCache() // Extracted for tests
		{
		return LruCache(.children?, .childrenCacheSize)
		}

	Create(library)
		{
		TreeModel.Create(library)
		.database("alter " $ library $
			" create (text, lib_invalid_text, lib_modified) key (name, group)")
		SvcDisabledLibraries.ResetCache()
		}

	database(q)
		{ Database(q) }

	libNumFactor: 100000
	Children(parent)
		{
		if parent is 0
			{
			items = Object()
			for lib in .libs.Members()
				items.Add(Object(name: .libs[lib], group: true, num: .libNum(lib)))
			return items
			}
		else
			{
			lib = .lib(parent)
			parent %= .libNumFactor
			children = .lts[lib].Children(parent)
			for i in children.Members()
				{
				item = children[i]
				item.num += .libNum(lib)
				item.parent += .libNum(lib)
				}
			return children
			}
		}

	libNum(lib)
		{
		return (lib + 1) * .libNumFactor
		}

	Children?(parent)
		{
		return .cache.Get(parent)
		}

	children?(parent)
		{
		// Return true if item has children
		return .lts.Member?(lib = .lib(parent))
			? .lts[lib].Children?(parent % .libNumFactor)
			: false
		}

	Container?(item)
		{
		// Return true if item is a container (i.e. a folder)
		if item % .libNumFactor is 0
			// Library roots are always folders...
			return true
		return .lts[.lib(item)].Container?(item % .libNumFactor)
		}

	Editable?(item)
		{
		return not .Container?(item)
		}

	Static?(item)
		{
		if item % .libNumFactor is 0
			// Library roots are always static
			return true
		return .lts[.lib(item)].Static?(item % .libNumFactor)
		}

	Get(num, origText? = false)
		{
		if num <= 0
			return false

		lib = .lib(num)
		if 0 is num %= .libNumFactor // Library root item
			{
			table = name = .libs[lib].Tr('()')
			return [:name, :table, group:]
			}

		if false is item = .lts[lib].Get(num)
			return false

		item.keyNum = item.num 		// Original record num (key value)
		item.num += .libNum(lib)	// Mangled record num (syncs with the TreeView)
		item.parent += .libNum(lib)
		if item.group is false
			.setItemText(item, origText?)
		return item
		}

	DisplayName(table, name)
		{
		table = table.Tr('()')
		if '' is name = name.Tr('()')
			name = table // Library folders do not have "name" set
		return Libraries().Has?(table) ? name : '(' $ name $ ')'
		}

	setItemText(item, origText?)
		{
		// Add a beforeText value on select to compare code changes
		if item.lib_before_text is '' and item.lib_modified is '' and
			item.lib_committed isnt ''
			item.lib_before_text = item.text
		if not origText?
			item.text = item.lib_current_text
		}

	Save(x)
		{ .codeState.Save(x) }

	Update(x)
		{
		.cache.Reset()
		table = .TableName(x.num)
		x.num %= .libNumFactor
		svcTable = SvcTable(table)
		QueryApply1(table, num: x.num)
			{
			it.lib_invalid_text = x.lib_invalid_text
			svcTable.Update(it, newText: x.GetDefault(#text, false), t: it.Transaction())
			}
		x.num = .MangleNum(table, x.num)
		return true
		}

	// NOTE: Calling code should verify that the newParent is not a child of the moved
	// record (Reference ExplorerMultiTreeControl.Move for an example)
	Move(x, newParent, t = false)
		{
		fromTable = SvcTable(.TableName(x.num))
		toTable = SvcTable(.TableName(newParent))
		.convertNums(x)
		.move(x, newParent % .libNumFactor, fromTable, toTable, t)
		.convertNums(x, toTable.Table())
		}

	convertNums(x, libName = false)
		{
		if libName isnt false
			{
			x.num = .MangleNum(libName, x.num)
			x.parent = .MangleNum(libName, x.parent)
			}
		else
			{
			x.num %= .libNumFactor
			x.parent %= .libNumFactor
			}
		}

	move(x, newParent, fromTable, toTable, t) // Recursive
		{
		DoWithTran(t, update:)
			{ |t|
			x.group
				? .moveFolder(x, newParent, fromTable, toTable, t)
				: .moveItem(x, newParent, fromTable, toTable, t)
			}
		}

	moveFolder(x, parent, fromTable, toTable, t)
		{
		children = .childrenToMove(fromTable, origNum = x.num, t)
		x.group = x.parent = parent
		x = TreeModel.EnsureUnique(x, toTable.Table())
		if tableMove = fromTable.Table() isnt toTable.Table()
			{
			x.num = NextTableNum(toTable.Table(), t)
			t.QueryOutput(toTable.Table(), x)
			}
		else
			t.QueryApply1(fromTable.Table(), num: x.num)
				{
				it.Merge(x)
				it.Update()
				}
		for child in children
			.move(child, x.num, fromTable, toTable, t)
		if tableMove
			t.QueryDo('delete ' $ fromTable.Table() $ ' where num is ' $ origNum)
		x.group = true
		}

	childrenToMove(svcTable, parent, t)
		{
		children = Object()
		t.QueryApply(svcTable.Table() $ ' where group >= -1 and parent is ' $ parent)
			{
			if not it.group = it.group > -1
				it.lib_before_path = svcTable.GetPath(it, t)
			children.Add(it)
			}
		return children
		}

	moveItem(x, newParent, fromTable, toTable, t)
		{
		fromTable.Table() isnt toTable.Table()
			? fromTable.MoveLibrary(x, toTable.Table(), newParent, t)
			: fromTable.Move(x, newParent, t)
		x.group = false
		}

	Rename(rec, newName)
		{
		libName = .TableName(rec.num)
		rec.num %= .libNumFactor
		rec.parent %= .libNumFactor
		if rec.group
			.renameFolder(rec, libName, newName)
		else
			SvcTable(libName).Rename(rec, newName)
		}

	renameFolder(rec, libname, newName)
		{
		children = Object()
		.collectChildren(libname, rec.num, children, svcTable = SvcTable(libname))
		QueryApply1(libname, num: rec.num)
			{ |folder|
			folder.name = newName
			folder.Update()
			children.Each({ svcTable.Move(it, it.parent, folder.Transaction()) })
			}
		}

	collectChildren(libname, parent, children, svcTable)
		{
		QueryApply(libname, :parent)
			{|rec|
			if rec.group > -1
				.collectChildren(libname, rec.num, children, svcTable)
			else if rec.group is -1
				{
				if rec.lib_before_path is ''
					rec.lib_before_path = svcTable.GetPath(rec)
				children.Add(rec)
				}
			}
		}

	Modified?(rec)
		{
		if rec.GetDefault(#group, true)
			return false
		return rec.lib_committed is '' and SvcSettings.Set?() or rec.lib_modified isnt ''
		}

	Valid?(data)
		{
		return data.lib_invalid_text is ''
		}

	Synced?(rec, savedRec)
		{
		return rec.lib_modified is savedRec.lib_modified and
			rec.lib_committed is savedRec.lib_committed and
			Adler32(rec.text) is Adler32(savedRec.text)
		}

	DeleteItem(num, name, group)
		{
		.cache.Reset()
		lib = .lib(num)
		if 0 is itemNum = num % .libNumFactor
			.deleteLibrary(lib, name)
		else
			.deleteItem(.TableName(num), itemNum, name, group, .lts[lib])
		}

	deleteItem(lib, num, name, group, tree)
		{
		svcTable = SvcTable(lib)
		if group
			{
			.deleteFolder(lib, num, tree)
			svcTable.Publish(#TreeChange)
			}
		else
			svcTable.StageDelete(name)
		}

	deleteFolder(lib, num, tree)
		{
		QueryAll(lib $ ' where parent is ' $ num $ ' and group >= -1').Each()
			{
			.deleteItem(lib, it.num, it.name, it.group isnt -1, tree)
			}
		QueryDo('delete ' $ lib $ ' where num is ' $ num)
		}

	deleteLibrary(lib, libname)
		{
		if #stdlib is libname = libname.Tr('()')
			return Alert('Cannot delete stdlib', 'Delete Library', 0, MB.ICONERROR)
		if Libraries().Has?(libname)
			ServerEval(#Unuse, libname)
		.database('drop ' $ libname)
		.libs.Delete(lib)
		.lts.Delete(lib)
		ResetCaches()
		LibraryTags.Reset()
		SvcTable.Publish(#TreeChange, type: 'lib')
		}

	Nextnum(parent)
		{
		lib = .lib(parent)
		return .lts[lib].Nextnum() + .libNum(lib)
		}

	NewItem(x)
		{
		.cache.Reset()
		lib = .lib(x.parent)
		x.parent %= .libNumFactor
		result = .lts[lib].NewItem(x)
		x.num += .libNum(lib)
		x.parent += .libNum(lib)
		SvcTable(.TableName(x.num)).Publish(#TreeChange)
		LibUnload(x.name)
		return result
		}

	lib(num)
		{
		return (num / .libNumFactor).Int() - 1
		}

	TreeSort(lParam1, lParam2)
		{
		// Alphabetical folders, then alphabetical items
		item1 = .Get(lParam1)
		item2 = .Get(lParam2)
		if item1.group and item2.group
			return item1.name < item2.name ? -1 : 1
		if item1.group
			return -1
		if item2.group
			return 1
		return item1.name < item2.name ? -1 : 1
		}

	TableName(num)
		{
		return .libs.GetDefault(.lib(num), false)
		}

	LibNum(num)
		{
		return .libNum(.lib(num))
		}

	MangleNum(libName, num)
		{
		return .libNum(.Libs().Find(libName)) + num
		}

	UnMangleNum(num)
		{
		return num % .libNumFactor
		}

	EnsureUnique(x)
		{
		lib = .lib(x.parent)
		parent = x.parent
		x.parent %= .libNumFactor
		origName = x.name
		x = .lts[lib].EnsureUnique(x)
		x.parent = parent
		if x.name isnt origName
			x.lib_committed = ''
		return x
		}
	}
