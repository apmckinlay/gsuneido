// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
/*
	Usage:		This function is to allow a buffered read -> encode -> write of a file.
				This prevents allocation of large sections of memory
	readFileName:	Name of file containing content to encode
	writable: 	something with a Write method (e.g. file, socket)
	102027:		102027 characters is roughly 99.63 kb (which remains under 100 kb).
				100kb is what PdfReader uses for its cache size.
				We used this number as concatenating encoded strings hinges upon being
				divisible by 3
*/
class
	{
	maxRead: 102027
	lineLength:	71 // divides ((102027 * 4 (encoding size)) / 3 (decoding size)) evenly
	CallClass (readFileName, writable)
		{
		File(readFileName)
			{
			.encode(it, writable)
			}
		}

	encode(readable, writable)
		{
		while false isnt text = readable.Read(.maxRead)
			Base64.Encode(text).MapN(.lineLength, { writable.Write(it $ '\r\n');; })
		}
	}
