// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Setup()
		{
		super.Setup()
		.TearDownIfTablesNotExist(UserSettings.Table)
		.masterTable = .MakeTable('(vl_sort_test_num, vl_sort_test_abbrev,
			vl_sort_test_name)
			key (vl_sort_test_num)
			key (vl_sort_test_name)
			index unique (vl_sort_test_abbrev)')
		.masterTable2 = .MakeTable('(vl_sort_test2_num, vl_sort_test2_abbrev,
			vl_sort_test2_name)
			key (vl_sort_test2_num)
			key (vl_sort_test2_name)
			index unique (vl_sort_test2_abbrev)')

		.table = .MakeTable('(a, b, c, vl_sort_test_num, vl_sort_test2_num) key(a)
			index (vl_sort_test_num) in ' $ .masterTable $
			' index (vl_sort_test2_num) in ' $ .masterTable2)

		.sf = SelectFields(#(vl_sort_test_num, vl_sort_test2_num))
		}

	Test_new()
		{
		sort = VirtualListSortModel(.table, .sf, .TempName(), loadAll?:)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' ')

		sort = VirtualListSortModel(.table, .sf, .TempName())
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' ')

		sort = VirtualListSortModel(
			.table $ ' sort  reverse  a', .sf, .TempName(), loadAll?:)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse a')

		sort = VirtualListSortModel(.table $ ' sort  reverse  a', .sf, .TempName())
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse a')

		Assert(sort.VirtualListSortModel_getSortStr(
			sortOb: [[col: 'b', dir: 1, id?: false]])
			is: ' sort b')

		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse a')
		}

	Test_no_save_name()
		{
		sort = VirtualListSortModel(.table $ ' sort a', .sf)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a')
		sort.SetSort('a')
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse a')
		sort.SetSort('c')
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort c')

		// testing multiple sorts
		sort = VirtualListSortModel(.table $ ' sort a, b', .sf)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a,b')

		sort = VirtualListSortModel(.table $ ' sort reverse a, b', .sf)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse a,b')
		sort.SetSort('c')
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort c')
		sort.SetSort('c')
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse c')
		sort.ResetSort()
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse a,b')
		}

	Test_BuildQuery()
		{
		sort = VirtualListSortModel(query = .table $ ' sort reverse a, b', .sf)
		Assert(sort.BuildQuery(query) is: .table $ ' sort reverse a,b')

		sort.SetSort('c')
		Assert(sort.BuildQuery(query) is: .table $ ' sort c')

		sort.SetSort('c')
		Assert(sort.BuildQuery(query) is: .table $ ' sort reverse c')

		// make sure BuildQuery does not modify interal sort object
		Assert(sort.BuildQuery(query, sortCol: 'b') is: .table $ ' sort b')
		Assert(sort.BuildQuery(query) is: .table $ ' sort reverse c')

		Assert(sort.BuildQuery(query, sortCol: 'c') is: .table $ ' sort c')
		Assert(sort.BuildQuery(query) is: .table $ ' sort reverse c')

		Assert(sort.BuildQuery(query, 'where b isnt ""', sortCol: 'c')
			is: .table $ ' where b isnt "" sort c')
		Assert(sort.BuildQuery(query, 'where b isnt ""', sortCol: 'b')
			is: .table $ ' where b isnt "" sort b')
		Assert(sort.BuildQuery(query) is: .table $ ' sort reverse c')
		Assert(sort.BuildQuery(query, 'where b isnt ""') is:
			.table $ ' where b isnt "" sort reverse c')
		}

	Test_UsingDefaultSort?()
		{
		sort = VirtualListSortModel(.table $ ' sort c', .sf)
		Assert(sort.UsingDefaultSort?(.table $ ' sort c'))
		Assert(sort.UsingDefaultSort?(.table $ ' sort reverse c') is: false)
		Assert(sort.UsingDefaultSort?(.table $ ' sort b') is: false)
		}

	Test_sortSaveName_load_all()
		{
		saveName = .TempName()
		UserSettings.Put(saveName, 'a')
		.AddTeardown({ UserSettings.Remove(saveName) })

		UserSettings.Put(saveName, 'xxxxx')
		sort = VirtualListSortModel(.table, .sf, saveName, loadAll?:)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' ')

		UserSettings.Put(saveName, 'xxxxx,a')
		sort = VirtualListSortModel(.table, .sf, saveName, loadAll?:)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a')

		UserSettings.Put(saveName, '-xxxxx,-a')
		sort = VirtualListSortModel(.table, .sf, saveName, loadAll?:)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse a')

		UserSettings.Remove(saveName)
		sort.SetDefaultSort()
		Assert(UserSettings.Get(saveName) is: #([dir: -1, col: 'a', id?: false]))

		UserSettings.Put(saveName, 'a')
		sort = VirtualListSortModel(.table, .sf, saveName, loadAll?:)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a')

		sort.SetSort('a')
		sort.SetDefaultSort()
		Assert(UserSettings.Get(saveName) is: #([dir: -1, col: 'a', id?: false]))

		sort.SetSort('b')
		sort.SetDefaultSort()
		Assert(UserSettings.Get(saveName) is: #([dir: 1, col: 'b', id?: false],
			[dir: -1, col: 'a', id?: false]))

		sort.SetSort('b')
		sort.SetDefaultSort()
		Assert(UserSettings.Get(saveName) is: #([dir: -1, col: 'b', id?: false],
			[dir: -1, col: 'a', id?: false]))

		sort.SetSort('a')
		sort.SetDefaultSort()
		Assert(UserSettings.Get(saveName) is: #([dir: 1, col: 'a', id?: false],
			[dir: -1, col: 'b', id?: false]))

		sort = VirtualListSortModel(.table, .sf, saveName, loadAll?:)
		Assert(sort.VirtualListSortModel_getSortStr() is:
			' extend reverse__b = true' $
			' sort a,reverse__b')

		sort = VirtualListSortModel(.table $ '  sort reverse a', .sf, saveName, loadAll?:)
		Assert(sort.VirtualListSortModel_getSortStr() is:
			' extend reverse__b = true' $
			' sort a,reverse__b')

		sort = VirtualListSortModel(.table $ '  sort reverse b', .sf, saveName, loadAll?:)
		Assert(sort.VirtualListSortModel_getSortStr() is:
			' extend reverse__b = true' $
			' sort a,reverse__b')

		sort.ResetSort()
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse b')
		Assert(UserSettings.Get(saveName) is: false)
		}

	Test_sortSaveName()
		{
		saveName = .TempName()
		UserSettings.Put(saveName, 'a')
		.AddTeardown({ UserSettings.Remove(saveName) })

		UserSettings.Put(saveName, 'xxxxx')
		sort = VirtualListSortModel(.table, .sf, saveName)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' ')

		UserSettings.Put(saveName, 'xxxxx,a')
		sort = VirtualListSortModel(.table, .sf, saveName)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a')

		UserSettings.Put(saveName, '-xxxxx,-a')
		sort = VirtualListSortModel(.table, .sf, saveName)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse a')

		UserSettings.Remove(saveName)
		sort.SetDefaultSort()
		Assert(UserSettings.Get(saveName) is: #([dir: -1, col: 'a', id?: false]))

		UserSettings.Put(saveName, 'a')
		sort = VirtualListSortModel(.table, .sf, saveName)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a')

		sort.SetSort('a')
		sort.SetDefaultSort()
		Assert(UserSettings.Get(saveName) is: #([dir: -1, col: 'a', id?: false]))

		sort.SetSort('b')
		sort.SetDefaultSort()
		Assert(UserSettings.Get(saveName) is: #([dir: 1, col: 'b', id?: false]))

		sort.SetSort('b')
		sort.SetDefaultSort()
		Assert(UserSettings.Get(saveName) is: #([dir: -1, col: 'b', id?: false]))

		sort.SetSort('a')
		sort.SetDefaultSort()
		Assert(UserSettings.Get(saveName) is: #([dir: 1, col: 'a', id?: false]))

		sort = VirtualListSortModel(.table, .sf, saveName)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a')

		sort = VirtualListSortModel(.table $ '  sort reverse a', .sf, saveName)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a')

		sort = VirtualListSortModel(.table $ '  sort reverse b', .sf, saveName)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a')

		sort.ResetSort()
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse b')
		Assert(UserSettings.Get(saveName) is: false)
		}

	Test_model_load_all()
		{
		sort = VirtualListSortModel(.table, .sf, loadAll?:)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' ')

		sort.SetSort('a')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a')

		sort.SetSort('a')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse a')

		sort.SetSort('a')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a')

		sort.SetSort('b')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort b,a')

		sort.SetSort('b')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is:
			' extend reverse__a = true ' $
			'sort reverse b,reverse__a')

		sort.SetSort('b')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort b,a')

		sort.SetSort('c')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort c,b')

		sort.SetSort('b')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort b,c')

		sort.SetSort('b')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is:
			' extend reverse__c = true' $
			' sort reverse b,reverse__c')

		sort.SetSort('c')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is:
			' extend reverse__b = true' $
			' sort c,reverse__b')
		}

	Test_model()
		{
		sort = VirtualListSortModel(.table, .sf)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' ')

		sort.SetSort('a')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a')

		sort.SetSort('a')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse a')

		sort.SetSort('a')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort a')

		sort.SetSort('b')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort b')

		sort.SetSort('b')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse b')

		sort.SetSort('b')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort b')

		sort.SetSort('c')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort c')

		sort.SetSort('b')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort b')

		sort.SetSort('b')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse b')

		sort.SetSort('c')
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort c')
		}

	Test_UsingDisplayCol()
		{
		sort = VirtualListSortModel(.table, .sf)
		Assert(sort.StripSort(.table) is: .table)
		Assert(sort.VirtualListSortModel_getSortStr() is: ' ')

		sort.SetSort('a', 'fred')
		Assert(sort.GetPrimarySort()
			is: [displayCol: 'a', dir: 1, col: "fred", id?: false])
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort fred')

		sort.SetSort('a', 'fred')
		Assert(sort.GetPrimarySort()
			is: [displayCol: 'a', dir: -1, col: "fred", id?: false])
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse fred')

		sort.SetSort('b')
		Assert(sort.GetPrimarySort() is: [dir: 1, col: "b", id?: false])
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort b')

		sort.SetSort('b')
		Assert(sort.GetPrimarySort() is: [dir: -1, col: "b", id?: false])
		Assert(sort.VirtualListSortModel_getSortStr() is: ' sort reverse b')
		}

	Test_sort_num_name_abbrev()
		{
		saveName = .TempName()
		UserSettings.Put(saveName, 'vl_sort_test_num')
		.AddTeardown({ UserSettings.Remove(saveName) })

		sort = VirtualListSortModel(.table, .sf, saveName)
		Assert(query = sort.BuildQuery(.table)
			is: '/* tableHint: ' $ .table $ ' */ ' $ .table $
				" leftjoin by(vl_sort_test_num) (" $
				.masterTable $ " project vl_sort_test_num, vl_sort_test_name, " $
				"vl_sort_test_abbrev)  sort vl_sort_test_name")

		sort.SetSort('vl_sort_test_num')
		Assert(query = sort.BuildQuery(query)
			is:  '/* tableHint: ' $ .table $ ' */ ' $ .table $
				" leftjoin by(vl_sort_test_num) (" $
				.masterTable $ " project vl_sort_test_num, vl_sort_test_name, " $
				"vl_sort_test_abbrev)  sort reverse vl_sort_test_name")

		sort.SetSort('vl_sort_test_num')
		Assert(query = sort.BuildQuery(query)
			is:  '/* tableHint: ' $ .table $ ' */ ' $ .table $
				" leftjoin by(vl_sort_test_num) (" $
				.masterTable $ " project vl_sort_test_num, vl_sort_test_name, " $
				"vl_sort_test_abbrev)  sort vl_sort_test_name")

		sort.SetSort('vl_sort_test2_num')
		Assert(query = sort.BuildQuery(query)
			is:  '/* tableHint: ' $ .table $ ' */ ' $ .table $
				" leftjoin by(vl_sort_test2_num) (" $
				.masterTable2 $ " project vl_sort_test2_num, vl_sort_test2_name, " $
				"vl_sort_test2_abbrev)  sort vl_sort_test2_name")

		sort.SetSort('vl_sort_test2_num')
		Assert(query = sort.BuildQuery(query)
			is:  '/* tableHint: ' $ .table $ ' */ ' $ .table $
				" leftjoin by(vl_sort_test2_num) (" $
				.masterTable2 $ " project vl_sort_test2_num, vl_sort_test2_name, " $
				"vl_sort_test2_abbrev)  sort reverse vl_sort_test2_name")

		sort.SetSort('c')
		Assert(query = sort.BuildQuery(.table)
			is: .table $ " sort c")

		sort.SetSort('vl_sort_test2_num')
		query = .table $ " leftjoin by(vl_sort_test2_num) (" $ .masterTable2 $
			" project vl_sort_test2_num, vl_sort_test2_name, vl_sort_test2_abbrev)" $
			"  sort reverse vl_sort_test2_name"
		Assert(query = sort.BuildQuery(query)
			is: .table $ " leftjoin by(vl_sort_test2_num) (" $ .masterTable2 $
			" project vl_sort_test2_num, vl_sort_test2_name, vl_sort_test2_abbrev)" $
			" sort vl_sort_test2_name")

		sort.SetSort('c')
		Assert(query = sort.BuildQuery(query)
			is: .table $ " leftjoin by(vl_sort_test2_num) (" $ .masterTable2 $
			" project vl_sort_test2_num, vl_sort_test2_name, vl_sort_test2_abbrev)" $
			" sort c")
		}

	Test_sort_name_on_join()
		{
		saveName = .TempName()
		UserSettings.Put(saveName, 'vl_sort_test_name')
		.AddTeardown({ UserSettings.Remove(saveName) })

		baseQuery = "/* tableHint: " $ .table $ " */ " $ .table $
			" join by(vl_sort_test_num) " $ .masterTable
		sort = VirtualListSortModel(baseQuery, .sf, saveName)
		Assert(query = sort.BuildQuery(baseQuery)
			is: baseQuery $ " sort vl_sort_test_name")

		sort.SetSort('vl_sort_test_abbrev')
		Assert(query = sort.BuildQuery(query)
			is: baseQuery $ " sort vl_sort_test_abbrev")

		sort.SetSort('vl_sort_test_num')
		Assert(query = sort.BuildQuery(query)
			is: baseQuery $ " sort vl_sort_test_name")

		sort.SetSort('vl_sort_test2_num')
		Assert(query = sort.BuildQuery(query)
			is: baseQuery $
				' leftjoin by(vl_sort_test2_num) (' $
				.masterTable2 $
				' project vl_sort_test2_num, vl_sort_test2_name, vl_sort_test2_abbrev)' $
				'  sort vl_sort_test2_name')
		}

	Test_SortInMemory()
		{
		m = VirtualListSortModel.SortInMemory
		data = Object(
			x21 = #(a: 2, b: 1),
			x32 = #(a: 3, b: 2),
			x12 = #(a: 1, b: 2),
			x22 = #(a: 2, b: 2),
			x11 = #(a: 1, b: 1),
			x31 = #(a: 3, b: 1)
			)

		m(data, 'a, b')
		Assert(data
			is: Object(x11, x12, x21, x22, x31, x32))

		m(data, 'reverse a, b')
		Assert(data is: Object(x32, x31, x22, x21, x12, x11))

		m(data, 'a, reverse__b')
		Assert(data is: Object(x12, x11, x22, x21, x32, x31))

		m(data, 'reverse a, reverse__b')
		Assert(data is: Object(x31, x32, x21, x22, x11, x12))

		m(data, 'reverse__a, b')
		Assert(data is: Object(x31, x32, x21, x22, x11, x12))

		m(data, 'reverse reverse__a, b')
		Assert(data is: Object(x12, x11, x22, x21, x32, x31))
		}

	Test_StripSort()
		{
		q = `/* tableHint: ` $ .table $ ` */ ` $ .table $
			` where a isnt "" ` $
			`/* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/ ` $
			`leftjoin by(vl_sort_test_num) ` $
			`(` $ .masterTable $ ` rename vl_sort_test_name to vl_sort_test_name_ren)`
		sort = ` sort vl_sort_test_num`
		result = VirtualListSortModel.StripSort(q $ sort)
		Assert(result
			is: `/* tableHint: ` $ .table $ ` */ ` $ .table $
			` where a isnt "" /* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/ ` $
			`leftjoin by(vl_sort_test_num) ` $
			`(` $ .masterTable $ ` rename vl_sort_test_name to vl_sort_test_name_ren)`,
			msg: 'join false')

		cl = VirtualListSortModel { VirtualListSortModel_join?: true }
		result = cl.StripSort(q)
		Assert(result
			is: `/* tableHint: ` $ .table $ ` */ ` $ .table $
			` where a isnt "" ` $
			`/* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/`, msg: 'join true')

		query = 'tables /* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/ ' $
			'extend reverse__nrows sort table'
		result = VirtualListSortModel.StripSort(query)
		Assert(result is: 'tables /* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/',
			msg: 'reverse__')
		}
	}