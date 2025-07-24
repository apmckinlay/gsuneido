// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
class
	{
	MakeDD(field, prompt, baseType)
		{
		name = "Field_" $ field
		base = "Field_" $ baseType

		if .getConfigRec(name) is false
			.outputField(name, prompt, base)
		else
			{
			try
				dd = Global("Field_" $ field)
			catch
				dd = Field_string
			if not dd.Base?(Global(base)) or dd.GetDefault('Heading', '') isnt prompt
				{
				QueryDo('delete configlib where name is ' $ Display(name))
				.outputField(name, prompt, base)
				}
			}
		return name
		}

	getConfigRec(name)
		{
		if not TableExists?('configlib')
			return false

		return Query1('configlib', :name, group: -1)
		}

	// promptInfo should be like [:baseField, :promptMethod, :prefix, :suffix]
	BuildParam(field, param_field_suffix, promptInfo, checkOnly = false)
		{
		paramField = field.RemoveSuffix("?") $ param_field_suffix
		paramRecName = "Field_" $ paramField

		cur = .getConfigRec(paramRecName)
		// field has been defined in code, don't overwrite it in configlib
		if paramRecName.GlobalName?() and not Uninit?(paramRecName) and cur is false
			return paramField

		if checkOnly is true
			return paramField

		if field.BeforeFirst('_') in ('total', 'min', 'max', 'average')
			field = field.AfterFirst('_')
		base = 'Field_' $ field
		if Uninit?(base)
			base = 'Field_string'
		rec = [name: paramRecName,
			text: base $ '\n' $
				'\t{' $
				'\n\tPromptInfo: ' $ Display(promptInfo) $ '\n\t}']

		if cur is false
			.outputRecord(rec)
		else if cur.text isnt rec.text
			.outputRecord(rec, delete?:)
		return paramField
		}

	outputRecord(rec, delete? = false)
		{
		Transaction(update:)
			{ |t|
			if delete? is true
				t.QueryDo('delete configlib where name is ' $ Display(rec.name))

			OutputLibraryRecord('configlib', rec, :t)
			}
		}

	outputField(name, prompt, base)
		{
		if Uninit?(base)
			base = 'Field_string'
		.outputRecord(Record(:name,
			text: base $ '\n' $
				'\t{' $
				'\n\tHeading: ' $ Display(prompt) $
				'\n\t}'))
		}

	GetCalcDD(type)
		{
		if type.Prefix?("Number")
			return Field_number
		else if type.Has?("Date")
			return Field_date
		else if type.Prefix?("Checkmark")
			return Field_boolean_checkmark
		return false
		}

	GetCalcDef(type, c)
		{
		if type.Prefix?("Number")
			{
			fmt = CustomFieldTypes.GetFormat(c)
			return "Field_number { Format: " $ Display(fmt) $ " }"
			}
		else if type.Has?("Date")
			{
			fmt = CustomFieldTypes.GetFormat(c)
			return "Field_date { Format: " $ Display(fmt) $ " }"
			}
		else
			return "Field_string { }"
		}

	GetWidth(sf, fld)
		{
		if not sf.HasPrompt?(fld)
			return Reporter.DefaultColWidth

		dd = Datadict(sf.PromptToField(fld))
		width = Object?(dd.Format) and dd.Format.Member?('width')
			? dd.Format.width : Reporter.DefaultColWidth
		if dd.Format[0] is 'Image'
			{
			maxWidth = 8.5.InchesInTwips()
			width = (dd.Format.width / maxWidth) * Reporter.LandscapeChars
			}
		return width
		}
	}