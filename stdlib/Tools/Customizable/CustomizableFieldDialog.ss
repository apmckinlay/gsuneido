// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
// TODO: add "mandatory" checkbox
Controller
	{
	originalPeditor: false
	originalData: false
	control: false
	CallClass(fieldobject = false, text = 'Edit Properties')
		{
		OkCancel(Object(this, fieldobject, text), text)
		}
	New(fieldobject = false, .title = '')
		{
		super(.makecontrols(fieldobject))
		.peditorslot = .FindControl('peditorslot')
		peditor = .getPeditor()
		if false isnt peditor
			{
			if false is type = .findType(.type)
				type = #()
			.peditorslot.Append(Object(peditor, typePlugin: type))
			.peditor = .FindControl('peditor')
			.peditor.Set(fieldobject.flddef)
			.originalPeditor = .peditor.Copy()
			.originalData = .originalPeditor.Get()
			}
		if fieldobject isnt false
			{
			if false isnt colpro = .FindControl('colpro')
				colpro.Set(fieldobject.flddef.Prompt)
			if false isnt ctllbl = .FindControl('ctllbl')
				ctllbl.Set(.type)
			}
		.Data.AddObserver(.Observer)
		}

	getPeditor()
		{
		peditor = 'CustomizableFieldDialogPropertiesEditor_' $ .ctlnme $ 'Control'
		try
			Global(peditor)
		catch (e)
			{
			if not e.Prefix?(`can't find`) // control does not have an editor class
				SuneidoLog('ERROR: (CAUGHT) ' $ e, calls:,
					caughtMsg: 'no properties added to edit')
			return false
			}
		return peditor
		}
	makecontrols(.fieldobject)
		{
		.types = CustomFieldTypes()

		.ctlnme = ''
		.colnme = ''
		.fldbse = ''
		.custpe = ''

		//Add
		if false is fieldobject
			return .controls(editName:, editType:)
		//Update
		.ctlnme = fieldobject.flddef.Control[0]
		.colnme = fieldobject.colnme
		.fldbse = fieldobject.fldbse
		.custpe = fieldobject.custpe
		.type = false

		ctrl = .types.FindIf(.matchingControl)

		if ctrl is false
			return .noPropertiesCtrl()

		typeOb = .types[ctrl]
		.type = typeOb.name
		editType = CustomFieldTypes.GetCompatible(typeOb)

		.peditor = .getPeditor()

		if typeOb.GetDefault('oneWay', false)
			return .noPropertiesCtrl()

		return false is .peditor and editType is false
			? .noPropertiesCtrl()
			: .controls(editName:, :editType, type: typeOb.base)
		}

	matchingControl(c)
		{
		return Datadict(c.base).Control[0] is .ctlnme
		}

	warning: 'Changing the name will not update Reporter Reports,\r\n' $
		'Custom Field Fill-Ins, and some other screens'
	warnStatic()
		{
		return Object('StyledStatic', .warning,	alwaysBold:, textStyle: 'note')
		}
	noPropertiesCtrl()
		{
		return Object('Record'
			Object('Vert'
				#Skip
				Object('Pair' #(Static 'Name')
				Object('Field' name: 'colpro', mandatory:, set: .custpe)),
				.warnStatic()
			))
		}

	controls(editName = false, editType = false, type = false)
		{
		ob = Object()
		ob.Add('Record')
		vert = Object('Vert' name: 'Vert')
		vert.Add('Skip')

		if editName is true
			{
			vert.Add(Object('Pair' #(Static 'Name')
				Object('Field' name: 'colpro', mandatory:, set: .custpe)))
			if type isnt false
				vert.Add(.warnStatic())
			}
		else
			vert.Add(Object('Pair' #(Static 'Name') #(Static name: 'colpro')))

		if editType isnt false
			vert.Add(Object('Pair' #(Static 'Type') Object('ChooseFieldType'
				filterBy: type, mandatory:, name: 'ctllbl')))
		else
			vert.Add(Object('Pair' #(Static 'Type') #(Static name: 'ctllbl')))

		vert.Add('Skip')
		vert.Add(#(GroupBox 'Field Properties'
			(Vert name: 'peditorslot' )
			xmin: 400
			ymin: 200
			xstretch: 1
			ystretch: 1
		))
		ob.Add(vert)
		return ob
		}

	Observer(member)
		{
		if member isnt 'ctllbl' or .peditorslot is false
			return
		.peditorslot.Remove(0)
		type = .Data.GetField('ctllbl')
		if false is t = .findType(type)
			return
		try
			{
			.control = t.Member?('control') ? t.control[0] : Datadict(t.base).Control[0]
			peditor = 'CustomizableFieldDialogPropertiesEditor_' $ .control
			if .originalPeditor?(peditor)
				peditor = .originalPeditor.Copy()

			.peditorslot.Append(Object(peditor, typePlugin: t))
			.peditor = .FindControl('peditor')
			if .fieldobject isnt false
				.peditor.Set(.fieldobject.flddef)
			}
		}

	originalPeditor?(peditor)
		{
		if .originalPeditor is false
			return false
		return peditor is 'CustomizableFieldDialogPropertiesEditor_' $ .ctlnme $ 'Control'
		}

	findType(name)
		{
		return .types.GetDefault(.types.FindIf({|c| c.name is name }), false)
		}

	OK()
		{
		if .notChanged?()
			{
			.On_Cancel()
			return false
			}
		data = .Data.Get()
		if .valid?(data)
			{
			if .hasConversionFunction?()
				.handleConversion(data)
			data.colnme = .colnme
			if false isnt ctllbl = .FindControl('ctllbl')
				data.ctllbl = ctllbl.Get()
			else
				data.ctllbl = .type
			data.colpro = .FindControl('colpro').Get()
			if false isnt t = .findType(data.ctllbl)
				data.fldbse = 'Field_' $ t.base
			return data
			}
		else
			Beep()
		return false
		}

	notChanged?()
		{
		if .Data.Valid(forceCheck:) isnt true
			return false
		peditor = .FindControl('peditor')
		if peditor isnt false and peditor.Valid?() isnt true
			return false
		return not .Data.Dirty?() and (peditor is false or .originalData is peditor.Get())
		}

	// when .originalPeditor is false, we are not editing a field
	// when .control is false, we have not selected a new Field Type, not converting
	// .control should be the field type we are converting to
	hasConversionFunction?()
		{
		return .originalPeditor isnt false and .control isnt false and
			.originalPeditor.Member?('ConvertFieldType_' $ .control)
		}

	handleConversion(data)
		{
		if false is requiredData = (.originalPeditor['ConvertFieldType_' $ .control])(
			.fieldobject, .originalData)
			return false

		if false is data.GetDefault('options', false)
			data.options = Object(control: Object(), format: Object())

		for mem in requiredData.control.Members()
			data.options.control.Add(requiredData.control[mem], at: mem)

		for mem in requiredData.format.Members()
			data.options.format.Add(requiredData.format[mem], at: mem)
		}

	valid?(data)
		{
		data.options = false
		pe = .FindControl('peditor')
		if  false isnt pe
			if false is pe.Valid?()
				return false
			else
				data.options = pe.Get()

		// calls function Custom_PromptInUse that checks if the prompt has been
		// used already. This function wipes out the prompt field in the data
		// which we didn't want so used fake field "custom_xxx", only way to
		// tell if prompt was invalid is to check if fake field was cleared.
		field = "custom_xxx"
		ob = Object(custom_xxx: field)
		data[field] = "VALID"
		Custom_PromptInUse(field, ob, data.colpro, .Window.Hwnd, data, exclude_custom?:)
		if data[field] is ''
			return false

		if "" isnt prompt_valid = CustomPromptValid(data.colpro, description: "Name")
			{
			.AlertError(.title, prompt_valid)
			return false
			}

		return .Data.Valid() is true
		}
	}
