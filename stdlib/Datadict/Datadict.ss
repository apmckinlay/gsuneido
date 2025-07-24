// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(fieldname, getMembers = false)
		{
		name = "Field_" $ fieldname
		try
			dd = Global(name)
		catch
			return .getDDForUndefinedField(fieldname)

		return getMembers isnt false
			? .ddValues(dd, getMembers)
			: .inject?(dd)
				? .Cache(name)
				: dd
		}

	ddValues(dd, members)
		{
		if members.Has?('Control') or members.Has?('Format')
			throw 'Cannot get Control or Format. ' $
				'Get whole datadict to ensure any injects are included'
		ob = Object()
		for mem in members
			if dd.Member?(mem)
				ob.Add(dd[mem] at: mem)
		return ob
		}

	Cache(name)
		{
		if not Suneido.Member?('DatadictCache')
			.init()
		return Suneido.DatadictCache.Get(name)
		}

	cacheSize: 250
	init()
		{
		LibUnload.AddObserver(#Datadict, Datadict.Unload)
		Suneido.DatadictCache = LruCache(.injectMembers, .cacheSize)
		}

	Unload(name)
		{
		if name.Prefix?("Field_") and
			false isnt cache = Suneido.GetDefault(#DatadictCache, false)
			cache.Reset()
		}

	getDDForUndefinedField(fieldname)
		{
		for func in #(total max min average)
			if fieldname.Prefix?(func $ "_")
				{
				d = .getDDForSummarizeField(fieldname, func)
				return .changePrompt(d,
					{ |origPrompt| func.Capitalize() $ " " $ origPrompt })
				}
		if fieldname =~ "_\d+$"
			return Datadict(fieldname.BeforeLast('_'))
		if fieldname.Suffix?('_lower!')
			{
			d = Datadict(fieldname.BeforeLast('_'))
			return .changePrompt(d, { |origPrompt| origPrompt $ '*' })
			}
		return Field_string
		}

	changePrompt(d, fn)
		{
		d = Class?(d) ? new d : d.Copy()
		for m in #(Prompt, Heading, SelectPrompt)
			if d.Member?(m)
				d[m] = fn(d[m])
		return d
		}

	getDDForSummarizeField(fieldname, func)
		{
		dd = Datadict(fieldname.AfterFirst("_"))
		if func in ('total', 'average') and not dd.Base?(Field_number)
			{
			newDD = new Field_number
			for m in #(Prompt, Heading, SelectPrompt)
				if dd.Member?(m)
					newDD[m] = dd[m]
			dd = newDD
			}
		return dd
		}

	inject?(dd)
		{
		return Class?(dd) and dd.Members(all:).HasIf?(
			{ it.Prefix?('Control_') or it.Prefix?('Format_') or
				it.Prefix?('SelectControl_') })
		}

	injectMembers(name)
		{
		dd = new (Global(name))
		ddMembers = dd.Members(all:)
		for type in #(Format, Control, SelectControl)
			{
			if dd.Member?(type)
				dd[type] = dd[type].Copy()
			for mem in ddMembers.Filter({ it.Prefix?(type $ "_") })
				dd[type][mem.AfterFirst(type $ "_")] = dd.Val_or_func(mem)
			}
		return dd
		}

	/* Field Methods */
	FieldMap(fields, func)
		{
		map = Object()
		for field in fields
			map[field] = this[func](field)
		return map
		}

	Prompt(field, excludeTags = #(Internal), logAsError = true)
		{
		ddMembers = Object('Prompt').Append(excludeTags)
		ddVals = .CallClass(field, ddMembers)
		.checkAndLogTags(ddVals, field, excludeTags, logAsError)
		if ddVals.Member?("Prompt")
			return ddVals.Prompt

		.logError('no Prompt for: ', field, logAsError)
		return field
		}

	logError(err, field, logAsError = true)
		{
		// Capital/digit/multiple words prefix presumed intentional shortcut around datadicts
		if not (field =~ "^[[:upper:]|[:digit:]]") and not field.Has?(' ')
			{
			if logAsError
				.programmerError(err $ field)
			else
				.suneiodlog('INFO: ' $ err $ field)
			}
		}
	programmerError(msg) // for tests
		{
		ProgrammerError(msg)
		}
	suneiodlog(msg) // for tests
		{
		SuneidoLog.Once(msg, calls:)
		}

	checkAndLogTags(ddVals, field, excludeTags, logAsError = true)
		{
		if false isnt tag = .checkTags(ddVals, excludeTags)
			{
			if .recentlyDeletedCustomField?(tag, field)
				return

			msg = field $ ' should have been excluded due to tag: ' $ tag
			.logError(msg, '', :logAsError)
			}
		}

	checkTags(ddVals, excludeTags)
		{
		for ddTag in excludeTags
			if ddVals.Member?(ddTag) and ddVals[ddTag] is true
				return ddTag
		return false
		}

	recentlyDeletedCustomField?(tag, field)
		{
		return #(Internal, ExcludeSelect).Has?(tag) and
			Customizable.CustomField?(field, includeDeleted:) and
			false isnt TableExists?('suneidolog') and
			not QueryEmpty?('suneidolog where sulog_message is ' $
				Display(Customizable.DeletedCustomFieldMessage(field)) $
				' and sulog_timestamp > ' $ Display(Date().Minus(days: 1)))
		}

	PromptOrHeading(field, excludeTags = #(Internal), logAsError = true)
		{
		ddMembers = Object('Prompt', 'Heading', 'SelectPrompt').Append(excludeTags)
		ddVals = .CallClass(field, ddMembers)
		.checkAndLogTags(ddVals, field, excludeTags, logAsError)

		promptOrHeading = ddVals.Member?("Prompt") and ddVals.Prompt isnt ""
			? ddVals.Prompt
			: ddVals.Member?("Heading")
				? ddVals.Heading
				: ddVals.Member?('SelectPrompt')
					? ddVals.SelectPrompt
					: field

		if promptOrHeading is field or promptOrHeading is ''
			.logError('no PromptOrHeading for: ', field, logAsError)
		return TranslateLanguage(promptOrHeading)
		}

	Heading(field, excludeTags = #(Internal), logAsError = true)
		{
		ddMembers = Object('Heading', 'Prompt').Append(excludeTags)
		ddVals = .CallClass(field, ddMembers)
		.checkAndLogTags(ddVals, field, excludeTags, logAsError)

		heading = ddVals.Member?("Heading")
			? ddVals.Heading
			: ddVals.Member?("Prompt")
				? ddVals.Prompt
				: field

		if heading is field or heading is ''
			.logError('no Heading for: ', field, logAsError)
		return TranslateLanguage(heading)
		}

	SelectPrompt(field, excludeTags = #(Internal, ExcludeSelect), logAsError = true)
		{
		ddMembers = Object('SelectPrompt', 'Prompt', 'Heading').Append(excludeTags)
		ddVals = .CallClass(field, ddMembers)
		.checkAndLogTags(ddVals, field, excludeTags, logAsError)

		selectPrompt = ddVals.Member?("SelectPrompt")
			? ddVals.SelectPrompt
			: ddVals.Member?('Prompt') and ddVals.Prompt isnt ""
				? ddVals.Prompt
				: ddVals.Member?('Heading') and ddVals.Heading isnt ""
					? ddVals.Heading
					: field

		if selectPrompt is field or selectPrompt is ''
			.logError('no SelectPrompt for: ', field, logAsError)
		return selectPrompt
		}

	// The purpose of GetFieldPrompt is to handle duplicate custom field prompts.
	GetFieldPrompt(field, promptList = #(), excludeTags = #(Internal), logAsError = true)
		{
		if field.Prefix?('custom_')
			{
			prompt = .Prompt(field, :excludeTags, :logAsError)
			if not promptList.Has?(prompt)
				return prompt
			}
		return .SelectPrompt(field, :excludeTags, :logAsError)
		}

	GetPromptMap(columns, logAsError = true)
		{
		promptList = Object()
		for col in columns.Copy()
			{
			prompt = .GetFieldPrompt(col, :promptList, :logAsError)
			if promptList.Has?(prompt)
				{
				duplicateField = promptList.Find(prompt)
				if duplicateField.Prefix?('custom_')
					{
					selectPrompt = .SelectPrompt(duplicateField, :logAsError)
					if selectPrompt isnt prompt
						{
						promptList[duplicateField] = selectPrompt
						promptList[col] = prompt
						continue
						}
					}
				if not (col.Prefix?('custom_') and duplicateField.Prefix?('custom_'))
					.suneidologOnce('ERROR: Duplicate Prompt: ' $ col $ ' & ' $
						promptList.Find(prompt))
				}
			promptList[col] = prompt
			}
		return promptList
		}

	suneidologOnce(msg)
		{
		SuneidoLog.Once(msg, calls:)
		}
	}
