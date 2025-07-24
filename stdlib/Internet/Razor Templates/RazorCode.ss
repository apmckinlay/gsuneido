// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// TODO handle "else"
class
	{
	CallClass(dest, s)
		{
		Dbg(CODE: s)
		Assert(s[0] is '@')
		s = s[1 ..]
		switch s[0]
			{
		case '@':
			// @@ is a literal @
			dest.Html('@')
			return s[1 ..]
		case '(':
			return .expr(dest, s)
		case '{':
			return .code(dest, s)
		default:
			return .name(dest, s)
			}
		}
	members: '^([.]?[a-zA-Z_]+([.][a-zA-Z_]+)*)'
	name(dest, s)
		{
		name = s.Extract(.members, 0)
		if name is 'if' or name is 'for' or name is 'while'
			return .control(dest, s)
		else if name isnt false
			return .var(dest, s, name)
		else
			throw "Razor invalid @" $ s[.. 10]
		}
	var(dest, s, var)
		{
		expr = var
		s = s[var.Size() ..]
		while s[0] is '(' or s[0] is '['
			{
			expr $= s[0]
			closing = .closing[s[0]]
			s = .scanNested(false, {|t| expr $= t }, s, s[0])
			expr $= closing
			if false isnt members = s.Extract('^([.][a-zA-Z_]+([.][a-zA-Z_]+)*)')
				{
				expr $= members
				s = s[members.Size() ..]
				}
			}
		dest.Expr(expr)
		return s
		}
	control(dest, s)
		{
		prefix = s.BeforeFirst('(')
		dest.CodeFragment(prefix $ '(')
		s = s[prefix.Size() ..]
		s = .scanNested(dest, dest.CodeFragment, s, '(')
		dest.Code(')')
		s = s.LeftTrim()
		dest.Code('{')
		s = .scanNested(dest, dest.Code, s, '{')
		dest.Code('}')
		return s
		}
	expr(dest, s)
		{
		return .scanNested(dest, dest.Expr, s, '(')
		}
	code(dest, s)
		{
		return .scanNested(dest, dest.Code, s, '{')
		}
	closing: ('(': ')', '{': '}', '[': ']')
	scanNested(dest, fn, s, opening)
		{
		Assert(s[0] is opening)
		closing = .closing[opening]
		nest = 0
		prev = ''
		scan = Scanner(s)
		do
			{
			if scan is token = scan.Next()
				throw "Razor code missing " $ closing
			if token is opening
				++nest
			if token is closing
				--nest
			// find code before <tag> to run
			if opening is '{' and prev is '<' and scan.Type() is #IDENTIFIER
				{
				pos = scan.Position() - token.Size() - 1
				fn(s[1 :: pos - 2])  // read back from '<' to find code
				s = s[pos ..]
				s = RazorHtml(dest, s)
				Dbg(MORE_CODE: s)
				scan = Scanner(s)  // skip the code
				prev = ''
				}
			// find expression to run
			else if opening is '{' and token is '@'
				{
				pos = scan.Position()
				s = .name(dest, s[pos ..])
				Dbg(MORE_EXPR: s)
				scan = Scanner(s) // skip the expression
				prev = ''
				}
			else
				prev = token
			}
			while nest > 0
		pos = scan.Position()
		if pos > 2
			fn(s[1 :: pos - 2]) // run code after <tags>
		return s[pos ..]
		}
	}
