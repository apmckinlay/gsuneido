// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Type()
		{
		return FORMULATYPE.DATE
		}

	CallClass(dateStr, fmt)
		{
		.Validate(Object(dateStr.type), Object(fmt.type))
		date = DateControl.ConvertToDate(dateStr.value, convertDateCodes?:,
			format: fmt.value)
		if not Date?(date)
			throw "Formula: " $ .DisplayName $ " cannot convert " $
				Display(dateStr.value) $ " to a Date"
		return Object(type: .Type(), value: date.NoTime())
		}

	DisplayName: 'DATE'
	Validate(@args)
		{
		if args.Size() < 2
			throw "Formula: " $ .DisplayName $ " missing arguments"

		if args.Size() > 2
			throw "Formula: " $ .DisplayName $ " too many arguments"

		if args[0] isnt Object(FORMULATYPE.STRING)
			throw "Formula: " $ .DisplayName $ " Date must be a <String>"

		if args[1] isnt Object(FORMULATYPE.STRING)
			throw "Formula: " $ .DisplayName $ " Format must be a <String>"

		return Object(.Type())
		}
	}
