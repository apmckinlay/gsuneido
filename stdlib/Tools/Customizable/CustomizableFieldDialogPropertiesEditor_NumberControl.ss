// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'peditor' // Must be named like this to be found by CustomizableFieldDialog

	New()
		{
		.decimals = .FindControl('decimals')
		.digits = .FindControl('digits')
		.noFormat = .FindControl('noFormat')
		.example = .FindControl('example')
		.tooltip = .FindControl('tooltip')
		.Defer(.updateExample)
		}
	Controls:	(Border (Vert
					(CheckBox text: 'No Formatting', name: 'noFormat')
					Skip
					(Static 'Number of digits before the decimal')
					Skip
					(Spinner set: 9 rangefrom: 2, rangeto: 9 name: 'digits')
					Skip
					(Static 'Number of digits after the decimal')
					Skip
					(Spinner set: 2 rangeto: 5 name: 'decimals' )
					Skip
					(Pair (Static 'e.g.')
						(Field width: 12, justify: 'RIGHT', readonly:, name: 'example'))
					Skip
					(Static 'Tooltip')
					Skip
					(Field name: 'tooltip')
					Skip
				))
	maxDigits: 9
	maxDecimals: 5
	Valid?()
		{
		decimals = .decimals.Get()
		digits = .digits.Get()
		return Number?(decimals) and 0 <= decimals and decimals <= .maxDecimals and
			Number?(digits) and 2 <= digits and digits <= .maxDigits
		}

	noFormat: false
	NewValue(value, source)
		{
		.updateExample()
		if source is .noFormat
			{
			.digits.SetReadOnly(value)
			.decimals.SetReadOnly(value)
			}
		.Send(#NewValue, value, source)
		}
	example: false
	updateExample()
		{
		if .example isnt false
			{
			if .noFormat.Get() is true
				.example.Set('999999999.9999')
			else
				.example.Set(.getMask()[1..].Tr('#', '9'))
			}
		}

	// should return an object with options i.e. (list:('a' 'b') width:50 )
	Get()
		{
		m = .getMask()
		status = .tooltip.Get()
		format = m is false
			? Object(mask: m, width: 8)
			: Object(mask: m)
		return Object(control: Object(mask: m, :status), :format)
		}
	getMask()
		{
		if .noFormat.Get() is true
			return false

		dp = .decimals.Get()
		digits = .digits.Get()

		m = .getDigitMask(digits)

		if dp > 0
			m $= '.' $ '#'.Repeat(dp)
		return m
		}

	getDigitMask(digits)
		{
		mod = digits % 3
		m = '#'.Repeat(mod) $ ',###'.Repeat((digits / 3).Floor())
		if mod is 0
			m = m[1 ..]
		m = '-' $ m
		return m
		}

	Set(object)
		{
		mask = object.Control_mask
		if mask is false
			{
			.noFormat.Set(true)
			.digits.SetReadOnly(true)
			.decimals.SetReadOnly(true)
			}
		else
			{
			.noFormat.Set(false)
			decimals = mask.AfterFirst('.').Size()
			.decimals.Set(decimals)

			digits = mask.BeforeFirst('.').Tr(',').Size() - 1 // "- 1" for minus sign
			.digits.Set(digits)
			}
		.tooltip.Set(object.GetDefault('Control_status', ''))
		}
}
