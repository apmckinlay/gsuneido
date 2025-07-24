// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(name, book = 'imagebook')
		{
		if false is text = GetBookText(name, :book)
			return false
		return 'data:' $
			MimeTypes.GetDefault(name.AfterLast('.'), 'application/octet-stream') $
			';base64,' $ Base64.Encode(text)
		}
	}