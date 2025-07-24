// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
class
	{
	Type: [
		// Event: "EVENT_NAME"	// arg1: uniqueId	arg2: args
		Heartbeat: '',
		UpdateStatus: 0,		// arg1: member		arg2: value
		SuJsTimeOut: 1,			// arg1: id
		SyncVisibility: 2,		// arg1: state

		/* For server response */
		// Actions: Object		// arg1: eventId
		CONNECTED: 3			// arg1: connectId
		OVERLAY: 4				// arg1: msg		arg2: hide
		]

	FormatEvent(typeOrEvent, eventId, ack, arg1 = #(n_a:), arg2 = #(n_a:))
		{
		ob = Object(typeOrEvent, eventId, ack)
		if not Same?(arg1, #(n_a:))
			ob.Add(arg1)
		if not Same?(arg2, #(n_a:))
			ob.Add(arg2)
		return ob
		}

	FormatResponse(typeOrActions, arg1 = #(n_a:), arg2 = #(n_b:))
		{
		ob = Object(typeOrActions)
		if not Same?(arg1, #(n_a:))
			ob.Add(arg1)
		if not Same?(arg2, #(n_b:))
			ob.Add(arg2)
		return ob
		}
	}
