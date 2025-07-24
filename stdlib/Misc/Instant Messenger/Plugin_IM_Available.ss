// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
#(
ExtensionPoints:
	(
	(im_available)
	)
Contributions:
	(
	('IM_Available', 'im_available', func: function ()
		{ return Sys.Client?()})
	)
)