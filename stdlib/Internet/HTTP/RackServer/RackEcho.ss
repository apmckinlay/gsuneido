// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// A simple app that just displays its environment
function (env)
	{
	return env.Assocs().Map!({ it[0] $ ': ' $ Display(it[1]) }).Join('\n')
	}