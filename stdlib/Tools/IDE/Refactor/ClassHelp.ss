// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Class?(text) // NOTE: only checks start - not complete syntax
		{
		return LibRecordType(text) is 'class'
		}
	global_pat: "^_?[[:upper:]]"
	SuperClass(text)
		{
		scan = ScannerWithContext(text)
		if scan is token = scan.Next()
			return false
		if token is 'class' and scan.Ahead() is ':'
			{
			scan.Next()
			if scan.Ahead() =~ .global_pat
				return scan.Ahead()
			}
		else if token =~ .global_pat
			return token
		return false
		}
	AddMethod(text, pos, method)
		{
		pos = .AfterMethod(text, pos)
		return .add_method(text, pos, method)
		}
	AddMethodAtEnd(text, method)
		{
		pos = text.FindLast('}')
		while text[pos - 1] is ' ' or text[pos - 1] is '\t'
			--pos
		return .add_method(text, pos, method)
		}
	add_method(text, pos, method)
		{
		return text[.. pos] $ '\t' $ method $ '\r\n' $ text[pos..]
		}
	AfterMethod(text, pos)
		{
		return .MethodRange(text, pos).to
		}
	AdvanceToNewline(text, pos)
		{
		return pos + text[pos..].Find('\n') + 1
		}
	MethodRange(text, pos)
		{
		nest = 0
		from = 0
		scanner = Scanner(text)
		while scanner isnt token = scanner.Next()
			{
			nest = .nesting(nest, token)
			if nest is 1 and scanner.Type() is #IDENTIFIER
				from = scanner.Position() - token.Size()
			if nest is 1 and scanner.Position() >= pos
				return Object(:from, to: .AdvanceToNewline(text, scanner.Position()))
			}
		return false
		}
	MethodName(text, pos)
		{
		if pos <= text.Find('{')
			return false
		range = ClassHelp.MethodRange(text, pos)
		token = Scanner(text[range.from..]).Next()
		if not String?(token) or token !~ '^[[:alpha:]][_[:alnum:]]*[?!]?$'
			return false
		return token
		}

	Locals(text, cond = function (unused) { true })
		{
		vars = Object()
		scanner = ScannerWithContext(text)
		while scanner isnt token = scanner.Next()
			if .local?(scanner, token) and cond(scanner)
				vars.AddUnique(token)
		return vars
		}
	local?(scanner, token)
		{
		return scanner.Type() is #IDENTIFIER and
			token[0].Lower?() and not scanner.Keyword?() and
			scanner.Ahead() isnt ':' and
			scanner.Prev() isnt '.' and scanner.Prev() isnt '#'
		}
	LocalsInputs(text)
		{
		return .Locals(text, .not_assignment?)
		}
	not_assignment?(scanner)
		{
		return scanner.Ahead() isnt '='
		}
	LocalsModified(text)
		{
		return .Locals(text, .modified?)
		}
	modified?(scanner)
		{
		prev = scanner.Prev()
		ahead = scanner.Ahead()
		return prev is '++' or prev is '--' or
			ahead is '++' or ahead is '--' or
			ahead is '=' or ahead =~ '[^=]='
		}
	LocalsAssigned(text)
		// pre: text includes the parenthesized parameters
		{
		params = Object()
		end = .RetrieveParamsList(text.AfterFirst('('), params)
		locals = .Locals(text[end..], .assignment?)
		return params.MergeUnion(locals)
		}
	RetrieveParamsList(text, vars)
		{
		scanner = Scanner(text)
		while scanner isnt token = scanner.Next()
			{
			if scanner.Type() is #IDENTIFIER
				vars.Add(token.Replace('^_').UnCapitalize())
			else
				{
				if token is '=' // default value
					token = .skipDefault(scanner, token)
				if token is ')'
					break
				}
			}
		return scanner.Position()
		}
	skipDefault(scanner, token)
		{
		nest = 0
		function? = false
		while scanner isnt token = scanner.Next()
			{
			if nest is 0 and token in (')', ',')
				return token
			nest = .nesting(nest, token)

			if '' is function? = .tokenIsFunction?(token, nest, function?)
				return ''
			}
		}

	tokenIsFunction?(token, nest, function?)
		{
		if token in (')', '}', ']')
			{
			if nest <= 0 and function? is false
				return ''
			if function? is true
				return false
			}
		else if token is 'function'
			return true
		return function?
		}

	assignment?(scanner)
		{
		// TODO: handle more than two block arguments
		return scanner.Ahead() is '=' or
			(scanner.Prev() is '|' or scanner.Ahead() is '|') or
			(scanner.Prev() is 'for' or scanner.Ahead() is 'in')
		}

	Methods(text)
		{
		list = Object()
		.foreach_member(text)
			{|scanner, token|
			if scanner.Ahead() is '('
				list.Add(token)
			else if scanner.Ahead() is ':'
				{
				scanner.Next()
				if scanner.Ahead() is 'function'
					{
					scanner.Next()
					list.Add(token)
					}
				}
			}
		return list
		}

	AllObjectMembers(text)
		{
		try //Filter out unnamed
			return text.SafeEval().Members().Filter(String?).Map({ it $ ':' })
		catch
			return Object()
		}

	//Only for classes - The block behaves odd in a function
	AllClassMembers(text)
		{
		classMembers = Object()
		.regularClassMembers(text, classMembers)
		.classMemberDotDeclarations(text, classMembers)
		return classMembers.MapMembers({ classMembers[it] })
		}

	//Only for classes - The block behaves odd in a function
	ClassMembers(text)
		{
		classMembers = Object()
		.regularClassMembers(text, classMembers)
		return classMembers
		}

	regularClassMembers(text, classMembers)
		{
		.foreach_member(text)
			{|scanner, token|
			if scanner.Ahead() is '('
				classMembers.Add(token)
			else if scanner.Ahead() is ':'
				{
				scanner.Next()
				if scanner.Ahead() is 'function'
					{
					scanner.Next()
					classMembers.Add(token)
					}
				else
					classMembers.Add(token $ ':')
				}
			}
		}

	classMemberDotDeclarations(text, classVariables, find = false)
		{
		methods = .MethodSizes(text)
		for method in methods
			{
			if method.to is ""
				endOfString = text.Size() - 1
			else
				endOfString = method.to

			methodCode = text[method.from .. endOfString]
			result = .classMemberDeclaresInParamsMethod(methodCode, classVariables, :find)
			if Number?(result)
				return result + method.from
			result = .classVariablesInMethodBody(methodCode, classVariables, :find)
			if Number?(result)
				return result + method.from
			}
		return classVariables
		}

	FindDotDeclarations(text, member)
		{
		pos = .classMemberDotDeclarations(text, Object(), find: member)
		return Number?(pos) ? pos : false
		}

	parametersText(text)
		{
		scan = Scanner(text)
		nest = 0
		token = scan.Next()
		do
			{
			if token in ('{', '[', '(')
				nest++
			else if token in ('}', ']', ')')
				nest--
			}
			while ((scan isnt token = scan.Next()) and nest isnt 0)
		return text[.. scan.Position()]
		}

	classMemberDeclaresInParamsMethod(text, classVariables, find = false)
		{
		period = false
		text = .parametersText(text)
		scan = Scanner(text)
		while scan isnt token = scan.Next()
			{
			if scan.Type() in (#WHITESPACE, #COMMENT)
				continue
			if period is true
				{
				if token is find
					return scan.Position()
				classVariables.Add(token $ ':')
				period = false
				}
			else
				{
				if token is '='
					token = .skipDefault(scan, token)
				if token is '.'
					period = true
				}
			}
		return classVariables
		}

	validTokensBeforePeriod: #('=', ';', '(', "is", "if", "isnt", "else", "{", ",", ":",
		"return", "or", "and")
	classVariablesInMethodBody(text, classVariables, find = false)
		{
		startPoint = text.Find('{') + 1
		text = text.AfterFirst('{')
		scan = Scanner(text)
		token = prevToken = prev2Token = prev3Token = prev4Token = ""
		prevPos = 0
		curKeyword? = prevKeyword? = true
		while scan isnt token = scan.Next()
			{
			curKeyword? = scan.Keyword?()
			if scan.Type() in (#WHITESPACE, #COMMENT)
				continue
			if (.classVariableOneLine(token, prevKeyword?, prev2Token, prev3Token) or
				.classVariableMultipleLines(token, prevKeyword?, prev2Token, prev3Token,
					prev4Token))
				{
				if find is prevToken
					return prevPos + startPoint
				classVariables.Add(prevToken $ ':')
				}

			prev4Token = prev3Token
			prev3Token = prev2Token
			prev2Token = prevToken
			prevToken = token
			prevPos = scan.Position()
			prevKeyword? = curKeyword?
			}
		return classVariables
		}

	classVariableOneLine(token, prevKeyword?, prev2Token, prev3Token)
		{
		return token is '=' and prevKeyword? is false and prev2Token is '.' and
			(prev3Token.Has?('\n') or .validTokensBeforePeriod.Has?(prev3Token))
		}

	classVariableMultipleLines(token, prevKeyword?, prev2Token, prev3Token, prev4Token)
		{
		return token is '=' and prevKeyword? is false and prev2Token.Has?('\n') and
			prev3Token is '.' and (.validTokensBeforePeriod.Has?(prev4Token) or
			prev4Token.Has?('\n'))
		}

	MethodRanges(text)
		{
		method = false
		list = Object()
		nest = 0
		scanner = ScannerWithContext(text)
		while scanner isnt token = scanner.Next()
			{
			nest = .nesting(nest, token)
			if token is '}' and nest is 1 and method isnt false
				{
				list.Add(method.Add(scanner.Position(), at: #to))
				method = false
				}

			if nest isnt 1 or scanner.Type() isnt #IDENTIFIER
				continue

			method = .methodPos(method, scanner, token)
			}
		return list
		}

	methodPos(method, scanner, token)
		{
		pos = scanner.Position()
		if scanner.Ahead() is '('
			return [name: token, from: pos]
		else if scanner.Ahead() is ':'
			{
			scanner.Next()
			if scanner.Ahead() is 'function'
				return [name: token, from: pos]
			}
		return method
		}

	MethodSizes(text)
		{
		if LibRecordType(text) is 'function'
			return Object([lines: .nonWhiteLineCount(text),
				from: ScannerFind(text, 'function')])
		return .MethodRanges(text).Each()
			{|x|
			x.lines = .nonWhiteLineCount(text[x.from .. x.to])
			}
		}
	nonWhiteLineCount(code)
		{
		nonblank = function (line)
			{ not line.Blank?() }
		return .removeComments(code).Lines().Filter(nonblank).Count()
		}
	removeComments(src)
		{
		dst = ""
		scan = Scanner(src)
		while scan isnt tok = scan.Next2()
			if tok is #COMMENT
				// need to keep newline in case there was text
				// on the same line before /* or after */
				dst $= scan.Text().Tr('^\n')
			else
				dst $= scan.Text()
		return dst
		}

	nesting(nest,  token)
		{
		if #('{', '(', '[').Has?(token)
			nest++
		else if #('}', ')', ']').Has?(token)
			nest--
		return nest
		}

	foreach_member(text, block)
		{
		nest = 0
		scanner = ScannerWithContext(text)
		while scanner isnt token = scanner.Next()
			{
			nest = .nesting(nest, token)
			if nest is 1 and scanner.Type() is #IDENTIFIER
				block(scanner, token)
			}
		}
	FindMethod(text, name)
		{
		.foreach_member(text)
			{|scanner, token|
			if token is name
				return scanner.Position() - token.Size()
			}
		return false
		}
	FindBaseMethod(lib, text, method_name)
		{
		x = Object(:lib, :text)
		while false isnt base = .SuperClass(x.text)
			{
			if false is (x = .find_base(x.lib, base)) or
				not .Class?(text)
				break
			if false isnt .FindMethod(x.text, method_name)
				return Object(lib: x.lib, name: base)
			}
		return false
		}
	find_base(lib, name)
		{
		libs = .libraries()
		if not lib.Suffix?('webgui')
			libs.RemoveIf({ it.Suffix?('webgui') })
		if name.Prefix?('_')
			{
			name = name[1..]
			libs = libs[.. libs.Find(lib)]
			}
		for lib in libs.Reverse!()
			if false isnt x = .getRecord(lib, name)
				return Object(:lib, text: x.text)
		return false
		}
	getRecord(lib, name)
		{
		return Query1(lib $
			' where group = -1 and name = ' $ Display(name))
		}
	libraries()
		{
		return Libraries()
		}

	PrivateMembers(text)
		{
		list = []
		nest = 0
		scanner = ScannerWithContext(text)
		while scanner isnt token = scanner.Next()
			{
			nest = .nesting(nest, token)
			if not token.LocalName?() or scanner.Type() isnt #IDENTIFIER
				continue
			if nest is 1
				list.Add(token)
			else if .privateMem?(nest, scanner)
				list.Add(token)
			}
		return list.Sort!().Unique!()
		}

	privateMem?(nest, scanner)
		{
		return nest > 1 and scanner.Prev() is '.' and not scanner.Prev2().LocalName?() and
			scanner.Ahead() is '='
		}

	PublicMembers(text) // includes inherited
		{
		list = []
		.foreach_member(text)
			{|unused, token|
			if token[0].Upper?()
				list.Add(token)
			}
		if false isnt base = .SuperClass(text)
			try
				.add_members(list, Global(base), base)
		return list.Sort!().Unique!()
		}
	PublicMembersOfName(name)
		{
		if name is "Suneido"
			return Suneido.Members(all:).Sort!().Unique!() // includes methods
		list = []
		word = name.Has?('.') ? name.AfterLast('.') : name
		try
			.add_members(list, Global(name), word)
		return list.Sort!().Unique!()
		}
	add_members(list, x, name) // recursive
		{
		name $= '_'
		for m in x.Members()
			if not m.Prefix?(name) and m !~ "^Class[0-9]+_"
				list.Add(m)
		if false isnt base = BaseClass(x)
			.add_members(list, base, Display(base).BeforeFirst(' '))
		}
	}
