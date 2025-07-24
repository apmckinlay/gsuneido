// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: "ParamsSelect"
	ComponentName: 'ParamsSelect'
	New(field, set = false, .emptyValue = 'all', paramPrompt = false, .excludeOps = #())
		{
		super(.makecontrols(field, emptyValue, paramPrompt))
		.Send("Data")
		.field = field
		.horz = .Pair.WndPane.Horz
		.op = .horz.operation
		.val = .horz.GetChildren()[.val_i]
		if not .convertDateCodes?
			{
			.horz.Append(#(Skip, small:))
			.horz.Append(#(Static, name: datePreview))
			.horz.Append(#(Skip, small:))
			.horz.Append(Object('HelpButton', imageColor: CLR.Inactive,
				tip: 'Date Shortcuts', noAccels?:, size: '-3'))
			}
		.val.SetReadOnly(true)
		.Top = .Pair.Top
		.Left = .Pair.Left
		.datePreview = .FindControl('datePreview')
		if set isnt false
			{
			.Set(set)
			.Send("NewValue", .Get())
			}
		if .Name is ''
			.SetReadOnly(true)
		}
	string?: true
	useFieldForCompare?: false
	makecontrols(field, emptyValue, paramPrompt)
		{
		control = Object()
		prompt = ''
		if String?(field)
			{
			prompt = paramPrompt isnt false ? paramPrompt : Prompt.WithInfo(field)
			dd = Datadict(field)
			if not dd.Base?(Field_string) or dd.Control[0] is 'ChooseDates'
				.string? = false
			control = dd.Member?('SelectControl') ? dd.SelectControl : dd.Control
			.Name = field
			}
		else if Object?(field)
			{
			prompt = field.Member?('prompt') ? field.prompt : field.name
			control = field
			.Name = field.name
			.string? = .IsStringControl?(control)
			}
		else
			throw "bad field type"

		.useFieldForCompare? = .getUseFieldForCompare?(control)

		control = .prepareParamControl(control, field)
		.initializeOps(control, emptyValue)

		return .layout(control, prompt, emptyValue)
		}

	layout(control, prompt, emptyValue)
		{
		horz = Object('Horz'
			Object('ParamsSelectChooseButton', emptyValue, .ops_trans,
				width: 20, :emptyValue, name: 'operation'),
			#(Skip 2), control, overlap:, xstretch: 1)

		controls = Object("WndPane" horz windowClass: "SuBtnfaceArrow")
		return prompt is ""
			? Object('Horz', controls, name: 'Pair')
			: Object('Pair' Object('Static' prompt) controls)
		}

	// default to true
	IsStringControl?(ctrl)
		{
		if ctrl.Empty?() or not ctrl.Member?(0)
			return true

		// is Date type
		if .ControlName(ctrl).Has?('Date')
			return false

		if false is ctrlClass = GetControlClass.FromControl(ctrl)
			return true

		// is Number, Dollar or CheckBox type
		return not (ctrlClass.Base?(NumberControl) or ctrlClass.Base?(CheckBoxControl))
		}

	getUseFieldForCompare?(control)
		{
		if .string? is false
			return false
		if false is ctrlClass = GetControlClass.FromControl(control)
			return false
		return not ctrlClass.Base?(ChooseControl)
		}

	ControlName(ctrl)
		{
		if ctrl.Empty?()
			return ''
		ctrlName = ctrl[0]
		return String?(ctrlName) ? ctrlName : ''
		}

	initializeOps(control, emptyValue)
		{
		.id = control[0] is 'Id'
		if .id
			.ops_org = Object(emptyValue, "equals", "not equal to", "empty",
				"not empty", "in list", "not in list")
		else if control[0] is 'ChooseDates'
			.ops_org = Object(emptyValue, "equals", "not equal to", "empty", "not empty",
				"contains", "does not contain")
		else if control[0] is 'DateRange'
			.ops_org = Object(emptyValue, "empty", "not empty")
		else if .string?
			.ops_org = Object(emptyValue,
				"less than", "less than or equal to",
				"greater than", "greater than or equal to",
				"equals", "not equal to", "empty", "not empty",
				"contains", "does not contain", "starts with", "ends with",
				"matches", "does not match", "range", "not in range",
				"in list", "not in list")
		else if control[0] is 'CheckBox'
			.ops_org = Object(emptyValue, "equals", "not equal to")
		else
			.ops_org = Object(emptyValue, "less than", "less than or equal to",
				"greater than", "greater than or equal to", "equals", "not equal to",
				"empty", "not empty", "range", "not in range", "in list", "not in list")

		.ops_org.Remove(@.excludeOps)
		.ops_trans = Object()
		for op in .ops_org
			.ops_trans.Add(TranslateLanguage(op))
		}
	prepareParamControl(control, field)
		{
		control = .convertControl(control.Copy())
		control.mandatory = true
		.control = control
		control.readonly = false
		ctrlName = not control[0].Suffix?('Control')
			? control[0] $ 'Control'
			: control[0]
		if control[0] is 'Key' or Global(ctrlName).Base?(KeyControl)
			{
			control.filterOnEmpty = false
			control.saveInfoName = field
			control.fillin = false
			}
		else if .control[0].Has?('Date')
			control.checkCurrent = control.checkValid = false
		control.weight = ''
		return control
		}

	showDate()
		{
		if .datePreview is false
			return

		.datePreview.Set(.dateString())
		}

	dateString()
		{
		if not .Valid?()
			return ''

		date1 = .val.Get()
		if .op.Get().Has?('range')
			{
			date2 = .val2.Get()
			if Date?(date1) and Date?(date2)
				return ''
			return ' ' $ .convertDate(date1)  $ ' to ' $ .convertDate(date2)
			}
		if String?(date1) and (.op.Get().Has?('equal') or .op.Get().Has?('than'))
			return ' ' $ .convertDate(date1)
		return ''
		}

	convertDate(text)
		{
		showTime = .control.showTime
		if false is date = DateControl.ConvertToDate(text, convertDateCodes?:, :showTime)
			return ''
		return showTime ? date.ShortDateTime() : date.ShortDate()
		}

	convertControl(control)
		{
		// switch Editor to Field to save space on screen,
		// Editor is overkill for params anyway
		if control[0] in ("Editor", "ScintillaAddonsEditor", "CourierEditor")
			control = Object('Field', width: 30, xstretch: 1)
		if control[0] is 'RadioButtons'
			control = .convertRadioButtonsToChooseList(control)
		if control[0].Has?('Static')
			control = .replaceStaticField(control[0])
		return control
		}
	replaceStaticField(static)
		{
		if static.Has?('Time')
			return Object('ChooseDateTime')
		if static.Has?('Date')
			return Object('ChooseDate')
		return Object('Field')
		}
	convertRadioButtonsToChooseList(control)
		{
		choices = control.ProjectValues(control.Members(list:)).Delete(0)
		return Object("ChooseList", choices)
		}
	val_i: 2
	val2_i: 3
	Valid?()
		{
		operation = .getOperation()

		if operation is ''
			return .validateEmptyOperation()

		if operation in ("empty", "not empty")
			return true

		if not .op.Valid?() or not .val.Valid?()
			return false

		if operation.Has?('range')
			return .rangeValid?()

		if operation.Has?('match')
			return Regex?(.val.Get())

		return true
		}

	getOperation()
		{
		operation = .op.Get()
		if operation isnt '' and .op.Valid?()
			operation = .ops_org[.ops_trans.Find(operation)]
		return operation
		}

	validateEmptyOperation()
		{
		// control from ListEditWindow does not commit if invalid
		if .Controller.Base?(ListEditWindow)
			return true
		cur = .Send('GetField', .Name)
		if not Object?(cur)
			return true
		// validate between set value vs display for saved data
		return cur.GetDefault('operation', false) in ('', .emptyValue)
		}

	rangeValid?()
		{
		val2_valid? = .val2.Valid?()
		if not .val.Method?("ValidateRange?") or .val.ValidateRange?()
			return .val2.Get() >= .val.Get() and val2_valid?
		return val2_valid?
		}

	find_operation(ob)
		{
		if ob.operation is '' or false is (pos = .ops_org.Find(ob.operation))
			return ''
		return .ops_trans[pos]
		}
	Set(ob)
		{
		.op.Set("")
		.val.Set("")
		if Object?(ob)
			{
			if (ob.Member?('operation'))
				{
				.op.Set(.find_operation(ob))
				// need to trigger control arranging code
				.arrange_controls(ob.operation)
				}
			if ob.Member?('value')
				.val.Set(ob.value)
			if ob.Member?('value2') and .op.Get().Has?('range')
				.val2.Set(ob.value2)
			.showDate()
			}
		else
			.arrange_controls('')
		.setReadOnly(Object?(ob) and ob.Member?('operation') ? ob.operation : "")
		}
	Get()
		{
		if (('' isnt operation = .op.Get()) and .op.Valid?())
			operation = .ops_org[.ops_trans.Find(operation)]

		ob = .id or .valType is #inList or not operation.Has?('range')
			? Object(:operation, value: .val.Get(), value2: "")
			: Object(:operation, value: .val.Get(), value2: .val2.Get())
		if .convertDateCodes? is false and .hasDateCode?(ob)
			ob.dateConversionInfo = Object(showTime: .control.showTime)
		return ob
		}

	hasDateCode?(ob)
		{
		return String?(ob.value) or ob.value2 isnt "" and String?(ob.value2)
		}

	GetFieldName()
		{
		return .Name
		}
	// valType can be '', 'inList' or 'fieldOp'
	valType: ''
	// resend newvalue so that this controller becomes the source and
	// thus responsible for Get method.
	NewValue(value, source)
		{
		// prevent issues from trying to get values from destroyed controls
		// sometimes they send NewValue in the process of being destroyed
		if source.Destroyed?()
			return
		if source is .op and .op.Valid?()
			{
			op = value is '' ? '' : .ops_org[.ops_trans.Find(value)]
			.arrange_controls(op)
			}
		if source.Name.Has?('ChooseDate')
			.showDate()
		if source is .op
			.showDate()
		.Send('NewValue', .Get())
		}
	Dirty?(dirty = '')
		{
		dirty = .val.Dirty?(dirty)
		dirty2 = .val2 isnt false ? .val2.Dirty?(dirty) : false
		return dirty or dirty2
		}
	Data() // block Data message
		{
		}
	arrange_controls(op)
		{
		.handleSpecialOp(op)

		// clear out the control's values where necessary
		if (not .ops_with_val1.Has?(op))
			.val.Set("")
		.handleRange(op)

		.setReadOnly(op)
		}
	val2: false
	handleRange(op)
		{
		if not op.Has?("range") and .val2 isnt false
			{
			.horz.Remove(.val2_i)
			.val2 = false
			}
		else if op.Has?('range') and .val2 is false
			{
			.horz.Insert(.val2_i,
				.valType is #fieldOp ? .fieldCtrl : .control)
			.val2 = .horz.GetChildren()[.val2_i]
			}
		}
	handleSpecialOp(op)
		{
		if .handleInList(op) is true
			return
		if .handleStringOps(op) is true
			return
		if .handleCompareOp(op) is true
			return
		.restoreControl()
		}

	handleInList(op)
		{
		if not op.Has?("in list") or .field is ""
			return false

		if .valType isnt #inList
			{
			.changeControl(Object('ParamsChooseList', .field, values: Object(),
				readonly: false))
			.valType = #inList
			}
		return true
		}

	fieldCtrl: #('Field', width: 20, mandatory:, xstretch: 1)
	handleStringOps(op)
		{
		if .string? isnt true or not .stringOps.Any?({ op.Has?(it) })
			return false

		if .valType isnt #fieldOp
			{
			.changeControl(.fieldCtrl)
			.valType = #fieldOp
			}
		return true
		}
	handleCompareOp(op)
		{
		if .useFieldForCompare? isnt true or not .compareOps.Any?({ op.Has?(it) })
			return false

		if .valType isnt #fieldOp
			{
			.changeControl(.fieldCtrl)
			.valType = #fieldOp
			}
		return true
		}

	restoreControl()
		{
		if .valType isnt ''
			{
			.changeControl(.control)
			.valType = ''
			}
		}

	stringOps: ('contain', 'starts with', 'ends with', 'match')
	compareOps: ('greater than', 'less than', 'range')
	changeControl(control)
		{
		.horz.Remove(.val_i)
		.horz.Insert(.val_i, control)
		.val = .horz.GetChildren()[.val_i]
		}
	ops_with_val1: #("greater than", "greater than or equal to",
		"less than", "less than or equal to", "equals", "not equal to",
		"contains", "does not contain", "starts with", "ends with",
		"does not match", "matches", "range", "not in range")
	setReadOnly(op)
		{
		if .valType isnt #inList
			.val.SetReadOnly(not .ops_with_val1.Has?(op))
		}
	SetReadOnly(readonly)
		{
		if .field is ''
			readonly = true
		super.SetReadOnly(readonly)
		if readonly isnt true
			.setReadOnly(.op.Get())
		}
	SetValid(valid? = true)
		{
		op = .op.Get()
		if .stringOps.Any?({ op =~ it }) or .compareOps.Any?({ op =~ it })
			.val.SetValid(valid?)
		}
	FocusValue()
		{
		.val.SetFocus()
		}
	HelpButton_HelpPage()
		{
		return '/General/Reference/Keyboard Shortcuts#keyboarddates'
		}
	convertDateCodes?: true
	DateControl_ConvertDateCodes(showTime = false)
		{
		result = .Send('DateControl_ConvertDateCodes')
		if not (.convertDateCodes? = result is 0)
			.control.showTime = showTime
		return result
		}
	Destroy()
		{
		.ClearFocus() // make sure afterfield processing is done before destroyed
		.Send("NoData")
		super.Destroy()
		}
	}