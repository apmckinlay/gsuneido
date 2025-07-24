// Copyright (C) 2014 Axon Development Corporation All rights reserved worldwide.
class
	{
	Secret()
		{
		return .getKey('amazon_secret_key')
		}
	Access()
		{
		return .getKey('amazon_access_key')
		}

	getKey(field)
		{
		keyRec = .getKeyRec()
		magic = keyRec[field]
		if magic is ''
			throw field $ ' must be defined in ' $ .table $ ' table'
		return StringXor(Base64.Decode(magic), keyRec.amazon_xor_key)
		}

	table: 'amazon_access_keys'
	getKeyRec()
		{
		Database('ensure ' $ .table $
			'(amazon_access_key, amazon_secret_key, amazon_xor_key) key()')
		keyRec = Query1Cached(.table)
		return keyRec is false ? [] : keyRec
		}
	}
