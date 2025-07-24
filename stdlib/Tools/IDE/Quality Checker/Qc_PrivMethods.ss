// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(recordData, minimizeOutput? = false)
		{
		warnings = Object()
		lineWarnings = Object()
		if not recordData.recordName.Prefix?("Rule_") and
			LibRecordType(recordData.code) is "class"
			{
			allMembers = ClassHelp.AllClassMembers(recordData.code)
			methods = ClassHelp.MethodRanges(recordData.code)
			privMethods = .getPrivMethods(recordData, allMembers, methods)
			privCalls = .findAllPrivMethodCalls(recordData.code, methods)
			privMethodCalls = .pruneOutClassVariables(privCalls, allMembers, privMethods)
			unusedPrivMethods = .calculateUnusedPrivMethods(privMethods, privMethodCalls)
			undefinedPrivMethods = .calculateUndefinedPrivMethods(privMethods,
				privMethodCalls)
			.addWarnings(recordData, warnings, undefinedPrivMethods, 'undefined',
				lineWarnings, :minimizeOutput?)
			.addWarnings(recordData, warnings, unusedPrivMethods, 'unused',
				lineWarnings, :minimizeOutput?)
			}
		desc = .getDescription(warnings, minimizeOutput?)
		rating = .getRating(warnings)
		return Object(:warnings, :desc, :rating, :lineWarnings)
		}

	pruneOutClassVariables(privMethodCalls, allPrivMembers, privMethods)
		{
		privMethodsList = privMethods.MapMembers({ privMethods[it].method })
		return privMethodCalls.Filter()
			{
			method = it.method
			allPrivMembers.Member?("get_" $ method) or
				allPrivMembers.Member?("getter_" $ method) or
				allPrivMembers.Member?(method $ ':') is false or
				privMethodsList.Member?(method)
			}
		}

	getRating(warnings)
		{
		return warnings.Size() isnt 0 ? 0 : 5
		}

	getDescription(warnings, minimizeOutput?)
		{
		if not warnings.Empty?()
			return 'There are unused or undefined private methods'

		if not minimizeOutput?
			return 'There are no undefined or unused private methods'
		return ''
		}

	findPrivMethodCallsInMethod(code, offset)
		{
		privMethodCalls = Object()
		scan = Scanner(code)
		prevToken = prev2Token = prev3Token = ""
		isPrevTokenIdentifier? = isPrev2TokenIdentifier? = ""
		isPrev3TokenIdentifier? = ""

		while scan isnt token = scan.Next()
			{
			if scan.Type() in (#WHITESPACE, #COMMENT)
				continue
			if .privateMethod?(token, prevToken, prev2Token, prev3Token,
				isPrev2TokenIdentifier?, isPrev3TokenIdentifier?)
				{
				privMethodCalls.Add(Record(method: token,
					line: code[.. scan.Position()].LineCount() + offset))
				}

			isPrev3TokenIdentifier? = isPrev2TokenIdentifier?
			isPrev2TokenIdentifier? = isPrevTokenIdentifier?
			isPrevTokenIdentifier? = scan.Type() is #IDENTIFIER
			prev3Token = prev2Token
			prev2Token = prevToken
			prevToken = token
			}
		return privMethodCalls
		}

	validTokensBeforePeriod: #("return", "if", "in", "try", "while", "throw")
	privateMethod?(token, prevToken, prev2Token, prev3Token, isPrev2TokenIdentifier?,
		isPrev3TokenIdentifier?)
		{
		if .invalidToken?(token)
			return false

		return .privMethodOneLine(prevToken, prev2Token, isPrev2TokenIdentifier?) or
			.privMethodMultipleLines(prevToken, prev2Token, prev3Token,
				isPrev3TokenIdentifier?)
		}

	privMethodMultipleLines(prevToken, prev2Token, prev3Token, isPrev3TokenIdentifier?)
		{
		return prevToken.Has?('\n') and prev2Token is '.' and
			prev3Token not in (')', ']') and
			(prev3Token.Has?('\n') or not isPrev3TokenIdentifier? or
			.validTokensBeforePeriod.Has?(prev3Token))
		}

	privMethodOneLine(prevToken, prev2Token, isPrev2TokenIdentifier?)
		{
		return prevToken is '.'  and prev2Token not in (')', ']') and
			(prev2Token.Has?('\n') or not isPrev2TokenIdentifier? or
			.validTokensBeforePeriod.Has?(prev2Token))
		}

	invalidToken?(token)
		{
		return token.Capitalized?() or not token[0].Alpha?() or token.Has?('\n') is true
		}

	findAllPrivMethodCalls(code, methods)
		{
		return methods.FlatMap({ .findPrivMethodCallsInMethod(code[it.from .. it.to],
			code[.. it.from].LineCount() - 1) })
		}

	getPrivMethods(recordData, allMembers, methods)
		{
		privMethods = Object()
		for method in methods
			{
			if method.name.Capitalized?()
				continue
			line = recordData.code[.. method.from].LineCount()
			if method.name is 'function' and allMembers.Member?(
				recordData.code.Lines()[line - 1].BeforeFirst(':').Trim())
				continue
			privMethods.Add([method: method.name, :line])
			}
		return privMethods
		}

	calculateUnusedPrivMethods(privMethods, privMethodCalls)
		{
		unusedPrivMethods = Object()
		for privMethodRec in privMethods
			{
			method = privMethodRec.method
			if method =~ `^contrib\d?_`
				continue // ignore getters
			if not privMethodCalls.HasIf?({ .privmeth(method, it.method) })
				unusedPrivMethods.Add(privMethodRec)
			}
		return unusedPrivMethods
		}

	calculateUndefinedPrivMethods(privMethods, privMethodCalls)
		{
		undefinedPrivMethods = Object()
		for privMethodCallRec in privMethodCalls
			{
			call = privMethodCallRec.method
			if not privMethods.HasIf?({ .privmeth(it.method, call) })
				undefinedPrivMethods.Add(privMethodCallRec)
			}
		return undefinedPrivMethods
		}
	privmeth(method, call)
		{
		return call is method or
			"get_" $ call is method or "getter_" $ call is method
		}

	addWarnings(recordData, warnings, methods, type, lineWarnings,
		minimizeOutput? = false)
		{
		err = type is 'undefined' ? 'ERROR' : 'WARNING'
		for rec in methods
			{
			name = minimizeOutput?
				? recordData.lib $ ':' $ recordData.recordName $ ':' $ rec.line $ ' ' $
					err $ ': '
				: ''
			name $= '.' $ rec.method $ " is " $ type
			warnings.Add([:name])
			if minimizeOutput?
				lineWarnings.Add(Object(rec.line))
			}
		}
	}
