// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(connections)
		{
		unlabeledConnections = Object()
		for client in connections
			if not .valid?(client)
				unlabeledConnections.Add(client)
		return unlabeledConnections
		}
	valid?(client)
		{
		return client.Has?('(') or client.Has?('@') or // named connection
			client.Suffix?('-thread') or // named threads
			client.Has?('SocketServer-thread') or
			client.Has?('suneido-thread') or
			client.Has?('ThreadTotal') or
			client =~ 'SocketServer-.*(HTTP|Rack) Server$' or
			client.Suffix?('LibLocateList') // libview
		}
	}