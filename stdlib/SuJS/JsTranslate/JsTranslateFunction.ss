// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
JsTranslate // inherit to get output stuff
	{
	CallClass(ast, outerName = false, inMethod = false)
		{
		.NameBegin(outerName, .CompilerNameSuffix[ast.type])
		(new this)(ast, inMethod)
		.NameEnd()
		}
	tmpi: 0
	pre: '.call(this, '
	Call(ast, inMethod = false)
		{
		callableType = ast.type
		Assert(callableType in ('Function', 'Block'))
		_inMethod = .isInMethod(callableType, inMethod)
		// break/continue/return need to know if they are in a block
		_inblock = callableType is 'Block'
		_inLoop = false
		// return from a block (possibly nested) needs to know its function id
		// function needs to know if it had any block that had a return
		fn = [id: .NextNum(), hasBlockReturn: false]
		if callableType is 'Function'
			_fn = fn
		// need to know if there are any assignments to dynamics
		_setsDynamic = [false]
		// need to know if there are any selfref in a block
		_blockThis = [false]
		_const = Object()
		_nextTmp = { '$tmp' $ .tmpi++ }
		.vars = Object(this: false)
		if not _inblock
			_var = {|var| if (not .vars.Member?(var)) .vars[var] = true }
		.Print('(') // start of IIFE
		.func()
		start = .PrintPos()
		params = ast.params
		.Indented()
			{
			if false isnt atParam = .atParam(params)
				{
				.atFuncImp(atParam)
				fwd = .atFuncFwd
				}
			else
				{
				.funcImp(params)
				fwd = .funcFwd
				}
			.Indented()
				{
				before = .PrintPos()
				.Println('su.maxargs(' $ params.size $ ', arguments.length);')
				.body(ast) // body
				after = .PrintPos()
				}
			.Println('};') // opening is in func imp
			.Println('$f.$callableType = "' $
				(inMethod ? 'Method' : callableType).Upper() $ '";')
			.Println('$f.$callableName = "' $ .Get('curName') $ '";')
			if _blockThis[0]
				.Println('$f.$blockThis = this;')
			fwd(:params)
			.printParams(params)
			.Println('return $f;')
			.Print(_blockThis[0] ? '}).call(this)' : '})()') // end of IIFE
			.Indented()
				{
				// must be after const to see results of child functions/blocks
				.addHandlers(fn, before, after)
				}
			.const(start)
			}
		}
	isInMethod(callType, inMethod)
		{
		if inMethod is true
			return true
		inMethod = false
		if callType is 'Block'
			try inMethod = _inMethod
		return inMethod
		}
	func(with = '()')
		{
		.Print('function ')
		if String?(with)
			.Print(with)
		else
			with()
		.Println(' {')
		}

	// @param
	atParam(params)
		{
		return params.size is 1 and params[0].name.Prefix?('@')
			? params[0].name[1..]
			: false
		}
	atFuncImp(atParam)
		{
		atParam = .ConvertQMark(atParam)
		.vars[atParam] = false
		.Print('let $f = ')
		.func('(' $ atParam $ ' = su.mandatory())')
		}
	atFuncFwd()
		{
		.Println("$f.$callAt = $f;")
		.Print("$f.$call = ")
		.func("(...args)")
		.Indented()
			{
			.Println("return $f" $ .pre $ "su.mkObject2(args));")
			}
		.Println("};")
		.Print("$f.$callNamed = ")
		.func("(named, ...args)")
		.Indented()
			{
			.Println("return $f" $ .pre $ "su.mkObject2(args, su.obToMap(named)));")
			}
		.Println("};")
		}

	// normal params
	funcImp(params)
		{
		paramNames = .mapParamName(params)
		list = paramNames.names.Join(', ')
		.Print("let $callNamed = ")
		.func("(named" $ Opt(', ', list) $ ")")
		.Indented()
			{
			.Println("su.maxargs(" $ (paramNames.names.Size() + 1) $
				", arguments.length);")
			if paramNames.normalNames.NotEmpty?()
				{
				name = paramNames.normalNames.Map({ it $ ' = ' $ it }).Join(', ')
				.Println("({ " $ name $ " } = named);")
				}
			for m, v in paramNames.qMarkNames
				{
				.Println('if (named["' $ m $ '"] !== undefined)')
				.Indented()
					{
					.Println(v $ ' = named["' $ m $ '"];')
					}
				}
			.Println(.fix("return $f" $ .pre $ "" $ list $ ");"))
			}
		.Println("};")

		.Print('let $f = ')
		.func({ .params(params) })
		.Indented()
			{
			.handleDotParams(params)
			}
		}
	mapParamName(params)
		{
		names = Object()
		normalNames = Object()
		qMarkNames = Object()
		for (i = 0; i < params.size; i++)
			{
			param = params[i]
			if param.name.Suffix?('?')
				qMarkNames.Add(name = .getParamName(param),
					at: .getParamName(param, keepQMark:))
			else
				normalNames.Add(name = .getParamName(param, keepQMark:))
			names.Add(name)
			}
		return Object(:normalNames, :qMarkNames, :names)
		}
	getParamName(param, keepCapital = false, keepQMark = false)
		{
		name = param.name.Replace('^[._]+')
		if keepQMark is false
			name = .ConvertQMark(name)
		return keepCapital ? name : name.UnCapitalize()
		}
	fix(s)
		{
		return s.Replace(`\(this, \)`, '(this)')
		}
	params(params)
		{
		.Print('(')
		sep = ''
		for (i = 0; i < params.size; i++)
			{
			.Print(sep); sep = ', '
			.param(params[i])
			}
		.Print(')')
		}
	param(p)
		{
		origParamName = .getParamName(p, keepQMark:)
		paramName = .getParamName(p)
		.vars[paramName] = false
		.Print(paramName $ ' = ')
		isDynParam = .isDynParam(p)
		if p.hasdef is false
			{
			if isDynParam
				.Print('su.dynparam("_' $ origParamName $ '")')
			else
				.Print("su.mandatory()")
			}
		else
			{
			c = p.defval
			if isDynParam
				.Print('su.dynparam("_' $ origParamName $ '", ')
			if Type(c) isnt 'AstNode'
				.Value(c)
			else
				.Const(c)
			if isDynParam
				.Print(')')
			}
		}
	isDynParam(param)
		{
		return '_' is (.isDotParam(param) ? param.name[1] : param.name[0])
		}
	isDotParam(param)
		{
		return param.name.Prefix?('.')
		}
	handleDotParams(params)
		{
		for (i = 0; i < params.size; i++)
			{
			param = params[i]
			if not .isDotParam(param)
				continue
			// Note: depends on the class implementation
			paramName = .getParamName(param)
			memberName = .getParamName(param, keepCapital:, keepQMark:)
			if _inMethod
				memberName = JsTranslateClass.PrivatizeName(memberName)
			.Println('su.put(this, ' $ Display(memberName) $ ', ' $ paramName $ ');')
			}
		}
	funcFwd(params/*unused*/)
		{
		.Println("$f.$call = $f;")
		.Println("$f.$callNamed = $callNamed;")
		.Println("$f.$callAt = function (args) {")
		.Indented()
			{
			.Println("return $callNamed" $ .pre $ "su.mapToOb(args.map), ...args.vec);")
			}
		.Println("};")
		}

	const(start)
		{
		if _const.Empty?()
			return
		.At(start)
			{
			.Indented()
				{
				.Println("let $const = [")
				.Indented()
					{
					for c in _const
						{
						JsTranslateConstant(c)
						.Println(',')
						}
					}
				.Println('];')
				}
			}
		}
	body(stmts)
		{
		pos = .PrintPos()
		last = stmts.size - 1
		.superInit(stmts)
		for (i = 0; i <= last; ++i)
			{
			.superChecks(i, stmts[i])
			JsTranslateStatement(stmts[i], i is last)
			}
		// declare variables
		.At(pos)
			{
			s = Opt('var ', Join(', ',
					not _inblock
						? .vars.Copy().Remove(false).Members().Sort!().Join(', ') : '',
					Seq(.tmpi).Map!({ '$tmp' $ it }).Join(', ')),
					';')
			if s isnt ''
				.Println(s)
			}
		}
	superInit(stmts)
		{
		if .isInNewMethod() and not .hasExplicitSuperCall(stmts)
			.Println('su.invokeBySuper($super, "New", this)')
		}

	hasExplicitSuperCall(stmts)
		{
		if stmts.size isnt 0
			return .isSuperNewCall(stmts[0])
		return false
		}

	superChecks(i, stmt)
		{
		if .isSuperNewCall(stmt)
			{
			if not .isInNewMethod()
				throw "super call only allowed in New"
			if i isnt 0
				throw "call to super must come first"
			}
		}

	isSuperNewCall(stmt)
		{
		return stmt.type is 'ExprStmt' and stmt.expr.type is 'Call' and
			stmt.expr.func.type is 'Ident' and stmt.expr.func.name is 'super'
		}

	isInNewMethod()
		{
		return  _inMethod and .Get(#curName).Suffix?(.MethodSeparaor $ 'New')
		}

	addHandlers(fn, before, after)
		{
		if fn.hasBlockReturn or _setsDynamic[0]
			{
			// insert block return handler
			// suffix first to avoid invalidating before position
			.At(after)
				{
				.Print('} ')
				if fn.hasBlockReturn
					// blockReturnHandler will rethrow if not block return
					.Println('catch (_e) { ' $
						'return su.blockReturnHandler(_e, ' $ fn.id $ '); }')
				if _setsDynamic[0]
					.Println('finally { su.dynpop(); }')
				}
			.At(before)
				{
				if _setsDynamic[0]
					.Println('try { su.dynpush();')
				else
					.Println('try {')
				}
			}
		if _blockThis[0]
			.At(before)
				{
				.Println('var blockThis = su.getBlockThis(this);')
				}
		}

	printParams(params)
		{
		ob = Object()
		for (i = 0; i < params.size; i++)
			ob.Add(.printParam(params[i]))
		.Println("$f.$params = '" $ ob.Join(', ') $ "';")
		}
	printParam(p)
		{
		s = .getParamName(p, keepQMark:)
		if p.hasdef
			{
			s $= '='
			defval = p.defval
			// Note: only print simple default values
			// since complex ones (e.g. function) are hard to print
			if Type(defval) isnt 'AstNode'
				s $= String?(defval)
					? '"' $ defval.Escape() $ '"'
					: Display(defval)
			else if not Same?(defval, val = AstFoldExpr(defval))
				s $= String?(val)
					? '"' $ val.Escape() $ '"'
					: Display(val)
			}
		return s
		}
	}
