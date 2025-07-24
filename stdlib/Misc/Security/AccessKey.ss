// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	MakeKey(nusers, expiry, mac = false, oldKey = false)
		{
		if mac is false
			mac = ServerEval('GetMacAddress')
		Assert(.valid?(nusers, expiry))
		expiry = expiry.Format('yyyyMMdd')
		if oldKey is true
			return (expiry $ nusers).Xor(mac).Base64Encode().RightTrim('=')
		return (expiry $ nusers) $ .getKeyHash(expiry, nusers, mac)
		}

	MaxUsers: 150
	valid?(nusers, expiry)
		{
		return Number?(nusers) and 1 <= nusers and nusers <= .MaxUsers and
			Date?(expiry) and #20060101 < expiry
		}

	invalidKey: #(0, #17000101, mac: false)
	dateDigits: 8
	hashLength: 27
	SplitKey(key, mac = false)
		{
		// spliting will give a error if is not a valid key
		if key.Has?('-') or key is ''
			return .invalidKey

		for mac in .getMacAddresses(mac)
			{
			try
				{
				s = key[.. -.hashLength]
				expiry = Date(s[.. .dateDigits], 'yMd')
				nusers = s[.dateDigits ..]
				if .getKeyHash(expiry, nusers, mac) isnt key[-.hashLength..]
					{
					nusers = false
					expiry = false
					}
				}
			catch (err)
				{
				SuneidoLog('ERROR: (CAUGHT) ' $ err $ '. Unable to authenticate ' $
					key $ ' with ' $ Display(mac.ToHex()),
					caughtMsg: 'invalidkey returned')
				nusers = false
				expiry = false
				}

			nusers = Numberable?(nusers) ? Number(nusers) : 0
			if true is .valid?(nusers, expiry)
				return Object(nusers, expiry, :mac)
			}
		return .invalidKey
		}

	IsNewKey?(key)
		{
		return key.Size() > .hashLength and Date?(Date(key[.. .dateDigits], 'yMd'))
		}

	getMacAddresses(mac)
		{
		if mac isnt false and not Object?(mac)
			mac = Object(mac)
		return mac is false ? ServerEval('GetMacAddresses') : mac
		}

	getKeyHash(expiry, nusers, mac)
		{
		if Date?(expiry)
			expiry = expiry.Format('yyyyMMdd')
		return Sha1(expiry $ nusers $ mac.ToHex()).Base64Encode().RightTrim('=')
		}
	}
