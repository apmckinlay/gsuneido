// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	operators: #(
		'add				+',
		'subtract			-',
		'divide 			/',
		'multiply			*',
		'greater than		>',
		'greater than or equal to		>=',
		'less than			<',
		'less than or equal to			<=',
		'equal to			is',
		'not equal to		isnt',
		'and				and',
		'or					or',
		'not				not',
		'concatenate text	$',
		'group				( )'
		)
	formulaBtns()
		{
		formulaReserved = GetContributions('FormulaReserved')
		funcButtons = formulaReserved.Members().Filter({
			false isnt formulaReserved[it].GetDefault(#func, false) }).
			Sort!({ |x,y| formulaReserved[x].func < formulaReserved[y].func})
		return Object('Pair' #(Skip name: 'add_field_skip')
			Object('Horz'
				#(MenuButton 'Add'
					#('Field', 'Text', 'Number',  'Date', 'Quantity', 'Rate'))
				#(Skip small:)
				Object('MenuButton' 'Add an Operator', .operators)
				#(Skip small:)
				Object(#MenuButton, 'Add a Function', funcButtons))
			)
		}
	reserved()
		{
		formulaReserved = GetContributions('FormulaReserved')
		return formulaReserved.Members().Append(.operators.Map({ it.AfterLast('\t') })).
			Append(Object('true', 'false')).Sort!({ |x, y| x.Size() > y.Size() })
		}

	AddFormulaButtons()
		{
		return .formulaBtns()
		}

	Add_a_Field(ctrl, selectFields, hwnd)
		{
		.askReplaceSection(ctrl)
			{
			Reporter.ChooseField(selectFields.Prompts(), hwnd)
			}
		}
	Add_a_Text(ctrl)
		{
		.askReplaceSection(ctrl)
			{
			text = Ask('Add Text', 'Customize Fields', ctrl: #(Editor height: 3))
			text is false ? false : Display(text.Replace('\n', '\r\n'))
			}
		}
	Add_a_Number(ctrl)
		{
		.askReplaceSection(ctrl)
			{
			number = Ask('Add a Number', 'Customize Fields', ctrl: #(Number mask: false))
			number is false ? false : String(number)
			}
		}
	Add_a_Date(ctrl, selectFields)
		{
		.askReplaceSection(ctrl)
			{
			FormulaAddFunction('DATE', selectFields)
			}
		}
	Add_a_Quantity(ctrl, selectFields)
		{
		.askReplaceSection(ctrl)
			{
			FormulaAddFunction('QUANTITY', selectFields)
			}
		}
	Add_a_Rate(ctrl, selectFields)
		{
		.askReplaceSection(ctrl)
			{
			FormulaAddFunction('RATE', selectFields)
			}
		}
	Add_a_Function(ctrl, fnName,selectFields)
		{
		.askReplaceSection(ctrl)
			{
			FormulaAddFunction(fnName, selectFields)
			}
		}

	Add_an_Operator(ctrl, option)
		{
		if false is ctrl
			return
		operator = ' ' $ option.AfterLast('\t') $ ' '
		ctrl.SetFocus()
		ctrl.ReplaceSel(operator)
		}

	askReplaceSection(formula, block)
		{
		if false is formula
			return

		formula.SetFocus()
		if false isnt result = block()
			{
			formula.ReplaceSel(result)
			formula.SetFocus()
			}
		}

	ConstructFormula(sf, prompt, type, formula, field = false, quiet = false,
		skipReturnTypeCheck? = false, checkOnly = false)
		{
		if field is false
			field = sf.PromptToField(prompt)
		cf = CustomFieldTypes(reporter:)
		baseType = cf[cf.FindIf({ it.name is type })].base
		ddName = checkOnly is true
			? 'Field_' $ field
			: ReporterModel.MakeDD(field, prompt, baseType)
		formulaOb = CustomizeField.TranslateFormula(sf, formula, field, true, :quiet,
			:skipReturnTypeCheck?)
		formulaOb.ddName = ddName
		return formulaOb
		}

	HighlightSelection(formula, selectFields)
		{
		sel = formula.GetSel()
		if sel[0] isnt sel[1]
			return

		fields = selectFields.ScanFields(formula.Get().Replace('\n', '\r\n'), .reserved())
		for field in fields
			if sel[0] > field.pos and sel[0] < field.end
				formula.SetSel(field.pos, field.end)
		}
	}