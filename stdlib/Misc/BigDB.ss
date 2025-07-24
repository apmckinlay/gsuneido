// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if TableExists?('tmp')
		Database("drop tmp")
	Database("create tmp (k, f) key(k)")
	numTempRecords = 40000
	bigText = 'helloworld'.Repeat(10000) /* = repeat to create large text */
	for i in .. numTempRecords
		QueryOutput('tmp', [k: i, f: bigText])
	}