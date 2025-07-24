// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	TranslateFormula(selectFields, formula, field, useParams = false, quiet = false,
		skipReturnTypeCheck? = false)
		{
		if formula.Blank?()
			return Object(fields: '', formulaCode: '')
		code = ''
		params = ''

		selectFields.ScanFormula(formula,
			{|f| code $= f },
			{|s| code $= s })

		_useParams = useParams
		_fields = Object()

		try
			code = .interpret(code, field, skipReturnTypeCheck?)
		catch (err)
			{
			if .isFormulaError?(err)
				return Object(fields: '', formulaCode: Object(:err))
			if not quiet
				SuneidoLog('ERROR: (CAUGHT) ' $ err, calls:,
					caughtMsg: 'unexpected issue converting formula. ' $
						'reverting to old format.')
			return Object(fields: '', formulaCode: false)
			}

		isFormat = ''
		if useParams is true
			{
			params = _fields.Join(', ')
			isFormat = ', isFormat:'
			}

		formulaCode = Join('\r\n', 'function(' $ params $ ')'
			'\t' $ '{'
			'\t\t' $ 'return FormulaReturn({'$ code $ '}, ' $
				Display(field) $ isFormat $ ')'
			'\t' $ '}')
		return Object(fields: _fields.Join(','), :formulaCode)
		}

	isFormulaError?(err)
		{
		return err.Prefix?("Formula:")
		}

	interpret(code, field, skipReturnTypeCheck?)
		{
		// Tdop throws "ununitialized member if it finds a symbol that
		// does not belong to the new Formula type. Assume it's an old formula
		try
			res = Tdop(code, type: 'expression', symbols: FormulaSymbols())
		catch(err, "Unexpected Symbol")
			throw "Formula: " $
				err.Replace('Unexpected Symbol', 'Invalid FormulaSymbol:')
		s = ""

		TdopTraverse2(res,
			{ |node| s $= .preVisit(node) },
			{ |node| s = .postVisit(node, s) })

		if not skipReturnTypeCheck?
			.validateReturnType(res, field)

		return s.RemoveSuffix(',')
		}

	preVisit(node)
		{
		if node.Position is -1
			return ''

		for fn in Object(.handleCall, .handleUnaryOp, .handleBinaryOp,
			.handleNumberOrString, .handleBoolean, .handleIdentifier, .handleRvalue)
			if false isnt res = fn(node)
				return (node.GetDefault(#block, false) ? '{' : '') $ res
		return ''
		}

	handleCall(node)
		{
		if not node.Match(TDOPTOKEN.CALL)
			return false

		// if IDENTIFIER is not found in FormulaReserved
		// it does not belong to the new Formula type. Assume it's an old formula
		call = node[0].ToWrite()
		if not node[0].Match(TDOPTOKEN.IDENTIFIER) or
			not GetContributions('FormulaReserved').Members().Has?(call)
			throw "Formula: Invalid Function " $ call

		if call is 'IF'
			{
			args = node[2]
			args.Children.Each()
				{ |arg|
				arg[0][0].block = true
				}
			}

		return 'Formula' $ call.Capitalize() $ '('
		}

	handleUnaryOp(node)
		{
		if not node.Match(TDOPTOKEN.UNARYOP)
			return false
		token = node[0].Match(TDOPTOKEN.SUB) ? 'NEG' : node[0].Token
		return 'Formula' $ token.Capitalize() $ '('
		}

	handleBinaryOp(node)
		{
		if node.Match(TDOPTOKEN.BINARYOP)
			return 'Formula' $ node[1].Token.Capitalize() $ '('
		return false
		}

	handleNumberOrString(node)
		{
		if node.Match(TDOPTOKEN.NUMBER) or node.Match(TDOPTOKEN.STRING)
			{
			node.formulaTypes = Object(node.Token)
			return 'Object(type: ' $ Display(node.Token) $
				', value: ' $ node.ToWrite() $ '),'
			}
		return false
		}

	handleBoolean(node)
		{
		if node.Match(TDOPTOKEN.TRUE) or node.Match(TDOPTOKEN.FALSE)
			{
			node.formulaTypes = Object(FORMULATYPE.BOOLEAN)
			return 'Object(type: "BOOLEAN", value: ' $ node.ToWrite() $ '),'
			}
		return false
		}

	handleIdentifier(node, _fields, _useParams)
		{
		if not node.Match(TDOPTOKEN.IDENTIFIER)
			return false

		if GetContributions('FormulaReserved').Members().Has?(f = node.ToWrite())
			return ''

		if node.Base?(TdopSymbolReserved)
			throw "Formula: Cannot use the reserved keyword " $ node.Value

		fields.AddUnique(f)
		type = .type(f)
		node.formulaTypes = Object(type)

		value = useParams is true ? f : '.' $ f
		return 'Object(type: ' $ Display(type) $ ', value: ' $ value $ '),'
		}

	handleRvalue(node)
		{
		if node.Match(TDOPTOKEN.RVALUE)
			return '('
		return false
		}

	type(fieldName)
		{
		dd = Datadict(fieldName)
		control = GetControlClass.FromField(fieldName)
		if control.Base?(UOMControl)
			return dd.Control.GetDefault(#div, '') is '/'
				? FORMULATYPE.UOM_RATE
				: FORMULATYPE.UOM
		if dd.Base?(Field_number)
			return FORMULATYPE.NUMBER
		if dd.Base?(Field_date)
			return FORMULATYPE.DATE
		if dd.Base?(Field_boolean)
			return FORMULATYPE.BOOLEAN
		return FORMULATYPE.STRING
		}

	postVisit(node, s)
		{
		if node.Position is -1
			return s

		if not (node.Match(TDOPTOKEN.BINARYOP) or node.Match(TDOPTOKEN.RVALUE) or
			node.Match(TDOPTOKEN.UNARYOP) or node.Match(TDOPTOKEN.CALL))
			return node.GetDefault(#block, false) ? s.RemoveSuffix(',') $ '},' : s

		.setNodeTypes(node)
		return s.RemoveSuffix(',') $ ')' $ (node.GetDefault(#block, false) ? '},' : ',')
		}

	setNodeTypes(node)
		{
		if node.Match(TDOPTOKEN.RVALUE)
			node.formulaTypes = node[1].formulaTypes
		else if node.Match(TDOPTOKEN.BINARYOP)
			node.formulaTypes = .validateBinaryOp(node)
		else if node.Match(TDOPTOKEN.UNARYOP)
			node.formulaTypes = .validateUnaryOp(node)
		else if node.Match(TDOPTOKEN.CALL)
			node.formulaTypes = .validateCall(node)
		}

	validateBinaryOp(node)
		{
		fn = Global('Formula' $ node[1].Token.Capitalize())
		op1Types = node[0].formulaTypes
		op2Types = node[2].formulaTypes
		resTypes = Object()
		for op1Type in op1Types
			for op2Type in op2Types
				resTypes.AddUnique(fn(.generateTestValue(op1Type),
					.generateTestValue(op2Type)).type)
		return resTypes
		}

	validateUnaryOp(node)
		{
		token = node[0].Match(TDOPTOKEN.SUB) ? 'NEG' : node[0].Token
		fn = Global('Formula' $ token.Capitalize())
		return node[1].formulaTypes.Map({ fn(.generateTestValue(it)).type }).
			Sort!().Unique!()
		}

	validateCall(node)
		{
		args = node[2].Children.Map({ it[0][0].formulaTypes })
		fn = Global('Formula' $ node[0].ToWrite().Capitalize())
		return fn.Validate(@args)
		}

	FormulaTestUnit: 'formula_test_uom'
	generateTestValue(type)
		{
		value = ''
		switch (type)
			{
		case FORMULATYPE.BOOLEAN:
			value = true
		case FORMULATYPE.DATE:
			value = Date()
		case FORMULATYPE.NUMBER:
			value = 1
		case FORMULATYPE.STRING:
			value = ''
		case FORMULATYPE.UOM, FORMULATYPE.UOM_RATE:
			value = '1 ' $ .FormulaTestUnit
			}
		return FormulaBase.GenerateElement(:type, :value)
		}

	validateReturnType(resNode, field)
		{
		dd = Datadict(field)
		for type in resNode.formulaTypes
			{
			value = .generateTestValue(type)
			try
				FormulaReturn.ProcessValue(value.value, dd)
			catch (e)
				if false isnt match = e.Match('Invalid <(.+?)> value')
					throw 'Formula: ' $ PromptOrHeading(field) $ ' cannot assign <' $
						type.Capitalize() $ '> to <' $
						e[match[1][0]::match[1][1]] $ '>'
			}
		}

	ValidateCode(code)
		{
		if code is false
			return "Invalid or missing operators in formula."

		if Object?(code)
			return code.err

		if '' isnt msg = .ExtraCheck(code)
			return msg

		if code is ''
			return ''

		if CheckCode(code) is false
			return "Invalid Formula."
		return ''
		}

	ExtraCheck(code)
		{
		Plugins().ForeachContribution("Formulas", 'checks')
			{
			if "" isnt msg = (it.Fn)(code)
				return msg
			}
		return ""
		}

	SetFormulas(customKey, record, protectField)
		{
		formulas = .GetCustomFieldFormulas(customKey)
		for field in formulas.Members()
			{
			formula = formulas[field]
			if not .protectedByScreen?(record, protectField, field)
				{
				record.AttachRule(field, formula.rule)
				record.AttachRule(field $ '__protect', formula.protectRule)
				if formula.fields isnt ''
					record.SetDeps(field, formula.fields)
				}
			}
		}

	GetCustomFieldFormulas(customKey)
		{
		if customKey is false
			return #()

		formulas = Customizable.CacheByKey(customKey, "CustomizedFieldFormulas")
			{ |name|
			if Customizable.NotCustomizableScreen?()
				return #()

			ob = Object()
			QueryApply('customizable_fields where custfield_name is ' $ Display(name) $
				' and custfield_formula isnt ""')
				{
				protectRule = 'function()
					{
					return #(' $ it.custfield_formula_fields $ ').Any?({this[it] isnt ""})
					}'
				ob[it.custfield_field] = Object(
					rule: 			it.custfield_formula_code.Compile(),
					protectRule: 	protectRule.Compile(),
					fields: 		it.custfield_formula_fields)
				}
			ob
			}
		return formulas
		}

	HasCustomFieldFormula?(customKey)
		{
		return not .GetCustomFieldFormulas(customKey).Empty?()
		}

	protectedByScreen?(rec, protectField, field)
		{
		if protectField is false
			return false
		protect = rec[protectField]
		if not String?(protect) and not Boolean?(protect) and not Object?(protect)
			{
			.logError('invalid return type from protect rule')
			return true
			}

		if protect is true or (String?(protect) and protect isnt '')
			return true

		if Object?(protect)
			{
			allbut? = protect.GetDefault(0, false) is 'allbut'
			if ((allbut? and not protect.Member?(field)) or
				(not allbut? and protect.Member?(field)))
				return true
			}

		return false
		}

	logError(str, prefix = "ERROR: ")
		{
		SuneidoLog(prefix $ str)
		}

	IsCustomized?(accessCustomKey)
		{
		if accessCustomKey isnt false and TableExists?('customizable_fields')
			{
			if not QueryEmpty?('customizable_fields
				where custfield_name is "' $ accessCustomKey $ '"')
				return true
			}
		return false
		}
	ConsideredEmpty?(field, data)
		{
		if data[field] is ""
			return true
		ctrl = GetControlClass.FromField(field)
		return ctrl.Method?('ConsideredEmpty?')
			? ctrl.ConsideredEmpty?(data[field])
			: false
		}

	CheckCustomFields(customFields, recordCtrl, protectField)
		{
		data = recordCtrl.Get()
		if customFields is false
			return ''

		// check invalid manadatory custom fields, because they are not constructed
		errCustom = Object()
		mandatories = customFields.Members().Filter({|x|
			customFields[x].GetDefault('mandatory', false) is true })
		for field in mandatories
			if .ConsideredEmpty?(field, data) and
				not .customized_protected?(field, data, recordCtrl, protectField)
				errCustom.Add(PromptOrHeading(field))
		return Opt('Required: ', errCustom.Join(', '))
		}

	customized_protected?(field, data, recordCtrl, protectField)
		{
		ctrl = recordCtrl.GetControl(field)
		if ctrl isnt false and ctrl.Method?('GetReadOnly')
			return ctrl.GetReadOnly()
		return FieldProtected?(field, data, protectField)
		}
	}