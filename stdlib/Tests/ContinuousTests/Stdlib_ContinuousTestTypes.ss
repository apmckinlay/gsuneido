// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
#(
	(type: "Using just stdlib"
		libs: #(),
		testGroup: 'Standalone',
		skipBookCheck?:)

	(type: "Client Server With just stdlib"
		libs: #(),
		testGroup: "Client/Server",
		skipBookCheck?:,
		currency: '',
		goFile: 'server.go')
)
