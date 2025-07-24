// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// nrecords - number of records to use in the table
// minutes, seconds - how long to run the test
// max_delay - maximum delay in ms between each transaction
class
	{
	hundred: 100
	CallClass(nrecords, seconds = 0, minutes = 0, max_delay = 0)
		{
		.create(nrecords)
		ops = 0
		stop_time = Date().Plus(:minutes, :seconds)
		while (Date() < stop_time)
			{
			++ops
			Transaction(update:)
				{|t|
				switch (Random(10000) % 6) /*= random 0 to 5 */
					{
				case 0 :
					.output_delete(nrecords, t)
				case 1, 2, 3 : /*= 50% chance for update record */
					.update(nrecords, t)
				case 4 : /*= 16.67% chance for deleting output */
					.delete_output(nrecords, t)
				case 5 : /*= 16.67% chance for readonly update */
					.readonly_update(nrecords, t)
					}

				if (Random(.hundred) < 20) /*= 20% chance for small sleep */
					Thread.Sleep(Random(5)) /*= 5 ms sleep */

				if (Random(.hundred) < 10) /*= 10% chance of rolling back */
					try
						t.Rollback()
					catch (unused, 'cannot Rollback completed Transaction')
						{}
				else
					try
						t.Complete()
					catch (unused, '*conflict')
						{}
				}
			if (max_delay > 0)
				Thread.Sleep(Random(max_delay))
			}
		Print(:ops)
		}
	create(nrecords)
		{
		if not TableExists?('serverstress')
			{
			Database("create serverstress (num, text) key(num)")
			for (i = 0; i < nrecords; ++i)
				QueryOutput('serverstress', Object(num: i, text: 'Number:' $ i))
			Print('created')
			}
		}
	output_delete(nrecords, t)
		{
		.output(nrecords, t)
		x = .input(nrecords, t)
		try
			x.Delete()
		catch (unused, "*transaction conflict")
			{}
		}
	input(nrecords, t)
		{
		for (i = 0; i < .hundred; ++i)
			if false isnt x = t.Query1("serverstress" , num: Random(nrecords * 2))
				return x
		Print("too many input loops")
		}
	find_missing(nrecords, t)
		{
		for (i = 0; i < .hundred; ++i)
			if false is t.Query1("serverstress" , num: num = Random(nrecords * 2))
				return num
		Print("too many find_missing loops")
		}
	output(nrecords, t)
		{
		num = .find_missing(nrecords, t)
		try
			t.QueryOutput('serverstress',
				[:num, text: 'Number:' $ num])
		catch (unused, '*duplicate key')
			{}
		}
	delete(x)
		{
		try
			x.Delete()
		catch (unused, "*transaction conflict")
			{}
		}
	update(nrecords, t)
		{
		x = .input(nrecords, t)
		x.text $= 'u'
		try
			if (Random(7) isnt 0) /*= 85% chance for normal update*/
					x.Update()
			else // update key
				{
				x.num = .find_missing(nrecords, t)
				.update_x(x)
				}
		catch (unused, "*transaction conflict")
			{}
		}
	update_x(x)
		{
		// could still fail if another process took the key
		try
			x.Update()
		catch (unused, "*duplicate key")
			{}
		}
	delete_output(nrecords, t)
		{
		x = .input(nrecords, t)
		.delete(x)
		.output(nrecords, t)
		}
	readonly_update(nrecords, t)
		{
		.input(nrecords, t)
		.input(nrecords, t)
		}
	}