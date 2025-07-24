// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
JsTranslate // to get output stuff
	{
	CallClass(ast, last = false)
		{
		this[ast.type](ast, :last)
		}
	ExprStmt(ast, last)
		{
		if last
			.Print('return ')
		JsTranslateExpression(ast.expr, discard: not last)
		.Println(';')
		}

	MultiAssign(ast)
		{
		size = ast.size
		Assert(ast[size - 1].type is: 'Call')
		.Print('[')
		for i in ..size-1
			{
			Assert(ast[i].type is: 'Ident')
			if i > 0
				.Print(', ')
			name = .ConvertQMark(ast[i].name)
			_var(name)
			.Print(name)
			}
		.Print('] = (')
		JsTranslateExpression(ast[size - 1])
		.Println(');')
		}

	Return(ast)
		{
		if _inblock
			{
			_fn.hasBlockReturn = true
			.Print('throw su.blockreturn(' $ _fn.id)
			if ast.expr isnt false
				{
				.Print(', ')
				JsTranslateExpression(ast.expr)
				}
			.Println(');')
			}
		else
			{
			.Print('return')
			if ast.size > 1 // multiple returns
				{
				.Print(' [')
				for i in ..ast.size
					{
					if i > 0
						.Print(', ')
					JsTranslateExpression(ast[i], discard:)
					}
				.Print(']')
				}
			else if ast.expr isnt false
				{
				.Print(' ')
				JsTranslateExpression(ast.expr)
				}
			.Println(';')
			}
		}
	Forever(ast)
		{
		.Print('while (true)')
		.loopBody(ast.body)
		}
	While(ast)
		{
		.Print('while (')
		JsTranslateExpression.Bool(ast.cond)
		.Print(')')
		.loopBody(ast.body)
		}
	For(ast)
		{
		.Print('for (')
		.exprList(ast.init)
		.Print('; ')
		if ast.cond isnt false
			JsTranslateExpression.Bool(ast.cond)
		.Print('; ')
		.exprList(ast.inc)
		.Print(')')
		.loopBody(ast.body)
		}
	DoWhile(ast)
		{
		.Print('do')
		.loopBody(ast.body, newline?: false)
		.Print(' while (')
		JsTranslateExpression.Bool(ast.cond)
		.Println(');')
		}
	ForIn(ast)
		{
		if ast.size is 3/*=range children*/
			{
			.forRange(ast)
			return
			}
		if ast.var2 isnt ''
			{
			.forIn2(ast)
			return
			}
		var = ast.var
		_var(var)
		iter = _nextTmp()
		.Print('for (' $ iter $ ' = su.iter(')
		JsTranslateExpression(ast.expr)
		.Print('); ' $ iter $ ' !== (' $ var $ ' = su.next(' $ iter $ ')); )')
		.loopBody(ast.body)
		}
	forRange(ast)
		{
		var = ast.var
		if var isnt ''
			_var(var)
		else
			var = _nextTmp()
		from = _nextTmp()
		to = _nextTmp()
		.Print(from $ ' = (')
		JsTranslateExpression(ast.expr)
		.Println(');')
		.Print(to $ ' = (')
		JsTranslateExpression(ast.expr2)
		.Println(');')
		.Print('for (' $ var $ ' = ' $ from $ '; su.lt(' $ var $ ', ' $ to $
			'); ' $ var $ ' = su.inc(' $ var $ '))')
		.loopBody(ast.body)
		}
	forIn2(ast)
		{
		var = ast.var
		var2 = ast.var2
		_var(var)
		_var(var2)
		tmp = _nextTmp()
		iter = _nextTmp()
		.Print('for (' $ iter $ ' = su.iter(su.invoke(')
		JsTranslateExpression(ast.expr)
		.Print(', "Assocs")); ' $ iter $ ' !== (' $ tmp $ ' = su.next(' $ iter $ ')); )')
		.Println(' {')
		.Println(var $ ' = ' $ tmp $ '.get(0);')
		.Println(var2 $ ' = ' $ tmp $ '.get(1);')
		.loopBody(ast.body)
		.Println('}')
		}
	loopBody(body, newline? = true)
		{
		_inLoop = true
		if newline?
			.bodyln(body)
		else
			.body(body)
		}
	Break(ast /*unused*/)
		{
		.Println(_inblock and not _inLoop
			? 'throw su.exception("block:break");'
			: 'break;')
		}
	Continue(ast /*unused*/)
		{
		.Println(_inblock and not _inLoop
			? 'throw su.exception("block:continue");'
			: 'continue;')
		}
	If(ast)
		{
		.Print('if (')
		JsTranslateExpression.Bool(ast.cond)
		.Print(')')
		.body(ast.t)
		if ast.f is false or
			ast.f.type is 'Compound' and ast.f.size is 0
			.Println()
		else
			{
			.Print(' else')
			if ast.f.type is 'If'
				{
				.Print(' ')
				JsTranslateStatement(ast.f) // recursive
				}
			else
				.bodyln(ast.f)
			}
		}
	Switch(ast)
		{
		.Print('switch (')
		JsTranslateExpression(ast.expr)
		.Println(') {')
		.cases(ast)
		.Println('}')
		}
	cases(ast)
		{
		for (i = 0; i < ast.size; i++)
			{
			c = ast[i] // c is Case
			for (j = 0; j < c.size; j++)
				{
				.Print('case ')
				JsTranslateExpression(c[j])
				.Println(':')
				}
			.Indented()
				{
				if c.body.size > 0
					{
					body = c.body
					for (k = 0; k < body.size; k++)
						{
						if body[k].type is 'Compound'
							for (p = 0; p < body[k].size; p++)
								JsTranslateStatement(body[k][p])
						else
							JsTranslateStatement(body[k])
						}
					if body[body.size - 1].type not in ("Break", "Return", "Throw")
						.Println('break;')
					}
				else
					.Println('break;')
				}
			}
		if false isnt def = ast.def
			{
			.Println('default:')
			for (k = 0; k < def.size; k++)
				JsTranslateStatement(def[k])
			if def.size is 0
				.Println('break;')
			}
		else
			{
			.Println('default:')
			.Indented()
				{
				.Println('throw su.exception("unhandled switch case");')
				}
			}
		}


	TryCatch(ast)
		{
		// JavaScript catch variable is scoped to catch but Suneido's isnt
		// In Suneido the catch variable is optional, in JavaScript it's not
		// Suneido allows a catch filter
		.Print('try')
		.body(ast.try)
		.Println(' catch (_e) {')
		if _fn.hasBlockReturn
			.Indented()
				{
				.Println('su.rethrowBlockReturn(_e);')
				}
		if false is catchPart = ast.catch
			{
			.Println('}')
			return
			}
		.Indented()
			{
			if ast.catchvar isnt false
				{
				filter = ast.catchpat is false ? '' : ast.catchpat
				.Println('var ' $ ast.catchvar.name $
					' = su.catchMatch(_e, "' $ filter $ '");')
				}
			if catchPart.type is 'Compound'
				for (i = 0; i < catchPart.size; i++)
					JsTranslateStatement(catchPart[i])
			else
				JsTranslateStatement(catchPart)
			}
		.Println('}')
		}
	Throw(ast)
		{
		.Print('throw su.exception(')
		JsTranslateExpression(ast.expr)
		.Println(');')
		}
	exprList(list) // list is an object of exprs
		{
		for (i = 0; i < list.Size(); ++i)
			{
			if i > 0
				.Print(', ')
			JsTranslateExpression(list[i], discard:)
			}
		}
	bodyln(ast)
		{
		.body(ast)
		.Println()
		}
	body(ast)
		{
		.Println(' {')
		.Indented()
			{
			if ast.type is 'Compound'
				for (i = 0; i < ast.size; i++)
					JsTranslateStatement(ast[i])
			else
				JsTranslateStatement(ast)
			}
		.Print('}')
		}
	}
