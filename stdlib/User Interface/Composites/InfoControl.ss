// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// Note: have to use WndPane so tab order is not messed up
// by removing and inserting controls
Controller
	{
	ComponentName: 'Info'
	ImagePadding: .25
	New(def = 0, .mandatory = false, .allowOnlyType = false, .tabover = false,
		types = false, controls = false, typeWidth = 7, .controlWidth = 25,
		readonly = false, .hidden = false)
		// pre: def is 0 -> .types.Size()
		{
		super(.layout(def, types, controls, typeWidth))
		.innerHorz = .FindControl('innerHorz')
		.outerHorz = .FindControl('outerHorz')
		.typeField = .FindControl('ChooseList')
		.valueField = .FindControl('Field')
		.allLabels = GetContributions(.Name.BeforeFirst('_').Capitalize() $ 'InfoLabels')
		.valueField.AddContextMenuItem("", "", .contextMenuEnabled)
		.valueField.AddContextMenuItem('Labels', .On_Labels, .contextMenuEnabled)
		.typeField.Set(.types[.def])
		.Top = .valueField.Top
		.Left = .typeField.Xmin - 1
		.defaultTip = .valueField.Status
		// initial SetReadOnly call has to happen before .readonly is set because the
		// SetReadOnly method checks .readonly and doesn't do anything if true
		.SetReadOnly(readonly)
		.readonly = readonly
		.Send('Data')
		}

	controls: (Phone, Phone, Phone, Phone, Phone, MailLink, HttpLink, Field, Field)

	layout(def, types, controls, typeWidth)
		{
		.types = types is false ? InfoTypes : types
		.controls = controls is false ? .controls : controls
		.def = Function?(def) ? def() : def
		return Object('Horz',
			Object("WndPane"
				Object("Horz",
					Object('ChooseList', .types, width: typeWidth,
						mandatory:, tabover: .tabover, hidden: .hidden),
					.control(.def),
					xstretch: false, overlap:, name: 'innerHorz'),
				windowClass: 'SuBtnfaceArrow'),
			name: 'outerHorz')
		}

	control(i)
		{
		.curControl = .controls[i]
		return [.controls[i], mandatory: .mandatory, name: "Field", width: .controlWidth,
			tabover: .tabover, hidden: .hidden]
		}

	// LabelDelimiter needs a space at the beginning for Biz_GetInfoFields to work
	LabelDelimiter: ' \x85Label\x85 '
	allLabels: ()
	labels: ''
	GetInfoLabels(value)
		{
		return value.AfterFirst(.LabelDelimiter).Split(',').Map(#Trim)
		}

	On_Labels()
		{
		typeLabels = .allLabels.GetDefault(.curControl, #()).Sort!()

		if .curControl is "MailLink" and not .valueField.ValidLink?(.valueField.Get()) and
			.getValidLabels().Blank?()
			{
			.AlertWarn('Labels', 'Can not add labels to an invalid email address')
			return
			}

		if false is labels = OkCancel(Object('ChooseManyList', typeLabels,
			value: .getValidLabels()), title: 'Labels', hwnd: .Window.Hwnd)
			return

		.labels = Opt(.LabelDelimiter, labels)
		.setLabelTagAndToolTip()
		.NewValue(.Get(), .valueField)
		}

	Set(value)
		{
		type_end = value.Find(':')
		field_start = type_end + 2
		label_start = value.Find(.LabelDelimiter)
		.labels = value[label_start ..]

		if value is ""
			{
			.typeField.Set(.types[.def])
			.setFieldControl(.types[.def])
			.valueField.Set("")
			}
		else
			{
			type = value[.. type_end + 1]
			.typeField.Set(type)
			.setFieldControl(type)
			.valueField.Set(value[field_start :: label_start - field_start])
			}
		}

	Get()
		{
		if '' isnt value = .valueField.Get()
			value = ' ' $ value
		return value isnt '' or .allowOnlyType ? .typeField.Get() $ value $ .labels : ''
		}

	NewValue(value, source)
		{
		if source is .typeField and .typeField.Valid?()
			.typeValueChanged(value)
		else if source is .valueField and .valueField.Valid?()
			.valueChanged()
		.Send('NewValue', .Get()) // send AFTER formatting
		}

	typeValueChanged(value)
		{
		prev_val = .valueField.Get()
		.setFieldControl(value)
		// remove invalid labels when switching types
		labelOb = .getValidLabels().Split(',').Intersect(.GetInfoLabels(.labels))
		.labels = Opt(.LabelDelimiter, labelOb.Join(','))
		.valueField.Set(prev_val)
		.valueField.KillFocus() // formatting for new control if necessary
		}

	valueChanged()
		{
		if .curControl is "MailLink" and .valueField.Get().Blank?()
			{
			.labels = ''
			.setLabelTagAndToolTip()
			}
		}

	Dirty?(dirty = "")
		{
		cd = .typeField.Dirty?(dirty)
		fd = .valueField.Dirty?(dirty)
		return cd or fd
		}

	ValidData?(@args)
		{
		if false is i = args.GetDefault('types', InfoTypes).
			FindIf({ args[0].Prefix?(it) })
			return false
		ctrl = Global(args.GetDefault('controls', .controls)[i] $ 'Control')
		args[0] = args[0].AfterFirst(': ')
		return ctrl.Method?('ValidData?')
			? ctrl.ValidData?(@args) : true
		}

	Valid?()
		{
		if .readonlyState is true
			return true

		if .curControl is "MailLink" and not .valueField.ValidLink?(.valueField.Get()) and
			not .getValidLabels().Blank?()
			return false

		return .typeField.Valid?() and .valueField.Valid?()
		}

	SetValid(valid = true)
		{
		if .valueField.Method?('SetValid')
			.valueField.SetValid(valid)
		}

	setFieldControl(type)
		{
		hasfocus? = .valueField.HasFocus?()
		if false is i = .types.Find(type)
			{
			i = .types.Find("Other:")
			.typeField.SetValid(false)
			}
		else
			.typeField.SetValid(true)

		if .controls[i] is .curControl
			{
			.setLabelTagAndToolTip()
			return
			}

		.innerHorz.Remove(1)
		.innerHorz.Insert(1, .control(i))
		.valueField = .FindControl('Field')
		.valueField.SetReadOnly(.readonlyState)
		.defaultTip = .valueField.Status
		.valueField.AddContextMenuItem("", "", .contextMenuEnabled)
		.valueField.AddContextMenuItem('Labels', .On_Labels, .contextMenuEnabled)
		if hasfocus?
			.setFieldFocus()
		.setLabelTagAndToolTip()
		}

	setFieldFocus()
		{
		// plain SetFocus doesn't work - control not finished construction
		PostMessage(.Window.Hwnd, WM.APP_SETFOCUS, 0, .valueField.Hwnd)
		}

	labelTag: false
	setLabelTagAndToolTip()
		{
		if .allLabels.Empty?()
			return

		validLabels = .getValidLabels()
		.valueField.SetStatus(validLabels.Blank?() ? .defaultTip : validLabels)

		if validLabels.Blank?()
			{
			if .labelTag isnt false
				.outerHorz.Remove(.outerHorz.Tally() - 1)
			}
		else if .labelTag is false
			.outerHorz.Append(
				Object('EnhancedButton',
					image: 'L', command: 'LabelsBtn', imageColor: CLR.Highlight,
					imagePadding: .ImagePadding, name: 'Labels'))

		if false isnt .labelTag = .FindControl('Labels')
			.labelTag.ToolTip(validLabels)
		}

	getValidLabels()
		{
		typeLabels = .allLabels.GetDefault(.curControl, #())
		selectedLabels = .labels.AfterFirst(.LabelDelimiter).Split(',')
		return selectedLabels.Intersect(typeLabels).Join(',')
		}

	On_LabelsBtn()
		{
		enabledOb = .contextMenuEnabled()
		if enabledOb.addToMenu is true and enabledOb.state is MFS.ENABLED
			.On_Labels()
		}

	contextMenuEnabled()
		{
		addToMenu? = .typeField.Valid?() and .valueField.Valid?() and
			not .allLabels.GetDefault(.curControl, #()).Empty?()
		state = not .readonlyState ? MFS.ENABLED : MFS.DISABLED
		return Object(addToMenu: addToMenu?, :state)
		}

	HandleTab()
		{
		if .typeField.HasFocus?()
			{
			.valueField.SetFocus()
			return true
			}
		return false
		}

	readonlyState: false
	readonly: false
	SetReadOnly(readonly)
		{
		if .readonly
			return

		.readonlyState = readonly
		super.SetReadOnly(readonly)
		}

	Status(status)
		{
		.Send('Status', status)
		}

	Field_SetFocus()
		{
		.Send('Field_SetFocus')
		}

	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}