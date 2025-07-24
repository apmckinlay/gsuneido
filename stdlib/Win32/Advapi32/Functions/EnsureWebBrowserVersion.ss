// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// REF: https://msdn.microsoft.com/en-us/library/ee330730%28v=vs.85%29.aspx
	IEVersion: 11001 // IE 11
	CallClass()
		{
		try
			{
			phkResult = Object()
			if 0 isnt err = .open(phkResult)
				{
				if err isnt .notFound
					return .logErr('open', err)
				if 0 isnt createErr = .createKey(phkResult = Object())
					return .logErr('create', createErr)
				}
			res = .ensureIEVersion(phkResult.x)
			RegCloseKey(phkResult.x)
			Assert(res)
			}
		catch(e)
			{
			.logErr('ensure', e)
			}
		}

	browserEmulationKey: `SOFTWARE\Microsoft\Internet Explorer` $
		`\Main\FeatureControl\FEATURE_BROWSER_EMULATION`
	open(phkResult)
		{
		return RegOpenKeyEx(
			REG_HKEY.CURRENT_USER,
			.browserEmulationKey,
			0,
			REG_KEY.QUERY_VALUE | REG_KEY.SET_VALUE,
			phkResult)
		}

	createKey(phkResult)
		{
		createErr = .createKeyEx(phkResult)
		if createErr is 0 or createErr isnt 1021 /*= parent key is volatile */
			return createErr

		return .createKeyEx(phkResult, volatile:)
		}

	createKeyEx(phkResult, volatile = false)
		{
		return RegCreateKeyEx(
			REG_HKEY.CURRENT_USER,
			.browserEmulationKey,
			reserved: 0,
			lpClass: 0,
			dwOptions: volatile ? 1 : 0,
			samDesired: REG_KEY.ALL_ACCESS,
			lpSecurityAttributes: 0,
			:phkResult,
			lpdwDisposition: 0
			)
		}

	notFound: 2
	ensureIEVersion(hKey)
		{
		exeName = ExeName()
		result = .retrieveVersion(hKey, exeName, lpData = Object())
		if result isnt 0 and result isnt .notFound
			return .logErr('query', result)

		if result is .notFound or lpData.x isnt .IEVersion
			if 0 isnt err = .setVersion(hKey, exeName)
				return .logErr('set', err)

		return true
		}

	retrieveVersion(hKey, exeName, lpData)
		{
		return RegQueryValueEx(
			hKey,
			exeName,
			0,
			Object()
			lpData,
			lpcbData: Object(x: 4))
		}

	setVersion(hKey, exeName)
		{
		RegSetValueEx(hKey,
			exeName,
			0,
			REG.DWORD,
			Object(x: .IEVersion),
			cbData: 4)
		}

	logErr(action, err)
		{
		address = LoginSessionAddress()
		ErrorLog('ERROR: EnsureWebBrowserVersion unable to ' $
			action $ ' registry from ' $ address $ ' - ' $ String(err))
		return false
		}
	}
