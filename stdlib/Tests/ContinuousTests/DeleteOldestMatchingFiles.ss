// Copyright (C) 2014 Axon Development Corporation All rights reserved worldwide.
function (path, filePattern, filesToKeep = 10)
	{
	Dir(path $ filePattern, files:, details:).
		Sort!(By(#date))[.. -filesToKeep].
		Each({ DeleteFile(path $ it.name) })
	}
