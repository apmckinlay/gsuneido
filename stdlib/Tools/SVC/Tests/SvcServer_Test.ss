// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_stateless()
		{
		_request = Object()
		svc = SvcServer
			{
			PLACEHOLDER(){ /* Fake method */ }
			}
		inst = new svc
		inst.SvcServer_state = 'Verify'
		fn = inst.SvcServer_stateless

		cmd = 'InvalidMethod'
		Assert({ fn([cmd]) }
			throws: 'connection must be verified prior to command: ' $ cmd)

		inst.SvcServer_state = 'Open'
		Assert({ fn([cmd]) } throws: 'invalid command: ' $ cmd)

		cmd = 'GET'
		Assert({ fn([cmd]) } throws: cmd $ ': missing argument')

		cmd = 'PLACEHOLDER'
		fn([cmd]) // Shouldn't throw

		Assert({ fn(cmd) } throws: 'invalid request, arguments are not in an Object')
		}
	}