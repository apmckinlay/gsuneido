// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Ceiling()
		{
		i = .Int()
		return this > i ? i + 1 : i
		}

	PercentToDecimal()
		{
		return (this * .01) /* = convert % to decimal */
		}

	DecimalToPercent(round = 0)
		{
		return (this * 100).Round(round) /* = convert decimal to % */
		}

	Floor()
		{
		i = .Int()
		return this < i ? i - 1 : i
		}
	Abs()
		{
		return this < 0 ? -this : this
		}

	ToWords()
		{
		written = #('ZERO', 'ONE', 'TWO', 'THREE', 'FOUR',
			'FIVE', 'SIX', 'SEVEN', 'EIGHT', 'NINE',
			'TEN', 'ELEVEN', 'TWELVE', 'THIRTEEN',
			'FOURTEEN', 'FIFTEEN', 'SIXTEEN', 'SEVENTEEN',
			'EIGHTEEN', 'NINETEEN', 'TWENTY', 30: 'THIRTY',
			40: 'FORTY', 50: 'FIFTY', 60: 'SIXTY',
			70: 'SEVENTY', 80: 'EIGHTY', 90: 'NINETY')

		result = function (num, divisable, desc)
			{
			joinwith = divisable is 100 ? 'AND ' : ""
			result = (num / divisable).Int()
			remainder = num % divisable
			return result.ToWords() $ desc $
				(remainder > 0 ? " " $ joinwith $ remainder.ToWords() : "")
			}

		num = .Int()
		if num <= 20
			return written[num]
		else if num < 100
			{
			remainder = num % 10
			return written[num - remainder] $
				(remainder > 0 ? " " $ written[remainder] : "")
			}
		else if num < 1000
			return result(num, 100, ' HUNDRED')
		else if num < 1000000
			return result(num, 1000, ' THOUSAND')
		else
			return result(num, 1000000, ' MILLION')
		}
	EnFrancais()
		{
		return LanguageFrench.NumberToWords(this)
		}
	ToWordsFrench()
		{
		return LanguageFrench.NumberToWords(this)
		}
	ToWordsDutch()
		{
		return LanguageDutch.NumberToWords(this)
		}
	ToWordSpanish()
		{
		return LanguageSpanish.NumberToWords(this)
		}
	ToWordsItalian()
		{
		return LanguageItalian.NumberToWords(this)
		}
	ToWordsSimple()
		{
		/* contributed by Johan Samyn
		Returns a string (in the current language), delimited by "* " and " *",
		summing up every digit in the number, from left to right,
		seperated by a space, and mentioning the decimal sign if applicable.
		The method also accepts negative numbers, and then returns a result
		starting with '* minus '.
		This simple notation is accepted in Belgium on official documents
		as a transformation in words of a number (e.g. invoice, check, ...).
		However, this method is ready to be used in any language available in
		translatelanguage.
		It can handle any number Suneido can handle, just like ToWordsDutch().
		*/
		words = #('zero', 'one', 'two', 'three', 'four',
			'five', 'six', 'seven', 'eight', 'nine')
		decimal = TranslateLanguage('decimalpoint')
		sNum = String(this.Abs())
		if sNum[0] is '.'
			sNum = '0' $ sNum
		sResult = "*"
		if this < 0
			sResult $= " " $ TranslateLanguage('minus')
		for (i = 0; i < sNum.Size(); ++i)
			sResult $= " " $ (sNum[i] is '.' ? decimal :
				TranslateLanguage(words[Number(sNum[i])]))
		return sResult $ " *"
		}

	Pad(minSize, char = "0")
		{
		Assert(char.Size() is 1)
		str = String(this.Abs()) // ignores sign
		n = minSize - str.Size()
		return (n > 0) ? char.Repeat(n) $ str : str
		}
	Factorial()
		{
		if (this > 0)
			return this * (this - 1).Factorial()
		else if (this is 0)
			return 1
		else
			throw "Factorial: can't handle negative values"
		}
	Int?()
		{
		return this is .Int()
		}
	IntDigits()
		{
		num = .Abs().Int()
		return num is 0 ? 0 : String(num).Size()
		}
	FracDigits()
		{
		return String(.Abs().Frac()).Size() - 1
		}
	EuroFormat(mask)
		{
		return .Format(mask.Tr(',.', '.,')).Tr(',.', '.,')
		}
	DollarFormat(mask)
		{
		if String?(mask) and mask.Prefix?('-')
			mask = '(' $ mask.Replace('^-', '') $ ')'
		return '$' $ .Format(mask)
		}
	Even?()
		{ return this % 2 is 0 }
	Odd?()
		{ return this % 2 isnt 0 }
	Sign()
		{ return this < 0 ? -1 : 1 }
	RoundToPrecision(p)
		{
		if this is 0 or this is (1/0) or this is (-1/0)
			return this
		Assert(p greaterThan: 0)
		return .Round(-(.Abs().Log10().Floor() - (p - 1)))
		}
	RoundToNearest(nearest = 1)
		{
		number = this.Int()
		if nearest <= 0
			return number
		rem = number % nearest
		result = rem < nearest / 2
			? number - rem
			: number + (nearest - rem)
		return result.Int()
		}
	ToRGB()
		{
		r =  this & 0xff
		g = (this & 0xff00) >> 8
		b = (this & 0xff0000) >> 16
		return Object(r, g, b)
		}
	Times(block)
		{
		for (i = 0; i < this; ++i)
			block()
		}
	MinutesInMs()
		{
		return this * 60 * 1000
		}
	SecondsInMs()
		{
		return this * 1000
		}
	SecondsInHours()
		{
		return (this / 3600).Round(1)
		}
	HoursInSeconds()
		{
		return (this * 3600).Round(0)
		}
	HoursInMinutes()
		{
		return (this * 60).Round(0)
		}
	SecondsInMinutes()
		{
		return (this / 60).Round(0)
		}
	InchesInTwips()
		{
		return this * 1440 /* = twips in an inch */
		}
	TwipsInInch()
		{
		return this / 1440 /* = inch in twips */
		}
	InchesInCanvasUnit()
		{
		return this * 1440 / 17 /*= 1440 is twips in an inch,
			15 twips is 1 pixel in 96 DPI,
			not sure why we use 17, but it has been used in the calculatioin*/
		}
	Kb()
		{
		byteFactor = 1024.0
		return this * byteFactor
		}
	Mb()
		{
		byteFactor = 1024.0
		return this * byteFactor * byteFactor
		}
	Gb()
		{
		byteFactor = 1024.0
		return this * byteFactor * byteFactor *	byteFactor
		}
	SafeEval()
		{
		return this
		}
	Of(block)
		{
		if not Function?(block)
			{
			ob = Object()
			for ..this
				ob.Add(block)
			return ob
			}
		return Nof(this, block)
		}
	OfStr(block)
		{
		s = ""
		for ..this
			s $= block()
		return s
		}
	}