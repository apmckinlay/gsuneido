// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: 'Repeat'
	RowName: 'Row'
	New(controls, .no_minus = false, .minusAtFront = false, .plusAtBottom = false,
		.project = false, .maxRecords = false, .noPlus = false, .headerFields = #(),
		.disableFieldProtectRules = false, .saveKey = false, .protectField = false)
		{
		super(.layout(controls))
		.vert = .FindControl('rows')
		for rc in .vert.GetChildren()
			.prepareRowRecordControl(rc)
		if .project isnt false
			.project()
		.Send(#Data)
		}
	emptyRecord: ()
	layout(controls)
		{
		ob = ['Horz' controls #(Skip 4)]
		plusButton = #(RepeatButton Plus 'plus.emf', tip: 'Add another row')
		if .plusAtBottom is false and .noPlus is false
			ob.Add(plusButton)
		if .no_minus is false
			ob.Add(#(RepeatButton Minus 'minus.emf', tip: 'Delete this row'),
				at: .minusAtFront isnt false ? 1 : ob.Size())
		.row = ['Record', ob, .disableFieldProtectRules]
		contentVert = ['Vert', ['Vert', .row, name: 'rows' overlap:], name: 'content']
		if .plusAtBottom isnt false and .noPlus is false
			contentVert.Add(plusButton)
		return ['WndPane', contentVert, windowClass: "SuBtnfaceArrow", name: 'pane']
		}
	project()
		{
		emptyRecord = []
		.project.Each({ emptyRecord[it] }) // reference fields to kick in rule
		.emptyRecord = emptyRecord.Project(.project).Remove("")
		}
	On_Plus(source)
		{
		if .AboveMax?()
			return false

		newRec = []
		.addHeaderFieldsToRecord(newRec)

		if .plusAtBottom
			return .onPlusAtBottom(newRec)
		.handleFocus()
		// have to add new rows at the end to keep tab order correct
		// move data to get the effect of inserting new rows in the middle
		.insertRow(.Tally())
		rows = .rows()
		j = .FindRow(source) + 1
		for (i = .Tally() - 1; i > j; --i)
			rows[i - 1].MoveStateTo(rows[i])
		.BeforeRowSet(rows[j], newRec)
		rows[j].Set(newRec)
		rows[j].SetFocus()
		.Send('Repeat_RowsChanged', newRow: j)
		}

	// tests if adding one record will put us over max, NOT if already over max
	AboveMax?(quiet = false)
		{
		if .maxRecords isnt false and .Tally() >= .maxRecords
			{
			if quiet is false
				.AlertWarn('Add New ' $ .RowName,
					'The maximum number of ' $ .RowName $ ' is ' $ .maxRecords)
			return true
			}
		return false
		}

	handleFocus()
		{
		// can't use SetFocus(NULL) - see suggn 18482
		.rows()[0].SetFocus()
		}

	insertRow(index)
		{
		.vert.Insert(index, .row)
		.prepareRowRecordControl(.rows()[index])
		}
	prepareRowRecordControl(row)
		{
		row.AddObserver(.RepeatRecord_Changed)
		row.SetProtectField(.protectField)
		.Send('Repeat_InsertRow', row)
		}
	removeRow(index)
		{
		.rows()[index].RemoveObserver(.RepeatRecord_Changed)
		.vert.Remove(index)
		}
	UpdateLayout(layout, .project)
		{
		while .Tally() > 0
			.removeRow(0)

		.layout(layout)
		if .project isnt false
			.project()
		}
	onPlusAtBottom(newRec) // when each row does not have same control
		{
		.handleFocus()
		length = .Tally()
		.insertRow(length)
		rows = .rows()
		.BeforeRowSet(rows[length], newRec)
		rows[length].Set(newRec)
		rows[length].SetFocus()
		.Send('Repeat_RowsChanged', newRow: length)
		}

	On_Minus(source)
		{
		.Send('Repeat_BeforeMinus', origSource: source)
		.handleFocus()
		if .vert.Tally() is 1
			.On_Plus(source)
		.removeRow(.FindRow(source))
		.Send(#NewValue, .Get())
		.dirty? = true
		.Send('Repeat_RowsChanged')
		}

	FindRowByRecordControl(rc)
		{
		return .rows().FindIf({ Same?(it, rc) })
		}

	FindRow(source)
		{
		return .FindRowByRecordControl(source.Parent.Parent)
		}

	Tally()
		{
		return .vert.Tally()
		}

	Get()
		{
		value = .rows().Map(#Get)
		if .project is false
			value.Map!({ it.Copy().Remove("") })
		else
			{
			value.Map!({ it.Project(.project).Remove("") })
			value.Remove(.emptyRecord)
			}
		value.Remove(#())
		if value.Empty?()
			return ''

		if .saveKey isnt false
			.assignUniqueIds(value)

		return value
		}

	assignUniqueIds(value)
		{
		for rec in value
			if rec.key is ''
				rec.key = Display(Timestamp()).Tr('#.')
		}

	Set(data)
		{
		.addHeaderObserver()
		if data.Size() is 0
			data = [[]]
		while .Tally() > data.Size()
			.removeRow(0)
		while .Tally() < data.Size()
			{
			i = .Tally()
			.insertRow(i)
			.rows()[i].SetReadOnly(.readonly)
			}
		rows = .rows()
		for (i = 0; i < data.Size(); ++i)
			{
			rowData = data[i].Copy()
			.BeforeRowSet(rows[i], rowData)
			rows[i].Set(rowData)
			}
		.dirty? = false
		}

	hdr_data: false
	addHeaderObserver()
		{
		if .headerFields.Empty?() or not Record?(.hdr_data = .getHeaderRecord())
			return
		.hdr_data.Observer(.Observer_HeaderData)
		}

	getHeaderRecord()
		{
		if Record?(hdr_data = .Send("GetData"))
			return hdr_data
		if 0 isnt ctrl = .Send('GetRecordControl')
			return ctrl.Get()
		return false
		}

	Observer_HeaderData(member)
		{
		if .headerFields.Has?(member)
			for rc in .rows()
				{
				rowData = rc.Get()
				rowData[member] = .hdr_data[member]
				}
		}

	BeforeRowSet(rowCtrl /*unused*/, rowData)
		{
		.addHeaderFieldsToRecord(rowData)
		}

	addHeaderFieldsToRecord(rec)
		{
		if .headerFields.Empty?() or not Record?(.hdr_data)
			return

		for field in .headerFields
			rec[field] = .hdr_data[field]
		}

	dirty?: false
	Dirty?(state = '')
		{
		if state isnt ''
			{
			for c in .rows()
				c.Dirty?(state)
			return state
			}
		else
			return .dirty? or .rows().Any?(#Dirty?)
		}
	Valid?(checkAll = false)
		{
		if not checkAll
			return .rows().Every?({|c| c.Valid(forceCheck:) is true })

		allValid? = true
		.GetRows().Each()
			{|c|
			if c.Valid(forceCheck:) isnt true
				allValid? = false
			}
		return allValid?
		}
	readonly: false
	SetReadOnly(readonly)
		{
		super.SetReadOnly(readonly)
		.readonly = readonly
		}

	GetReadOnly()
		{
		return .readonly
		}

	RepeatRecord_Changed(member)
		{
		if member.Suffix?('_protect')
			return

		if .readonly is true
			return

		.Send(#NewValue, .Get())
		.dirty? = true
		}

	rows()
		{
		return .vert.GetChildren()
		}

	GetRows()
		{
		return .rows()
		}

	AppendRow(rowRec) // Used by AttachmentRepeatControl
		{
		row = .Tally()
		.insertRow(row)
		rows = .rows()
		rows[row].Set(rowRec)
		}

	Destroy()
		{
		.Send(#NoData)
		super.Destroy()
		}
	}
