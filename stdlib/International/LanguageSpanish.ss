// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// contributed by Luis Alfredo
class
	{
	written: #('cero', 'uno', 'dos', 'tres', 'cuatro',
		'cinco', 'seis', 'siete', 'ocho', 'nueve',
		'diez', 'once', 'doce', 'trece',
		'catorce', 'quince', 'dieciseis', 'diecisiete',
		'dieciocho', 'diecinueve', 'veinte', 'veintiuno',
		'veintidos', 'veintitres', 'veinticuatro', 'veinticinco', 'veintiseis',
		'veintisiete', 'veintiocho', 'veintinueve',
		30: 'treinta', 40: 'cuarenta', 50: 'cincuenta', 60: 'sesenta',
		70: 'setenta', 80: 'ochenta', 90: 'noventa')

	variaciones: #(21: 'veintiun', 31: 'treinta y un' 41: 'cuarenta y un',
		51: 'cincuenta y un', 61: 'sesenta y un', 71: 'setenta y un',
		81: 'ochenta y un', 91: 'noventa y un')

	NumberToWords(number)
		{
		num = number.Int()
		if num <= 30 /*= from map */
			return .written[num]
		else if num < .hundred
			{
			remainder = num % 10 /*= units */
			return .written[num - remainder] $
				(remainder > 0 ? " y " $ .written[remainder] : "")
			}
		else
			return .greaterThan100(num)
		}

	hundred: 100
	greaterThan100(num)
		{
		if num < .thousand
			if num is .hundred
				return 'cien'
			else if (num < 200) /*= speical */
				{
				vnum = num - .hundred
				return 'ciento' $ .result(vnum, .hundred, "")
				}
			else if ((num >= 500) and (num < 600)) /*= speical */
				{
				vnum = num - 500 /*= speical */
				return 'quinientos' $ .result(vnum, .hundred, "")
				}
			else
				return .greaterThan600(num)
		else
			return .greaterThan1000(num)
		}

	greaterThan600(num)
		{
		if ((num >= 700) and (num < 800)) /*= speical */
			{
			vnum = num - 700 /*= speical */
			return 'setecientos' $ .result(vnum, .hundred, "")
			}
		else if (num > 900) /*= speical */
			{
			vnum = num - 900 /*= speical */
			return 'novecientos' $ .result(vnum, .hundred, "")
			}
		else
			return .result(num, .hundred, 'cientos')
		}

	thousand: 1000
	million: 1000000
	greaterThan1000(num)
		{
		if num < .million
			if num is .thousand
				return 'mil'
			else if (num < 2000) /*= speical */
				{
				vnum = num - .thousand
				return 'mil' $ .result(vnum, .thousand, "")
				}
			else
				return .greaterThan2000(num)
		else
			return .greaterThanMillion(num)
		}

	firsts: (21, 31, 41, 51, 61, 71, 81, 91)
	greaterThan2000(num)
		{
		vnum = (num / .thousand).Int()
		remainder = num % .thousand
		vremainder = 0
		if .firsts.Has?(vnum)
			return .variaciones[vnum] $ ' mil' $ .result(remainder, .thousand, "")
		else
			if vnum > .hundred
				vremainder = vnum % .hundred
			if .firsts.Has?(vremainder)
				if vnum - vremainder is .hundred
					return 'ciento ' $ .variaciones[vremainder] $
						' mil' $ .result(remainder, .thousand, "")
				else
					return (vnum - vremainder).ToWordSpanish() $ ' ' $
						.variaciones[vremainder] $ ' mil' $
						.result(remainder, .thousand, "")
			else
				return vnum.ToWordSpanish() $ ' mil' $
					.result(remainder, .thousand, "")
		}

	greaterThanMillion(num)
		{
		if num is .million
			return 'un millon'
		else if num < 2000000 /*= 2 million */
			return .result(num - .million, .million, 'un millon')
		else if num < 1000000000000 /*= speical */
			return .greaterThan2Million(num)
		else
			return "Esa cifra es demasiado alta"
		}
	greaterThan2Million(num)
		{
		vnum = (num / .million).Int()
		remainder = ((num / .million) - vnum) * .million
		vremainder = 0
		if .firsts.Has?(vnum)
			return .variaciones[vnum] $ ' millones' $
				.result(remainder, .million, "")
		else
			{
			if vnum > .hundred
				vremainder = vnum % .hundred
			else if vnum > .thousand
				vremainder = vnum % .thousand
			else if vnum > 10000 /*= speical */
				vremainder = vnum % 10000 /*= last four digits */
			else if vnum > 100000 /*= speical */
				vremainder = vnum % 100000 /*= last five digits */
			if .firsts.Has?(vremainder)
				if vnum-vremainder is .hundred
					return 'ciento ' $ .variaciones[vremainder] $
						' millones' $ .result(remainder, .million, "")
				else
					return (vnum - vremainder).ToWordSpanish() $ ' ' $
						.variaciones[vremainder] $ ' millones' $
						.result(remainder, .million, "")
			else
				return vnum.ToWordSpanish() $ ' millones' $
					.result(remainder, .million, "")
			}
		}

	result(num, divisable, desc)
		{
		result = (num / divisable).Int()
		remainder = num % divisable
		if result < 1
			return desc $ (remainder > 0 ? " " $ remainder.ToWordSpanish() : "")
		else
			return result.ToWordSpanish() $ desc $
				(remainder > 0 ? " " $ remainder.ToWordSpanish() : "")
		}
	}