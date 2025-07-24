// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// forces the list's current selection to repaint ???
	CallClass(list, record, member)
		{
		cl = new this(list)
		cl.Observe(record, member)
		}

	New(.list)
		{
		.model = .list.GetModel()
		}

	Observe(record, member)
		{
		if .list.Empty?() or member is "invalidCols" or member.Prefix?("vl_") or
			member is 'List_InvalidData'
			return

		.makeDirtyQueryColumnChanged(member, record)
		if .list.GetColumns().Has?(member)
			.list.RepaintSelectedRows()
		}

	makeDirtyQueryColumnChanged(member, record, _committing = false)
		{
		if not .model.Columns().Has?(member)
			return

		if Record?(record.vl_origin) and record.vl_origin[member] is record[member]
			.model.EditModel.ClearMemberChange(record, member)
		else
			.model.EditModel.AddChanges(record, member, record[member])
		mandatoryAndEmpty = ListCustomize.MandatoryAndEmpty?(
			record,
			member,
			.model.ColModel.GetCustomFields(),
			.model.EditModel.ProtectField)
		valid = not mandatoryAndEmpty and
			ControlValidData?(record, member)
		if valid
			{
			.model.EditModel.RemoveInvalidCol(record, member)
			.model.UpdateStickyField(record, member)
			}
		else
			.model.EditModel.AddInvalidCol(record, member)

		if member isnt committing
			ListControl.SetInvalidFieldData(record, member, '')

		.list.RefreshValid(record)
		.list.Send("VirtualList_RecordChange", member, record)

		.model.ColModel.Plugins_Execute(data: record, :member,
			hwnd: .list.Parent.Window.Hwnd,	query: .model.GetKeyQuery(record),
			pluginType: 'Observers')
		}
	}
