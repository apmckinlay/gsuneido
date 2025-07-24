// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	title: 'UserColumns_Test'
	Setup()
		{
		.TearDownIfTablesNotExist('usercolumns')
		}
	dummylist: class
		{
		GetColWidth(i)
			{ return .Wids[i] }
		SetColWidth(i, w)
			{ .Wids[i] = w }
		GetColumns()
			{ return .Cols }
		SetColumns(cols, reset /*unused*/ = false)
			{ .Cols = cols.Copy(); .Wids = Object().AddMany!(9, cols.Size()) }
		SetHeaderChanged(status /*unused*/)
			{}
		}
	Test_main()
		{
		list = new .dummylist
		list.Cols = original_cols = #(a, b, c, d)
		list.Wids = original_wids = #(1, 2, 3, 4)

		UserColumns.Save(.title, list, original_cols) // create
		Assert(Query1(.query()) is: Record(
			usercolumns_user: Suneido.User,
			usercolumns_title: .title,
			usercolumns_order: original_cols,
			usercolumns_sizes: original_wids))

		list.Cols = new_cols = #(b, a, c, d)
		list.Wids = new_wids = #(5, 6, 0, 7)
		UserColumns.Save(.title, list, original_cols) // update
		Assert(Query1(.query()) is: Record(
			usercolumns_user: Suneido.User,
			usercolumns_title: .title,
			usercolumns_order: new_cols,
			usercolumns_sizes: new_wids))

		// load customized columns
		list.Cols = original_cols
		list.Wids = original_wids
		UserColumns.Load(original_cols, .title, list, load_visible?:)
		Assert(list.Cols is: #(b, a, d))
		Assert(list.Wids is: #(5, 6, 7))

		// load after column added
		list.Cols = original_cols
		list.Wids = original_wids
		UserColumns.Load(#(a,b,nu,c,d), .title, list, load_visible?:)
		Assert(list.Cols is: #(b, a, d, nu))
		Assert(list.Wids is: #(5, 6, 7, false))

		// load after column added with option to hide added columns turned on
		UserColumns.Load(#(a,b,nu,c,d), .title, list, load_visible?:,
			hideColumnsNotSaved?:)
		Assert(list.Cols is: #(b, a, d))
		Assert(list.Wids is: #(5, 6, 7))

		// load after column added
		list.Cols = original_cols
		list.Wids = original_wids
		UserColumns.Load(#(a,b,custom_999999,c,d), .title, list, load_visible?:)
		Assert(list.Cols is: #(b, a, d))
		Assert(list.Wids is: #(5, 6, 7))

		// load after column deleted
		list.Cols = original_cols
		list.Wids = original_wids
		UserColumns.Load(#(a,b,d), .title, list, load_visible?:)
		Assert(list.Cols is: #(b, a, d))
		Assert(list.Wids is: #(5, 6, 7))

		.delete_user_columns()
		}

	Test_Duplicates()
		{
		list = new .dummylist
		list.Cols = #(a, b, c, a, d)
		list.Wids = #(1, 2, 3, 1, 4)
		UserColumns.Save(.title, list, list.Cols) // create
		Assert(Query1(.query()) is: Record(
			usercolumns_user: Suneido.User,
			usercolumns_title: .title,
			usercolumns_order: #(a, b, c, d),
			usercolumns_sizes: #(1, 2, 3, 4)))

		.delete_user_columns()
		}

	Test_ignoreExtraCols()
		{
		list = new .dummylist
		list.Cols = #(select, a, b, c, d)
		list.Wids = #(75, 1, 2, 3, 4)
		UserColumns.Save(.title, list, list.Cols, ignoreCols: #(select)) // create
		Assert(Query1(.query()) is: Record(
			usercolumns_user: Suneido.User,
			usercolumns_title: .title,
			usercolumns_order: #(a, b, c, d),
			usercolumns_sizes: #(1, 2, 3, 4)))

		UserColumns.Load(#(a, b, c, d), .title, list,
			extraCols: #((field: 'select', pos: 0, width: 75)))
		Assert(list.GetColumns() is: #(select, a, b, c, d))

		// no default saved so just returns list.Columns
		UserColumns.Reset(list, .title, list.Cols)
		Assert(list.GetColumns() is: #(select, a, b, c, d))

		// create default columns
		UserColumns.SaveDefaultColumns(.title, #(a, b, c), #(a, b, c, d))
		UserColumns.Reset(list, .title, list.Cols,
			extraCols: #((field: 'select', pos: 0, width: 75)))
		Assert(list.GetColumns() is: #(select, a, b, c, d))

		.delete_user_columns()
		}

	Test_default_Columns()
		{
		list = new .dummylist
		list.Cols = #(a, b, c, a, d)
		list.Wids = #(1, 2, 3, 1, 4)
		UserColumns.Save(.title, list, list.Cols) // create
		.delete_user_columns()

		UserColumns.SaveDefaultColumns(.title, #(a, b, c), #(a, b, c, d, e))
		Assert(Query1(.query('')) is: Record(
			usercolumns_title: .title,
			usercolumns_order: #(a, b, c, d, e),
			usercolumns_sizes: #(100, 100, 100, 0, 0)))

		UserColumns.SaveDefaultColumns(.title, #(a, b), #(a, b))
		Assert(Query1(.query('')) is: Record(
			usercolumns_title: .title,
			usercolumns_order: #(a, b),
			usercolumns_sizes: #(100, 100)))

		UserColumns.SaveDefaultColumns(.title, #(a, b), #(a, b, c))
		Assert(Query1(.query('')) is: Record(
			usercolumns_title: .title,
			usercolumns_order: #(a, b, c),
			usercolumns_sizes: #(100, 100, 0)))

		.delete_user_columns()
		}

	Test_update_default_columns()
		{
		list = new .dummylist
		list.Cols = #(a, b, c)
		list.Wids = #(1, 2, 3)
		.delete_user_columns()

		custom_fields = #(e, d)
		UserColumns.AddCustomFields(.title, list, custom_fields,
			#(a, b, c), #(a, b, c, d, e))
		Assert(Query1(.query('')) is: Record(
			usercolumns_title: .title,
			usercolumns_order: #(a, b, c, e, d),
			usercolumns_sizes: #(100, 100, 100, 100, 100)))

		Assert(Query1(.query()) is: Record(
			usercolumns_sizes: #(1, 2, 3, false, false),
			usercolumns_title: .title,
			usercolumns_order: #(a, b, c, e, d),
			usercolumns_user: Suneido.User))

		.delete_user_columns()

		// Test duplicate custom_fields
		custom_fields = #(e, d, c)
		UserColumns.AddCustomFields(.title, list, custom_fields,
			#(a, b, c), #(a, b, c, d, e))
		Assert(Query1(.query('')) is: Record(
			usercolumns_title: .title,
			usercolumns_order: #(a, b, c, e, d),
			usercolumns_sizes: #(100, 100, 100, 100, 100)))

		Assert(Query1(.query()) is: Record(
			usercolumns_sizes: #(1, 2, 3, false, false),
			usercolumns_title: .title,
			usercolumns_order: #(a, b, c, e, d),
			usercolumns_user: Suneido.User))

		.delete_user_columns()
		}

	Test_hidden_columns()
		{
		list = new .dummylist
		list.Cols = #(a, b, c)
		list.Wids = #(1, 2, 0)
		.delete_user_columns()

		custom_fields = #(e, d)
		UserColumns.AddCustomFields(.title, list, custom_fields,
			#(a, b), #(a, b, c, d, e))

		Assert(Query1(.query('')) is: Record(
			usercolumns_title: .title,
			usercolumns_order: #(a, b, e, d, c),
			usercolumns_sizes: #(100, 100, 100, 100, 0)))

		Assert(Query1(.query()) is: Record(
			usercolumns_sizes: #(1, 2, 0, false, false),
			usercolumns_title: .title,
			usercolumns_order: #(a, b, c, e, d),
			usercolumns_user: Suneido.User))

		.delete_user_columns()
		}

	Test_hide_and_add_custom_columns()
		{
		list = new .dummylist
		list.Cols = #(a, b, c)
		// current user resizes the 3rd column to 0
		list.Wids = #(1, 2, 0)
		.delete_user_columns()

		custom_fields = #(e, d)
		UserColumns.AddCustomFields(.title, list, custom_fields,
			#(a, b, c), #(a, b, c, d, e))

		// other users still should see the 3rd column
		Assert(Query1(.query('')) is: Record(
			usercolumns_title: .title,
			usercolumns_order: #(a, b, c, e, d),
			usercolumns_sizes: #(100, 100, 100, 100, 100)))

		// current user should not see the 3rd column
		Assert(Query1(.query()) is: Record(
			usercolumns_sizes: #(1, 2, 0, false, false),
			usercolumns_title: .title,
			usercolumns_order: #(a, b, c, e, d),
			usercolumns_user: Suneido.User))

		.delete_user_columns()
		}

	Test_columns()
		{
		list = new .dummylist
		list.Cols = original_cols = #(a, b, c, d)
		list.Wids = #(1, 2, 3, 4)
		.delete_user_columns()

		UserColumns.SaveDefaultColumns(.title, #(a, b), #(a, b, c, d))
		Assert(Query1(.query('')) is: Record(
			usercolumns_title: .title,
			usercolumns_order: #(a, b, c, d),
			usercolumns_sizes: #(100, 100, 0, 0)))

		UserColumns.Load(#(a, b, c, d), .title, list, load_visible?:)
		Assert(list.Cols is: #(a, b))
		Assert(list.Wids is: #(false, false))

		list.Wids = #(5, 6)
		list.Cols = #(a, b)
		UserColumns.Save(.title, list, original_cols)
		Assert(Query1(.query()) is: Record(
			usercolumns_user: Suneido.User,
			usercolumns_title: .title,
			usercolumns_order: original_cols,
			usercolumns_sizes: #(5, 6, 0, 0)))

		UserColumns.Load(#(a, b, c, d, f), .title, list, load_visible?:)
		Assert(list.Cols is: #(a, b, f))
		Assert(list.Wids is: #(5, 6, false))

		UserColumns.Reset(list, .title, #(a, b, c, d, f), load_visible?:)
		Assert(list.Cols is: #(a, b, f))
		Assert(list.Wids is: #(false, false, false))

		.delete_user_columns()
		}

	Test_invalid_saved_columns()
		{
		list = new .dummylist

		// c and d no longer exist in columns.
		// No default columns
		// Displays Columns not saved
		list.Cols = #(c, d)
		list.Wids = #(1, 2)
		UserColumns.Save(.title, list, #(c, d))
		UserColumns.Load(#(a, b), .title, list, load_visible?:,
			hideColumnsNotSaved?: false)
		Assert(list.Cols is: #(a, b))
		Assert(list.Wids is: #(false, false))
		Assert(Query1(.query()) isnt: false)

		// c and d no longer exist in columns.
		// No default columns
		// Hides Columns not saved
		list.Cols = #(c, d)
		list.Wids = #(1, 2)
		UserColumns.Save(.title, list, #(c, d))
		UserColumns.Load(#(a, b), .title, list, load_visible?:,
			hideColumnsNotSaved?:)
		Assert(list.Cols is: #(a, b))
		Assert(list.Wids is: #(false, false))
		Assert(Query1(.query()) is: false)

		// c and d no longer exist in columns.
		// default columns: display a, hide b
		// Displays Columns not saved
		list.Cols = #(c, d)
		list.Wids = #(1, 2)
		UserColumns.Save(.title, list, #(c, d))
		UserColumns.SaveDefaultColumns(.title, #(a), #(a, b))
		UserColumns.Load(#(a, b), .title, list, load_visible?:,
			hideColumnsNotSaved?: false)
		Assert(list.Cols is: #(a, b))
		Assert(list.Wids is: #(false, false))
		Assert(Query1(.query()) isnt: false)
		Assert(Query1(.query("")) isnt: false)

		// c and d no longer exist in columns.
		// default columns: display a, hide b
		// Hides Columns not saved
		list.Cols = #(c, d)
		list.Wids = #(1, 2)
		UserColumns.Save(.title, list, #(c, d))
		UserColumns.SaveDefaultColumns(.title, #(a), #(a, b))
		UserColumns.Load(#(a, b), .title, list, load_visible?:,
			hideColumnsNotSaved?:)
		Assert(list.Cols is: #(a))
		Assert(list.Wids is: #(false))
		Assert(Query1(.query()) is: false)
		Assert(Query1(.query("")) isnt: false)

		// c and d no longer exist in columns.
		// default columns: display c, hide d
		// Displays Columns not saved
		list.Cols = #(c, d)
		list.Wids = #(1, 2)
		UserColumns.Save(.title, list, #(c, d))
		UserColumns.SaveDefaultColumns(.title, #(c), #(c, d))
		UserColumns.Load(#(a, b), .title, list, load_visible?:,
			hideColumnsNotSaved?: false)
		Assert(list.Cols is: #(a, b))
		Assert(list.Wids is: #(false, false))
		Assert(Query1(.query()) isnt: false)
		Assert(Query1(.query("")) isnt: false)

		// c and d no longer exist in columns.
		// default columns: display c, hide d
		// Hides Columns not saved
		list.Cols = #(c, d)
		list.Wids = #(1, 2)
		UserColumns.Save(.title, list, #(c, d))
		UserColumns.SaveDefaultColumns(.title, #(c), #(c, d))
		UserColumns.Load(#(a, b), .title, list, load_visible?:,
			hideColumnsNotSaved?:)
		Assert(list.Cols is: #(a, b))
		Assert(list.Wids is: #(9, 9))
		Assert(Query1(.query()) is: false)
		Assert(Query1(.query("")) is: false)

		// c and d no longer exist in columns.
		// default columns: display b, hide c and d
		// Displays Columns not saved
		list.Cols = #(b, c, d)
		list.Wids = #(1, 2, 3)
		UserColumns.Save(.title, list, #(b, c, d))
		UserColumns.SaveDefaultColumns(.title, #(b, c), #(b, c, d))
		UserColumns.Load(#(a, b), .title, list, load_visible?:,
			hideColumnsNotSaved?: false)
		Assert(list.Cols is: #(b, a))
		Assert(list.Wids is: #(1, false))
		Assert(Query1(.query()) isnt: false)
		Assert(Query1(.query("")) isnt: false)

		// c and d no longer exist in columns.
		// default columns: display b, hide c and d
		// Hides Columns not saved
		list.Cols = #(b, c, d)
		list.Wids = #(1, 2, 3)
		UserColumns.Save(.title, list, #(b, c, d))
		UserColumns.SaveDefaultColumns(.title, #(b, c), #(b, c, d))
		UserColumns.Load(#(a, b), .title, list, load_visible?:,
			hideColumnsNotSaved?:)
		Assert(list.Cols is: #(b))
		Assert(list.Wids is: #(1))
		Assert(Query1(.query()) isnt: false)
		Assert(Query1(.query("")) isnt: false)

		// c and d no longer exist in columns.
		// default columns: display b, hide c
		// Displays Columns not saved
		list.Cols = #(c, d)
		list.Wids = #(1, 2)
		UserColumns.Save(.title, list, #(c, d))
		UserColumns.SaveDefaultColumns(.title, #(b, c), #(b, c))
		UserColumns.Load(#(a, b), .title, list, load_visible?:,
			hideColumnsNotSaved?: false)
		Assert(list.Cols is: #(a, b))
		Assert(list.Wids is: #(false, false))
		Assert(Query1(.query()) isnt: false)
		Assert(Query1(.query("")) isnt: false)

		// c and d no longer exist in columns.
		// default columns: display b, hide c
		// Displays Columns not saved
		list.Cols = #(c, d)
		list.Wids = #(1, 2)
		UserColumns.Save(.title, list, #(c, d))
		UserColumns.SaveDefaultColumns(.title, #(b, c), #(b, c))
		UserColumns.Load(#(a, b), .title, list, load_visible?:,
			hideColumnsNotSaved?:)
		Assert(list.Cols is: #(b))
		Assert(list.Wids is: #(false))
		Assert(Query1(.query()) is: false)
		Assert(Query1(.query("")) isnt: false)
		}

	Test_load_all_columns()
		{
		list = new .dummylist
		list.Cols = original_cols = #(a, b, c, d)
		list.Wids = original_wids = #(1, 2, 3, 4)

		UserColumns.Save(.title, list) // create
		Assert(Query1(.query()) is: Record(
			usercolumns_user: Suneido.User,
			usercolumns_title: .title,
			usercolumns_order: original_cols,
			usercolumns_sizes: original_wids))

		list.Cols = new_cols = #(b, a, c, d)
		list.Wids = new_wids = #(5, 6, 0, 7)
		UserColumns.Save(.title, list) // update
		Assert(Query1(.query()) is: Record(
			usercolumns_user: Suneido.User,
			usercolumns_title: .title,
			usercolumns_order: new_cols,
			usercolumns_sizes: new_wids))

		// load customized columns
		list.Cols = original_cols
		list.Wids = original_wids
		UserColumns.Load(original_cols, .title, list)
		Assert(list.Cols is: #(b, a, c, d))
		Assert(list.Wids is: #(5, 6, 0, 7))

		// load after column added
		list.Cols = original_cols
		list.Wids = original_wids
		UserColumns.Load(#(a,b,nu,c,d), .title, list)
		Assert(list.Cols is: #(b, a, c, d, nu))
		Assert(list.Wids is: #(5, 6, 0, 7, false))

		// load after column added
		list.Cols = original_cols
		list.Wids = original_wids
		UserColumns.Load(#(a,b,custom_999999,c,d), .title, list)
		Assert(list.Cols is: #(b, a, c, d, custom_999999))
		Assert(list.Wids is: #(5, 6, 0, 7, false))

		// load after column deleted
		list.Cols = original_cols
		list.Wids = original_wids
		UserColumns.Load(#(a,b,d), .title, list)
		Assert(list.Cols is: #(b, a, d))
		Assert(list.Wids is: #(5, 6, 7))

		.delete_user_columns()
		}

	Test_removeDeletedAndNonPermissableCustom()
		{
		f = UserColumns.UserColumns_removeDeletedAndNonPermissableCustom
		origCols = #(field1, custom1, custom_nonpermissible, custom2, customDeleted,
			field2)
		origWids = #(1, 2, 3, 4, 5, 6)
		visibleCols = origCols.Copy()
		visibleWids = origWids.Copy()
		f('', visibleCols, visibleWids)
		Assert(visibleCols is: origCols)
		Assert(visibleWids is: origWids)

		.SpyOn(Customizable.GetNonPermissableFields).Return(#(custom_nonpermissible))
		f('', visibleCols, visibleWids)
		Assert(visibleCols is: #(field1, custom1, custom2, customDeleted, field2))
		Assert(visibleWids is: #(1, 2, 4, 5, 6))

		visibleCols = origCols.Copy()
		visibleWids = origWids.Copy()
		.SpyOn(Customizable.DeletedField?).Return(false, false, false, false, true, false)
		f('', visibleCols, visibleWids)
		Assert(visibleCols is: #(field1, custom1, custom2, field2))
		Assert(visibleWids is: #(1, 2, 4, 6))
		}

	query(user = false)
		{
		return 'usercolumns
			remove usercolumns_TS
			where usercolumns_title is ' $ Display(.title) $ ' and
			usercolumns_user is ' $ Display(user is false ? Suneido.User : user)
		}

	delete_user_columns()
		{
		QueryDo('delete usercolumns where usercolumns_title is ' $ Display(.title))
		}

	Teardown()
		{
		.delete_user_columns()
		super.Teardown()
		}
	}
