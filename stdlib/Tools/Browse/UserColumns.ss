// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// can be used with List, Browse
// requires list have:
//		HeaderChanged?(),
//		GetColWidth(i), SetColWidth(i, w)
//		GetColumns(), SetColumns(cols, reset = false)
class
	{
	EnsureTable()
		{
		Database('ensure usercolumns (usercolumns_user, usercolumns_title,
			usercolumns_order, usercolumns_sizes, usercolumns_TS)
			index (usercolumns_title)
			key(usercolumns_user, usercolumns_title)')
		}
	deletedFlagWidth: 17
	defaultColWidth: 100
	GetDefaultColumns(title)
		{
		.EnsureTable()
		if false isnt rec = Query1(.query(title, ''))
			{
			sizes = Object()
			for(si =0; si < rec.usercolumns_sizes.Size(); si++)
				sizes.Add(rec.usercolumns_order[si] isnt 'listrow_deleted'
					? rec.usercolumns_sizes[si] is 0 ? 0 : false
					: .deletedFlagWidth)
			rec.usercolumns_sizes = sizes
			return rec
			}
		return false
		}
	SaveDefaultColumns(title, cols, available_cols, deletecol = false, t = false)
		{
		.EnsureTable()
		query = .query(title, '')

		remainder = available_cols.Difference(cols)
		new_cols = cols.Copy().Append(remainder)
		sizes = Object()
		new_cols.Remove('listrow_deleted')
		for c in new_cols
			sizes.Add(cols.Has?(c) ? .defaultColWidth : 0)
		if deletecol is true
			{
			new_cols.Add('listrow_deleted', at: 0)
			sizes.Add(.deletedFlagWidth, at: 0)
			}
		DoWithTran(t, update:)
			{ |t|
			if false is t.Query1(query)
				t.QueryOutput('usercolumns', Record(
					usercolumns_user: "",
					usercolumns_title: .truncate_title(title),
					usercolumns_order: new_cols,
					usercolumns_sizes: sizes))
			else
				t.QueryApply(query)
					{|x|
					x.usercolumns_order = new_cols
					x.usercolumns_sizes = sizes
					x.Update()
					}
			}
		}
	AddCustomFields(title, list, custom_fields, original_columns, available_cols,
		deletecol = false)
		{
		.EnsureTable()
		.append_default_cols(title, custom_fields, original_columns,
			available_cols, deletecol)
		.append_cols(title, list, available_cols, custom_fields, deletecol)
		}
	append_default_cols(title, custom_fields, original_columns, available_cols,
		deletecol)
		{
		Transaction(update:)
			{ |t|
			if false is rec = t.Query1(.query(title, ''))
				cols = original_columns.Copy().Append(custom_fields).UniqueValues()
			else
				{
				cols = Object()
				for(i = 0; i < rec.usercolumns_order.Size(); i++)
					if rec.usercolumns_sizes[i] > 0
						cols.Add(rec.usercolumns_order[i])
				cols = cols.Append(custom_fields).UniqueValues()
				}
			.SaveDefaultColumns(title, cols, available_cols, deletecol, t)
			}
		}
	append_cols(title, list, available_cols, custom_fields, deletecol)
		{
		result = .cols_and_sizes(list, available_cols, custom_fields)
		.save_columns(title, result)

		saved = .get_saved_rec(title)
		.set_visible_columns(saved, list, available_cols, deletecol)
		list.SetHeaderChanged(true)
		}
	Save(title, list, columns = false, ignoreCols = #())
		{
		.EnsureTable()
		if list.Method?('HeaderChanged?') and not list.HeaderChanged?()
			return

		if columns is false
			columns = list.GetColumns()
		result = .cols_and_sizes(list, columns, :ignoreCols)
		.save_columns(title, result)
		}
	save_columns(title, colsOb)
		{
		RetryTransaction()
			{ |t|
			if (false is x = .get_saved_rec(title, t))
				t.QueryOutput('usercolumns', Record(
					usercolumns_user: Suneido.User,
					usercolumns_title: .truncate_title(title),
					usercolumns_order: colsOb.cols,
					usercolumns_sizes: colsOb.sizes))
			else
				{
				x.usercolumns_order = colsOb.cols
				x.usercolumns_sizes = colsOb.sizes
				x.Update()
				}
			}
		}
	cols_and_sizes(list, available_cols, custom_fields = #(), ignoreCols = #())
		{
		visibleWids = Object()
		visibleCols = Object()
		listCols = list.GetColumns()
		for(i = 0; i < listCols.Size(); i++)
			if not ignoreCols.Has?(listCols[i]) and
				not visibleCols.Has?(listCols[i])
				{
				visibleCols.Add(listCols[i])
				visibleWids.Add(list.GetColWidth(i))
				}

		visibleCols.MergeUnion(custom_fields)
		visibleWids.AddMany!(false, visibleCols.Size() - visibleWids.Size())

		hiddenCols = available_cols.UniqueValues().
			Difference(visibleCols).Difference(ignoreCols)
		visibleWids.AddMany!(0, hiddenCols.Size())
		cols = visibleCols.MergeUnion(hiddenCols)

		return Object(sizes: visibleWids, :cols)
		}
	Reset(list, title, columns, deletecol = false, load_visible? = false,
		extraCols = #(), permissableQuery = false, hideColumnsNotSaved? = false)
		{
		if list is false
			return

		.EnsureTable()
		saved = .GetDefaultColumns(title)
		visible = .set_visible_columns(saved, list, columns.Copy(), deletecol,
			load_visible?, :extraCols, :permissableQuery, :hideColumnsNotSaved?)
		if visible.Empty?()
			{
			SuneidoLog('INFO: Deleted default usercolumns with no visible columns',
				params: saved, calls:)
			QueryDo('delete usercolumns
				where usercolumns_title is ' $ Display(title) $
				' and usercolumns_user is ""')
			list.SetColumns(columns)
			list.SetHeaderChanged(true)
			return
			}

		if saved is false
			for ci in columns.Members()
				list.SetColWidth(ci, false)
		list.SetHeaderChanged(true)
		}
	Load(columns, title, list, deletecol = false, initialized? = false,
		load_visible? = false, extraCols = #(), permissableQuery = false,
		hideColumnsNotSaved? = false)
		{
		if list is false
			return

		.EnsureTable()
		columns = columns.Copy()
		if deletecol is false
			deletecol = columns.Has?('listrow_deleted')

		if false is saved = .get_saved_rec(title)
			saved = .GetDefaultColumns(title)

		visibleCols = .set_visible_columns(saved, list, columns, deletecol, load_visible?,
			:extraCols, :permissableQuery, :hideColumnsNotSaved?)

		if visibleCols.Empty?()
			{
			QueryDo('delete usercolumns
				where usercolumns_title is ' $ Display(title) $
				' and usercolumns_user is ' $ Display(Suneido.User))
			.Reset(list, title, columns, deletecol, load_visible?, extraCols,
				permissableQuery, :hideColumnsNotSaved?)
			return
			}

		if saved is false
			.hide_columns(list, columns, initialized?)
		}
	set_visible_columns(saved, list, columns, deletecol, load_visible? = false,
		extraCols = #(), permissableQuery = false, hideColumnsNotSaved? = false)
		{
		if saved is false
			{
			list.SetColumns(columns)
			return columns
			}

		visibleCols = Object()
		visibleWids = Object()
		savedCols = saved.usercolumns_order
		savedWids = saved.usercolumns_sizes
		extraColFields = extraCols.Map({ it.field })
		for(i = 0; i < savedCols.Size(); i++)
			{
			col = savedCols[i]
			if col isnt 'listrow_deleted' and columns.Has?(col) and
				.visible?(col, visibleCols, savedWids[i], load_visible?, extraColFields)
				{
				visibleCols.Add(savedCols[i])
				visibleWids.Add(savedWids[i])
				}
			}

		added = .getAddedColumns(hideColumnsNotSaved?, columns, savedCols, extraColFields,
			load_visible?)

		visibleCols.Add(@added)
		visibleWids.AddMany!(false, added.Size())

		if deletecol
			{
			visibleCols.Add('listrow_deleted', at: 0)
			visibleWids.Add(.deletedFlagWidth, at: 0)
			}
		.addExtraCols(extraCols, visibleCols, visibleWids)
		.removeDeletedAndNonPermissableCustom(permissableQuery, visibleCols, visibleWids)
		.setColumns(list, visibleCols, visibleWids)
		return visibleCols
		}

	visible?(col, visibleCols, width, load_visible?, extraColFields)
		{
		return not visibleCols.Has?(col) and
			(width is false or width > 0 or load_visible? is false) and
			not extraColFields.Has?(col)
		}

	getAddedColumns(hideColumnsNotSaved?, columns, savedCols, extraColFields,
		load_visible?)
		{
		added = Object()
		if not hideColumnsNotSaved?
			added = columns.Difference(savedCols).Difference(extraColFields)

		if load_visible? is true
			added = added.Filter({|x| x !~ "^custom_[0-9]+$" })

		return added
		}

	addExtraCols(extraCols, visibleCols, visibleWids)
		{
		for col in extraCols
			{
			pos = not col.Member?('pos') or col.pos is 'end'
				? extraCols.Size() - 1 : col.pos
			width = not col.Member?('width') ? .defaultColWidth : col.width
			visibleCols.Add(col.field, at: pos)
			visibleWids.Add(width, at: pos)
			}
		}

	setColumns(list, visibleCols, visibleWids)
		{
		list.SetColumns(visibleCols)
		for(i = 0; i< visibleCols.Size(); i++)
			list.SetColWidth(i, visibleWids[i])
		}

	removeDeletedAndNonPermissableCustom(permissableQuery, visibleCols, visibleWids)
		{
		nonPermissable = Customizable.GetNonPermissableFields(permissableQuery)
		for field in visibleCols.Copy()
			if ((Customizable.DeletedField?(field) or nonPermissable.Has?(field)) and
				false isnt pos = visibleCols.Find(field))
				{
				visibleCols.Delete(pos)
				visibleWids.Delete(pos)
				}
		}

	hide_columns(list, columns, initialized? = false)
		{
		if not initialized?
			for col in columns.Members()
				if Customizable.CustomField?(col)
					list.SetColWidth(col, 0)
		}
	truncate_title(title)
		{
		// truncate title so key on usercolumns doesn't get too large
		return TruncateKey(title, replace: "", length: 100)
		}
	query(title, user = false)
		{
		title = .truncate_title(title)

		if title is ""
			SuneidoLog('ERROR: UserColumns has empty title', calls:)

		if user is false
			user = Suneido.User
		return 'usercolumns
			where usercolumns_user is '	$ Display(user) $
			' and usercolumns_title is ' $ Display(title)
		}
	get_saved_rec(title, t = false)
		{
		DoWithTran(t)
			{|t|
			result = t.Query1(.query(title))
			}
		return result
		}
	UpdateSavedWithNewColumns(columns, title)
		{
		QueryApplyMulti(`usercolumns where usercolumns_title is ` $ Display(title),
			update:)
			{
			for col in columns
				{
				it.usercolumns_order.Add(col)
				it.usercolumns_sizes.Add(0)
				}
			it.Update()
			}
		}
	ReplaceColumn(title, oldCol, newCol)
		{
		QueryApplyMulti(`usercolumns where usercolumns_title is ` $ Display(title),
			update:)
			{
			it.usercolumns_order = it.usercolumns_order.Replace(oldCol, newCol)
			it.Update()
			}
		}
	}
