// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(env)
		{
		IM_MessengerManager.LogConnectionsIfChanged()
		args = IM_MessengerManager.ParseQuery(env.body)
		userName = args.user

		if not Suneido.Member?('IM_Messages')
			return Json.Encode(#())

		if not Suneido.IM_Messages.Member?(userName)
			return Json.Encode(#())

		return Json.Encode(.getMessagesFromGlobal(userName))
		}
	MaxMessageFetchLimit: 500 /* This controls the actual trimming of the queue, the value
								in IM_MessengerControl is just for the user to see in a
								popup message, change this to set a maximum fetch limit
								on messages*/
	getMessagesFromGlobal(userName)
		{
		ob_messages = Object()
		if not Suneido.IM_Messages[userName].Empty?()
			ob_messages.Merge(Suneido.IM_Messages[userName])
		Suneido.IM_Messages.Delete(userName)

		if ob_messages.Size() > .MaxMessageFetchLimit
			{
			ob_messages = ob_messages[.. .MaxMessageFetchLimit]
			ob_messages["error"] = true
			/* in the future error codes can be used to throw specific errors*/
			}
		return ob_messages // The output format is a 2d nested object
		}
	}
