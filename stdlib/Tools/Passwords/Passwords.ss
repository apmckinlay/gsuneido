// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	maxLoopAttempts: 50
	Generate(pwLength, option)
		{
		loop = 0
		while '' isnt .validateStrength(pass = .genPassword(pwLength, option))
			if ++loop > .maxLoopAttempts
				{
				pass = .lastTryPassword(pwLength)
				break
				}

		return pass
		}

	lastTryPassword(pwLength)
		{
		return .GenerateSimple(pwLength)
		}

	allowedRanges:
		#('Easy to read': 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ',
		'Easy to spell': 'abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789~!@#$%^' $
			'&*(){}[];:,<>./?\\=-_'
		'All Characters': 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567' $
			'89`~!@#$%^&*()-_=[]{}\\|;:\'",.<>/? '
		'Pass Phrase': '')

	genPassword(pwLength, option)
		{
		allowedRange = .allowedRanges[option]
		if allowedRange.Blank?()
			return .getPassPhrase(pwLength)
		return pwLength.Of({ allowedRange.RandChar() }).Join()
		}

	getPassPhrase(wordCount)
		{
		return wordCount.Of({ PasswordDictList.RandVal().Capitalize() }).Join()
		}

	ValidateStrength(pw)
		{
		return .validateStrength(Base64.Decode(pw))
		}

	validateStrength(pw)
		{
		if pw is ''
			return ''

		if '' isnt msg = .passwordComplexity(pw)
			return msg
		return .pwnedPassword?(pw)
		}

	charCount: 26
	numCount: 10
	symbolCount: 33
	/// 8 length, lower+upper+symbol = min allowed complexity
	minAllowedComplexity: 7e11
	passwordComplexity(pw)
		{
		uniqueTypes = .getUniqueTypes(pw)
		baseSet = (uniqueTypes.cap ? .charCount : 0) +
			(uniqueTypes.lower ? .charCount : 0) +
			(uniqueTypes.number ? .numCount : 0) +
			(uniqueTypes.symbol ? .symbolCount : 0)
		possibilities = baseSet.Pow(pw.Size())

		return .minAllowedComplexity <= possibilities
			? ''
			: 'This password is not secure.\n'$
				'Please make it longer or include more types of characters\n' $
				'(lower case, upper case, numbers, symbols)'
		}

	getUniqueTypes(pw)
		{
		uniqueTypes = Object(cap: false, lower: false, number: false, symbol: false)
		for char in pw
			{
			if char.Upper?()
				uniqueTypes.cap = true
			else if char.Lower?()
				uniqueTypes.lower = true
			else if char.Numeric?()
				uniqueTypes.number = true
			else if char.AlphaNum?() is false
				uniqueTypes.symbol = true
			}
		return uniqueTypes
		}
	GetUniqueTypes(pw)
		{
		return .getUniqueTypes(pw)
		}

	prefixLength: 5
	pwnedPassword?(pw)
		{
		hashed = Sha1(pw).ToHex()
		suffix = hashed[.prefixLength ..].Upper()
		prefix = hashed[.. .prefixLength]
		try
			content = .getPwnedList(prefix)
		catch (err)
			{
			.logError('connection err: ' $ err $ '; check for one time error')
			return ''
			}
		if content.Blank?()
			{
			.logError('returned empty content list, unable to compare for prefix: ' $
				prefix $ '; check for problem with this prefix')
			return ''
			}

		suffixOb = content.Lines().Map({ it.BeforeFirst(':') })
		return suffixOb.HasIf?({ it is suffix })
			? 'This password is not secure.\nPlease try a different password.'
			: ''
		}

	// extracted for tests
	getPwnedList(prefix)
		{
		return Https.Get('https://api.pwnedpasswords.com/range/' $ prefix,
			timeoutConnect: 10) /*= do not want the users to wait long if the connection
								fails */
		}
	logError(msg)
		{
		ErrorLog('ERROR (CAUGHT): api.pwnedpasswords ' $ msg $
			' (password allowed without pwned checking)')
		}

	symbols: '~!@#$^*/%&'
	numbers: '0123456789'
	GenerateSimple(pwLength = 15, symbols = 2, numbers = 2)
		{
		letters = Max(0, pwLength - symbols - numbers)
		chars = .nof(symbols, .symbols) $ .nof(numbers, .numbers) $
			.nof(letters, .allowedRanges['Easy to read'])
		return chars.Shuffle()
		}
	nof(n, chars)
		{
		return n.OfStr({ chars.RandChar() })
		}
	}