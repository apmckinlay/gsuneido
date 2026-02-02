// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
// NOTE: separate record, so EmailAttachment won't be referenced on suneido.js
function ()
	{
	// Amazon SES limit is 10mb
	// but we use 7 here because mime base64 encoding increases the size by 3:4
	return 7 /* = max size in mb */
	}