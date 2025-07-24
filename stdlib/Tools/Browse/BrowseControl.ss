// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// Note: Default method forwards to list
// i.e. you can call ListControl methods on a BrowseControl
PassthruController
	{
	Name:				"Browse"
	Xstretch:			1
	Ystretch:			1
	value:				false		// from Set
	query:				false
	columns:			false
	query_columns: 		()
	linkField:			false
	validationField:	false
	protectField:		false
	headerData:			false
	report:				false
	dataMember:			"browse_data"
	allDataMember:		"all_browse_data"

	New(query, columns = false, linkField = false, stickyFields = false,
		validField = false, protectField = false, .title = false,
		statusBar = false, .headerFields = #(), dataMember = "browse_data",
		.noShading = false, .noHeaderButtons = false, notifyLast = false,
		.columnsSaveName = '', .buttonBar = false, .headerSelectPrompt = false,
		.mandatoryFields = #(), stretch = false, .alwaysReadOnly = false, name = false,
		extra = false, hideContents = false, .noSaveOnDestroy = false,
		.loadRecordNotification = false, .addons = #())
		{
		super(.createControls(statusBar, stretch, extra))
		.addons = .addons.Copy().Append(GetContributions('FormatContextMenuItems'))
		.Addons = AddonManager(this, .addons)
		// REFACTOR: move this into .addons
		.addon = BrowseAddon(this, options: false)

		.initObserverListRow()

		columns = .initVisibleBrowseComponents(query, statusBar, columns)

		if ReadOnlyAccess(this) or .alwaysReadOnly
			.SetReadOnly(true)
		.stickyFields = StickyFields(stickyFields)
		.SetValidationField(validField)
		.SetProtectField(protectField)
		.SetLinkField(linkField)
		.setCustomization(name, query)
		.setDataMember(dataMember)
		.SetQuery(query, columns)
		.Send("Data")

		.initValidationNotifyAndObserver(linkField, notifyLast)
		.SetMenu()
		if hideContents
			{
			.list.HideContent(true)
			.list.SetEnabled(false)
			}

		if not .Destroyed?() and linkField is false
			.list.SetListFocus()
		if query is ''
			.setAttachmentsManager(query, #())
		}
	setCustomization(name, query)
		{
		.SetCustomKey(.buildCustomKey(name, query))
		}
	buildCustomKey(name, query)
		{
		if .linkField is false
			return ListCustomize.BuildCustomKeyFromQueryTitle(query, .title)

		access_custom_key = .Send('GetAccessCustomKey')
		return String?(access_custom_key)
			? access_custom_key $ ' | ' $ name
			: false
		}
	SetCustomKey(customKey)
		{
		.customKey = customKey
		.customFields = ListCustomize.GetCustomizedFields(.customKey)
		if .customFields isnt false
			.list.SetCustomFields(.customFields)
		}

	ResetCustomKey(customKey)
		{
		.SetCustomKey(customKey)
		.setRecordCustomDeps()
		}

	GetCustomFields()
		{
		return .customFields
		}

	initValidationNotifyAndObserver(linkField, notifyLast)
		{
		if (linkField is false)
			.Window.AddValidationItem(this)
		if notifyLast is true
			.Send("AccessObserver", .AccessChanged, 10) /* = 10 so this is last notified*/
		else
			.Send("AccessObserver", .AccessChanged)
		}
	initObserverListRow()
		{
		// observer_ListRow is used as a straight function (not as a Browse method)
		// therefore "this" is the record that was modified
		// so ._browse is coming from the record
		.observer_ListRow = function (member)
			{ ._browse.Observer_ListRow2(this, member) }
		}
	initVisibleBrowseComponents(query, statusBar, columns)
		{
		if .Title is ""
			.Title = query.Tr("\r\n\t", " ") $ " - Browse"
		.statusBar = statusBar ? .Vert.Status : false

		if (.statusBar isnt false) .statusBar.SetValid(true)
		.list = .Member?("Vert") ? .Vert.List : .List
		.list.SetMultiSelect(true)
		return columns
		}
	createControls(statusBar, stretch, extra)
		{
		Assert(String?(.title) or (.title is false))
		Assert(Boolean?(statusBar))

		list = Object(stretch is true ? "ListStretch" : "List", defWidth: false,
			noShading: .noShading, noHeaderButtons: .noHeaderButtons,
			headerSelectPrompt: .headerSelectPrompt)
		if statusBar or .title isnt false
			{
			controls = Object("Vert")
			.addStatusBarAndTitle(controls, extra, list, statusBar)
			}
		else
			controls = list
		return controls
		}
	addStatusBarAndTitle(controls, extra, list, statusBar)
		{
		if .title isnt false
			controls.Add(Object('CenterTitle', .title))
		if extra isnt false
			{
			controls.Add(extra)
			controls.Add(#(Skip 5))
			}
		controls.Add(list)
		if .buttonBar isnt false
			controls.Add(.addon.ButtonBar())
		if statusBar isnt false
			controls.Add("Status")
		}
	SetMenu(menu = false)
		{
		if menu is false
			.addon.BuildMenu(.linkField, .columnsSaveName, .readOnly)
		else
			.addon.SetMenu(menu)
		}
	setAttachmentsManager(query, key)
		{
		.addon.SetAttachmentsManager(query, key)
		}
	Set(value)
		{
		orig_value = .value
		.value = value
		if (.linkField isnt false)
			{
			.headerData = .Send("GetData")
			.SetQuery(.query, .columns)
			// only clear sticky values if the value changed (could just be reload)
			if orig_value isnt .value
				.stickyFields.ClearStickyFieldValues()
			}
		// used in BrowseFlipForm, needs new records to be loaded in browse
		.Send('Browse_AfterSet')
		}
	Get()
		{
		return .value
		}
	GetBrowseData()
		{
		ob = Object()
		for (row in .list.Get())
			if (row.listrow_deleted isnt true)
				ob.Add(row)
		return ob
		}
	GetAllBrowseData()
		{
		return .list.Get()
		}
	SetQuery(query, columns = false, header_data = false, t = false,
		max_records = false, max_records_msg = '', exclude_custom_columns = false)
		{
		.list.FinishEdit()
		if (query is '')
			return
		Assert(String?(query), "query must be a string")

		if exclude_custom_columns is false
			columns = ListCustomize.AddCustomColumns(query, columns)

		// save column widths before resetting columns
		if .columns isnt false and .GetColumnsSaveName() isnt false
			UserColumns.Save(.GetColumnsSaveName(), .list, .columns)

		result = .buildListData(query, t, max_records, max_records_msg)
		data = result.data
		status_msg = result.status_msg

		.setListDataAndProperties(columns, data, header_data)
		.setBrowseStatusBar(status_msg)
		.setAttachmentsManager(query, .key)
		if .linkField isnt false
			.Send('RegisterLinkedBrowse', this, .Name)
		}
	setBrowseStatusBar(status_msg, valid = true)
		{
		if (.statusBar isnt false)
			{
			.statusBar.Set(status_msg)
			.statusBar.SetValid(valid)
			}
		}
	setListDataAndProperties(columns, data, header_data)
		{
		.SetColumns(columns)
		.list.Set(data)
		.SetHeaderData(header_data isnt false ? header_data: .headerData)
		.setRecordCustomDeps()
		.Send('DoWithoutDirty')
			{
			.Send("SetField", .dataMember, data)
			.Send("SetField", .allDataMember, .list.Get())
			}
		}
	setRecordCustomDeps()
		{
		if not CustomizeField.HasCustomFieldFormula?(.customKey)
			return
		data = .list.Get()
		for row in data
			CustomizeField.SetFormulas(.customKey, row, .protectField)
		}
	base_query: false
	SetBaseQuery(.base_query)
		{
		}
	MaxRecordsLoaded?: false
	buildListData(query, t, max_records, max_records_msg)
		{
		.MaxRecordsLoaded? = false
		.query = query
		if (.linkField isnt false)
			query = QueryHelper.AddWhere(query,
				" where " $ .linkField $ " = " $ Display(.value))
		data = Object()
		count = 0
		status_msg = ''
		.key = ShortestKey(.getBaseQuery(query)).Split(',')
		DoWithTran(t)
			{|tran|
			tran.Query(query)
				{|q|
				.query_columns = q.Columns()
				while false isnt x = q.Next()
					{
					x.Browse_OriginalRecord = x.Copy()
					x._browse = this
					x.Observer(.observer_ListRow)
					if .loadRecordNotification
						.Send('Browse_LoadRecord', x)
					data.Add(x)
					++count
					if max_records isnt false and count >= max_records
						{
						status_msg = max_records_msg is "" and
							max_records isnt false
							? 'Only the first ' $ max_records $
								' transactions were loaded'
							: max_records_msg
						.MaxRecordsLoaded? = true
						break
						}
					}
				}
			}
		return Object(:data, :status_msg)
		}
	getAvailableCols()
		{
		return ListCustomize.AddCustomColumns(.query, .getOriginalCols())
		}
	getOriginalCols()
		{
		return .columns is false and .query isnt ''
			? QueryHelper.AvailableColumns(.query)
			: .columns
		}
	GetQuery()
		{ return .query }
	GetList()
		{ return .list }
	SetLinkField(linkField)
		{
		Assert(String?(linkField) or (linkField is false))
		if (linkField is .linkField)
			return
		.linkField = linkField
		if (.query isnt false)
			.SetQuery(.query)
		}
	GetLinkField()
		{ return .linkField }
	SetHeaderData(headerData)
		{
		Assert(Record?(headerData) or (headerData is false))
		data = .list.Get()
		if ((.headerData = headerData) is false)
			return
		.headerData.Observer(.Observer_HeaderData)
		for (row in data)
			for (field in .headerFields)
				row.PreSet(field, headerData[field])
		}
	SetRecordHeaderData(record)
		{
		if (.headerData isnt false)
			{	// copy-in header data
			for (field in .headerFields)
				record.PreSet(field, .headerData[field])
			}
		}
	SetColumns(columns = false)
		{
		initialized? = .columns isnt false

		columns = columns is false ? QueryHelper.AvailableColumns(.query) : columns.Copy()
		.setting_columns(columns)

		if .GetColumnsSaveName() is false
			{
			.list.SetColumns(columns)
			return
			}
		UserColumns.Load(.columns, .GetColumnsSaveName(), .list, true,
			initialized?, load_visible?:)
		.appendMissingMandatoryCols(columns)
		}
	appendMissingMandatoryCols(columns)
		{
		visibleCols = .list.GetColumns()
		missingCols = columns.Difference(visibleCols).Filter(.MandatoryColumn?)
		if not missingCols.Empty?()
			{
			.list.AppendColumns(missingCols)
			.list.SetHeaderChanged(true)
			}
		}
	setting_columns(columns)
		{
		Assert(Object?(columns) or (columns is false))
		columns.Remove("listrow_deleted")
		columns.Add("listrow_deleted", at: 0)

		.columns = columns
		}
	GetColumns()
		{ return .columns }
	SetProtectField(protectField)
		{
		Assert(String?(protectField) or (protectField is false))
		.protectField = protectField
		}
	GetProtectField()
		{ return .protectField }
	SetValidationField(validationField)
		{
		Assert(String?(validationField) or (validationField is false))
		.validationField = validationField
		}
	GetValidationField()
		{ return .validationField }
	readOnly: false
	SetReadOnly(readOnly)
		{
		if readOnly isnt true and .alwaysReadOnly
			return
		Assert(Boolean?(readOnly))
		.list.SetReadOnly(readOnly)
		.readOnly = readOnly
		}
	setDataMember(dataMember)
		{
		.dataMember = dataMember
		.allDataMember = 'all_' $ .dataMember
		}
	GetDataMember()
		{
		return .dataMember
		}
	GetFields()
		{
		return .query_columns
		}
	// need GetField/SetField in browse because
	// ListEditWindow checks for them (DONT REMOVE)
	GetField(field)
		{
		.list.GetField(field)
		}
	SetField(field, value, idx = false, invalidate = false)
		{
		.list.SetField(field, value, idx)
		if invalidate
			{
			.Send("InvalidateFields", Object(.dataMember, .allDataMember))
			.Send("BrowseData_Changed", .GetBrowseData(), idx, field, false)
			}
		}
	SetMainRecordField(field, value)
		{
		.SetField(field, value)
		}
	GetTransQuery()
		{
		return .query
		}
	AddRecord(rec, idx = false, validateData = false, useDefaultsIfEmpty? = false)
		{
		rec.appAddedRecord? = true
		if (idx is false)
			idx = .list.GetNumRows()
		rec.Browse_ValidateData = validateData
		if (false isnt newrec = .list.CheckAndInsertRow(idx, rec, :useDefaultsIfEmpty?))
			{
			newrec.Browse_RecordDirty = true
			newrec.Browse_ValidateData = validateData
			.Refresh_BrowseData()
			}
		return newrec
		}
	ForceDeleteAll()
		{
		data = .GetAllBrowseData()
		size = data.Size()
		for (i = size - 1; i >= 0; --i)
			.ForceDeleteOne(data[i], i)
		}
	ForceDeleteOne(rec, row)
		{
		// skip records that are already marked as deleted
		if rec.listrow_deleted is true
			return
		.AllowNextDelete()
		.DeleteRows(row)
		}
	// DeleteRows is not handled by the Default method because it can potentially
	// have a large number of row arguments to delete. Doing @+1 args as the
	// Default method does causes Suneido to not optimize the arguments and
	// you can end up getting "value stack overflow" errors
	DeleteRows(@args)
		{
		.list.DeleteRows(@args)
		}
	FindLinkedRecord(keyValue, linkField, includeDeleted? = false)
		{
		browse_data = .GetAllBrowseData()
		if false is (row = browse_data.FindIf(
			{ it[linkField] is keyValue and
				(includeDeleted? or it.listrow_deleted isnt true) }))
			return false
		rec = browse_data[row]
		return Object(:rec, :row)
		}
	Save(tran = false)
		{
		try
			.list
		catch
			{
			SuneidoLog("ERROR: Browse.Save called but no .list", calls:)
			return true
			}

		.list.FinishEdit()
		.Send('Browse_BeforeValid')
		if not .Valid?(evalRule?:, selectFirstInvalid?:)
			return false

		.save_result = true
		KeyException.TryCatch(
			block:
				{
				DoWithTran(tran, .save, update:)
				}
			catch_block:
				{|e|
				KeyException(e)
				// if error occurred after list deletions in the save method,
				// nextSaveRec may not be in list, must check
				if (.nextSaveRec isnt false and
					false isnt (i = .list.Get().Find(.nextSaveRec)))
					.list.SetSelection(i)
				return false
				}
			)

		return .save_result
		}
	QueueDeleteAttachmentFile(newFile, oldFile, name, action)
		{
		sel = .list.GetSelection()
		rec = .list.GetRow(sel[0])
		.addon.QueueDeleteAttachmentFile(newFile, oldFile, rec, name, action)
		}
	RestoreAttachmentFiles()
		{
		.addon.CleanupAttachments(true)
		}
	deletes: #()
	save(tran)
		{
		// make sure deletes starts as a fresh object
		// (no remaining deletes from previous failed save)
		.deletes = Object()

		// keep track of the record we are trying to save, used to highlight the line in
		// the browse if save fales
		.nextSaveRec = false
		if false is .Send("Browse_BeforeSave", tran)
			return .saveResult(false)
		data = .list.Get()
		if (data.Size() is 0)
			return .saveResult(true)

		query = QueryHelper.StripSort(.query)
		deletes = Object()
		outputs = Object()
		oldrecs = Object()

		// do ALL query deletes first to prevent duplicate keys when a key is
		// deleted, then inserted again.
		if false is .handleDeletedItems(data, tran, query, deletes, oldrecs)
			return .saveResult(false)

		// process the rest of the data
		.handleUpdatedItems(data, outputs, deletes, oldrecs)

		// do outputs last so we don't get conflicts with updates
		.handleOutputItems(outputs, tran, query)

		.Send("Browse_AfterSave", tran)

		// list deletions and clearing record flags must be last
		if .linkField is false
			.after_save_processing(data, deletes)
		else
			.deletes = deletes // save deletes for save complete access notification
		return .saveResult(true)
		}
	saveResult(result)
		{
		.save_result = result
		return result
		}

	handleDeletedItems(data, tran, query, deletes, oldrecs)
		{
		valid = true
		for row in data.Members()
			{
			.nextSaveRec = record = data[row]

			if not .existingRecordChanged?(record)
				continue

			cur = tran.Query1(.keyQuery(.getBaseQuery(query),
				record.Browse_OriginalRecord))
			if false is result = .handleConflictItem(record, cur, row, valid)
				return false
			if Object?(result)
				{
				valid = result.valid
				continue
				}

			if record.listrow_deleted is true
				{
				if .deleteRecord(record, tran, deletes, row, cur) is false
					return false
				}
			else // record dirty
				{
				oldrecs[row] = cur
				}
			}
		if not valid
			.setBrowseStatusBar("Another user has modified records, " $
				"restore or discard changes", false)
		return valid
		}

	handleConflictItem(record, cur, row, valid)
		{
		if .addon.RecordConflict?(
			record, cur, .query_columns, .linkField, quiet?: valid is false)
			{
			if .linkField isnt false
				return false
			// if valid is true, then this is the first record conflict
			if valid
				{
				.list.SetSelection(row)
				valid = false
				}
			.list.AddHighlight(row, CLR.ErrorColor)
			return Object(:valid)
			}
		return true
		}

	handleUpdatedItems(data, outputs, deletes, oldrecs)
		{
		for (row in data.Members())
			{
			record = data[row]
			.nextSaveRec = record
			if (record.listrow_deleted is true)
				continue
			if (record.Member?('Browse_NewRecord'))
				{
				// only save new records that are dirty
				if (record.Browse_RecordDirty is true)
					{
					// have to queue up outputs and do last so that
					// they don't conflict with updates
					outputs.Add(record)
					}
				else
					deletes.Add(row)
				}
			else if oldrecs.Member?(row)
				{
				oldrecs[row].Update(record)
				}
			}
		}

	handleOutputItems(outputs, tran, query)
		{
		if not outputs.Empty?()
			{
			tran.Query(query)
				{|q|
				for (record in outputs)
					{
					.nextSaveRec = record
					q.Output(record)
					}
				}
			}
		}

	getBaseQuery(query)
		{
		return .base_query is false ? query : .base_query
		}

	existingRecordChanged?(record)
		{
		return ((record.Member?('Browse_RecordDirty') and
			not record.Member?('Browse_NewRecord'))	or
			record.listrow_deleted is true)
		}

	deleteRecord(record, tran, deletes, row, cur)
		{
		result = .list.Send('Browse_AllowDelete', record, tranFromSave: tran)
		if result is 0 or result is true
			{
			deletes.Add(row)
			.addon.DeleteRecordAttachments(record)
			cur.Delete()
			return true
			}
		if String?(result)
			AlertDelayed(result, 'Delete Record')
		return false
		}

	after_save_processing(data, deletes)
		{
		.list.SetForceDelete()
		.list.DeleteRows(@deletes)

		// clear records' flags and set original record
		for (record in data)
			if (record.Member?('Browse_NewRecord') or
				record.Member?('Browse_RecordDirty'))
				.clearRecordFlags(record)

		if .save_result is true
			.addon.CleanupAttachments()
		}
	clearRecordFlags(record)
		{
		record.Delete('Browse_NewRecord')
		record.Delete('Browse_RecordDirty')
		record.Delete('Browse_OriginalRecord')
		record.Delete('List_InvalidData')

		// have to keep the following lines separate due to the fact that
		// object.Delete returns false if the member is not present in the object
		original_record = record.Copy()
		original_record.Delete("_browse")
		record.Browse_OriginalRecord = original_record
		}
	access_after_save()
		{
		data = .list.Get()
		.after_save_processing(data, .deletes)
		.deletes = Object()
		}

	accessInvalid()
		{
		if .firstInvalidRow isnt false
			.list.SetSelection(.firstInvalidRow)
		}
	Dirty?(dirty/*unused*/ = "")
		{ return false }

	BrowseDataDirty?(data = false)
		{
		if data is false
			data = .GetAllBrowseData()
		for rec in data
			if rec.listrow_deleted is true or rec.Browse_RecordDirty is true
				return true
		return false
		}

	firstInvalidRow: false
	Valid?(evalRule? = false, selectFirstInvalid? = false)
		{
		.firstInvalidRow = false
		if not .BrowseDataDirty?()
			return true

		if false is .Send("Browse_ExtraValid")
			return false

		data = .GetAllBrowseData()
		for (row in data.Members())
			{
			rec = data[row]
			if .skipValid?(rec)
				continue
			if not .check_valid_field(rec, evalRule?)
				{
				.firstInvalidRow = row
				if selectFirstInvalid?
					.list.SetSelection(.firstInvalidRow)
				return false
				}
			}
		return true
		}

	skipValid?(rec)
		{
		if rec.listrow_deleted is true
				return true
		// bypass new non-dirty records
		if (rec.Member?('Browse_NewRecord') and
			rec.Browse_RecordDirty isnt true)
			return true
		// bypass old non-dirty records
		if (not rec.Member?('Browse_NewRecord') and
			not rec.Member?("Browse_RecordDirty"))
			return true
		return false
		}

	keyQuery(query, x)
		{
		keyRestrictions = Object()
		for fld in .key
			keyRestrictions.Add(fld $ ' = ' $ Display(x[fld]))
		return QueryHelper.StripSort(query) $
			Opt(" where ", keyRestrictions.Join(' and '))
		}
	Browse_InvalidateFields(@fields)
		{
		// first remove named members (if sent as a message, will have "source")
		fields = fields.Values(list:)
		for (rec in .list.Get())
			rec.Invalidate(@fields)
		}
	// interface (list)
	List_Selection(selection)
		{
		validSelect? = data = false
		if (selection isnt false)
			validSelect? = .check_valid_field(data = .list.GetRow(selection[0]))
		if validSelect? and data isnt false
			.check_browse_status_message(data)
		.Send("List_Selection" selection)
		}

	check_browse_status_message(data)
		{
		if data.Member?('browse_status_message') and data.browse_status_message isnt ''
			.statusBar.Set(data.browse_status_message)
		}

	List_CellEdit(col, row, data/*unused*/)
		{
		// store the previous row value in BrowseControl_PreviousData
		record = .list.GetRow(row)
		field = .list.GetCol(col)
		.FieldPreChange(record)
		if (false is .Send("Browse_CellEdit", col, row, record))
			return false
		// check the field against protection field
		if (FieldProtected?(field, record, .protectField))
			return false	// disallow change if protection on
		return true
		}
	FieldPreChange(record) // call this before modifying
		{
		record.Browse_RecordDirty = true
		}
	List_AfterEdit(col, row, data, valid?)
		{
		record = .list.GetRow(row)
		.Send('Browse_AfterEdit', :col, :row, :data, :valid?, :record)
		}
	List_InvalidDataChanged(rec)
		{
		rec.Browse_RecordDirty = true
		.Refresh_BrowseData() // will make Access record dirty if linked
		}
	List_AfterEditWindowCommit(col, row, data, valid?, valueChanged?)
		{
		field = .list.GetCol(col)
		record = .list.GetRow(row)
		if valueChanged?
			{
			if not valid?
				{
				.list.AddInvalidCell(col, row)
				.stickyFields.RemoveInvalidStickyValue(field, data)
				}
			else
				.list.RemoveInvalidCell(col, row)
			}
		for column in .list.GetColumns().Members()
			if ListCustomize.MandatoryAndEmpty?(
				record, .list.GetCol(column), .customFields, .protectField)
				.list.AddInvalidCell(column, row)

		.check_valid_field(record)
		.Send('List_AfterEditWindowCommit', col, row, data, valid?)
		}
	check_valid_field(record, evalRule? = false)
		{
		.setBrowseStatusBar('')

		if .readOnly or record.listrow_deleted is true or
			(record.Member?('Browse_NewRecord') and record.Browse_RecordDirty isnt true)
			return true

		// if any cells are invalid, do not bother with valid rule since the
		// record may not have valid data types (causing errors in valid rule)
		if .list.RowHasInvalidCell?(record)
			return false

		if (.validationField isnt false)
			{
			msg = .validationMsg(evalRule?, record)
			.setBrowseStatusBar(msg, msg is "")
			return msg is ""
			}
		return true
		}
	validationMsg(evalRule?, record)
		{
		try
			msg = evalRule? is true
				? record.Eval(Global('Rule_' $ .validationField))
				: record[.validationField]
		catch (err)
			{
			msg = "There was a problem validating the record"
			SuneidoLog("ERROR: (CAUGHT) " $ err, caughtMsg: 'user message: ' $ msg)
			}
		return msg
		}
	Refresh_BrowseData()
		{
		.Send("SetField", .dataMember, .GetBrowseData())
		.Send("InvalidateFields", Object(.dataMember))
		}
	List_Move(fromRow, toRow)
		{
		.Send('List_Move', fromRow, toRow)
		.Refresh_BrowseData()
		}
	List_CellValueChanged(col, row, data)
		{
		// mark record dirty
		record = .list.Get()[row]
		record.Browse_RecordDirty = true

		// send notifications of data changes to controller
		.Send("Browse_CellValueChanged", data, .list.GetCol(col), :record)
		browse_data = .GetBrowseData()
		.Send("SetField", .dataMember, browse_data)
		.Send("SetField", .allDataMember, .list.Get())
		.Send("InvalidateFields", Object(.dataMember, .allDataMember))
		if (browse_data.Member?(row))
			.Send("BrowseData_Changed", browse_data, row, .list.GetCol(col), false)
		}
	allowNextDelete: false
	AllowNextDelete()
		{
		.allowNextDelete = true
		}
	List_DeleteRecord(record)
		{
		if .allowDelete?(record)
			{
			.allowNextDelete = false
			row = .list.Get().Find(record)
			if false is .handleDeletedLine(record)
				return false
			.list.RepaintRow(row)
			.Send("Browse_DeleteRecord", record)
			data = .GetBrowseData()
			.Send("SetField", .dataMember, data)
			.Send("SetField", .allDataMember, .list.Get())
			.Send("InvalidateFields", Object(.dataMember, .allDataMember))
			.Send("BrowseData_Changed", data, row, false, .list.Get()[row])
			if false isnt newRecord? = record.Member?('Browse_NewRecord')
				.addon.DeleteNewRecordAttachments(record)
			return newRecord?
			}
		return false
		}

	allowDelete?(record)
		{
		return .allowNextDelete or ProtectRuleAllowsDelete?(record, .protectField,
			record.Member?('Browse_NewRecord'))
		}

	handleDeletedLine(record)
		{
		if (record.listrow_deleted isnt true)
			{
			if false is result = .Send("Browse_AllowDelete", record, tranFromSave: false)
				return false
			if String?(result)
				{
				.AlertInfo('Delete Record', result)
				return false
				}

			// listrow_deleted member is needed for the Open_Account_Balances
			// screen functionality
			record.listrow_deleted =  true
			// clear status bar error when error line is deleted
			if (.statusBar isnt false)
				{
				.statusBar.Set('')
				.statusBar.SetValid(true)
				}
			}
		else
			{
			if (false is .Send("Browse_AllowUnDelete", record))
				return false
			record.listrow_deleted =  false
			// need to validate the undeleted record but also mark dirty otherwise the
			// browse will not validate it when Valid? is called
			record.Browse_RecordDirty = true
			.check_valid_field(record)
			}
		return true
		}
	List_Deletions(deletions/*unused*/)
		{
		.Send("SetField", .allDataMember, .list.Get())
		.Send("InvalidateFields", Object(.allDataMember))
		}
	List_WantNewRow(prevRow/*unused*/, record, useDefaultsIfEmpty? = false)
		{
		if false is .Send("Browse_AllowNewRecord")
			return false

		record.Browse_NewRecord = true
		record._browse = this

		.SetRecordHeaderData(record)
		ListCustomize.HandleCustomizableFields(.customKey, record, .protectField,
			:useDefaultsIfEmpty?)
		if (.linkField isnt false)
			record[.linkField] = .value

		// make copy of record to evaluate the protect rule on because we don't
		// want this to trigger certain rules that depend on the record being in
		// list already
		if (.readOnly or
			(not .appAddedRecord?(record) and
				false is AllowInsertRecord?(record.Copy(), .protectField)))
			return false	// disallow adding a protected record

		record.Observer(.observer_ListRow)
		.stickyFields.SetRecordStickyFields(record)
		record.Browse_RecordDirty = false
		.Refresh_BrowseData()

		// send new record to controller to allow processing
		.Send("Browse_AddRecord", record)

		// return the new record to the list
		return record
		}
	appAddedRecord?(record)
		{
		return record.Member?('appAddedRecord?') and record.appAddedRecord? is true
		}
	List_AllowCellEdit(col, row/*unused*/)
		{
		return .list.GetCol(col) isnt "listrow_deleted"
		}
	List_NewRowAdded(row)
		{
		// mark fields that require entry as invalid
		rec = .list.GetRow(row)
		for col in .list.GetColumns().Members()
			{
			field = .list.GetCol(col)
			if ListCustomize.MandatoryAndEmpty?(
				rec, field, .customFields, .protectField) or
				false is .controlValidData?(rec, field)
				.list.AddInvalidCell(col, row)
			}
		// TODO: change Browse_NewRowAdded to pass record as well
		.Send("Browse_NewRowAdded", row)

		// mark record dirty if customize fields set default value on the row.
		// in case they have a default value on the foreign key field,
		// and it may fillin all required fields and record should save
		if rec.Member?('CustomizableSetDefaultValues')
			{
			rec.Browse_RecordDirty = true
			.Send("SetField", .allDataMember, .list.Get())
			.Send("InvalidateFields", Object(.allDataMember))
			.Refresh_BrowseData()
			}
		}

	getControlOption(field, option)
		{
		if field is 'listrow_deleted'
			return false
		return ListCustomize.GetControlOption(.customFields, field, option)
		}

	controlValidData?(record, field)
		{
		// don't check empty since mandatory checking is done separately
		if record.Browse_ValidateData isnt true or field is 'listrow_deleted'
			return true

		if ControlValidData?(record, field) is false
			{
			try
				{
				ctrl = Global(Datadict(field).Control[0] $ 'Control')
				if ctrl.Method?('GetUnvalidated')
					{
					.list.SetInvalidFieldData(record, field, record[field])
					record[field] = ''
					}
				}
			catch (err /*unused*/, "can't find")
				return false
			return false
			}
		return true
		}

	List_AllowZoom(row/*unused*/, col/*unused*/)
		{
		return .readOnly
		}
	List_DoubleClick(row, col)
		{
		if col is -1 // no columns set yet
			return 0 // allow edit
		field = .list.GetCol(col)
		if (field is "listrow_deleted")
			{
			.Send("Browse_DoubleClick_DeletedColumn")
			return false
			}
		.Send("List_DoubleClick", row, col)
		return 0
		}
	List_EditFieldReadonly(col, row)
		{
		name = .list.GetCol(col)
		record = .list.GetRow(row)
		if false isnt (result = FieldProtected?(name, record, .protectField))
			return result
		return .customFields isnt false and .customFields.Member?(name) and
			.customFields[name].GetDefault('readonly', false) is true
		}
	List_AllowHeaderResize(col)
		{
		return .list.GetCol(col) isnt "listrow_deleted"
		}
	List_AllowHeaderReorder(col)
		{
		return .list.GetCol(col) isnt "listrow_deleted"
		}
	List_HeaderContextMenu(x, y)
		{
		if .query is false or .query is ''
			return false
		.addon.HeaderMenu(this, x, y)
		}
	List_ContextMenu(x, y)
		{
		if .query is false or .query is ''
			return
		.pt = Object(:x, :y)
		.addon.ListMenu(this, x, y, .readOnly)
		}
	List_Tabover?(field)
		{
		// skip columns field with tabover option
		return .getControlOption(field, 'tabover')
		}
	List_NextCellNotAvailable?(row)
		{
		nums = .list.GetNumRows()
		if row > nums
			return false

		if row is nums
			row = row - 1
		return .list.GetRow(row).Browse_NewRecord is true
		}
	mandatoryColMinWidth: 50
	Header_TrackMinWidth(col)
		{
		if .MandatoryColumn?(.list.GetColumns()[col + 1])
			return .mandatoryColMinWidth
		return 0
		}
	On_New()
		{
		.On_Context_New()
		}
	On_Current_Restore()
		{
		.On_Context_Restore()
		}
	On_Current_DeleteUndelete()
		{
		.On_Context_DeleteUndelete()
		}
	On_Global_Save()
		{
		.On_Context_Save_All()
		}
	On_Global_Restore()
		{
		.On_Context_Restore_All()
		}
	On_Global_Summarize()
		{
		.On_Context_Summarize()
		}
	On_Global_Reporter()
		{
		.On_Context_Reporter()
		}
	On_Global_CrossTable()
		{
		.On_Context_CrossTable()
		}
	On_Global_Export()
		{
		.On_Context_Export()
		}
	On_Context_Print()
		{
		.addon.Print(.title, .query)
		}
	On_Context_Reset_Columns()
		{
		cols = .GetColumnsSaveName() isnt false ? .getAvailableCols() : .getOriginalCols()
		.setting_columns(cols)
		if .GetColumnsSaveName() isnt false
			{
			UserColumns.Reset(.list, .GetColumnsSaveName(), .columns, deletecol:,
				load_visible?:)
			.appendMissingMandatoryCols(cols)
			}
		else
			.list.SetColumns(.columns, reset:)
		}
	On_Context_Customize_Columns()
		{
		availableCols = .getAvailableCols()
		mandatoryFields = availableCols.Filter(.MandatoryColumn?)
		CustomizeColumnsDialog(.Window.Hwnd, .list, availableCols,
			.GetColumnsSaveName(), mandatoryFields, .headerSelectPrompt)
		}
	MandatoryColumn?(column)
		{
		return ListCustomize.MandatoryColumn?(column, .mandatoryFields, .customFields)
		}
	On_Context_Customize()
		{
		if false is result = ListCustomize.Customize(
			.query, .columns, .customKey, .linkField, this)
			return
		.handle_customize_result(result.dirty, result.custom_fields)
		}

	handle_customize_result(dirty, custom_fields)
		{
		Assert(.GetColumnsSaveName() isnt false)
		if Object?(dirty) and dirty.screen and
			custom_fields isnt false and custom_fields.NotEmpty?()
			{
			UserColumns.AddCustomFields(.GetColumnsSaveName(), .list, custom_fields,
				.getOriginalCols(), .getAvailableCols(), deletecol: )

			.setting_columns(.columns.Append(custom_fields))
			UserColumns.Load(.columns, .GetColumnsSaveName(), .list, deletecol:,
				load_visible?:)
			}
		if not Object?(dirty) or dirty.fields or dirty.screen
			.Send('BookRefresh')
		}
	On_Context_New()
		{
		.list.ContextNew()
		}
	On_Context_Edit_Field()
		{
		.list.ContextEdit(.pt)
		}
	On_Context_DeleteUndelete()
		{
		.list.DeleteSelection()
		}
	On_Context_Reason_Protected()
		{
		sel = .list.GetSelection()
		if sel is #()
			return
		row = sel[0]
		record = .list.GetRow(row)

		ListCustomize.ReasonProtected(record, .protectField, .Window.Hwnd)
		}
	On_Context_Save_All()
		{
		if (0 is .Send("Access_Save"))
			.Save()
		}
	On_Context_Restore()
		{
		.addon.Restore(.recQuery, .observer_ListRow)
		}
	recQuery(record)
		{
		return .keyQuery(.query, record.Browse_OriginalRecord)
		}
	On_Context_Restore_All()
		{
		.addon.RestoreAll(.query, .columns)
		}
	On_Context_Summarize()
		{
		.addon.Summarize()
		}
	On_Context_Reporter()
		{
		.addon.Reporter()
		}
	On_Context_CrossTable()
		{
		.addon.CrossTable(.query, .title)
		}
	On_Context_Export()
		{
		.addon.Export(.query)
		}
	On_Context_Go_To_QueryView()
		{
		.addon.Go_To_QueryView(.list, .query, .linkField, .value, .keyQuery)
		}
	On_Context_Cut_Record()
		{
		.addon.CutRecord()
		}
	On_Context_Paste_Record_Before()
		{
		.addon.PasteRecord()
		}
	On_Context_Paste_Record_After()
		{
		.addon.PasteRecord(before: false)
		}
	Observer_HeaderData(member)
		{
		if (.Empty?() or // after destroy
			not .headerFields.Has?(member))
			return
		repaint = false
		data = .list.Get()
		for (row in data)
			{
			if (row[member] isnt .headerData[member])
				{
				repaint = true
				row[member] = .headerData[member]
				}
			}
		if (repaint)
			.list.Repaint()
		}
	Browse_SetFieldValid(member, valid)
		{
		selection = .list.GetSelection()
		if selection.Empty?()
			return
		row = selection[0]
		col = .list.GetColumns().Find(member)
		if (col is false)
			return

		if valid is true
			.list.RemoveInvalidCell(col, row)
		else
			.list.AddInvalidCell(col, row)

		record = .list.Get()[row]
		.check_valid_field(record)
		}
	Observer_ListRow2(ob, member)
		// forces the list's current selection to repaint
		{
		if .skipObserver?(member)
			return

		.makeDirtyQueryColumnChanged(member, ob)
		.stickyFields.UpdateStickyField(ob, member, ob.Member?("Browse_NewRecord"))
		col = .list.GetColumns().Find(member)
		if (col is false)
			return

		data = .list.Get()
		selection = .list.GetSelection()
		if (selection isnt #() and data[selection[0]] is ob)
			row = selection[0]
		else if (false is row = data.Find(ob))
			return

		.refreshAndNotifyChange(data, row, member, col)
		}

	skipObserver?(member)
		{
		if .Destroyed?() or member is "listrow_deleted" or member is 'List_InvalidData'
			return true

		// due to sequence of events, browse could be readonly but have changes
		// happening to the record. To ignore changes when linked Access is not in edit
		// mode we need to ask the AccessControl if it's in edit mode or not, we
		// cannot rely on the readOnly property of the browse
		editMode? = .Send('EditMode?')
		if .linkField isnt false and editMode? is false
			return true

		return false
		}

	makeDirtyQueryColumnChanged(member, record)
		{
		if (.query_columns.Has?(member) and
			(not Record?(record.Browse_OriginalRecord) or
			record.Browse_OriginalRecord[member] isnt record[member]))
			{
			record.Browse_RecordDirty = true
			.Refresh_BrowseData() // will make Access record dirty if linked
			}
		}
	refreshAndNotifyChange(data, row, member, col, _committing = false)
		{
		record = data[row]
		if .list.HasInvalidCell?(record, member) and
			.list.InvalidCellValue(record, member) isnt record[member]
			.list.RemoveInvalidCell(col, row)

		// if field is modified by other field like rules, fill-in, assuming is valid
		if member isnt committing
			.list.SetInvalidFieldData(record, member, '')

		.list.RepaintRow(row)

		.Send("Browse_AfterField", member, record)

		// "fix" for header rules not always updated from browse
		.Send("InvalidateFields", Object(.dataMember, .allDataMember))
		.Send("BrowseData_Changed", .GetBrowseData(), row, member, false)

		// the above Sends can make other header fields change
		// but the observer will NOT be called recursively
		if (record.Browse_RecordDirty isnt true and
			.query_columns.Has?(member) and
			(not Record?(record.Browse_OriginalRecord) or
			record.Browse_OriginalRecord[member] isnt record[member]))
			{
			record.Browse_RecordDirty = true
			.Refresh_BrowseData() // will make Access record dirty if linked
			}
		}
	// interface (access)
	AccessChanged(@args)
		{
		switch (args[0])
			{
		case 'restore' :	.restore()
		case 'delete' :		return .AccessDelete()
		case 'save' :		return .Save(args[1])
		case 'after_save' : .access_after_save()
		case 'accessInvalid' : .accessInvalid()
		default:
			}
		return true
		}
	restore(t = false)
		{
		.RestoreAttachmentFiles()
		.SetQuery(.query, .columns, :t)
		}
	AccessDelete()
		// checks if OK to delete, empties the list and deletions list
		// Note: 	There should be a foreign key from lines to header
		//			with cascading deletes enabled.
		{
		for rec in .GetBrowseData()
			if false is ProtectRuleAllowsDelete?(rec, .protectField,
				rec.Member?('Browse_NewRecord'), true)
				return false
		if .checkConflicts()
			return false
		for record in .list.Get()
			{
			record.RemoveObserver(.observer_ListRow)
			.addon.DeleteRecordAttachments(record)
			}
		.addon.CleanupAttachments()
		.list.Clear()
		return true
		}
	checkConflicts()
		{
		data = .list.Get()
		for row in data.Members()
			{
			record = data[row]
			if record.Browse_NewRecord is true
				continue
			cur = Query1(.keyQuery(.getBaseQuery(.query), record.Browse_OriginalRecord))
			if .addon.RecordConflict?(record, cur, .query_columns, .linkField)
				return true
			}
		return false
		}
	Refresh(t = false) // Refresh method does not save!
		{
		.restore(t)
		}
	Default(@args)
		{
		method = args[0]
		if .list.Method?(method)
			return .list[method](@+1 args)

		if method.Prefix?('On_Context_') and args.Member?('item')
			{
			.Send('Browse_ContextMenuItemClicked', args.item)
			.Addons.Send(method)
			return 0
			}

		throw 'method not found in BrowseControl: ' $ method
		}
	GetColumnsSaveName()
		{
		if .columnsSaveName isnt ''
			return .columnsSaveName

		if .customKey isnt false
			return .customKey

		.addon.LogColumnsSaveWarning(.Title)
		return .Title
		}
	SetColumnSaveName(name)
		{
		.columnsSaveName = name
		}
	// interface (destroy)
	ConfirmDestroy()
		// returns true iff invalid data has not been entered into this
		{
		// this will be used in CloseWindowConfirmation
		// to log if the changes are discarded
		Suneido.AccessRecordDestroyed = .title is false ? .Title : .title

		// TODO: check for overwrite conflicts
		return .noSaveOnDestroy ? true : .Save()
		}
	Destroy()
		{
		Suneido.AccessRecordDestroyed = ''

		// if linked, the Access will trigger the save
		if (.linkField is false)
			.Window.RemoveValidationItem(this)
		if .GetColumnsSaveName() isnt false
			UserColumns.Save(.GetColumnsSaveName(), .list, .columns)
		.Send("RemoveAccessObserver", .AccessChanged)
		.Send("NoData")
		super.Destroy()
		}
	}