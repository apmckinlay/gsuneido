// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Generate(user, sessionId, tmpPass? = false)
		{
		if Sys.Client?()
			return ServerEval('TwoFAManager.Generate', user, sessionId)

		ob = .getVars()
		.removeExpired(ob)

		for ..10 /*=max retry*/
			{
			otp = .GenerateOTP(tmpPass?)
			.Synchronized()
				{
				if not ob.tokens.Member?(otp)
					{
					ob.tokens[otp] = Object(:user, :sessionId, invalid: false,
						expire: .expire())
					ob.expires.Add(otp)
					return otp
					}
				}
			}
		.log('Generate OTP failed (.tokens size: ' $ ob.tokens.Size() $ ')',
			params: [:user, :sessionId])
		return false
		}

	expire()
		{
		return Date().Plus(minutes: 10)
		}

	SystemSessionId: 'system'
	Auth(user, sessionId, otp)
		{
		if Sys.Client?()
			return ServerEval('TwoFAManager.Auth', user, sessionId, otp)

		ob = .getVars()
		.removeExpired(ob)

		if false is token = ob.tokens.GetDefault(otp, false)
			return false

		return token.invalid is false and token.user is user and
			(token.sessionId is .SystemSessionId or token.sessionId is sessionId)
		}

	Invalidate(otp)
		{
		if Sys.Client?()
			return ServerEval('TwoFAManager.Invalidate', otp)

		ob = .getVars()

		if false is token = ob.tokens.GetDefault(otp, false)
			return

		token.invalid = true
		}

	getVars()
		{
		Suneido.GetInit('TwoFAManager',
			{ Object(tokens: Object(), expires: Object(), inRemove?: false) })
		}

	codeLength: 6 /*=two FA code works with users existing password so it can be shorter*/
	codeChars: '0123456789'
	passLength: 8 /*= seems sufficiently secure for a temporary password with the delay
			protection against brute force attacks in validation */
	passChars: "abcdefghijklmnopqrstuvwxyz"
	GenerateOTP(tmpPass? = false)
		{
		length = tmpPass? ? .passLength : .codeLength
		chars = tmpPass? ? .passChars : .codeChars
		return length.Of({ chars.RandChar() }).Join()
		}

	removeExpired(ob)
		{
		if ob.inRemove? is true
			return

		ob.inRemove? = true
		i = 0
		n = ob.expires.Size()
		now = Date()
		while i < n
			{
			if ob.tokens[ob.expires[i]].expire < now
				ob.tokens.Delete(ob.expires[i])
			else
				break
			i++
			}

		ob.expires = ob.expires[i..]
		ob.inRemove? = false
		}

	log(msg, params = "")
		{
		SuneidoLog('TwoFAManager: ERROR - ' $ msg, calls:, :params)
		}

	SendEmail(otp, from, to, recipient, type = 'code')
		{
		if false is email = GetContributions('TwoFAEmailMessage').GetDefault(type, false)
			return false
		msg = email.message.
			Replace('<recipient>', recipient).
			Replace('<otp>', otp)
		return BookSendEmail(0, from, to, MimeText(msg).Subject(email.subject), quiet?:)
		}

	IsAuthEmail?(subject)
		{
		return GetContributions('TwoFAEmailMessage').Any?({ it.subject is subject })
		}
	}
