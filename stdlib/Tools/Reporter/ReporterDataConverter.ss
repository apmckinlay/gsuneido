// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	ToPrompts1(data, sf)
		{
		if data.GetDefault(#usingFieldsInSave?, false) isnt true
			return

		map = .prepairMap(sf)
		converter = { |field| map.GetDefault(field, field) }
		.convertSummarize(data, converter)
		}

	ToPrompts2(data, sf, _failedConverts = false)
		{
		if data.GetDefault(#usingFieldsInSave?, false) isnt true
			return

		map = .prepairMap(sf)
		converter = {
			|field|
			result = map.GetDefault(field, field)
			if result is field and Object?(failedConverts)
				failedConverts.AddUnique(field)
			result
			}
		.convertColOptions(data, converter, .convertColOptionsHeadingToPrompts)
		.convertColumns(data, converter)
		.convertSort(data, converter)
		.convertSelect(data, converter)
		// only formula converts need to be handled
		if Object?(failedConverts) and not failedConverts.Empty?()
			failedConverts.RemoveIf({ not it.Blank?() })
		.convertFormulasToPrompts(data, converter)
		data.Delete(#usingFieldsInSave?)
		}

	prepairMap(sf)
		{
		fields = sf.Fields
		map = Object()
		for prompt, field in fields
			map[field] = prompt
		return map
		}

	ToFields(data, sf)
		{
		if data.GetDefault(#usingFieldsInSave?, false) is true
			return

		converter = { |prompt| .toField(prompt, sf) }

		data.columns = data.columns.DeepCopy()
		if Object?(data.select)
			data.select = data.select.DeepCopy()

		.convertColOptions(data, converter, .convertColOptionsHeadingToFields)
		.convertColumns(data, converter)
		.convertSort(data, converter)
		.convertSelect(data, converter)
		.convertSummarize(data, converter)
		.convertFormulasToFields(data, sf)
		data.usingFieldsInSave? = true
		}

	convertFormulasToPrompts(data, converter)
		{
		for row in data.formulas
			{
			if row.Member?(#fields)
				{
				row.formula = row.formula.Replace('FORMULA_#(\d+)#',
					{ |s| converter(row.fields[Number(s[9/*=prefix*/..-1])]) })
				row.Delete(#fields)
				}
			}
		}

	convertFormulasToFields(data, sf)
		{
		for row in data.formulas
			if not row.Member?(#fields) and not row.formula.Trim().Blank?()
				{
				code = ''
				fields = Object()
				sf.ScanFormula(row.formula,
					{|f| code $= 'FORMULA_#' $ fields.Size() $ '#'; fields.Add(f) },
					{|s| code $= s })
				row.formula = code
				row.fields = fields
				}
		}

	convertColOptions(data, converter, extraConvert)
		{
		coloptions = Object()
		for m in data.coloptions.Members()
			{
			option = data.coloptions[m]
			converted = String?(m) ? converter(m) : m
			extraConvert(option, m, converted)
			coloptions[converted] = option
			}
		data.coloptions = coloptions
		}

	convertColOptionsHeadingToPrompts(option, original/*unused*/, converted)
		{
		if Object?(option) and not option.Member?(#heading)
			option.heading = converted
		}

	convertColOptionsHeadingToFields(option, original, converted/*unused*/)
		{
		if Object?(option) and option.GetDefault(#heading, false) is original
			option.Delete(#heading)
		}

	convertColumns(data, converter)
		{
		if not Object?(data.columns.GetDefault(0, #()))
			{
			data.columns.Map!(converter)
			return
			}
		for column in data.columns
			column.text = converter(column.text)
		}

	convertSort(data, converter)
		{
		for i in .. Reporter.SortRows
			{
			if data['sort' $ i] isnt ''
				data['sort' $ i] = converter(data['sort' $ i])
			}
		}

	convertSelect(data, converter)
		{
		if not Object?(data.select)
			return

		for i in .. Select2Control.Numrows
			{
			fieldStr = 'fieldlist' $ i
			if data.select[fieldStr] isnt ''
				data.select[fieldStr] = converter(data.select[fieldStr])
			}
		}

	convertSummarize(data, converter)
		{
		if data.Member?(#summarize_by)
			data.summarize_by = data.summarize_by.Split(',').Map(converter).Join(',')

		for i in .. Reporter.MaxSummarizeFields
			{
			fieldStr = 'summarize_field' $ i
			if data[fieldStr] isnt ''
				data[fieldStr] = converter(data[fieldStr])
			}
		}

	toField(prompt, sf)
		{
		if false isnt field = sf.PromptToField(prompt)
			return field
		return prompt
		}
	}
