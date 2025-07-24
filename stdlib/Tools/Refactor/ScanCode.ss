// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// keeps track of "context" = 'code', 'class', 'constant'
class
	{
	New(code, begin = false, end = false)
		{
		.scan = ScannerWithContext(code)
		.stack = Stack()
		.attached = Object()
		.begin = begin is false ? function(unused){} : begin
		.end = end is false ? function(unused){} : end
		}
	context: ''
	next_context: false
	context_after: false
	nest: 0
	nest_by: ''
	params: ''
	inParams: false
	Next()
		{
		if .context_after isnt false
			{
			.context = .context_after
			.context_after = false
			}
		token = .scan.Next()
		if token is .scan
			return this
		ahead = .scan.Ahead()

		if .context is 'class' and token is '('
			{
			.push()
			.context = 'constant'
			.nest_by = '()'
			}
		// fall through

		if token is 'dll'
			{
			.push()
			.nest_by = '()'
			.context = 'dll'
			}
		else if token is 'function'
			{
			.push()
			.nest_by = '{}'
			.context = 'code'
			}
		else if .method?()
			{
			.push()
			.nest_by = '{}'
			.params = '()'
			.context = 'class'
			.context_after = 'code'
			}
		else if .object?(token)
			{
			.push()
			.context = 'constant'
			.nest_by = ahead is '(' ? '()' : '{}'
			}
		else if .class?(token)
			{
			.push()
			.context = 'code'
			.next_context = 'class'
			.nest_by = '{}'
			}
		else if token is .params[0]
			.inParams = true
		else if token is .params[1]
			{
			.inParams = false
			.params = ""
			}
		else if token is .nest_by[0] and not .inParams
			{
			if .next_context isnt false
				{
				.context = .next_context
				.next_context = false
				}
			++.nest
			}
		else if token is .nest_by[1] and not .inParams
			{
			if --.nest is 0
				.pop()
			}
		return token
		}
	push()
		{
		.stack.Push([context: .context, nest: .nest, nest_by: .nest_by,
			attached: .attached, params: .params, inParams: .inParams])
		.nest = 0
		.nest_by = ''
		.context = ''
		.params = ''
		.inParams = false
		.attached = Object()
		(.begin)(.attached)
		}
	pop()
		{
		(.end)(.attached)
		top = .stack.Pop()
		.nest = top.nest
		.nest_by = top.nest_by
		.params = top.params
		.inParams = top.inParams
		.context = top.context
		.attached = top.attached
		}
	object?(token)
		{
		ahead = .scan.Ahead()
		return token is '#' and (ahead is '(' or ahead is '{')
		}
	class?(token)
		{
		if .next_context is 'class'
			return false // e.g. class : Global {
		if token is 'class'
			return true
		prev = .scan.Prev()
		special = #('.':, '#':, 'if':, 'switch':, 'while':, 'in':, 'is':, 'isnt':)
		return .scan.Type() is #IDENTIFIER and
			(token.GlobalName?() or (token[0] is '_' and token[1].Upper?())) and
			.scan.Ahead() is '{' and
			not special.Member?(prev)
		}
	method?()
		{
		return .context is 'class' and
			.scan.Type() is #IDENTIFIER and
			.scan.Ahead() is '('
		}
	Context()
		{
		return .context
		}
	Attached()
		{
		return .attached
		}
	Global?()
		{
		return .Context() is 'code' and
			.Type() is #IDENTIFIER and
			not .Keyword?() and
			.Prev() isnt '.' and .Prev() isnt '#' and
			(.Ahead() isnt ':' or .AheadPos() isnt .Position() + 1) and
			.Token().GlobalName?()
		}
	Default(@args)
		{
		return .scan[args[0]](@+1args)
		}
	}