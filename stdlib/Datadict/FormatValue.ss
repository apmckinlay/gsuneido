// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(value, field)
		{
		dd = Datadict(field)
		if Boolean?(value) or dd.Base?(Field_boolean)
			return .formatBoolean(value)
		// check for '' AFTER checking boolean since '' treated as false for boolean
		if value is ''
			return ''
		if Object?(value)
			return Display(value)
		if false isnt fmtValue = .FormatDataToString(dd, value)
			return fmtValue
		if Date?(value)
			return .fmtDate(value)
		return .formatBasedOnDD(dd, field, value)
		}

	formatBoolean(value)
		{
		return value is true ? 'yes' : 'no'
		}

	FormatDataToString(dd, value)
		{
		try
			{
			fmt = Report.Construct(dd.Format)
			if fmt isnt false and fmt.Method?('DataToString')
				return fmt.DataToString(value, [])
			}
		return false
		}

	fmtDate(value)
		{
		return value is value.NoTime() ? value.ShortDate() : value.ShortDateTime()
		}

	formatBasedOnDD(dd, field, value)
		{
		if dd.Base?(Field_dollars) or dd.Base?(Field_dollar)
			return .fmtDollars(field, value)
		if dd.Base?(Field_scintilla_rich)
			return ScintillaRichStripHTML(value)
		if dd.Base?(Field_info)
			return StripInfoLabel(value)
		return value
		}

	fmtDollars(field, value)
		{
		mask = OptionalNumberFormat.EvalMask(Datadict(field).Format.mask)
		return value is 0 ? '$0.00' : '$' $ value.Format(mask)
		}
	}
