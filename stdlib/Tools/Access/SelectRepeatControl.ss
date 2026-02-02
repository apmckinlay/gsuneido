// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'SelectRepeat'
	MaxRecords: 50
	New(.sf, select_vals, .name, .option = '', .title = '',
		.noUserDefaultSelects? = false, .fromFilter = false, selChanged = false)
		{
		.filters = .FindControl('conditions')
		.scroll = .FindControl('Scroll')
		if fromFilter
			{
			minRows = Max((.Send('GetDefaultSelect').Size() + 1)/2, 1) // inital two rows
			.setScrollYmin(minRows)
			}
		// fill in values from previous select
		.Data.SetField('conditions', [].Merge(select_vals))
		.loadButton = .FindControl('loadButton') // false when from filter
		if selChanged
			.SetSelectApplied(false)
		}

	Controls()
		{
		initial = Object('Clear All', Object('Clear All'))
		if not .noUserDefaultSelects?
			initial.Add('Default Settings', #('Set My Default...', 'Change My Default...',
				'Reset To My Default'))
		return Object('Record',
			.fromFilter
				? .layoutFromFilter(initial)
				: .defaultLayout(initial)
			)
		}

	defaultLayout(initial)
		{
		return Object('Scroll', Object('Vert',
			Object('Horz',
				#(Static "Checkmark the rows you want to apply"),
				'Fill',
				.option is '' ? 'Skip'
					: Object('Presets', .option, .title, :initial)
				)
			'Skip',
			Object('ChooseFilters', useCheckBoxes:, fieldOptional:,
				maxRecords: .MaxRecords, name: 'conditions')))
		}

	layoutFromFilter(initial)
		{
		Assert(.option isnt: '')
		if 0 is extraLayout = .Send('Select_ExtraLayout')
			extraLayout = #()
		buttons = Object('HorzEqual',
			Object('EnhancedButton', 'Select', command: 'Select', buttonStyle:,
				mouseEffect:, weight: 'bold', pad: 40,
				name: 'loadButton', xstretch: 0), name: 'buttons')
		buttons.Add(#(Skip, small:))
		buttons.Add(Object('Button', 'Count', tip: SelectControl.CountBtnTip))
		buttons.Add(#(Skip, small:))
		buttons.Add(Object('Presets', .option, .title, :initial, xstretch: 0
			extraMenu: #('', 'Open Select Window'), alignHorz:))
		buttons.Add('Fill')
		VirtualListTopLayoutControl.BuildLayout(extraLayout, buttons)
		buttons.name = 'buttons'
		return Object(#Vert
			Object('Scroll'
				Object('ChooseFilters', useCheckBoxes:, fieldOptional:,
					maxRecords: .MaxRecords, name: 'conditions'))
			buttons)
		}

	setScrollYmin(rows)
		{
		.scroll.Ymin = (.filters.Ymin * rows).Round(0) + 4 /*= border */
		}

	Recv(@args)
		{
		source = args.source
		if source.Parent.Name is 'buttons' and args[0].Prefix?('On_')
			.Controller.Send(@args)
		return 0
		}

	On_Count()
		{
		if not .valid?()
			return
		.Send('On_Count')
		}

	On_Select()
		{
		if not .valid?()
			return
		if false is .Send('Select_Apply')
			return
		.SetSelectApplied(true)
		}

	SetSelectApplied(applied = false)
		{
		if .loadButton is false
			return
		.Send('SelectControl_SetSelectApplied', .valid?)
		.loadButton.SetTextColor(applied ? CLR.BLACK : CLR.RED)
		}

	Record_NewValue(field, value /*unused*/)
		{
		if .fromFilter is true and field is 'conditions'
			.SetSelectApplied(false)
		}
	SelectChanged?()
		{
		return .loadButton.GetTextColor() is CLR.RED
		}

	FieldPrompt_GetSelectFields()
		{
		return .sf
		}

	Where()
		{
		if .valid?()
			return .BuildWhere(.sf, .Get().conditions)
		return false
		}

	valid?(quiet = false)
		{
		if '' isnt msg = .filters.ForceValid()
			{
			if not quiet
				.AlertInfo('Select', msg)
			return false
			}
		return true
		}

	BuildWhere(sf, conditions)
		{
		joinflds = Object()
		whereConditions = conditions.DeepCopy()
		for condition in whereConditions
			{
			if .skip?(condition)
				continue
			op = condition[condition.condition_field].operation
			if op is ''
				continue
			fld = Select2.Empty_field(condition.condition_field, Object(op), sf)
			if fld isnt condition.condition_field
				{
				condition[fld] = condition[condition.condition_field]
				condition.Delete(condition.condition_field)
				condition.condition_field = fld
				}
			else if false isnt joinNums = GetForeignNumsFromNameAbbrevFilter(fld, sf,
				condition[fld].operation, condition[fld].value, condition[fld].value2)
				{
				condition[joinNums.numField] = joinNums.nums.Empty?()
					? [operation: 'less than', value: '', value2: '']
					: [operation: 'in list', value: joinNums.nums, value2: '']
				condition.Delete(fld)
				fld = condition.condition_field = joinNums.numField
				}
			joinflds.Add(fld)
			}
		fields = Object()
		where = ChooseFiltersControl.BuildWhereFromFilter(
			whereConditions, useCheckBoxes:, conditionFields: fields)
		return Object(:where, errs: '', :joinflds, :fields)
		}

	skip?(condition)
		{
		if condition.check isnt true or condition.condition_field is ''
			return true
		if not Object?(condition[condition.condition_field])
			{
			SuneidoLog('ERROR: (CAUGHT) condition field must be an object', calls:,
				params: condition, caughtMsg: 'Skipping invalid condition.  Check to ' $
					'see if bad value is saved.')
			return true
			}
		return false
		}

	Get() // from SelectControl
		{
		// in some cases conditions may not be present in the data,
		// like when Presets->Delete is done
		data = .Data.Get()
		if data.conditions is ''
			data.conditions = Object()
		return data
		}

	Set(data)  // from Presets
		{
		.Data.Set(data)
		}

	On_Presets_Clear_All_Clear_All()
		{
		.Data.SetField('conditions', Object())
		}

	On_Presets_Default_Settings_Set_My_Default()
		{
		if not .valid?()
			return
		data = .Get().conditions.Filter({ it.check is true })
		.chooseDefaultSelect(data)
		}

	chooseDefaultSelect(data)
		{
		if false is defaultFilters = DateCodeChooseFilterControl(data, .sf,
			"Default Selects")
			return

		defaultName = .name $ '~default'
		.deleteDefaultUserSelect(defaultName)
		QueryOutput("userselects", Record(
			userselect_user: Suneido.User,
			userselect_title: defaultName,
			userselect_selects: defaultFilters))
		}

	deleteDefaultUserSelect(defaultName)
		{
		QueryDo("delete userselects
			where userselect_user is " $ Display(Suneido.User) $
			" and userselect_title is " $ Display(defaultName))
		}

	On_Presets_Default_Settings_Change_My_Default()
		{
		if false is cur = Query1('userselects', userselect_user: Suneido.User
			userselect_title: .name $ '~default')
			.On_Presets_Default_Settings_Set_My_Default()
		else
			.chooseDefaultSelect(cur.userselect_selects)
		}

	On_Presets_Default_Settings_Reset_To_My_Default()
		{
		defaultName = .name $ '~default'
		if false isnt sel = Query1("userselects", userselect_user: Suneido.User,
			userselect_title: defaultName)
			{
			.ProcessPresets([conditions: sel.userselect_selects])
			.SetSelectApplied(false)
			}
		}

	newRow: false
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		if .newRow isnt false
			{
			.scrollToBottom(.newRow)
			.newRow = false
			}
		}
	scrollToBottom(newRow)
		{
		.Defer(uniqueID: 'SelectRepeat_ScrollToBottom')
			{
			.scroll.VSCROLL(MAKELONG(SB.PAGEDOWN, 0))
			.filters.FocusRow(newRow)
			}
		}

	HighlightLastRow(row)
		{
		.scroll.VSCROLL(MAKELONG(SB.BOTTOM, 0))
		.filters.FocusRow(row, value?:)
		}

	Repeat_RowsChanged(.newRow = false)
		{
		// could not delay and scroll to bottom here,
		// since the Repeat insert also triggers delay
		}

	On_Open_Select_Window()
		{
		.Send('Select_OpenDialog')
		}

	LoadPresets(value)
		{
		if value is false
			.SetSelectApplied(false)
		}

	ProcessPresets(params)
		{
		if not Object?(params) or not Object?(params.GetDefault('conditions', false))
			return true
		newConditions = params.conditions.DeepCopy()
		if '' is curConditions = .Data.Get().conditions
			curConditions = Object()
		curConditions.Each({ it.check = false })
		for c in newConditions
			curConditions.RemoveIf({ it.condition_field is c.condition_field })
		newConditions.Append(curConditions)
		newConditions = newConditions[.. .MaxRecords]
		.Set([conditions: newConditions])
		return true
		}

	GetPresetsSaveData()
		{
		conditions = .Data.Get().conditions
		conditions = conditions.Filter({ it.check is true })
		if conditions.Empty?()
			{
			.AlertWarn('Presets', 'No row checkmarked, cannot save Presets')
			return false
			}
		return [:conditions]
		}

	DateControl_ConvertDateCodes()
		{
		return false
		}
	}
