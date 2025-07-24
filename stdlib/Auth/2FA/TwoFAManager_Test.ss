// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		.SpyOn(Sys.Client?).Return(false)

		ob = Object(tokens: Object(), expires: Object(), inRemove?: false)
		.SpyOn(TwoFAManager.TwoFAManager_getVars).Return(ob)

		user = .TempName()
		sessionId = .TempName()

		otp = TwoFAManager.Generate(user, sessionId)
		Assert(otp isnt: false)
		Assert(TwoFAManager.Auth(.TempName(), sessionId, otp) is: false)
		Assert(TwoFAManager.Auth(user, .TempName(), otp) is: false)
		Assert(TwoFAManager.Auth(user, sessionId, .TempName()) is: false)
		Assert(TwoFAManager.Auth(user, sessionId, otp), msg: 'auth generate')

		TwoFAManager.Invalidate(otp)
		Assert(TwoFAManager.Auth(user, sessionId, otp) is: false)
		Assert(ob.tokens[otp].invalid, msg: 'invalid')
		ob.tokens.Delete(all:)
		ob.expires.Delete(all:)

		.SpyOn(TwoFAManager.TwoFAManager_expire).Return(
			Date().Plus(minutes: -10), Date().Plus(minutes: 10))
		otp = TwoFAManager.Generate(user, sessionId)
		Assert(TwoFAManager.Auth(user, sessionId, otp) is: false)
		Assert(ob.tokens isSize: 0) // expired is removed

		.SpyOn(TwoFAManager.GenerateOTP).Return('000000')
		spy = .SpyOn(TwoFAManager.TwoFAManager_log).Return('')
		Assert(TwoFAManager.Generate(user, sessionId) is: '000000')
		Assert(TwoFAManager.Generate(user, sessionId) is: false)
		Assert(spy.CallLogs() is: [[msg: 'Generate OTP failed (.tokens size: 1)',
			params: [:user, :sessionId]]])
		}
	Test_GenerateOTP()
		{
		pass = TwoFAManager.GenerateOTP(tmpPass?: false)
		Assert(pass, isSize: 6)
		Assert(pass.Numeric?(), msg: 'is number')

		pass = TwoFAManager.GenerateOTP(tmpPass?: true)
		Assert(pass, isSize: 8)
		Assert(not pass.Numeric?(), msg: 'no number')
		}
	}