// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (filename)
	{
	return throw RetryBool(maxretries: 3, min: 250)
		{
		result = DeleteFileApi(filename)
		if String?(result) and result.Has?("does not exist")
			result = true // no need to retry
		result
		}
	}