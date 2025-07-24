// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
JsTranslate // to get output stuff
	{
	// TODO class literals

	CallClass(ast, discard = false, parens = false)
		{
		if not Same?(ast, val = AstFoldExpr(ast))
			{
			.Value(val)
			return
			}
		this[ast.type](ast, :parens, :discard)
		}

	Ident(ast)
		{
		name = ast.name
		switch
			{
		case .selfref?(name):
			.selfref()
		case .local?(name):
			name = .ConvertQMark(name)
			_var(name)
			.Print(name)
		case .global?(name):
			.Print('su.global("' $ name $ '")')
		case .dynamic?(name):
			.Print('su.dynget("' $ name $ '")')
		//TODO _Name library overload backref
			}
		}
	selfref(_inblock = false)
		{
		if inblock
			{
			_blockThis[0] = true
			.Print('blockThis')
			}
		else
			.Print('this') // assuming JS 'this' is used
		}

	selfref?(name)
		{
		return name is 'this'
		}
	// Note: these do not do a full check, they assume a valid identifier
	local?(name)
		{
		return name[0].Lower?()
		}
	global?(name)
		{
		return name[0].Upper?()
		}
	dynamic?(name)
		{
		return name[0] is '_' and name[1].Lower?()
		}
	backref?(name)
		{
		return name[0] is '_' and name[1].Upper?()
		}

	Unary(ast, parens, discard)
		{
		op = ast.op
		if discard and op in ('PostInc', 'PostDec')
			op = op.RemovePrefix('Post') // optimization, avoid temp var
		switch ast.op
			{
		case 'Add':
			.unary(ast, 'uadd')
		case 'Sub':
			.unary(ast, 'usub')
		case 'BitNot':
			.unary(ast, 'bitnot')
		case 'PostInc', 'PostDec':
			.postIncDec(ast, parens)
		case 'Inc', 'Dec':
			.preIncDec(ast, parens)
		case 'Not':
			.boolNot(ast, parens)
		case 'LParen':
			JsTranslateExpression(ast.expr, parens)
			}
		}

	unary(ast, op)
		{
		.Print('su.' $ op $ '(')
		JsTranslateExpression(ast.expr)
		.Print(')')
		}

	boolNot(ast, parens)
		{
		.doWithParens(parens)
			{
			.Print('! ')
			.Bool(ast.expr, parens:)
			}
		}

	Nary(ast)
		{
		switch ast.op
			{
		case 'And':
			.and_or(ast, ' && ')
		case 'Or':
			.and_or(ast, ' || ')
		case 'Mul':
			.mul(ast, ast[ast.size - 1], ast.size - 1)
		case 'Add':
			.add(ast, ast[ast.size - 1], ast.size - 1)
		default:
			.nary(ast, ast.size - 1)
			}
		}

	and_or(ast, op)
		{
		.Bool(ast[0], parens:)
		for (i = 1; i < ast.size; ++i)
			{
			.Print(op)
			.Bool(ast[i], parens:)
			}
		}

	mul(ast, op2, i)
		{
		if i is 0
			JsTranslateExpression(op2, parens:)
		else
			{
			op = 'mul'
			op1 = ast[i - 1]
			if op2.type is 'Unary' and op2.op is 'Div'
				{
				op = 'div'
				op2 = op2.expr
				}
			.Print('su.' $ op $ '(')
			.mul(ast, op1, i - 1)
			.Print(', ')
			JsTranslateExpression(op2, parens:)
			.Print(')')
			}
		}

	add(ast, op2, i)
		{
		if i is 0
			JsTranslateExpression(op2, parens:)
		else
			{
			op = 'add'
			op1 = ast[i - 1]
			if op2.type is 'Unary' and op2.op is 'Sub'
				{
				op = 'sub'
				op2 = op2.expr
				}
			.Print('su.' $ op $ '(')
			.add(ast, op1, i - 1)
			.Print(', ')
			JsTranslateExpression(op2, parens:)
			.Print(')')
			}
		}

	nary(ast, i)
		{
		if i is 0
			JsTranslateExpression(ast[0], parens:)
		else
			{
			.Print(.xlat(ast))
			.nary(ast, i - 1)
			.Print(', ')
			JsTranslateExpression(ast[i], parens:)
			.Print(')')
			}
		}

	Trinary(ast)
		{
		.Bool(ast.cond, parens:)
		.Print(' ? ')
		JsTranslateExpression(ast.t, parens:)
		.Print(' : ')
		JsTranslateExpression(ast.f, parens:)
		}
	Bool(ast, parens = false)
		{
		boolResult = .boolResult(ast)
		.doWithParens(parens and boolResult)
			{
			if not boolResult
				.Print('su.toBool(')
			JsTranslateExpression(ast)
			if not boolResult
				.Print(')')
			}
		}
	boolResult(ast)
		{
		if ast.type is 'Unary'
			{
			if ast.op is 'LParen'
				return .boolResult(ast.expr)
			else if ast.op is 'Not'
				return true
			else
				return false
			}
		return ast.type is 'Nary' and ast.op in ('And', 'Or') or
			ast.type is 'Binary' and ast.op in ('Is', 'Isnt', 'Lt', 'Lte', 'Gt', 'Gte',
				'Match', 'MatchNot')
		}

	preIncDec(ast, parens)
		{
		post = .lvalue(ast.expr, parens)
		.Print(.xlat(ast))
		JsTranslateExpression(ast.expr)
		.Print(')' $ post)
		}
	postIncDec(ast, parens)
		{
		.doWithParens(parens)
			{
			post = .lvalue(ast.expr, parens:)
			tmp = _nextTmp()
			.Print(.xlat(ast) $ tmp $ ' = ')
			JsTranslateExpression(ast.expr)
			.Print(post $ '), ' $ tmp)
			}
		}

	Binary(ast, parens)
		{
		if ast.op is 'Eq'
			.Eq(ast, parens)
		else
			{
			post = ast.op.Suffix?('Eq')
				? .lvalue(ast.lhs, parens)
				: ''
			.Print(.xlat(ast))
			JsTranslateExpression(ast.lhs, parens:)
			.Print(', ')
			JsTranslateExpression(ast.rhs, parens:)
			.Print(')')
			if post isnt ''
				.Print(post)
			}
		}

	Eq(ast, parens)
		{
		post = .lvalue(ast.lhs, parens)
		JsTranslateExpression(ast.rhs, parens:)
		.Print(post)
		}

	lvalue(ast, parens) // returns what should be after the rvalue, '' or ')'
		{
		switch ast.type
			{
		case 'Ident':
			return .lvalue_id(parens,  ast)
		case 'Mem':
			.Print('su.put(')
			JsTranslateExpression(ast.expr)
			.Print(', ')
			if ast.mem.type is 'Constant' and String?(ast.mem.value)
				.Print('"' $ (JsTranslateClass.Privatize(ast.expr, ast.mem.value)) $ '"')
			else
				JsTranslateExpression(ast.mem)
			.Print(', ')
			return ')'
			}
		}
	lvalue_id(parens,  ast)
		{
		Assert(ast.type is 'Ident')
		name = ast.name
		switch
			{
		case .local?(name):
			name = .ConvertQMark(name)
			_var(name)
			if (parens) .Print('(')
			.Print(name)
			.Print(' = ')
			return parens ? ')' : ''
		case .dynamic?(name):
			_setsDynamic[0] = true
			.Print('su.dynset("' $ name $ '", ')
			return ')'
			}
		}

	Mem(ast)
		{
		.Print('su.get(')
		JsTranslateExpression(ast.expr)
		.Print(', ')
		if ast.mem.type is 'Constant' and String?(ast.mem.value)
			.Value(JsTranslateClass.Privatize(ast.expr, ast.mem.value))
		else
			JsTranslateExpression(ast.mem, parens:)
		.Print(')')
		}

	strLengthLimit: '536870888' // by JavaScript engines
	RangeTo(ast)
		{
		.Print('su.rangeto(')
		JsTranslateExpression(ast.expr)
		.Print(', ')
		if false is ast.from
			.Print('0')
		else
			JsTranslateExpression(ast.from, parens:)
		.Print(', ')
		if false is ast.to
			.Print(.strLengthLimit)
		else
			JsTranslateExpression(ast.to, parens:)
		.Print(')')
		}

	RangeLen(ast)
		{
		.Print('su.rangelen(')
		JsTranslateExpression(ast.expr)
		.Print(', ')
		if false is ast.from
			.Print('0')
		else
			JsTranslateExpression(ast.from, parens:)
		.Print(', ')
		if false is ast.len
			.Print(.strLengthLimit)
		else
			JsTranslateExpression(ast.len, parens:)
		.Print(')')
		}

	In(ast, parens)
		{
		.doWithParens(parens)
			{
			.Print('-1 != [')
			sep = ''
			for (i = 0; i < ast.size; i++)
				{
				.Print(sep)
				sep = ', '
				JsTranslateExpression(ast[i])
				}
			.Print('].indexOf(') // or does this need to be .some(su.is) ???
			JsTranslateExpression(ast.expr)
			.Print(')')
			}
		}

	Call(ast)
		{
		.call(ast)
			{|callType|
			if ast.func.type is 'Ident' and ast.func.name is 'super'
				{
				.Print('su.invoke' $ callType $ 'BySuper($super, ')
				.Print('"New", this')
				}
			else if ast.func.type is 'Mem'
				{
				m = ast.func
				if m.expr.type is 'Ident' and m.expr.name is 'super'
					{
					.Print('su.invoke' $ callType $ 'BySuper($super, ')
					.Print('"' $ m.mem.value $ '", this')
					}
				else if m.mem.type is 'Constant' and m.mem.value is '*new*'
					{
					.Print('su.instantiate' $ callType $ '(')
					JsTranslateExpression(m.expr)
					}
				else
					{
					.Print('su.invoke' $ callType $ '(')
					JsTranslateExpression(m.expr, parens:)
					.Print(', ')
					if m.mem.type is 'Constant' and String?(m.mem.value)
						.Print('"' $ JsTranslateClass.Privatize(m.expr, m.mem.value) $
							'"')
					else
						JsTranslateExpression(m.mem, parens:)
					}
				}
			else
				{
				.Print('su.call' $ callType $ '(')
				JsTranslateExpression(ast.func)
				}
			}
		}
	call(ast, block)
		{
		callType = ''
		if ast.size is 1 and ast[0].name in ('@', '@+1')
			callType = 'At'
		else if .hasNamed(ast)
			callType = 'Named'
		else
			callType = ''
		block(callType)
		.args(ast, callType)
		.Print(')')
		}
	hasNamed(ast)
		{
		for (i = 0; i < ast.size; i++)
			if ast[i].name isnt false
				return true
		return false
		}
	args(ast, callType)
		{
		if callType is 'At'
			{
			.Print(', ')
			JsTranslateExpression(ast[0].expr, parens:)
			.Print(', ' $ (ast[0].name is '@' ? '0' : ast[0].name[2..]))
			return
			}
		n = ast.size
		if callType is 'Named'
			.handleNamed(ast)
		for (i = 0; i < n; ++i)
			.unnamed(ast[i])
		}
	handleNamed(ast)
		{
		numberNamed = Object()
		stringNamed = Object()
		n = ast.size
		for (i = 0; i < n; ++i)
			{
			cur = ast[i]
			if cur.name is false
				continue

			if Number?(cur.name)
				numberNamed[cur.name] = cur.expr
			else
				stringNamed[cur.name] = cur.expr
			}

		if numberNamed.Empty?()
			.namedWithoutNumber(stringNamed)
		else
			.namedWithNumber(numberNamed, stringNamed)
		}
	namedWithoutNumber(stringNamed)
		{
		.Print(', {')
		for name, node in stringNamed
			{
			.Print('"' $ name $ '": ')
			JsTranslateExpression(node, parens:)
			.Print(', ')
			}
		.Print('}')
		}
	namedWithNumber(numberNamed, stringNamed)
		{
		tmp = _nextTmp()
		.Print(', (' $ tmp $ ' = new Map([')
		for name, node in numberNamed
			{
			.Print('[' $ name $ ', ')
			JsTranslateExpression(node, parens:)
			.Print('], ')
			}
		.Print(']), ')
		for name, node in stringNamed
			{
			.Print(tmp $ '.' $ name $ ' = ')
			JsTranslateExpression(node, parens:)
			.Print(', ')
			}
		.Print(tmp $ ')')
		}
	unnamed(arg)
		{
		if arg.name isnt false
			return false
		.Print(', ')
		JsTranslateExpression(arg.expr, parens:)
		return true
		}
	xlat(ast)
		{
		return 'su.' $ ast.op.Lower().Replace('^post|eq$') $ '('
		}
	Constant(ast)
		{
		this[ast.value.type](ast.value)
		}
	Object(ast)
		{
		.Const(ast)
		}
	Record(ast)
		{
		.Const(ast)
		}
	Function(ast)
		{
		.Const(ast)
		}
	Block(ast)
		{
		JsTranslateFunction(ast)
		}
	Class(ast)
		{
		.Const(ast)
		}

	doWithParens(parens, block)
		{
		if parens
			.Print('(')
		block()
		if parens
			.Print(')')
		}
	}
