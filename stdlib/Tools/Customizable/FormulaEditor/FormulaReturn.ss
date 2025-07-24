// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(code, field, isFormat = false)
		{
		value = ''
		try
			{
			value = code().value
			dd = Datadict(field)
			value = .ProcessValue(value, dd, :isFormat)
			}
		catch(err)
			{
			value = ''
			if err =~ '^Formula: Incompatible unit of measure|^Formula: Invalid Value'
				return value

			msg = "There is a problem calculating " $ Heading(field)
			if err.Prefix?("Formula: ")
				msg $= "\r\n" $ err.AfterFirst("Formula: ")
			SuneidoLog(err, params: Object(:msg, :field, :isFormat, gui: Sys.GUI?()))
			.handleError(msg, field, isFormat)
			}
		return value
		}

	ProcessValue(value, dd, isFormat = false)
		{
		value = .handleStrings(dd, value, isFormat)
		if dd.Method?('Encode')
			value = dd.Encode(value)
		value = .handleNumbers(dd, value, :isFormat)
		value = .handleBoolean(dd, value)
		value = .handleDate(dd, value)
		}

	handleStrings(dd, val, isFormat)
		{
		if not dd.Base?(Field_string)
			return val

		str = FormulaConvertToString(val)
		if isFormat or UnsortableField?('', dd)
			return str

		if str.Size() > FieldControl.MaxCharacters
			throw 'Formula: Value exceeds max size for single line text field'
		return str
		}

	handleNumbers(dd, val, isFormat)
		{
		if not dd.Base?(Field_number)
			return val

		if val is ''
			return 0

		if not Number?(val)
			throw "Formula: Invalid <Number> value: " $ Display(val)

		if IsInf?(val)
			return .handleThrow("Formula: Invalid <Number> value: " $ Display(val),
				isFormat)

		if false is mask = .getMask(dd, isFormat)
			return val

		if '#' is formattedVal = val.Format(mask)
			return .handleThrow('Formula: Invalid <Number> value: ' $ Display(val) $
				'. Maximum digits before decimal is ' $
				dd.Control.mask.BeforeFirst('.').Count('#'), isFormat)

		try
			return Number(formattedVal)
		catch
			throw 'Formula: The format for this field does not support displaying ' $
				Display(val)
		}

	handleThrow(msg, isFormat)
		{
		if isFormat
			return '#'
		throw msg
		}

	getMask(dd, isFormat)
		{
		ddMethod = isFormat ? dd.Format : dd.Control
		return ddMethod.GetDefault('mask', false)
		}

	handleBoolean(dd, value)
		{
		if not dd.Base?(Field_boolean)
			return value
		if value isnt '' and not Boolean?(value)
			throw "Formula: Invalid <Boolean> value: " $ Display(value)
		return value is true
		}

	handleDate(dd, value)
		{
		if not dd.Base?(Field_date) or value is ''
			return value
		if not Date?(value)
			throw "Formula: Invalid <Date> value: " $ Display(value)
		if dd.Base?(Field_date_time) or dd.Base?(Field_num)
			return value
		return value.NoTime()
		}

	handleError(msg, field, isFormat, _showFormulaError = false)
		{
		if isFormat or showFormulaError
			throw "SHOW: " $ msg
		else if Sys.GUI?()
			AlertDelayed(msg, "Formula Error", uniqueId: field)
		else
			SuneidoLog('Problem Calculating Formula: ' $ msg, params: Object(:field))
		}
	}
