// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Subscribe(event, block)
		{
		ps = Suneido.GetInit(#pubsub_subscribers, { Object().Set_default(Object()) })
		unused = Suneido.CompareAndSet(#pubsub_nextnum, 0)
		n = Suneido.pubsub_nextnum++
		ps[event][n] = block
		return (.unsub)(event, n)
		}

	unsub: class
		{
		New(.event, .n)
			{}
		Unsubscribe()
			{
			Suneido.pubsub_subscribers[.event].Delete(.n)
			return // avoid returning the whole Suneido object
			}
		}

	Publish(@args)
		{
		if not Suneido.Member?('pubsub_subscribers')
			return

		subscribers = Suneido.pubsub_subscribers[args[0]]
		if subscribers.Size() is 1
			return subscribers.Values().First()(@+1 args)
		// Copy to prevent "object modified during iteration"
		for func in subscribers.Copy()
			func(@+1args)
		}

	PublishWithArgs(@args)
		{
		if not Suneido.Member?('pubsub_subscribers')
			return

		// Copy to prevent "object modified during iteration"
		for func in Suneido.pubsub_subscribers[args[0]].Copy()
			func(@+1args)
		}

	/* READ ME: PublishConsolidate works as follows:
	1. The arguments used are stored in pubsub_delayed[Consolidate_<event>_<consolidate>]
	2. When the Defer call is triggered, the consolidated events are carried out
	Example:
	function ()
		{
		event = #LibraryTreeChange

		table = #stdlib
		PubSub.PublishConsolidate(event, name: #Rec1, :table)
		PubSub.PublishConsolidate(event, name: #Rec2, :table)

		table = #Test_lib
		PubSub.PublishConsolidate(event, name: #Rec3, :table, consolidate: '2')
		}
	This will trigger LibraryTreeChange subscribers twice with the following arguments:
		#( // consolidate: Consolidate_LibraryTreeChange
			(name: Rec1, table: stdlib)
			(name: Rec2, table: stdlib)
			)
		#( // consolidate: Consolidate_LibraryTreeChange_2
			(name: Rec3, table: Test_lib)
			)
	*/
	PublishConsolidate(@args)
		/*usage: event, args, consolidate: optional member to customize consolidation */
		{
		// gsport.exe does not have Defer. As a result, we cannot consolidate events
		// when not running in a GUI environment. Instead of consolidating, immediately
		// publish the event with the provided arguments as if it were consolidated.
		if not Sys.GUI?()
			.Publish(args.Extract(0), args: Object(args))
		else if Suneido.Member?('pubsub_subscribers')
			.consolidate(args,
				Suneido.GetInit('pubsub_delayed', Object),
				args.Extract('consolidate', ''))
		}

	consolidate(args, pd, consolidate)
		{
		event = args.Extract(0)
		if pd.Member?(member = 'Consolidate_' $ event $ Opt('_', consolidate))
			pd[member].AddUnique(args)
		else
			{
			pd[member] = Object(args)
			// Have to use Defer to ensure all events in the current stack are
			// consolidated before triggering Publish
			Defer({ .Publish(event, args: pd.Extract(member)) })
			}
		}

	Count()
		{
		return Suneido.GetDefault(#pubsub_subscribers, #()).SumWith(#Size)
		}
	}
