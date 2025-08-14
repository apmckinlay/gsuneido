// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Start/End a trial"
	CallClass()
		{
		return OkCancel(this, .Title, 0)
		}

	list: false
	New()
		{
		.model = new SvcModel()
		.model.SetSettings(SvcSettings.Get())
		.list = .FindControl('localList')
		.data = .FindControl('Data')
		.tag = .FindControl('tag')
		.data.SetField('option', 'Start')
		.loadList('Start')
		}

	Controls()
		{
		.trialTags = LastContribution('Svc_TrialTags')
		return Object('Vert'
			Object('Record'
				Object('Vert'
				Object('RadioButtons', 'Start', 'End', horz:, name: 'option')
				Object('Pair'
					Object('Static', 'Tag'),
					Object('ChooseList', .trialTags.Map2({ |m, v| m $ ' - ' $ v }),
						name: 'tag'))))
			#Skip
			#(List columns: #(svc_checked, svc_lib, svc_type, svc_date,
				svc_local_date, svc_name),
				noShading:,	name: 'localList', defWidth: false,
				columnsSaveName: 'svc_local', stretchColumn: 'svc_name',
				checkBoxColumn: 'svc_checked'))
		}

	maps: (name: #svc_name,
		lib: #svc_lib,
		type: #svc_type,
		modified: #svc_date,
		committed: #svc_local_date)
	Record_NewValue(field, value)
		{
		// init
		if .list is false or field isnt 'option'
			return

		.loadList(value)
		}

	loadList(option)
		{
		if .trialTags.Empty?()
			return

		if option is 'Start'
			{
			.tag.SetReadOnly(false)
			.model.SetTable('svc_all_changes')
			for rec in .model.LocalChanges
				{
				rec.svc_checked = false
				.list.AddRow(.formatRec(rec))
				}
			}
		else if option is 'End'
			{
			.tag.SetReadOnly(true)
			data = Object()
			libs = Libraries().Append(UnusedStandardLibraries())
			for lib in libs
				{
				data.Append(QueryAll(lib $ '
					rename lib_modified to modified,
						lib_committed to committed
					extend lib = ' $ Display(lib) $ ',
						svc_checked = false,
						tag = LibraryTags.GetTagFromName(name)' $ Opt('
					where ', .trialTags.Members().Map(
						{ 'tag.Suffix?(' $ Display('_' $ it) $ ')' }).Join(' or ')) $ '
					where group is -1
					project name, group, lib, modified, committed, svc_checked'))
				}
			.list.Set(data.Map(.formatRec))
			}
		}

	formatRec(rec)
		{
		return rec.MapMembers({ .maps.GetDefault(it, it) })
		}

	List_SingleClick(row, col, source)
		{
		if source.GetCol(col) is 'svc_checked' and row isnt false
			{
			data = source.GetRow(row)
			data.svc_checked = data.svc_checked isnt true
			source.RepaintRow(row)
			}
		return 0
		}

	OK()
		{
		selects = .list.Get().Filter({ it.svc_checked is true })
		if selects.Empty?()
			{
			.AlertError(.Title, 'Please select at least one record')
			return false
			}

		if .data.Valid(forceCheck:) isnt true
			return false

		option = .data.Get()
		if option.option is 'Start' and option.tag is ''
			{
			.tag.SetValid(false)
			return false
			}

		block = option.option is 'Start' ? .copyAndRestore : .renameAndDelete
		.forEachSelects(selects, block, option.tag)
		return true
		}

	copyAndRestore(select, svcTable, t, tag)
		{
		srcRec = svcTable.Get(select.svc_name, :t)
		destName = LibraryTags.SetTrialTag(select.svc_name, tag, .trialTags.Members())

		svcTable.Output([
			name: destName,
			parent: srcRec.parent,
			text: srcRec.text,
			lib_invalid_text: srcRec.lib_invalid_text
			], :t)
		svcTable.Restore(select.svc_name, :t)
		}

	renameAndDelete(select, svcTable, t)
		{
		srcRec = svcTable.Get(select.svc_name, :t)
		destName = LibraryTags.SetTrialTag(select.svc_name, '', .trialTags.Members())
		if false is destRec = svcTable.Get(destName, :t)
			svcTable.Rename(srcRec, destName, :t)
		else
			{
			destRec.lib_invalid_text = srcRec.lib_invalid_text
			svcTable.Update(destRec, newText: srcRec.text, :t)
			svcTable.StageDelete(select.svc_name, :t)
			}
		}

	forEachSelects(selects, block, tag = '')
		{
		lib = false
		svcTable = false
		Transaction(update:)
			{ |t|
			for select in selects
				{
				if select.svc_lib isnt lib
					svcTable = SvcTable(lib = select.svc_lib)

				block(select, svcTable, t, :tag)
				}
			}
		}
	}