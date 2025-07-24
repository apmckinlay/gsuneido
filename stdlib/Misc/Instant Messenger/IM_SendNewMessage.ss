// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(env)
		{
		IM_MessengerManager.LogConnectionsIfChanged()
		args = IM_MessengerManager.ParseQuery(env.body)
		msgOb = Object(im_to: args.to, im_message: String(args.msg),
			im_from: args.from, im_num: args.date, imchannel_num: args.imchannel_num)

		if not Suneido.Member?('IM_Messages')
			Suneido.IM_Messages = Object()
		if not Suneido.Member?('IM_CheckMessages')
			Suneido.IM_CheckMessages = Object()

		if #() isnt parsed_recipients = IM_MessengerManager.ParseRecipients(args.to)
			{
			if not msgOb.im_message.Prefix?(IM_MessengerManager.SysMsgTemplate)
				parsed_recipients.Remove(args.from)
			for recipient in parsed_recipients
				{
				msgOb.im_to = recipient
				.AddMessage(msgOb)
				}
			}
		return Json.Encode(#())
		}
	AllReadSystemMessage: "~~~Channel Maintenance~~~ #(type: 'all read'," $
				"token:'a0cb03a6f32a2448')"
	AddMessage(ob)
		{
		addMessageToGlobalQueue = true
		to = ob.im_to
		if not Suneido.IM_Messages.Member?(to)
			Suneido.IM_Messages[to] = Object()

		if not Suneido.Member?(#IM_MessagesRead)
			Suneido.IM_MessagesRead = Object()

		if ob.im_message.Prefix?(IM_MessengerManager.SysMsgTemplate)
			{
			if ob.im_message is .AllReadSystemMessage
				{
				Suneido.IM_MessagesRead[to] = Date()
				addMessageToGlobalQueue = false
				}
			else
				Suneido.IM_CheckMessages[to $ '_sys_msg'] = Date()
			}
		Suneido.IM_CheckMessages[to] = Date()
		if addMessageToGlobalQueue
			Suneido.IM_Messages[to].Add(ob)
		}
	}
