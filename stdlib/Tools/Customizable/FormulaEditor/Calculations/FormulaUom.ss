// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	DisplayName: "UOM"
	CallClass(rec)
		{
		.Validate(Object(rec.type))
		uom = Split_UOM(rec.value)
		return Object(type: FORMULATYPE.STRING, value: uom.uom.Trim())
		}

	Validate(@args)
		{
		if args.Size() < 1
			throw "Formula: " $ .DisplayName $ " missing arguments"

		if args.Size() > 1
			throw "Formula: " $ .DisplayName $ " too many arguments"

		if not args[0].Difference(Object(FORMULATYPE.UOM, FORMULATYPE.UOM_RATE)).Empty?()
			throw "Formula: " $ .DisplayName $ " Field must be a <Quantity> or <Rate>"

		return Object(FORMULATYPE.STRING)
		}
	}