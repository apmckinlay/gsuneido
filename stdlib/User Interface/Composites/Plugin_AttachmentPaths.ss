// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
#(
ExtensionPoints:
	(
	('subfolder')
	)
Contributions:
	(
	(AttachmentPaths, subfolder,
		getPath: function(subfolder)
			{
			return Paths.ToStd(subfolder).BeforeFirst("/") is subfolder
				? OpenImageSettings.Copyto()
				: false
			}
		)
	)
)
