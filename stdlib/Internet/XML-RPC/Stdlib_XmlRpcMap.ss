// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function ()
	// returns a list of valid functions for XmlRpc server
	{
	return Object(
		'test.echo': function (@args) { return args }
		'test.throw': function () { throw "test exception" }
		'examples.getStateName': function (i) { return States[i] }
		)
	}