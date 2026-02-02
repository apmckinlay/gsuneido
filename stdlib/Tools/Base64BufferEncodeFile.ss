// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
/*
	Usage:		This function is to allow a buffered read -> encode -> write of a file.
				This prevents allocation of large sections of memory
	readFileName:	Name of file containing content to encode
	writable: 	something with a Write method (e.g. file, socket)
*/
class
	{
	/*
	// maxRead characters is roughly 99.82 kb (which remains under 100 kb)
	// 100kb is what PdfReader uses for its cache size

	// Base64 encodes every 3 characters into 4 bytes
	102222 / 3 = 34074  // number of Base64 encoding chunks
	34074 * 4 = 136296  // total bytes after encoding
	136296 / 72 = 1893  // number of lines generated for each read
	// the total bytes (136296) needs to be divisible by 72 (per line)
	// so there is nothing left for the next chunk
	*/

	maxRead: 102222
	// Microsoft Exchange server has issue with line length that is not divisible by 4
	// and could end up with file corrupt
	// should be also consistent with Base64.EncodeLines
	lineLength:	72
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
