// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(src)
		{
		switch Type(src)
			{
		case 'Object': dst = Object()
		case 'Record': dst = Record()
			}
		Suneido.temp = dst // make it concurrent
		Suneido.Delete('temp')
		.handle(src, dst)
		return dst
		}

	handle(src, dst)
		{
		for m in src.Members()
			try
				dst[m] = src[m]
			catch (unused, "*cannot be set to concur")
				{
				switch Type(src[m])
					{
				case 'Object':
					.handle(src[m], dst[m] = Object()) // recursive
				case 'Record':
					.handle(src[m], dst[m] = Record()) // recursive
				default:
					dst[m] = Display(src[m])
					}
				}
		}
	}
