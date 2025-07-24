// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Limit: 1000
	CallClass(data, printFn, completion)
		{
		records? = .records?(data)
		nrecs = nlines = 0
		rand = false
		errors = Object()
		QueryApply(.query(data))
			{|x|
			if x.lines.Size() is 1 and String?(x.lines[0])
				{
				errors.Add(x.lib $ ':' $ x.name $ ':' $ x.lines[0])
				if errors.Size() > .Limit
					{
					errors.Each(printFn)
					printFn("TOO MANY ERRORS")
					return
					}
				continue
				}
			++nrecs
			if data.show is 'random record'
				{
				if .random1ofn?(nrecs)
					rand = x
				continue
				}
			if .forceStop?(data.findUid, printFn)
				return
			nlines += .print(x, records?, nlines, printFn)
			if .limitReached?(nrecs, nlines, printFn)
				return
			}
		if rand isnt false
			printFn(rand.lib $ ':' $ rand.name)
		completion(nrecs, nlines)
		errors.Each(printFn)
		}

	random1ofn?(n)
		{
		f = 13 // so we're not just always comparing to 0
		return Random(n * f) < f
		}

	forceStop?(findUid, printFn)
		{
		if stop? = ServerSuneido.GetAt(#workSpaceFindStop, findUid, false) is true
			printFn("FIND STOPPED")
		return stop?
		}

	print(x, records?, nlines, printFn)
		{
		nlines = 0
		if records?
			printFn(x.lib $ ':' $ x.name)
		else
			nlines = .printLines(printFn, x, nlines)
		return nlines
		}

	limitReached?(nrecs, nlines, printFn)
		{
		if limitReached? = (nrecs > .Limit or nlines > .Limit)
			printFn("TOO MANY MATCHES")
		return limitReached?
		}

	records?(data)
		{
		if data.show is 'records'
			return true
		return data.bytoken is "" and data.byexpression is "" and
			not data.textRepeat.Any?({ it.text isnt '' and it.exclude isnt true })
		}

	printLines(printFn, x, nlines)
		{
		for lineOb in x.lines
			{
			first? = true
			prefix = x.lib $ ':' $ x.name $ ':' $ (lineOb[0] + 1) $ ': '
			lines = lineOb.Map({ x.lib_current_text.NthLine(it) })
			sharedTabs = lines.
				Map({ it.Blank?()
					? 9999/*=big*/
					: (it.Size() - it.LeftTrim('\t').Size()) }).
				Min()
			for line in lines.Map({ it[sharedTabs..] })
				{
				printFn((first? ? prefix : ' '.Repeat(prefix.Size())) $ line)
				first? = false
				if ++nlines > .Limit
					break
				}
			}
		return nlines
		}

	query(data)
		{
		return .libsQuery(data) $
			.nameWhere(data) $
			.textExcludes(data) $
			.linesExtend(data) $
			' sort ' $ (data.GetDefault('sort', 'name') is 'name' ? 'name' : 'lib,name')
		}

	libsQuery(data)
		{
		allLibs = LibraryTables()
		libs = data.libs is '(All)' or data.libs is ''
			? Libraries().MergeUnion(allLibs) // searching lib in use first
			: data.libs is '(In Use)'
				? Libraries()
				: data.libs.Split(',')
		if data.exclude is true
			libs = allLibs.Difference(libs)
		return '(' $
			libs.Map({ '(' $ it $ ' extend lib=' $ Display(it) $ ', lib_current_text)' }).
				Join('\nunion\n') $
			' where group = -1)\n'
		}

	nameWhere(data)
		{
		expr = data.nameRepeat.
			Filter({ it.text isnt '' }).
			Map({ op = it.exclude is true ? ' !~ ' : ' =~ '
				'name' $ op $ Display(Find.Regex(it.text, it)) }).
			Join(' and ')
		return Opt(' where ', expr, '\n')
		}

	textExcludes(data)
		{
		expr = data.textRepeat.
			Filter({ it.text isnt '' and it.exclude is true }).
			Map({ 'lib_current_text !~ ' $ Display(Find.Regex(it.text, it))	}).
			Join(' and ')
		return Opt(' where ', expr, '\n')
		}

	linesExtend(data)
		{
		return Opt('extend lines = FindCombine(', .linesExpr(data),
			', ' $ Display(data.findUid) $
			')\nwhere not lines.Empty?()\n')
		}

	linesExpr(data)
		{
		list = data.textRepeat.
			Filter({ it.text isnt '' and it.exclude isnt true }).
			Map({ 'Object(#RegexMatchLines, lib_current_text, ' $
				Display(Find.Regex(it.text, it)) $ ')' })
		if data.bytoken isnt ""
			{
			tokens = FindByTokenScan(data.bytoken)
			regex = .quickCheck(tokens)
			list.Add('Object(#FindByToken, lib_current_text, ' $
				Display(tokens) $ ', ' $ Display(regex) $ ')')
			}
		if data.byexpression isnt ''
			{
			hint = AstSearch.GetHint(data.byexpression)
			list.Add('Object(#FindByExpression, name, lib_current_text, ' $
				Display(data.byexpression) $ ', ' $ Display(hint) $ ')')
			}
		return list.Join(',\n')
		}

	quickCheck(tokens)
		{
		return tokens.
			Filter({ it.Identifier?() and it.Size() > 1 and not .keyword?(it) }).
			Map({ '\<' $ it $ '\>' }).
			SortWith!({ -it.Size() }) // longest first (hopefully less common)
		}

	keyword?(tok)
		{
		scan = Scanner(tok)
		return scan isnt scan.Next2() and scan.Keyword?()
		}
	}
