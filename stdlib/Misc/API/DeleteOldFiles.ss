// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
class
	{
	failureLimit: 5
	CallClass(filePath, days = -14, logPrefix = 'ERROR')
		{
		Assert(days lessThan: 0)
		cutOff = Date().Plus(:days)
		failures = 0
		Dir(filePath $ '*', details:)
			{ |file|
			if file.date < cutOff
				{
				try
					.deleteItem(file, filePath)
				catch (err)
					{
					++failures
					SuneidoLog(logPrefix $ ': (CAUGHT) DeleteOldFiles: ' $
							'error deleting item: ' $ err,
						calls:, params: Object(:file, :filePath),
						caughtMsg: 'will continue deleting other items')
					if failures > .failureLimit
						throw "DeleteOldFiles: failed deletes exceeded limit"
					}
				}
			}
		}

	deleteItem(file, filePath)
		{
		pathToDel = Paths.Combine(filePath, file.name)
		if file.name.Suffix?('/')
			DeleteDir(pathToDel)
		else
			DeleteFile(pathToDel)
		}
	}