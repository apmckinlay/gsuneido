// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	GetInitSet(user/*unused*/)
		{
		return 'IDE'
		}

	ForgotPassword(user/*unused*/, forgotPassword)
		{
		if forgotPassword is true
			{
			return 'Forgot Password is not supported'
			}
		return ''
		}

	GetUserRec(user)
		{
		return Query1('users', :user)
		}

	TwoFA(userRec, data, hasCookies?/*unused*/, extra = #())
		{
		if '' is email = data.GetDefault(#email, '')
			return 'Missing email address.'

		sessionId = 16.Of({ Random(16/*=hex*/).Hex() }).Join()
		if false is otp = TwoFAManager.Generate(userRec.user, sessionId)
			return 'Generate Multi-factor Authentication code failed'

		from = BookEmailInfo().from
		recipient = userRec.user

		send? = TwoFAManager.SendEmail(otp, from, email, recipient) isnt false
		msg = 'Please enter the login code emailed to ' $ email
		if not send?
			return 'An email with the Login Code was not able to be sent to ' $ email
		return [:sessionId, :msg, :email, :from].Merge(extra)
		}

	AllowAuthFailure?(user/*unused*/)
		{
		return false
		}

	AfterLoginSuccess(user/*unused*/)
		{
		return ''
		}

	PostLogin(@unused)
		{
		return true
		}
	}