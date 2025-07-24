// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
#(
	ExtensionPoints:
		(
		(event) 		// source events
		(info_type)  	// Information Type from events
		(action) 		// perform action
		(permission)	// permission to change trigger
		)

	Contributions:
		(
		(ECA, info_type, name: 'Subject', type: 'string')
		(ECA, info_type, name: 'Message', type: 'string')
		)
)