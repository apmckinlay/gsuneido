// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	fields: ()
	data_changed: false
	New(.query, .keyfields, .headerFields = #(), headerData = false)
		{
		.data = Object()
		if headerData isnt false
			{
			.get_data_from_query()
			.SetHeaderData(headerData)
			}
		.deleted = Object()
		}
	get_data_from_query()
		// post: the data member is filled with the records from the
		// query.
		{
		if .query is ""
			return
		Transaction(read:)
			{ |t|
			q = t.Query(.query)
			while (false isnt x = q.Next())
				{
				x.PreSet("Explorer_PreviousData", x.Copy())
				.data.Add(x)
				}
			.fields = q.Columns()
			}
		}
	NewRecord(add = true)
		{
		r = Record()
		.setHeaderFields(r, .headerData)
		if add
			.data.Add(r)
		return r
		}
	SetHeaderData(data)
		{
		.headerData = data
		.headerData.Observer(.Observer_HeaderData)
		for r in .data
			.setHeaderFields(r, .headerData)
		}
	SetRecordHeaderData(r)
		{
		if .headerData isnt false
			.setHeaderFields(r, .headerData)
		}
	setHeaderFields(r, headerData)
		{
		for field in .headerFields
			r.PreSet(field, headerData[field])
		}
	Observer_HeaderData(member)
		{
		if not .headerFields.Has?(member)
			return
		for r in .data
			r.PreSet(member, .headerData[member])
		}
	GetData()
		{
		return .data
		}
	GetHeaderFields()
		{
		return .headerFields
		}
	GetFields()
		{
		return .fields
		}
	Output(rec)
		{
		rec.Explorer_Output_NewRecord = true
		.data_changed = true
		return true
		}
	Update(item)
		{
		i = .index_from_record(item)
		item.Explorer_UpdateRecord = true
		.data[i] = item
		.data_changed = true
		return true
		}
	DeleteItem(item)
		{
		if item.Member?('Explorer_Output_NewRecord')
			return true
		i = .index_from_item(item)
		if i is false
			throw "can't find item: " $ Display(item)
		.data[i].Explorer_RecordDeleted? = true
		.deleted.Add(.data[i])
		.data_changed = true
		return true
		}
	GetDeleted()
		{
		return .deleted
		}
	SetDeleted(deleted)
		{
		.deleted = deleted
		}
	Clear() // i.e. header delete
		{
		.data = Object()
		.deleted = Object()
		.data_changed = false
		}
	Restore()
		// post: model's data is set back to the original records from the query
		{
		.Clear()
		.get_data_from_query()
		}
	Save(tran = false)
		{
		.save_result = true
		KeyExceptionTransaction(tran, block: .SaveChanges)
		return .save_result
		}
	SaveChanges(t)
		{
		// Note: order of the save is import to avoid conflicts, clear flags is done
		// last so flags remain if there is a problem during the save
		if not .needToSave?()
			return .saveResult(true)

		.base_query = QueryStripSort(.query)
		.handleDeletes(t)
		outputs = Object()
		for x in .data
			{
			if x.Member?('Explorer_RecordDeleted?')
				continue
			if x.Member?('Explorer_Output_NewRecord')
				outputs.Add(x)
			else if x.Member?('Explorer_UpdateRecord') and
				false is .updateRecord(x, t)
				return .saveResult(false)
			}
		.handleOutputs(t, outputs)
		.clearAllRecordFlags()
		return .saveResult(true)
		}

	saveResult(result)
		{
		.save_result = result
		return result
		}

	needToSave?()
		{
		return ((.data.Size() > 0 or .deleted.Size() > 0) and .data_changed)
		}

	handleDeletes(t)
		{
		for x in .deleted
			t.QueryDo("delete " $ .Make_where(x))
		.deleted = Object()
		}

	updateRecord(record, t)
		{
		if false is updateRec = t.Query1(.Make_where(record))
			{
			AlertDelayed("List can't find record to update",
				title: 'Error', flags: MB.ICONERROR)
			return true
			}
		if RecordConflict?(record.Explorer_PreviousData, updateRec, .fields)
			return false
		updateRec.Update(record)
		return true
		}

	handleOutputs(t, outputs)
		{
		t.Query(.query)
			{ |q|
			for rec in outputs
				q.Output(rec)
			}
		}

	clearAllRecordFlags()
		{
		for record in .data
			if record.Member?('Explorer_Output_NewRecord') or
				record.Member?('Explorer_UpdateRecord')
				.clearRecordFlags(record)
		}

	clearRecordFlags(record)
		{
		record.Delete('Explorer_Output_NewRecord')
		record.Delete('Explorer_UpdateRecord')
		record.Delete('Explorer_PreviousData')
		record.PreSet("Explorer_PreviousData", record.Copy())
		}
	Make_where(x)
		{
		if x.Member?('Explorer_PreviousData')
			x = x.Explorer_PreviousData

		conditions = Object()
		for field in .keyfields
			conditions.Add(field $ " is " $ Display(x[field]))

		return .base_query $ " where " $ conditions.Join(" and ")
		}
	GetKey()
		{
		return .keyfields
		}
	GetQuery()
		{
		return .query
		}
	GetBaseQuery()
		{
		return .base_query
		}
	SetBaseQuery(base_query)
		{
		.base_query = base_query
		}
	Get(item)
		{
		return item
		}
	index_from_record(record)
		{
		if record.Member?('Explorer_PreviousData')
			record = record.Explorer_PreviousData
		for i in .data.Members()
			{
			r = .data[i]
			if r.Member?('Explorer_PreviousData')
				r = r.Explorer_PreviousData
			rec? = true
			for field in .keyfields
				if r[field] isnt record[field]
					rec? = false
			if rec?
				return i
			}
		return false
		}
	index_from_item(item)
		{
		for i in .data.Members()
			{
			r = .data[i]
			rec? = true
			for field in .keyfields
				if r[field] isnt DatadictEncode(field, item[field])
					rec? = false
			if rec?
				return i
			}
		SuneidoLog("ERROR: can't find item extra info",
			params: Object(data: .data, keyfields: .keyfields))
		return false
		}
	GetEntries()
		{
		return .data
		}
	GetDataChanged()
		{
		return .data_changed
		}
	ChangeQuery(query)
		{
		.Clear()
		.query = query
		.get_data_from_query()
		}
	}
