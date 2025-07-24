// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(recordData, minimizeOutput? = false)
		{
		lineWarnings = Object()
		if LibRecordType(recordData.code) isnt "class" or
			not Libraries().Has?(recordData.lib)
			return .invalidReturnValue(minimizeOutput?)
		curRecPubMethodCalls = .calculatePubMethodCalls(recordData.code)
		undefinedPubMethods = .calculateUndefinedPublicMethods(
			recordData, curRecPubMethodCalls)
		warnings = .addWarnings(
			undefinedPubMethods, recordData, minimizeOutput?, lineWarnings)
		desc = .getDescription(warnings, minimizeOutput?)
		return Object(:warnings, :desc, :lineWarnings)
		}

	invalidReturnValue(minimizeOutput?)
		{
		desc = minimizeOutput? ? "" : "Undefined public method checking aborted"
		return  Object(warnings: #(), :desc, lineWarnings: #())
		}

	getDescription(warnings, minimizeOutput?)
		{
		if not warnings.Empty?()
			return 'Undefined public methods were found'

		if not minimizeOutput?
			return 'No undefined public methods were found'
		return ''
		}

	addWarnings(undefinedPubMethods, recordData, minimizeOutput?, lineWarnings)
		{
		warnings = Object()
		for undefinedPubMethod in undefinedPubMethods
			{
			name = ''
			if minimizeOutput?
				{
				name $= recordData.lib $ ':' $ recordData.recordName $ ':' $
					undefinedPubMethod.line $ ' - '
				lineWarnings.Add(Object(undefinedPubMethod.line))
				}
			name $= '.' $ undefinedPubMethod.method $ ' is undefined'
			warnings.Add([:name])
			}
		return warnings
		}

	calculateUndefinedPublicMethods(recordData, curRecPubMethodCalls)
		{
		recName = recordData.recordName
		builtinMethods = Objects.Members().Add(@BasicMethods)
		try
			{
			cl = recordData.code.SafeEval()
			//If lower library overloads current record and changes the class to a function
			if Type(cl) is "Function"
				cl = RemoveUnderscoreRecordName(recName,
					QueryFirst(recordData.lib $
						' where name is ' $ Display(recName) $
						' sort group').lib_current_text).SafeEval()
			if cl.Method?("Default")
				return Object()
			// could try using .Base? if it works with builtin classes on gSuneido
			if cl.Base() is SocketServer
				builtinMethods.MergeUnion(#('Read', 'Readline', 'RemoteUser', 'Write',
					'Writeline'))
			}
		catch (e)
			{
			// do not need to log for _ inheritance or syntax errors
			if not e.Prefix?('invalid SafeEval') and not e.Prefix?(`can't find`)
				SuneidoLog('ERROR: (CAUGHT) ' $ e, calls:, caughtMsg: 'invalid code')
			return Object()
			}
		return curRecPubMethodCalls.Filter({ not cl.Method?(it.method) and
			not builtinMethods.Has?(it.method)})
		}

	calculatePubMethodCalls(code)
		{
		scan = Scanner(code)
		pubMethodCalls = Object()
		tokenTracker = Record()
		while scan isnt tokenTracker.token = scan.Next()
			{
			type = scan.Type()
			if type in ("WHITESPACE", "COMMENT")
				continue
			tokenTracker.identifier? = scan.Type() is 'IDENTIFIER'//not scan.Keyword?()
			if .pubMethodCall?(tokenTracker) and
				not .stringMethMultipleLines(tokenTracker)
				pubMethodCalls.Add(
					Record(method: tokenTracker.prevToken,
						line: code[.. scan.Position()].LineCount()))

			tokenTracker.prev4Token = tokenTracker.prev3Token
			tokenTracker.prev3Token = tokenTracker.prev2Token
			tokenTracker.prev2Token = tokenTracker.prevToken
			tokenTracker.prevToken = tokenTracker.token
			tokenTracker.prevIdentifier? = tokenTracker.identifier?
			}
		return pubMethodCalls
		}

	pubMethodCall?(tokenTracker)
		{
		if not .pubMethod?(tokenTracker)
			return false
		return .pubMethOneLine?(tokenTracker) or .pubMethMultipleLines?(tokenTracker)
		}

	pubMethod?(tokenTracker)
		{
		return tokenTracker.token is '(' and tokenTracker.prevIdentifier? and
			tokenTracker.prevToken.Capitalized?()
		}

	validTokensBeforePeriod: #('=', '(', '{')
	pubMethOneLine?(tokenTracker)
		{
		return tokenTracker.prev2Token is '.' and
			(tokenTracker.prev3Token.Has?('\n') or
			.validTokensBeforePeriod.Has?(tokenTracker.prev3Token))
		}

	pubMethMultipleLines?(tokenTracker)
		{
		return tokenTracker.prev2Token.Has?('\n') and
			tokenTracker.prev3Token is '.' and
			(tokenTracker.prev4Token.Has?('\n') or
			.validTokensBeforePeriod.Has?(tokenTracker.prev4Token))
		}

	stringMethMultipleLines(tokenTracker)
		{
		return tokenTracker.prev2Token is '.' and .hasQuote?(tokenTracker.prev3Token)
		}

	hasQuote?(s)
		{
		return String?(s) and (s.Has?(`"`) or s.Has?("`") or s.Has?(`'`))
		}
	}
