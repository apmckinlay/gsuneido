// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// TODO include name: and .name= members
// TODO change defaultMethods to most common instead of all basic
Addon_auto_complete
	{
	minWordSize: 3
	scanInterval: 5 // scan at most every this many seconds

	Init()
		{
		super.Init()
		.defaultMethods = Objects.Members().Add(@BasicMethods).Sort!().Unique!()
		LibLocateList.Start()
		.lastIdleAfterChange = Date()
		.paramsCache = Object()
		}
	AutoComplete(word)
		{
		if word.Has?('.')
			.autocomplete_method(word)
		else
			.autocomplete_word(word)
		}
	AutocShow(word/*unused*/, matches)
		{
		// override Addon_auto_complete and pass 0 as word length
		// this is to get around Scintilla wanting an exact prefix match
		.SCIAutocShow(0, matches)
		}
	autocomplete_method(word)
		{
		methods = Object()
		if word[0] is '.'
			methods = .thisMembers(word[1..])
		else if word[0].Upper?()
			methods = ClassHelp.PublicMembersOfName(word.BeforeLast('.'))
		else
			methods = .defaultMethods
		word = word.AfterLast('.')

		.AutoShow(word, .Matches(methods, word))
		}
	thisMembers(word)
		{
		if word is ""
			return .all_members
		else if word.LocalName?()
			return .private_members
		else if word.GlobalName?()
			return .public_members
		else
			return #()
		}

	autocomplete_word(word)
		{
		paramCandidates = .buildParamMatchCandidates(word)
		wordCandidates = .buildWordMatchCandidates(word)
		.AutoShow(word, paramCandidates.Append(wordCandidates))
		}

	buildParamMatchCandidates(word)
		{
		if word !~ '^[a-z]'
			return []

		if false is paramList = .getParamList()
			return []

		return paramList.Filter({ it.Prefix?(word) }).Map({ it $ ':' })
		}

	getParamList()
		{
		if false is caller = .findCaller(.Get(), .GetCurrentPos())
			return false

		if caller[0] is '.'
			{
			caller = caller[1..]
			name = .Send("CurrentName")
			if not caller.Capitalized?()
				caller = name $ '_' $ caller
			caller = name $ '.' $ caller
			}

		if .paramsCache.Member?(caller)
			return .paramsCache[caller]

		return .paramsCache[caller] = .buildParamList(caller)
		}

	buildParamList(caller)
		{
		try
			value = Global(caller)
		catch
			return false

		if Type(value) not in ("Class", "Method", "Function")
			return false

		if Type(value) is #Class
			{
			if false isnt c = value.MethodClass('CallClass')
				value = c.CallClass
			else if false isnt c = value.MethodClass('New')
				value = c.New
			else
				return false
			}

		params = value.Params()
		fakeFunc = 'function' $ params $ '{}'
		try
			{
			ast = Tdop(fakeFunc)
			list = Object()
			if ast[2].Token is TDOPTOKEN.PAREM_AT
				list.Add(ast[2][1].Value)
			else
				for param in ast[2].Children
					list.Add(param[1].Value)
			return list
			}
		catch
			return false
		}

	buildWordMatchCandidates(word)
		{
		if word.Size() < .minWordSize
			return []
		return .matching_words(word)
		}

	matching_words(word)
		{
		return word.Capitalized?()
			? LibLocateList.GetMatches(word, justName:)
			: .MatchesExcludingSelf(.text_names, word).Sort!().Unique!()
		}

	AutocSelection(scn)
		{
		// do our own fillin because Scintilla requires prefix to match fillin
		// these actions automatically cancel the Scintilla fillin
		word = ('.' $ .GetCurrentReference()).AfterLast('.')
		for unused in word
			.CharLeftExtend()
		.ReplaceSel(scn.text)
		}

	IdleAfterChange()
		{
		if Date().MinusSeconds(.lastIdleAfterChange) > .scanInterval
			{
			.idleAfterChangeThrottled()
			.lastIdleAfterChange = Date()
			}
		}

	text_names: ()
	all_members: ()
	public_members: ()
	private_members: ()
	idleAfterChangeThrottled()
		{
		text = .Get()
		.text_names = .getTextNames(text)
		if LibRecordType(text) is 'class'
			{
			.private_members = ClassHelp.PrivateMembers(text)
			.public_members = ClassHelp.PublicMembers(text)
			.all_members = .public_members.Copy().Add(@.private_members)
			}
		else
			.private_members = .public_members = .all_members = #()
		}
	getTextNames(text)
		{
		words = Object()
		for token in scan = Scanner(text)
			if scan.Type() is #IDENTIFIER and token.Size() > .minWordSize
				words.Add(token)
		return words.Sort!().Unique!()
		}

	findCaller(text, pos)
		{
		type = LibRecordType(text)
		if type not in (#class, #function)
			return false

		if type is #class
			{
			if false is range = ClassHelp.MethodRange(text, pos)
				return false
			text = text[range.from..range.to]
			pos -= range.from
			}

		return .find(text, pos)
		}

	find(text, pos)
		{
		stack = Stack()
		stack.Push(false)
		scan = ScannerWithContext(text)
		inBody? = false
		while scan isnt scan.Next() and scan.Position() <= pos
			{
			if inBody?
				{
				if scan.Ahead() is '('
					stack.Push(.handleLParen(scan))
				else if scan.Token() is ')'
					stack.Pop()
				}
			inBody? = inBody? or scan.Token() is '{'
			}

		return stack.Pop()
		}

	handleLParen(scan)
		{
		if scan.Type() isnt #IDENTIFIER or scan.Keyword?()
			return false
		caller = scan.Token()

		if scan.Prev_Token is '.'
			{
			caller = '.' $ caller
			if scan.Prev2_Type is #IDENTIFIER and not scan.Prev2_Keyword?
				caller = scan.Prev2_Token $ caller
			}
		return caller
		}
	}
