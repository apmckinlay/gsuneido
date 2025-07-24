// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_parse()
		{
		Assert(SchedAt("garbage") is: false)
		Assert(SchedAt("25:30") is: false)

		test = function (s, expected)
			{
			sched = SchedAt("at " $ s)
			Assert(sched.SchedAt_at is: expected)
			}
		test("1:00", "0100")
		test("19:30", "1930")
		}
	Test_due()
		{
		test = function (at, prev, cur, expected)
			{
			if prev isnt false
				prev = Date('Jan 1 2001 ' $ prev)
			cur = Date('Jan 1 2001 ' $ cur)
			// SchedAt requires a date as a param, but only uses hours/minutes
			Assert(SchedAt("at " $ at).Due?(prev, cur) is: expected)
			}
		// "at" is before prev and cur
		test('9:30', '9:35', '9:45', false)
		// "at" is after prev and cur
		test('10:30', '9:35', '9:45', false)
		// "at" is between prev and cur
		test('9:30', '9:15', '9:35', true)
		// "at" is same as cur
		test('9:35', '9:15', '9:35', true)
		// "at" is same as prev
		test('9:15', '9:15', '9:35', false)
		// "at" is before midnight, between prev and cur and day changed
		test('23:55', '23:50', '0:10', true)
		// "at" is after midnight, between prev and cur and day changed
		test('0:05', '23:50', '0:10', true)

		// "at" is after midnight, after prev and cur and day changed
		test('0:11', '23:50', '0:10', false)
		// "at" is after midnight, same time as cur and day changed
		test('0:10', '23:50', '0:10', true)
		// "at" is before midnight, between prev and cur and day changed
		test('23:55', '23:50', '0:10', true)
		// "at" is before midnight, before prev and cur and day changed
		test('21:55', '23:50', '0:10', false)
		// "at" is before midnight, on prev and day changed
		test('21:50', '23:50', '0:10', false)

		// prevcheck is false (first run)
		test('9:15', false '0:10', false)
		test('9:15', false '9:10', false)
		}
	Test_skipWeekend()
		{
		// ensure that skipWeekends with "" is handled
		Assert(SchedAt("at 11:00").Due?(#20170904.1059, #20170904.1101))
		Assert(SchedAt("at 11:00").Due?(#20170902.1059, #20170902.1101))

		test = function (at, prev, cur, expected)
			{
			Assert(SchedAt("at " $ at $ ' skip weekends').Due?(prev, cur) is: expected)
			}

		// 20170904 is a monday
		test('11:00', #20170904.1059, #20170904.1101, true)

		// 20170902 is a saturday
		test('11:00', #20170902.1059, #20170902.1101, false)

		// prevcheck friday, curtime saturday, "at" before midnight
		test('23:55', #20170901.2350, #20170902.0010, true)
		// prevcheck friday, curtime saturday, "at" after midnight
		test('00:05', #20170901.2350, #20170902.0010, false)

		// prevcheck sunday, curtime monday, "at" before midnight
		test('23:55', #20170903.2350, #20170904.0010, false)
		// prevcheck sunday, curtime monday, "at" after midnight
		test('00:05', #20170903.2350, #20170904.0010, true)

		// prevcheck friday, curtime saturday, "at" at midnight ("at" 00:00 saturday)
		test('00:00', #20170901.2350, #20170902.0010, false)
		// prevcheck sunday, curtime monday, "at" at midnight ("at" 00:00 monday)
		test('00:00', #20170903.2350, #20170904.0010, true)
		}
	}
