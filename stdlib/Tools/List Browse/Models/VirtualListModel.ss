// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// TODO: do not destroy cursor when seek or change direction
// REFACTOR: always initiate .ExpandModel to avoid false checking
class
	{
	ColModel:			false
	CheckBoxColModel:	false
	EditModel:			false
	ExpandModel:		false
	NextNum: 			false
	Created:			false

	Offset:			0
	VisibleRows: 	0

	baseQuery: 		false
	query:			false
	limit:			1000
	segment:		200

	New(query, startLast = false, columns = #(), columnsSaveName = false,
		headerSelectPrompt = false, where = false, hideCustomColumns? = false,
		mandatoryFields = #(), .checkBoxColumn = false, checkBoxAmountField = false,
		lockFields = #(),
		sortSaveName = false, .protectField = false, disableCheckSortLimit? = false,
		.extraSetupRecord = false, .observerList = false,
		validField = false, nextNum = false, excludeSelectFields = #(),
		.loadAll? = false, .extraFmts = false, customKey = false, .keyField = false,
		stickyFields = false, enableUserDefaultSelect = false, disableSelectFilter =false,
		stretchColumn = false, .enableMultiSelect = false, .saveQuery = false,
		hideColumnsNotSaved? = false, .linked? = false, useExpandModel? = false,
		option = false, defaultColumns = false, .asof = false, .useQuery = 'auto',
		warningField = false)
		{
		.Created = Timestamp()
		.baseQuery = .query = query

		sf = VirtualListColModel.BuildSelectFields(
			columns, excludeSelectFields, headerSelectPrompt)
		.sortModel = VirtualListSortModel(query, sf, sortSaveName, .loadAll?,
			disableCheckSortLimit?)
		.InitSelection()

		if where is false
			where = ''
		.query = .sortModel.BuildQuery(.query, where)

		if .checkBoxColumn isnt false
			{
			.CheckBoxColModel = VirtualListCheckBoxModel(.checkBoxColumn, keyField,
				checkBoxAmountField)
			mandatoryFields = mandatoryFields.Copy().Add(.checkBoxColumn)
			}

		// expecting all keys
		.EditModel = VirtualListEditModel(
			lockFields, protectField, validField, warningField, .queryKeys(.baseQuery))
		.stickyFields = StickyFields(stickyFields)

		.initStartLast = startLast
		.startLast = startLast
		.init()

		.ColModel = VirtualListColModel(columns, columnsSaveName,
			headerSelectPrompt, mandatoryFields, checkBoxColumn, excludeSelectFields,
			hideCustomColumns?, .query, .extraFmts, :customKey, :enableUserDefaultSelect,
			:disableSelectFilter, :stretchColumn, :hideColumnsNotSaved?, :option,
			:defaultColumns)

		if useExpandModel?
			.ExpandModel = VirtualListExpandModel()
		.NextNum = VirtualListNextNumModel(nextNum)

		if .loadAll?
			.loadAllIfSmall(false) // force to load
		else
			.setupCursors()

		.attachmentsManager = .EditModel.Editable?()
			? AttachmentsManager(.baseQuery, .keys(.baseQuery))
			: class { Default(@unused) { return false } }
		}

	UpdateStickyField(record, member)
		{
		.stickyFields.UpdateStickyField(record, member)
		}

	ClearStickyFieldValues()
		{
		.stickyFields.ClearStickyFieldValues()
		}

	curTop: false
	curBottom: false
	init()
		{
		.data = Object()
		.Offset = 0
		.RefreshColumns()
		.AllRead? = false
		}

	RefreshColumns()
		{
		.columns = QueryColumns(.query)
		}

	setupCursors()
		{
		// .curTop should be the first record's index
		.curTop = VirtualListModelCursor(.data, .query, .startLast,
			setupRecord: .setupRecord, asof: .asof, useQuery: .useQuery)
		// .curBottom should be the last record's index + 1
		.curBottom = VirtualListModelCursor(.data, .query, .startLast,
			setupRecord: .setupRecord, asof: .asof, useQuery: .useQuery)
		}

	setupRecord(rec)
		{
		if .CheckBoxColModel isnt false
			{
			if rec.Member?(.CheckBoxColModel.CheckBoxColumn)
				{
				if rec[.CheckBoxColModel.CheckBoxColumn]
					.CheckBoxColModel.SelectItem(rec)
				else
					.CheckBoxColModel.UnselectItem(rec)
				}
			else
				rec[.CheckBoxColModel.CheckBoxColumn] = .CheckBoxColModel.IsSelected(rec)
			}

		if .ColModel.HasCustomFormula?
			CustomizeField.SetFormulas(.ColModel.GetCustomKey(), rec, .protectField)

		.extraSetup(rec)

		if .EditModel isnt false and .EditModel.Editable?()
			{
			rec.vl_origin = rec.Copy()
			rec.vl_list = .observerList
			rec.Observer(VirtualListObserverOnChange)
			}
		}

	extraSetup(rec)
		{
		if .extraSetupRecord isnt false
			(.extraSetupRecord)(rec)
		}

	UsingCursor?()
		{
		.curTop.UsingCursor?()
		}

	GetBaseQuery()
		{
		return .baseQuery
		}

	GetQuery()
		{
		return .query
		}

	GetTableName()
		{
		return QueryGetTable(.GetBaseQuery(), orview:)
		}

	AllAvailableColumns(extra = "")
		{
		return QueryColumns(.sortModel.StripSort(.baseQuery) $ ' ' $ extra)
		}

	Columns()
		{
		return .columns
		}

	SetSort(displayCol, dataCol = false)
		{
		.sortModel.SetSort(displayCol, dataCol)
		.query = .sortModel.BuildQuery(.query)
		.Selection.ClearShiftStart()
		.resetQuery()
		}

	CheckSortable(field)
		{
		return .sortModel.CheckSortable(.baseQuery, field)
		}

	SetRecordsToTop(field, values)
		{
		Assert(.AllRead?)
		top = Object()
		.data.RemoveIf({
			if true is inlist? = values.Has?(it[field])
				top.Add(it)  // to keep the order
			inlist?
			})
		if top.Empty?()
			return false

		.data = top.Append(.data)
		return true
		}

	CheckSlowQuery(queryState, after)
		{
		whereSpecs = .ColModel.GetWhereSpecs(queryState.conditions, .AllAvailableColumns)
		query = .sortModel.BuildQuery(.baseQuery, whereSpecs.where, queryState.sortCol)
		if whereSpecs.where is '' and .sortModel.UsingDefaultSort?(query)
			return true
		allCols = .ColModel.GetAvailableColumns(query)
		return SlowQuery.Validate(query, allCols, after, queryState)
		}

	SetWhere(where)
		{
		if .initStartLast isnt .startLast // startLast could be changed on small table
			.startLast = .initStartLast

		if .sortModel.CheckAboveSortLimit?()
			.sortModel.SetOverSortLimit(.QueryAboveSortLimit?(.ColModel.GetSelectVals()))

		.query = .sortModel.BuildQuery(.baseQuery, where)
		.resetQuery()
		}

	RefreshData()
		{
		preOffset = .Offset
		preRows = .VisibleRows

		.destroyCursors()
		.init()
		.Selection = VirtualListGridSelection(this, .enableMultiSelect)
		if .loadAllIfSmall(preRows)
			{
			.UpdateOffset(preOffset, fromRefresh?:)
			return
			}
		.setupCursors()

		.Offset = 0
		.VisibleRows = 0

		.UpdateVisibleRows(preRows)
		if .Offset isnt preOffset
			.UpdateOffset(preOffset, fromRefresh?:)
		}

	GetInitStartLast()
		{
		return .initStartLast
		}

	GetStartLast()
		{
		return .startLast
		}

	SetStartLast(startLast)
		{
		if .startLast is startLast and not .curTop.Seeking
			return false

		.startLast = startLast
		.resetQuery()
		return true
		}

	resetQuery()
		{
		if .ExpandModel isnt false
			{
			// currently it does not retain the collapse states when resetting query
			recycled = .ExpandModel.RecycleExpands()
			.ExpandModel = VirtualListExpandModel()
			.ExpandModel.SetRecycledExpands(recycled)
			}

		.destroyCursors()
		.init()
		if .loadAllIfSmall(.VisibleRows)
			return
		.setupCursors()

		.Offset = 0
		pre_rows = .VisibleRows
		.VisibleRows = 0
		.UpdateVisibleRows(pre_rows)
		}

	chunkSize: 10
	ReadAllData(maxLoad = false)
		{
		count = 0
		do
			{
			lines = .startLast is false
				? .curBottom.ReadDown(.chunkSize)
				: .curTop.ReadUp(.chunkSize)

			count += lines * (.startLast is false ? 1 : -1)
			if maxLoad isnt false and count > maxLoad
				return false
			}
		while lines isnt 0
		.closeCursors()
		return true
		}

	UpdateOffset(offset, saveAndCollapse = false, _slowQueryLog = false,
		fromRefresh? = false)
		{
		if slowQueryLog is false
			_slowQueryLog = Object(logged: false, from: 'UpdateOffset')
		if .Begin?() and .End?()
			return 0

		if offset > 0 											// scrolls down
			{
			lines = offset + .endOffset() - .curBottom.Pos
			.curBottom.ReadDown(lines)
			// plus 1 for blank row at end
			offset = Max(0, Min(offset, .curBottom.Pos - .Offset - .VisibleRows + 1))
			}
		else													// scrolls up
			{
			lines = fromRefresh?
				? .curTop.Pos - offset
				: .curTop.Pos - offset - .Offset
			.curTop.ReadUp(lines)
			// .Offset should not exceed top cursor position
			offset = Max(offset, .curTop.Pos - .Offset)
			}
		.Offset = .Offset + offset
		if saveAndCollapse isnt false
			.recyling(offset, saveAndCollapse)
		.closeCursorsIfAllRead()
		return offset
		}

	recyling(offset, saveAndCollapse)
		{
		if .AllRead? or not .overLimit?()
			return

		if offset > 0
			{
			if .curTop.Pos + .segment < .Offset
				{
				if false is recToRelease = .collapseBeforeReleaseTop(saveAndCollapse)
					return
				.curTop.ReleaseDown(recToRelease)
				}
			}
		else
			{
			if .curBottom.Pos - .segment > .endOffset()
				{
				if false is recToRelease = .collapseBeforeReleaseBottom(saveAndCollapse)
					return
				.curBottom.ReleaseUp(recToRelease)
				}
			}
		}

	collapseBeforeReleaseTop(saveAndCollapse)
		{
		recToRelease = releaseCount = 0
		while releaseCount < .segment
			{
			if false is .curTop
				return false
			pos = .curTop.Pos + recToRelease
			recylingRec = rec = .data[pos]

			if rec.vl_expanded_rows isnt ''
				{
				if pos + rec.vl_expanded_rows >= .Offset
					break
				releaseCount += rec.vl_expanded_rows
				}
			if false is saveAndCollapse(recylingRec, pos - .Offset, model: this)
				return false
			releaseCount++
			recToRelease++
			}
		return recToRelease
		}

	collapseBeforeReleaseBottom(saveAndCollapse)
		{
		recToRelease = releaseCount = 0
		while releaseCount < .segment
			{
			if false is .curBottom
				return false
			pos = .curBottom.Pos - 1 - recToRelease
			recylingRec = rec = .data[pos]

			if rec.vl_expand? is true
				{
				if pos - rec.vl_rows <= .endOffset()
					{
					break
					}
				releaseCount += rec.vl_rows
				recylingRec = .data[pos - rec.vl_rows]
				Assert(recylingRec.vl_expanded_rows is: rec.vl_rows)
				}
			if false is saveAndCollapse(
				recylingRec, pos - rec.vl_rows - .Offset, model: this)
				return false

			releaseCount++
			recToRelease++
			}
		return recToRelease
		}

	overLimited: false
	overLimit?()
		{
		if ((.curBottom.Pos - .curTop.Pos).Abs() > .limit)
			{
			.overLimited = true
			return true
			}
		return false
		}

	firstTime: true
	UpdateVisibleRows(visibleRows, _slowQueryLog = false)
		{
		if slowQueryLog is false
			_slowQueryLog = Object(logged: false, from: 'UpdateVisibleRows')
		if not .updateVisibleRows?(visibleRows)
			return

		pre_rows = .VisibleRows
		ensureAtBottom = .keepAtBottom?()
		if not .startLast
			{
			lines = visibleRows + .Offset - .curBottom.Pos
			.curBottom.ReadDown(lines)
			.closeCursorsIfAllRead()
			}

		.VisibleRows = visibleRows
		if .startLast and pre_rows is 0 // for startLast first time
			{
			// plus 1 for blank row at end
			.UpdateOffset(-visibleRows + 1)
			}

		.fillEmpty(ensureAtBottom, visibleRows)
		}

	updateVisibleRows?(visibleRows)
		{
		if visibleRows <= 0 or visibleRows is .VisibleRows
			return false
		if .AllRead? is true
			return true
		if .firstTime
			{
			.firstTime = false
			if .loadAllIfSmall(visibleRows)
				{
				.VisibleRows = visibleRows
				return false
				}
			}
		return true
		}

	keepAtBottom?()
		{
		return (.VisibleRows isnt 0 and // not first time
			.initStartLast is true and // default to be end
			(.ExpandModel is false or
				.ExpandModel.GetExpanded().Size() is 0) and // for speed issue
			.End?())
		}

	fillEmpty(ensureAtBottom, visibleRows)
		{
		if .seeking? or .curBottom.Seeking
			{
			if 1 < emptyRows = .emptyRows()
				.UpdateOffset(-emptyRows)
			}
		else if ensureAtBottom
			.scrollToBottom(visibleRows)
		}

	scrollToBottom(visibleRows)
		{
		if not .End?()
			{
			if .AllRead? is true
				.UpdateOffset(.data.Size() - .Offset - visibleRows)
			else if .startLast
				.UpdateOffset(-(.Offset + visibleRows))
			}
		else if 1 < emptyRows = .emptyRows()
			.UpdateOffset(-emptyRows)
		}

	maxToLoadAll: 300
	loadAllLimit: 250
	loadAllIfSmall(visibleRows, _stopLoadAll = false)
		{
		if visibleRows is 0
			return false
		if stopLoadAll
			return false
		queryNoSort = QueryStripSort(.query)
		sort = QueryGetSort(.query)
		if .loadAll?
			{
			.loadAll(queryNoSort, sort, visibleRows)
			return true
			}

		if not .query.Has?('sort')
			return false

		blacklist = Suneido.GetInit('VirtualList_NotLoadAll', Object())
		if blacklist.GetDefault(queryNoSort, false)
			return false

		return VirtualListModelCursor.EstimatedSmall?(queryNoSort, .loadAllLimit)
			? .loadAll(queryNoSort, sort, visibleRows, .maxToLoadAll, :blacklist)
			: false
		}

	loadAll(queryNoSort, sort, visibleRows, maxLoad = false, blacklist = false)
		{
		prevStartLast = .startLast
		.startLast = false
		.destroyCursors()
		.curTop = VirtualListModelCursor(.data, queryNoSort, .startLast,
			setupRecord: .setupRecord, asof: .asof, useQuery: .useQuery)
		.curBottom = VirtualListModelCursor(.data, queryNoSort, .startLast,
			setupRecord: .setupRecord, asof: .asof, useQuery: .useQuery)
		if false is .ReadAllData(maxLoad)
			{
			blacklist[queryNoSort] = true
			.startLast = prevStartLast
			.destroyCursors()
			.init() // fall back
			.setupCursors()
			return false
			}

		VirtualListSortModel.SortInMemory(.data, sort)
		if prevStartLast isnt .startLast
			.UpdateOffset(.data.Size() - visibleRows +1)

		return true
		}
	SetFirstSelection()
		{
		if .Selection.NotEmpty?()
			return false
		if .data.Size() is 0
			return false
		focusedRow = .initStartLast
			? .startLast ? -1 : .data.Size() - 1
			: 0
		.Selection.SelectRows(false, false, focusedRow)
		return focusedRow
		}

	emptyRows()
		{
		return .Offset + .VisibleRows - .curBottom.Pos
		}

	AllRead?: false
	closeCursorsIfAllRead()
		{
		if .overLimited
			return
		if not .AllRead? and .curBottom.ReadDown(1) is 0 and .curTop.ReadUp(1) is 0
			.closeCursors()
		}

	closeCursors()
		{
		.AllRead? = true
		.curTop.Close()
		.curBottom.Close()
		}

	ValidateRow(row, returnBoundary? = false)
		{
		if .data.Size() is 0 or .curTop is false
			return false

		if row < .curTop.Pos
			return returnBoundary? ? .curTop.Pos : false

		if row >= .curBottom.Pos
			return returnBoundary? ? .curBottom.Pos - 1 : false

		return row
		}

	seeking?: false
	Seek(field, prefix)
		{
		.seeking? = true
		i = .seek(field, prefix)
		.seeking? = false
		return i
		}

	seek(field, prefix)
		{
		if .AllRead?
			{
			.UpdateOffset((i = .seekAllRead(field, prefix)) - .Offset)
			return i
			}

		.startLast = false
		.destroyCursors()
		.init()
		.setupCursors()

		.curTop.Seek(field, prefix, .data)
		.curBottom.Seek(field, prefix, .data)

		.curBottom.ReadDown( .VisibleRows )

		if .data.Size() < .VisibleRows
			.UpdateOffset(.data.Size() - .VisibleRows)
		return 0
		}

	// Based on BinarySearch with the following optimizations:
	// - Early returns for edge cases (removing the need to search entirely)
	// - Capable of searching objects with negative indexes
	seekAllRead(field, prefix)
		{
		indexes = .data.Members().Sort!()
		if indexes.Empty?()
			return 0

		// First value is greater than prefix, no further searching required
		if .data[index = indexes[start = 0]][field] >= prefix
			return index

		// Last value is less than prefix, no further searching required
		if .data[index = indexes[end = .data.Size() - 1]][field] <= prefix
			return index

		return .binarySeek(field, prefix, indexes, start, end)
		}

	binarySeek(field, prefix, indexes, start, end)
		{
		do
			{
			cur = ((end - start) / 2).RoundDown(0) + start
			if prefix is value = .data[index = indexes[cur]][field]
				return index // Exact match
			if value < prefix
				start = cur
			else
				end = cur
			} while end - start > 1 // Search cannot get any closer
		return indexes[end]
		}

	Begin?()
		{
		if .curTop.Pos < .Offset
			return false

		.curTop.ReadUp(1)
		return .curTop.Pos >= .Offset
		}

	End?()
		{
		if .curBottom.Pos > .endOffset()
			return false

		.curBottom.ReadDown(1)
		return .endOffset() >= .curBottom.Pos
		}

	endOffset()
		{
		return .Offset + .VisibleRows
		}

	GetPosition()
		{
		return .Begin?()
			? 'top'
			: .End?()
				? 'bottom'
				: 'middle'
		}

	GetRecordRowNum(record)
		{
		if false is index = .data.Find(record)
			{
			keys = .keys(.baseQuery)
			keyRec = false
			if keys.Size() is 1
				keyRec = .GetRecordByKeyPair(record[keys[0]], keys[0])
			keyRecPos = false
			if keyRec isnt false
				keyRecPos = .data.FindIf({
					Object?(it.vl_origin) and it.vl_origin[keys[0]] is record[keys[0]] })

			keyRecInSelection? = ''
			recInSelection? = ''
			if .Selection isnt false
				{
				keyRecInSelection? = .Selection.HasSelectedRow?(keyRec)
				recInSelection? = .Selection.HasSelectedRow?(record)
				}

			keyRecInExpand? = ''
			recInExpand? = ''
			if .ExpandModel isnt false
				{
				expanded = .ExpandModel.GetExpanded()
				keyRecInExpand? = expanded.Has?(keyRec)
				recInExpand? = expanded.Has?(record)
				}

			keyRecInChangedRec? = ''
			recInChangedRec? = ''
			if .EditModel isnt false
				{
				changedRecs = .EditModel.GetOutstandingChanges(all?:)
				keyRecInChangedRec? = changedRecs.Has?(keyRec)
				recInChangedRec? = changedRecs.Has?(record)
				}

			SuneidoLog('ERROR: (CAUGHT) VirtualList cannot find record', calls:,
				params: [:record, :keyRec, datasize: .data.Size(), :keyRecPos,
					:keyRecInSelection?, :recInSelection?,
					:keyRecInExpand?, :recInExpand?,
					:keyRecInChangedRec?, :recInChangedRec?],
				caughtMsg: 'additional logging to track down errors')
			SujsAdapter.CallOnRenderBackend(#DumpStatus, 'VirtualList cannot find record')
			index = 0
			}
		return index - .Offset
		}

	GetRecord(row)
		{
		return .data.GetDefault(row + .Offset, false)
		}

	GetRecordByKeyPair(val, field, str? = false)
		{
		return .data.FindOne({
			Object?(it.vl_origin) and
				(str? ? String(it.vl_origin[field]) : it.vl_origin[field]) is val
			})
		}

	ReplaceRecord(oldRec, newRec)
		{
		keys = .keys(.baseQuery)
		if false is rec = .data.FindOne({|x| keys.Every?({ x[it] is oldRec[it] }) })
			return false
		return .ReloadRecord(rec, force:, :newRec)
		}

	ReloadRecord(rec, force = false, newRec = false)
		{
		if .cannotReload(rec, force)
			return rec

		query = .GetKeyQuery(not force
			? rec.vl_origin
			: newRec isnt false
				? newRec
				: rec)
		if false is freshRec = Query1(query)
			return 'The current record has been deleted.'

		// need to ensure that fields not actually shown on the screen
		// (i.e. scrolled off the end of the page) have their associated rules kicked in
		// to ensure that dependencies get set correctly
		for fld in .ColModel.GetColumns()
			freshRec[fld]

		expandRows = rec.vl_expanded_rows
		freshRec.vl_expanded_rows = expandRows
		.setupRecord(freshRec)
		pos = .recPos(rec)
		Assert(pos isnt: false, msg: "Please consult with user to gather information")
		.data[pos] = freshRec
		if .ExpandModel isnt false
			.ExpandModel.SetExpandRecord(freshRec, rec)
		.Selection.ReloadRecord(rec, freshRec)
		if .CheckBoxColModel isnt false
			.CheckBoxColModel.ReloadRecord(rec, freshRec)
		return .data[pos]
		}

	cannotReload(rec, force)
		{
		return (not (.AutoSave? and rec.Member?('vl_origin')) and not force) or
			.EditModel.RecordLocked?(rec)
		}

	// vvvvvvvvvvvvvvvvvvvvvvvvv Temporary debugging for: 36029 vvvvvvvvvvvvvvvvvvvvvvvvv
	recPos(rec)
		{
		if false isnt pos = .data.Find(rec)
			return pos
		try
			{
			keyPairs = rec.Project(.keys(.baseQuery))
			if false is pos = .data.FindIf({ .matchesKeyPairs?(keyPairs, it) })
				.logDebugging(reason: 'Record no longer exists in .data', :keyPairs)
			else
				{
				diffs = Object()
				dataRec = .data[pos]
				rec.Members().MergeUnion(dataRec.Members()).Each()
					{
					dataVal = dataRec.GetDefault(it, '')
					recVal = rec.GetDefault(it, '')
					if dataVal isnt recVal
						diffs[it] = 'dataVal: ' $ Display(dataVal) $
							', recVal: ' $ Display(recVal)
					}
				.logDebugging(reason: 'Record out of sync', :diffs)
				}
			}
		catch (error)
			.logDebugging(reason: 'Error caught', :error)
		return false
		}

	matchesKeyPairs?(keyPairs, dataRec)
		{
		for m, v in keyPairs
			if dataRec[m] isnt v
				return false
		return true
		}

	logDebugging(@params)
		{
		SuneidoLog('WARNING: extra debugging, please include with suggestion/notes',
			:params)
		}
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	GetKeyQuery(rec, save? = false)
		{
		query = QueryHelper.StripSort(save? is true ? .GetSaveQuery() : .baseQuery)
		return query $ Opt(" where ",
			.keys(query).Map({ it $ ' is ' $ Display(rec[it]) }).Join(' and '))
		}

	GetSaveQuery()
		{
		return .saveQuery is false ? .baseQuery : .saveQuery
		}

	GetLoadedData()
		{
		return .data
		}

	GetLastVisibleRowIndex()
		{
		return Min(.Offset + .VisibleRows, .curBottom.Pos - 1)
		}

	InsertNewRecord(record = false, row_num = false, force = false)
		{
		newRec = .buildNewRec(record)
		if not force and false is AllowInsertRecord?(newRec.Copy(), .protectField)
			return false	// disallow adding a protected record
		if not .startLast
			{
			if row_num is false and .VisibleRows isnt 0 and
				.curBottom.Pos > .Offset + .VisibleRows
				row_num = .VisibleRows
			if row_num is false
				.data.Add(newRec, at: .curBottom.Pos)
			else
				{
				for (i = .curBottom.Pos; i > row_num + .Offset; i--)
					.data[i] = .data[i - 1]
				.data[row_num + .Offset] = newRec
				}
			.curBottom.Pos += 1
			}
		else
			{
			if row_num is false
				row_num = Min(-.Offset, .VisibleRows)
			for (i = .curTop.Pos; i <= row_num + .Offset - 1; i++)
				.data[i-1] = .data[i]
			.data[row_num + .Offset - 1] = newRec
			.curTop.Pos -= 1
			.Offset -= 1
			}
		.highlightInvalidFields(newRec)
		return newRec
		}

	buildNewRec(record = false)
		{
		newRec = record is false ? Record() : record
		for key in .keys(.baseQuery)
			newRec[key] // trigger rules on key fields
		.NextNum.ReserveNextNum(newRec, .baseQuery)
		ListCustomize.HandleCustomizableFields(
			.ColModel.GetCustomKey(), newRec, .protectField)
		.extraSetup(newRec)
		newRec.vl_list = .observerList
		newRec.Observer(VirtualListObserverOnChange)
		.stickyFields.SetRecordStickyFields(newRec)
		// want RecordChange to kick in with sticky fields
		if .stickyFields.GetStickyFields() isnt false
			.EditModel.ClearChanges(newRec)
		return newRec
		}

	RestoreNewRecord(rec)
		{
		newRec = .buildNewRec()
		.data.Replace(rec, newRec)
		if .ExpandModel isnt false and rec.vl_expanded_rows isnt ''
			{
			// Re-associate the expanded rows with the new record.
			// Otherwise, the expanded controls will be inoperable post restore
			newRec.vl_expanded_rows = rec.vl_expanded_rows
			.ExpandModel.SetExpandRecord(newRec, rec)
			}
		.highlightInvalidFields(newRec)
		return newRec
		}

	highlightInvalidFields(rec)
		{
		customFields = .ColModel.GetCustomFields()
		.ColModel.GetColumns().Each()
			{
			if ListCustomize.InvalidField?(rec, it, customFields, .EditModel.ProtectField)
				.EditModel.AddInvalidCol(rec, it)
			}
		}

	DeleteRecord(record)
		{
		if not .startLast
			{
			.removeRecord(record, .curTop.Pos, .curBottom.Pos)
			.curBottom.Pos -= 1
			}
		else
			{
			.removeRecord(record, .curBottom.Pos - 1, .curTop.Pos - 1)
			.curTop.Pos += 1
			.Offset += 1
			}
		.attachmentsManager.DeleteNewRecordFiles(record)
		.EditModel.ClearChanges(record)
		}

	removeRecord(record, startPos, endPos)
		{
		dir = .startLast ? -1 : 1
		found = false
		for i in startPos.Abs() .. endPos.Abs() - 1
			{
			pos = dir * i
			if .data[pos] is record
				found = true
			if found
				.data[pos] = .data[pos+dir]
			}
		.data.Delete(endPos - dir)
		}

	CheckRecord(rec, col = false, forceCheck = false)
		{
		if .CheckBoxColModel isnt false and rec isnt false
			.CheckRecordByKey(rec[.keyField], col, forceCheck)
		}

	CheckRecordByKey(key, col = false, forceCheck = false)
		{
		if .CheckBoxColModel is false
			return
		if col isnt .checkBoxColumn and forceCheck isnt true
			return
		if false is rec = .data.FindOne({ key is it[.keyField] })
			if false is rec = Query1(QueryStripSort(.baseQuery) $
				' where ' $ .keyField $ ' is ' $ Display(key))
				return
		.CheckBoxColModel.CheckRecord(rec)
		}

	CheckAll()
		{
		if .CheckBoxColModel is false
			return

		.CheckBoxColModel.SelectAll()
		for rec in .data
			rec[.checkBoxColumn] = true
		}

	UncheckAll()
		{
		if .CheckBoxColModel is false
			return

		.CheckBoxColModel.UnselectAll()
		for rec in .data
			rec[.checkBoxColumn] = false
		}

	GetCheckedRecords()
		{
		return .CheckBoxColModel.GetSelectedInfo()
		}

	destroyCursors()
		{
		if .curTop isnt false
			{
			.curTop.Close()
			.curTop = false
			}

		if .curBottom isnt false
			{
			.curBottom.Close()
			.curBottom = false
			}
		}

	keys(query)
		{
		keys = .EditModel.Editable?()
			? .EditModel.LockKeyField
			: ShortestKey(.queryKeys(query))
		return keys.Split(',')
		}

	queryKeys(query)
		{
		// strip the where to keep the keys more consistent
		return QueryKeys(query)
		}

	Limit()
		{
		// can only assume (.limit - .segment) is in memory
		return .limit - .segment
		}

	GoToQueryView()
		{
		GotoQueryView(.query)
		}

	SetRecordExpanded(row_num, rows)
		{
		rec = .GetRecord(row_num)
		if rec.vl_expanded_rows isnt ''
			return
		rec.vl_expanded_rows = rows

		if not .startLast
			{
			for(i = .curBottom.Pos-1; i >= row_num + .Offset; --i)
				.data[i+rows] = .data[i]
			for(i = 1; i <= rows; ++i)
				.data[row_num + .Offset + i] = .expandRec(rows, i -1)
			.curBottom.Pos += rows
			}
		else
			{
			for(i = .curTop.Pos; i <= row_num + .Offset; ++i)
				.data[i-rows] = .data[i]
			for i in ..rows
				.data[row_num + .Offset - i] = .expandRec(rows, i)
			.curTop.Pos -= rows
			.Offset -= rows
			}
		}

	expandRec(rows, index)
		{
		return [vl_expand?: true, vl_rows: rows, vl_expand_index: index]
		}

	SetRecordCollapsed(row_num, keepPosition? = false)
		{
		rec = .GetRecord(row_num)
		if '' is rows = rec.vl_expanded_rows
			return

		rec.Delete('vl_expanded_rows')
		if not .startLast
			{
			i = row_num + .Offset + 1
			for(; i < .curBottom.Pos-rows; ++i)
				.data[i] = .data[i+rows]
			for (j = .curBottom.Pos -1 ; j >= i; --j)
				.data.Delete(j)
			.curBottom.Pos -= rows
			}
		else
			{
			for(i = row_num + .Offset; i >= .curTop.Pos; --i)
				.data[i+rows] = .data[i]
			for i in ..rows
				.data.Delete(.curTop.Pos + i)
			.curTop.Pos += rows
			}
		.updateOffsetForCollapse(row_num, rows, :keepPosition?)
		}

	updateOffsetForCollapse(row_num, rows, keepPosition?)
		{
		if keepPosition? is true // so grid.updateExpand() does not move the record index
			{
			.Offset += .startLast ? rows : 0
			return
			}
		if not .startLast and row_num < 0
			.Offset = Max(.Offset - Min(rows, -row_num), 0)
		if .startLast and row_num + rows >= 0
			.Offset = Min(.Offset + Min(rows, row_num + rows), -1)
		}

	AutoSave?: false
	LockRecord(rec)
		{
		if not .AutoSave?
			return true
		if true isnt result = .EditModel.LockRecord(rec)
			return result
		if .ExpandModel isnt false
			.ExpandModel.SetExpandReadOnly(rec, readonly: false)
		.highlightInvalidFields(rec)
		return result
		}

	UnlockRecord(rec)
		{
		if not .AutoSave?
			return
		.handleNextNumTimer()
		.EditModel.UnlockRecord(rec)
		if .ExpandModel isnt false
			.ExpandModel.SetExpandReadOnly(rec, readonly:)
		}

	handleNextNumTimer()
		{
		recs = .EditModel.GetOutstandingChanges()
		if not recs.Empty?()
			.NextNum.CheckAndClearNewNextNums(recs[0])
		}

	OwnLock?(rec)
		{
		if not .AutoSave?
			return true
		return .EditModel.RecordLocked?(rec)
		}

	// TMP logging for 27831
	LogInvalidFocus(row)
		{
		if .invalidFocus?(row)
			{
			SuneidoLog('INFO: found invalid focused row number', calls:,
				params: [:row, size: .data.Size()])
			return true
			}
		return false
		}

	invalidFocus?(row)
		{
		return row isnt false and row < 0 and .AllRead? and
			not .GetStartLast() and .curTop isnt false and not .curTop.Seeking
		}

	GetPrimarySort()
		{
		return .sortModel.GetPrimarySort()
		}

	SetDefaultSort()
		{
		.sortModel.SetDefaultSort()
		}

	ResetSort()
		{
		.sortModel.ResetSort()
		.query = .sortModel.BuildQuery(.query)
		.resetQuery()
		}

	SaveSort?()
		{
		return .sortModel.SaveSort?()
		}

	QueueDeleteAttachmentFile(newFile, oldFile, rec, fieldName, action)
		{
		return .attachmentsManager.QueueDeleteFile(
			newFile, oldFile, rec, fieldName, action)
		}

	DeleteRecordAttachments(rec)
		{
		.attachmentsManager.QueueDeleteRecordFiles(rec)
		}

	CleanupAttachments(restore? = false)
		{
		return .attachmentsManager.ProcessQueue(restore?)
		}

	Selection: false
	InitSelection()
		{
		if .Selection is false or .linked?
			return .Selection = VirtualListGridSelection(this, .enableMultiSelect)

		return .Selection
		}

	CheckAboveSortLimit?()
		{
		return .sortModel.CheckAboveSortLimit?()
		}
	SetOverSortLimit(overSortLimit?)
		{
		.sortModel.SetOverSortLimit(overSortLimit?)
		}
	SortLimitChecked?()
		{
		return .sortModel.SortLimitChecked?()
		}

	sortLimit: 10
	QueryAboveSortLimit?(conditions)
		{
		if not .CheckAboveSortLimit?()
			return false
		sf = .ColModel.GetSelectFields()
		indexes = .SelectableIndexes()
		indexedConditions = conditions.Filter(
			{ indexes.Has?(it.condition_field) or
				indexes.Has?(sf.GetJoinNumField(it.condition_field)) })

		whereSpecs = .ColModel.GetWhereSpecs(indexedConditions, .AllAvailableColumns)
		return .sortModel.QueryAboveSortLimit?(.baseQuery, whereSpecs.where)
		}

	SelectableIndexes()
		{
		query = QueryStripSort(.baseQuery)
		allCols = .ColModel.GetAvailableColumns(query)
		return SlowQuery.SelectableIndexes(query, allCols)
		}

	OverSortLimit?()
		{
		return .sortModel.OverSortLimit?()
		}

	Destroy()
		{
		.handleNextNumTimer()
		.EditModel.Destroy()
		.sortModel.Destroy()
		.ColModel.Destroy()
		.destroyCursors()
		}
	}
