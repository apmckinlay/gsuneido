// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// TODO: anywhere that us calling Custom_PromptInuse can now call this instead
	CallClass(field, fieldMap, tagLabels, value, hwnd = 0, data = false,
		exclude_custom? = false, _customizable = false)
		{
		if data is false
			data = Object()
		source = customizable isnt false ? customizable.GetName() : false
		if value isnt "" and
			"" isnt msg = .PromptInUse(tagLabels, value, fieldMap[field],
				exclude_custom?, source)
			{
			Alert(msg, tagLabels[0].Tr('_', ' ') $ " In Use", hwnd, MB.ICONERROR)
			data[field] = ""
			return false
			}
		return true
		}
	exclCustomFieldsWhere: ' and not (name >= "Field_custom_000000"
		and name < "Field_custom_999999"
		and name =~ "^Field_custom_[0-9]+$")'

	PromptInUse(tagLabels, value, field, exclude_custom?, source = false)
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
					? #('Heading')  // TODO: handle calcFieldTags other than heading
					: tagLabels

				for type in types
					{
					dbPrompt = .prompt(x.text, type)
					if not .checkMatch(dbPrompt, value)
						continue

					if not calcFld? or
						(source isnt false and .rptFldMatch?(value, source))
						return tagLabels[0].Tr('_', ' ') $ ": " $ dbPrompt $
							" is already in use by the system."
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
