// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Presets Manager'
	CallClass(parent)
		{
		ModalWindow(Object(this, parent), keep_size:, closeButton?:,
			onDestroy:
				{|windowResult|
				if windowResult is false
					windowResult = 'Presets'
				parent.SetPresets(windowResult)
				})
		}


	New(.parent)
		{
		super(.controls())
		.list = .FindControl('VirtualList')
		}

	controls()
		{
		report = .parent.ReportDetails()
		.reportTitle = report.title
		.reportName = report.name
		.baseName = report.baseName
		return Object('Vert',
			Object('VirtualList',
				.parent.Query(),
				columns: #(preset_name, preset_accessible_to, preset_createdBy),
				enableMultiSelect:,
				columnsSaveName: 'ParamsManager',
				selectiveEditList:,
				headerSelectPrompt:,
				disableSelectFilter:,
				preventCustomExpand?:,
				sortSaveName: 'ParamsManager - Sort')
			#(Skip small:)
			#(Static, 'Double-click a preset to load it')
			#(Skip medium:)
			#(Horz,
				(Button, Rename)
				Skip
				(MenuButton, Delete, (Delete))
				Fill, Skip,
				(MenuButton, 'Accessible To', (Everyone, 'Only Me (and admin users)'))
				Fill, Skip,
				(MenuButton, Import, ('For Everyone', 'For Only Me')),
				Skip
				(Button, Export),
				)
			)
		}

	On_Rename()
		{
		if false is selected = .selectOne('Rename')
			return
		oldname = selected.preset_name
		if false isnt option = .parent.PresetsAsk('Rename "' $ oldname $ '" To')
			if true isnt result = .rename(oldname, option)
				.AlertInfo(.Title, result)
		}

	selectOne(option)
		{
		selected = .selected(suppressAlert?:)
		if selected.Size() isnt 1
			{
			.AlertInfo(.Title, 'You must select one Preset to ' $ option)
			return false
			}
		return selected[0]
		}

	selected(suppressAlert? = false)
		{
		selected = .list.GetSelectedRecords()
		if selected.Empty?() and not suppressAlert?
			.AlertInfo(.Title, 'No preset selected')
		return selected
		}

	rename(oldname, newname)
		{
		Transaction(update:)
			{ |t|
			if not t.QueryEmpty?(.parent.ReportQuery(newname))
				return 'Unable to rename. Another user has used the name'
			t.QueryDo('update ' $ .parent.ReportQuery(oldname) $
				' set report = ' $ Display(.parent.Report(newname)))
			}
		.list.Refresh()
		return true
		}

	On_Delete_Delete()
		{
		if .selected(suppressAlert?:).Size() > 1
			if not YesNo('You have mutliple presets marked for deletion.\r\n' $
				'This action cannot be undone.\r\n\r\n' $
				'Continue with delete?', title: .Title)
				return
		.forSelectedPresets(
			{
			QueryDo('delete ' $ .parent.ReportQuery(it))
			})
		}

	On_Accessible_To_Everyone()
		{
		.changeAccessForSelected(false)
		}

	changeAccessForSelected(private?)
		{
		illegalModify? = false
		.forSelectedPresets(
			{
			if false is .changeAccess(it, private?)
				illegalModify? = true
			})
		if illegalModify?
			.AlertInfo(.Title, 'You can only set a Preset to "Only Me" if you created it')
		}

	forSelectedPresets(block)
		{
		selected = .selected()
		selected.Each({ block(it.preset_name) })
		if not selected.Empty?()
			.list.Refresh()
		}

	changeAccess(presetName, private?)
		{
		modifyAllowed? = true
		user = .parent.UserInfo.user
		admin? = .parent.UserInfo.admin?
		QueryApply1(.parent.ReportQuery(presetName))
			{
			if modifyAllowed? = .changeAccessForPreset(it, user, admin?, private?)
				it.Update()
			}
		return modifyAllowed?
		}

	changeAccessForPreset(rec, user, admin?, private?)
		{
		if rec.user is '' or rec.user is user or admin?
			{
			if rec.user is '' or private?
				rec.user = user
			if not Object?(rec.report_options)
				rec.report_options = []
			rec.report_options.private? = private?
			}
		else if private?
			return false
		return true
		}

	On_Accessible_To_Only_Me_and_admin_users()
		{
		.changeAccessForSelected(true)
		}

	On_Import_For_Everyone()
		{
		.import(false)
		}

	On_Import_For_Only_Me()
		{
		.import(true)
		}

	type: 'presets'
	filter: "Options (*.opt)\x00*.opt"
	import(private?)
		{
		title = 'Import Presets'
		data = ImportExportObject.Import(title, .type, .filter, .Window.Hwnd)
		if data is false
			return
		if not .compatiblePresets?(data.report_name)
			{
			.AlertWarn(title,
				'Sorry, these presets belong to another option:\n\n' $
				'\t' $ data.report_title $ '\n\n' $
				'Please import them there.')
			return
			}
		if not data.Member?(#preset_name)
			return false
		// Handle importing old .opt files
		params = data.Member?(#params)
			? data.params	// new way (after 29508)
			: data			// old way (prior to 29508)
		.parent.OutputParam(data.preset_name, params, report_options: [:private?])
		.list.Refresh()
		}

	compatiblePresets?(reportName)
		{
		return .baseName isnt false
			? reportName.Has?(.baseName)
			: reportName is .reportName
		}

	On_Export()
		{
		if false is rec = .selectOne('Export')
			return
		rec.report_name = .reportName
		rec.report_title = .reportTitle
		ImportExportObject.Export('Export Presets', rec, .type, .filter, .Window.Hwnd,
			name: rec.preset_name)
		}

	VirtualList_DoubleClick(rec, col /*unused*/)
		{
		if rec is false
			return false
		QueryDo('update ' $ .parent.ReportQuery(rec.preset_name) $
			' set params_TS = ' $ Display(Timestamp()))
		.Parent.Result(rec.preset_name)
		return false
		}

	VirtualList_AddGlobalMenu?()
		{
		return false
		}
	}
