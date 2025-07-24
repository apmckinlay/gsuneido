// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// e.g	Window(#(Border (Record (Vert
//			(Presets name title)
//			Skip
//			(Field name: one)
//			(Field name: two)
//			xstretch: 1, ystretch: 1))), w: 300, h: 200)
// TODO gray out name of preset if user changes data
Controller
	{
	Name: 'Presets'
	Separator: '~presets~'
	New(report_name, report_title = false, initial = #(New), saveCtrlsOnly = false,
		.extraMenu = #(), alignHorz = false, .baseName = false)
		{
		super(['MenuButton', 'Presets', menu: false, tabover:,
			tip: "Save and Load different settings"])
		.mb = .GetChild()
		if not alignHorz
			.Top = .mb.Ymin - 3 /* = to align with Title */
		.report_name = report_name
		.report_title = report_title is false ? report_name : report_title
		.saveCtrlsOnly = saveCtrlsOnly
		.initial = initial
		if 0 is .rc = .Send(#GetRecordControl)
			throw 'PresetsControl must be used inside a RecordControl'
		}

	SetPresets(preset = 'Presets')
		{
		.presets = .list()
		.On_Presets(preset)
		}

	presets: false
	limitWarning: '(see all presets)'
	MenuButton_Presets()
		{
		if .presets is false
			.presets = .list()
		menu = .initial.Copy()
		menu.Add('Manage...' $ Opt(' ', .listTrimmed? ? .limitWarning : ''), '')
		if not .presets.Empty?()
			menu.Add(@.presets).Add('')
		if .presets.Has?(name = .mb.Get()) and
			false is .standardPresets.FindOne({ it.preset_name is name })
			menu.Add('Save', 'Save As...')
		else
			menu.Add('Save As...')
		menu.Add(@.extraMenu)
		return menu
		}

	On_Presets(option)
		{
		buttonText = false
		if .Member?(m = 'Presets_' $ option.Tr('. '))
			buttonText = this[m]()
		else if .extraMenu.Has?(option) or .initial.Has?(option)
			{
			.Send('On_' $ ToIdentifier(option))
			if .initial.Has?(option)
				buttonText = 'Presets'
			}
		else if option is 'Manage...' or option is 'Manage... ' $ .limitWarning
			PresetsManagerControl(this)
		else
			{
			.load(option)
			buttonText = option
			}
		if buttonText isnt false
			.mb.Set(buttonText)
		}

	ReportDetails()
		{
		return [title: .report_title, name: .report_name, baseName: .baseName]
		}

	On_Presets_Clear_All_Clear_All()
		{
		.Send('On_Presets_Clear_All_Clear_All')
		}

	On_Presets_Default_Settings(option)
		{
		.Send('On_Presets_Default_Settings_' $ ToIdentifier(option))
		}

	Presets_New()
		{
		.rc.Set([])
		return 'Presets'
		}

	Presets_Save()
		{
		.Save(.mb.Get())
		return false
		}

	saveAs: Controller
		{
		Controls()
			{
			return #(Vert
				(Pair
					(Static 'Save As')
					Field)
				(Pair
					('StyledStatic', 'NOTE:', textStyle: 'note')
					('StyledStatic', 'Only saves checkmarked rows', textStyle: 'note')))
			}
		Get()
			{
			return .FindControl('Field').Get()
			}
		}

	Presets_SaveAs()
		{
		if false is params = .Send('GetPresetsSaveData')
			return false
		askArgs = params is 0 ? #('Save As') : Object('', ctrl: .saveAs)
		if false is option = .PresetsAsk(@askArgs)
			return false
		.Save(option)
		.SetPresets(option)
		return option
		}

	PresetsAsk(prompt, ctrl = false)
		{
		if ctrl is false
			ctrl = Object("Field", mandatory:)
		saveName = Ask(prompt, 'Presets', hwnd: .Window.Hwnd, valid: .valid, :ctrl)
		if false is Query1(.ReportQuery(saveName))
			return saveName
		.AlertError(.Title, 'Preset name already used')
		return false
		}

	valid(s)
		{
		maxNameLength = 30
		if s =~ '[\r\n\t]'
			return 'Name cannot contain newlines or tabs'
		if s.Size() >= maxNameLength
			return 'Name must be less than ' $ maxNameLength $ ' characters long'
		if s.Blank?()
			return 'Name cannot be empty'
		return ''
		}

	Save(name)
		{
		if 0 is params = .Send('GetPresetsSaveData')
			params = .saveCtrlsOnly ? .rc.GetControlData() : .rc.Get().Copy()
		if params is false
			return
		Params.RemoveIgnoreFields(params)
		.OutputParam(name, params)
		}

	OutputParam(name, params, report_options = '')
		{
		report = .Report(name)
		RetryTransaction()
			{ |t|
			t.QueryDo('delete ' $ .ReportQuery(name))
			t.QueryOutput('params',
				[:report, :params, user: .UserInfo.user, :report_options])
			}
		}

	limit: 20
	listTrimmed?: false
	list()
		{
		userPresets = TableExists?('params')
			? QueryAll(.Query(sortTS?:)).Map!({ it.ProjectValues('preset_name')[0] })
			: []

		.listTrimmed? = userPresets.Size() > .limit
		stdPresets = .getStandardPresets()
		if not stdPresets.Empty?() and not userPresets.Empty?() // add divider
			stdPresets.Add('')

		return stdPresets.Append(userPresets[.. .limit].Sort!())
		}

	standardPresets: #()
	getStandardPresets()
		{
		presets = Object()
		.standardPresets = Object()
		for cont in GetContributions('StandardPresets')
			if cont.ReportName is .report_name
				{
				for item in cont.Presets
					{
					preset = item.Copy()
					preset.preset_name = preset.report.AfterLast('~')
					presets.Add(preset.preset_name)
					.standardPresets.Add(preset)
					}
				}
		return presets.Sort!()
		}

	Query(sortTS? = false)
		{
		q = 'params
			extend
				preset_name = report.AfterLast("' $ .Separator $ '"),
				preset_createdBy = user,
				preset_accessible_to,
				curUser = ' $ Display(.UserInfo.user) $ ',
				admin? = ' $ .UserInfo.admin? $ ',
				preset_accessible?
			where report > ' $ Display(.Report('')) $
			' and report < ' $ Display(.Report('~')) $
			' and preset_accessible? is true '
		q $= sortTS?
			? 'sort reverse params_TS'
			: 'sort preset_name'
		return q
		}

	Getter_UserInfo()
		{
		roles = Suneido.user_roles
		admin? = Object?(roles) or String?(roles) ? roles.Has?('admin') : false
		userInfo = [user: Suneido.User, :admin?]
		return .UserInfo = OptContribution('CurrentUserInfo', { userInfo })()
		}

	load(name)
		{
		if false is rec = .standardPresets.FindOne({ it.preset_name is name })
			{
			query = .ReportQuery(name)
			if QueryCount(query) > 1
				SuneidoLog.Once('ERROR: (CAUGHT) query returned multiple presets',
					params: [:name, :query])
			rec = false
			RetryTransaction()
				{ |t|
				t.QueryApply(query)
					{
					it.Update() // touching the record so param_TS updates
					rec = it
					}
				}
			if rec is false
				return
			}
		.Send(#LoadPresets, true)

		if 0 is .Send(#ProcessPresets, rec.params)
			{
			.rc.Set(Record())
			for m in rec.params.Members()
				if not .saveCtrlsOnly or .rc.GetControl(m) isnt false
					.rc.SetField(m, rec.params[m])
			}
		.Send(#LoadPresets, false)
		}

	ReportQuery(name)
		{
		return 'params where report is ' $ Display(.Report(name))
		}

	Report(name)
		{
		return .report_name $ .Separator $ name
		}

	BeforeCopy(copyRec)
		{
		if false isnt pos = copyRec.report.FindRxLast('_\d\d\d\d\d\d\d\d_(\d)+$')
			copyRec.report = copyRec.report[.. pos]
		copyRec.report $= .getSuffix()
		}

	getSuffix()
		{
		return String(Timestamp()).Tr('#.', '_')
		}
	}
