// Copyright (C) 2003 Axon Development Corporation All rights reserved worldwide.
class
	{
	Calc(fromValue, fromUom, toUom, rate? = false)
		{
		if fromUom.Lower() is toUom.Lower()
			return fromValue

		fromType = .GetUomRec(fromUom)
		toType = .GetUomRec(toUom)

		if fromType isnt false and toType isnt false and
			fromType.etauom_type is toType.etauom_type
			return fromValue * (rate? is false
				? toType.etauom_size / fromType.etauom_size
				: fromType.etauom_size / toType.etauom_size)

		return false
		}

	GetUomRec(uom)
		{
		return Uom_Conversions.GetDefault(uom.Lower(), false)
		}

	// add qty1 to qty2
	// expects qty1 and qty2 to be in format: #(value: value, uom: uom)
	// i.e. what Split_UOM returns
	Add(qty1, qty2)
		{
		if qty2 is ''
			return qty1
		 else
			if false is val = .Calc(qty1.value, qty1.uom, qty2.uom)
				return false
		qty2.value += val.Round(2)
		return qty2
		}
	}
