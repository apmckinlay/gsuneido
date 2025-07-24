// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(field, fieldMap, value, hwnd = 0, data = false,
		exclude_custom? = false, _customizable = false)
		{
		source = customizable isnt false ? customizable.GetName() : false
		if value isnt "" and
			"" isnt msg = .PromptInUse(value, fieldMap[field], exclude_custom?, source)
			{
			Alert(msg, "Prompt In Use", hwnd, MB.ICONERROR)
			data[field] = ""
			}
		}
	exclCustomFieldsWhere: ' and not (name >= "Field_custom_000000"
		and name < "Field_custom_999999"
		and name =~ "^Field_custom_[0-9]+$")'

	PromptInUse(prompt, field, exclude_custom?, source = false)
		{
		for lib in .libraries()
			{
			where = ' where name > "Field_a"
				and name < "Field_zzz"
				and name isnt "Field_' $ field $ '"'
			if lib is 'configlib' and exclude_custom?
				where $= .exclCustomFieldsWhere
			QueryApply(lib $ where $ " and group is -1")
				{ |x|
				calcFld? = .calcFld?(x.name)
				types = calcFld?
					? #('Heading')
					: #('Prompt', 'SelectPrompt')

				for type in types
					{
					dbPrompt = .prompt(x.text, type)
					if not .checkMatch(dbPrompt, prompt)
						continue

					if not calcFld? or
						(source isnt false and .rptFldMatch?(prompt, source))
						return "Prompt: " $ dbPrompt $ " is already in use by the system."
					}
				}
			}
		return ""
		}

	calcFld?(field)
		{
		return false isnt field.RemovePrefix('Field_').Match(ReporterModel.Calc_prefix)
		}

	rptFldMatch?(prompt, source)
		{
		QueryApply('params where report =~ "Reporter -"')
			{ |x|
			if not x.params.Member?('formulas') or not Object?(x.params.formulas)
				continue

			if not x.params.formulas.Any?({ .checkMatch(it.calc, prompt) })
				continue

			find = false
			Plugins().ForeachContribution('Reporter', 'queries')
				{
				if x.params.Source is it.name
					{
					find = it
					break
					}
				}

			if find isnt false and find.tables.Has?(source)
				return true
			}

		return false
		}

	checkMatch(dbPrompt, prompt)
		{
		if dbPrompt is false
			return false
		return dbPrompt.Lower() is prompt.Lower()
		}

	//extracted for test
	libraries()
		{
		return Libraries()
		}

	prompt(text, type)
		{
		if text.Size() is pos = text.FindRx('\<' $ type $ ':\>')
			return false

		scan = Scanner(text[pos+type.Size()+1..])
		do
			{
			tok = scan.Next2()
			}
		while tok in (#WHITESPACE, #NEWLINE, #COMMENT)
		return tok is #STRING ? scan.Value() : false
		}
	}
