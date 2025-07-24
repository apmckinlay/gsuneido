// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// values can be any of the basic types like string, number, boolean and date
	// values can also be a compare function
	CallClass(code, expressions, values = #(), findFirst? = false, skipFn? = false)
		{
		try
			codeAst = Type(code) is 'AstNode' ? code : Suneido.Parse(code)
		catch (e)
			{
			if e.Has?('gSuneido does not implement')
				return #()
			return 'Parse error: ' $ e
			}

		if Type(codeAst) isnt 'AstNode'
			return values.Any?({ .compareValue(codeAst, it) })
				? [Object(pos: 0, end: code.Size())]
				: #()

		if String?(searches = .buildSearches(expressions, values))
			return searches

		return .search(codeAst, searches, findFirst?, skipFn?)
		}

	buildSearches(expressions, values)
		{
		if String?(expressions)
			expressions = Object(expressions)

		searches = Object()
		for expr in expressions
			{
			try
				searches.Add(AstParseExpression(expr))
			catch (e)
				return 'Parse Search text ' $ Display(expr) $ ' error: ' $ e
			}

		searches.Append(values)
		return searches
		}

	search(codeAst, searches, findFirst?, skipFn?)
		{
		results = Object()
		.traverse(codeAst)
			{ |node, parents|
			status = 0 // 0 - continue, 1 - stop, 2 - skip children
			if skipFn? isnt false and skipFn?(node, parents) is true
				status = 2
			else
				{
				for search in searches
					{
					if .compare(node, search)
						{
						if false is res = .getPos(node, parents)
							return 'Cannot get pos info'

						results.Add(res)
						if findFirst? is true
							{
							status = 1
							break
							}
						}
					}
				}
			status
			}
		return results
		}

	traverse(target, block)
		{
		stack = Object()
		stack.Add(target)
		parents = Object()

		do
			{
			cur = stack[stack.Size() - 1]
			stack.Delete(stack.Size() - 1)
			if cur is -1
				{
				parents.Delete(parents.Size() - 1)
				continue
				}
			status = block(cur, parents)
			if status is 1
				break
			if status is 0
				{
				parents.Add(cur)
				stack.Add(-1)
				if cur.type is #Function
					.pushFunction(cur, stack)
				else if cur.type is #Param
					.pushParam(cur, stack)
				else
					{
					children = cur.children
					i = 0
					while false isnt child = children[i]
						{
						if Type(child) is 'AstNode'
							stack.Add(child)
						i++
						}
					}
				}
			}
		while stack.Size() isnt 0
		}

	pushFunction(cur, stack)
		{
		for i in .. cur.params.size
			{
			if Type(cur.params[i]) is 'AstNode'
				stack.Add(cur.params[i])
			}
		for i in .. cur.size
			if Type(cur[i]) is 'AstNode'
				stack.Add(cur[i])
		}

	pushParam(cur, stack)
		{
		if cur.hasdef and Type(cur.defval) is 'AstNode'
			stack.Add(cur.defval)
		}

	compare(node1, nodeOrValue)
		{
		if Type(nodeOrValue) isnt #AstNode
			{
			method = 'Compare_' $ node1.type $ '_To_Value'
			if not .Method?(method)
				return false
			return (this[method])(node1, nodeOrValue)
			}

		return .Compare2(node1, nodeOrValue)
		}

	Compare2(node1, node2, noSpecial? = false)
		{
		_noSpecial? = noSpecial?
		if .matchSpecial?(node1, node2)
			return true

		if node1.type isnt node2.type
			return false

		return (this['Compare_' $ node1.type])(node1, node2)
		}

	Compare_Object(node1, node2)
		{
		return .compareArray(node1, node2)
		}

	Compare_Record(node1, node2)
		{
		return .compareArray(node1, node2)
		}

	Compare_Class(node1, node2)
		{
		if node1.base isnt node2.base
			return false

		return .compareArray(node1, node2)
		}

	Compare_Member(node1, node2)
		{
		if node1.named isnt node2.named
			return false
		return (node1.named is false or node1.key is node2.key) and
			.compareValue(node1.value, node2.value)
		}

	Compare_Member_To_Value(node1, value)
		{
		return .compareValue(node1.value, value)
		}

	Compare_Function(node1, node2)
		{
		if .compare(node1.params, node2.params) is false
			return false
		return .compareArray(node1, node2)
		}

	Compare_Params(node1, node2)
		{
		return .compareArray(node1, node2)
		}

	Compare_Param(node1, node2)
		{
		if node1.name isnt node2.name or node1.hasdef isnt node2.hasdef
			return false
		return node1.hasdef is false or .compareValue(node1.defval, node2.defval)
		}

	Compare_Param_To_Value(node1, value)
		{
		return node1.hasdef isnt false and .compareValue(node1.defval, value)
		}

	Compare_Constant(node1, node2)
		{
		return .compareValue(node1.value, node2.value)
		}

	Compare_Constant_To_Value(node1, value)
		{
		return .compareValue(node1.value, value)
		}

	Compare_Ident(node1, node2)
		{
		return node1.name is node2.name
		}

	Compare_Unary(node1, node2)
		{
		return node1.op is node2.op and .compare(node1.expr, node2.expr)
		}

	Compare_Binary(node1, node2)
		{
		return node1.op is node2.op and .compare(node1.lhs, node2.lhs) and
			.compare(node1.rhs, node2.rhs)
		}

	Compare_Nary(node1, node2)
		{
		if node1.op isnt node2.op
			return false

		for (i = 0; i <= node1.size - node2.size; i++)
			{
			match = true
			for (j = 0; j < node2.size; j++)
				{
				if .compare(node1[i + j], node2[j]) is false
					{
					match = false
					break
					}
				}
			if match is true
				return true
			}
		return false
		}

	Compare_Mem(node1, node2)
		{
		if not .compare(node1.expr, node2.expr)
			return false

		// wildcard
		if node2.mem.type is 'Constant' and String?(node2.mem.value) and
			node2.mem.value =~ '^[[:alpha:]]$'
			return true

		return .compare(node1.mem, node2.mem)
		}

	Compare_Trinary(node1, node2)
		{
		return .compare(node1.cond, node2.cond) and
			.compare(node1.t, node2.t) and
			.compare(node1.f, node2.f)
		}

	Compare_RangeTo(node1, node2)
		{
		return node1.from is node2.from and node1.to is node2.to and
			.compare(node1.expr, node2.expr)
		}

	Compare_RangeLen(node1, node2)
		{
		return node1.from is node2.from and node1.len is node2.len and
			.compare(node1.expr, node2.expr)
		}

	Compare_In(node1, node2)
		{
		return .compare(node1.expr, node2.expr) and
			.compareArray(node1, node2)
		}

	Compare_Call(node1, node2)
		{
		return .compare(node1.func, node2.func) and
			.compareArgs(node1, node2)
		}

	compareArgs(node1, node2)
		{
		if .matchAnyArgs?(node2)
			return true

		args = Object()
		sharedArgs = Object()
		args2 = Object()
		sharedArgs2 = Object()
		names = Object()
		defs = Object()

		for i in .. node2.size
			{
			names[i] = node2[i].name isnt false
				? node2[i].name
				: 'ARG_' $ i
			args2[names[i]] = node2[i].expr
			}

		for i in .. node1.size
			{
			name = node1[i].name
			if name isnt false
				{
				if not args2.Member?(name)
					{
					defs.Add(node1[i].expr)
					args[name] = node1[i].expr
					}
				else
					{
					sharedArgs[name] = node1[i].expr
					sharedArgs2[name] = args2[name]
					args2.Delete(name)
					names.Remove(name)
					}
				}
			else
				args[i] = node1[i].expr
			}

		try
			.nameArgs(args, names, defs)
		catch (unused, "missing argument:")
			return false

		return sharedArgs2.Members().Every?({
			.compare(sharedArgs[it], sharedArgs2[it]) }) and
			args2.Members().Every?({
				.compare(args[it], args2[it]) })
		}

	matchAnyArgs?(node)
		{
		return node.size is 1 and node[0].name is '@' and node[0].expr.type is #Ident and
			node[0].expr.name is 'ANY_ARGS'
		}

	nameArgs(args, names, defs = #())
		{
		nlist = listCount = args.Size(list:)
		idefs = 0
		for (i = names.Size() - 1; i >= 0; --i)
			{
			if not args.Member?(names[i])
				if i < nlist
					{
					args[names[i]] = args[i]
					listCount--
					}
				else if names[i].Prefix?('ARG_') and idefs < defs.Size()
					args[names[i]] = defs[idefs++]
				else
					throw "missing argument: " $ names[i]
			args.Erase(i)
			}
		if listCount >  0
			throw "missing argument:"
		return args
		}

	Compare_Argument(node1, node2)
		{
		return node1.name is node2.name and .compare(node1.expr, node2.expr)
		}

	Compare_Block(node1, node2)
		{
		return .compareArray(node1, node2)
		}

	Compare_ExprStmt(node1, node2)
		{
		return .compare(node1.expr, node2.expr)
		}

	Compare_If(node1, node2)
		{
		return .compare(node1.cond, node2.cond) and
			.compare(node1.t, node2.t) and
			.compareFalseOrNode(node1.f, node2.f)
		}

	Compare_Switch(node1, node2)
		{
		return .compare(node1.expr, node2.expr) and
			.compareFalseOrNode(node1.def, node2.def) and
			.compareArray(node1, node2)
		}

	Compare_Case(node1, node2)
		{
		return .compareArray(node1, node2) and .compare(node1.body, node2.body)
		}

	Compare_Return(node1, node2)
		{
		return .compareArray(node1, node2)
		}

	Compare_Throw(node1, node2)
		{
		return .compare(node1.expr, node2.expr)
		}

	Compare_TryCatch(node1, node2)
		{
		return node1.catchpat is node2.catchpat and
			.compareFalseOrNode(node1.catchvar, node2.catchvar) and
			.compare(node1.try, node2.try) and
			.compare(node1.catch, node2.catch)
		}

	Compare_Forever(node1, node2)
		{
		return .compare(node1.body, node2.body)
		}

	Compare_ForIn(node1, node2)
		{
		if node1.size isnt node2.size
			return false

		if node1.var is node2.var and
			.compare(node1.expr, node2.expr) and
			.compare(node1.body, node2.body)
			{
			return node1.size is 3/*=range children*/
				? .compare(node1.expr2, node2.expr2)
				: true
			}

		return false
		}

	Compare_For(node1, node2)
		{
		if node1.init.Size() isnt node2.init.Size() or
			node1.inc.Size() isnt node2.inc.Size()
			return false
		for i in .. node1.init.Size()
			if .compare(node1.init[i], node2.init[i]) is false
				return false
		for i in .. node1.inc.Size()
			if .compare(node1.inc[i], node2.inc[i]) is false
				return false
		return .compareFalseOrNode(node1.cond, node2.cond) and
			.compare(node1.body, node2.body)
		}

	Compare_While(node1, node2)
		{
		return .compare(node1.cond, node2.cond) and
			.compare(node1.body, node2.body)
		}

	Compare_DoWhile(node1, node2)
		{
		return .compare(node1.cond, node2.cond) and
			.compare(node1.body, node2.body)
		}

	Compare_Break(@unused)
		{
		return true
		}

	Compare_Continue(@unused)
		{
		return true
		}

	Compare_MultiAssign(node1, node2)
		{
		return .compareArray(node1, node2)
		}

	compareFalseOrNode(node1, node2)
		{
		if node1 is false
			return node2 is false
		return node2 isnt false and .compare(node1, node2)
		}

	compareValue(value1, value2)
		{
		if Type(value1) isnt 'AstNode'
			{
			try
				return Type(value2) isnt 'AstNode' and
					(value1 is value2 or value2(value1))
			return false
			}
		return Type(value2) is 'AstNode' and .compare(value1, value2)
		}

	compareArray(node1, node2)
		{
		if node1.size isnt node2.size
			return false
		for i in .. node1.size
			if .compare(node1[i], node2[i]) is false
				return false
		return true
		}

	matchSpecial?(node1, node2, _noSpecial? = false)
		{
		if noSpecial? is true
			return false

		if .isExpression?(node1)
			return .isWildcard?(node2)

		if node1.type is 'Class' and
			node2.type is 'Ident' and
			node1.base.RemovePrefix('_') is node2.name
			return true

		return false
		}

	expressions: #(RangeLen:, Mem:, Binary:, Ident:, Call:, Nary:, RangeTo:, In:, Unary:,
		Constant:, Trinary:, Block:)
	isExpression?(node)
		{
		return .expressions.GetDefault(node.type, false)
		}

	isWildcard?(node)
		{
		return node.type is 'Ident' and node.name =~ '^[[:alpha:]]$'
		}

	getPos(node, parents)
		{
		if node.pos isnt false
			return Object(pos: node.pos, end: node.end)
		for (i = parents.Size() - 1; i >= 0; i--)
			if parents[i].pos isnt false
				return Object(pos: parents[i].pos, end: parents[i].end)
		return false
		}

	GetHint(search)
		{
		try
			{
			node = AstParseExpression(search)
			res = Object()
			.traverse(node)
				{ |node, parents/*unused*/|
				if node.type is 'Ident' and node.name !~ '^[[:alpha:]]$' and
					node.name not in ('ANY_ARGS', 'this')
					res.Add(node.name)
				else if node.type is 'Constant' and String?(node.value) and
					node.value !~ '^[[:alpha:]]$'
					res.Add(node.value)
				else if node.type is 'Member'
					{
					if node.named is true
						res.Add(node.key)
					if String?(node.value)
						res.Add(node.value)
					}
				0
				}
			if res.NotEmpty?()
				return res.MaxWith(#Size)
			}
		return false
		}
	}
