// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// NOTE: returns a binary string, you may need ToHex or Base64Encode
function (message, key)
	{
	return BuildHmac(message, key, Sha1)
	}