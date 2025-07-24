// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(ProtectRuleAllowsDelete?([], false))

		cl = ProtectRuleAllowsDelete? { EvalProtectRule(@unused) { return true } }
		Assert(cl([], 'test') is: false)

		cl = ProtectRuleAllowsDelete?
			{ EvalProtectRule(@unused) { return Object(allowHeaderDelete: true, abc:) } }
		Assert(cl([], 'test', header_delete:))

		cl = ProtectRuleAllowsDelete?
			{ EvalProtectRule(@unused) { return Object(abc:) } }
		Assert(cl([], 'test', header_delete:) is: false)

		cl = ProtectRuleAllowsDelete?
			{ EvalProtectRule(@unused) { return Object(neverDelete: true) } }
		Assert(cl([], 'test') is: false)

		cl = ProtectRuleAllowsDelete? { EvalProtectRule(@unused) { return true } }
		Assert(cl([], 'test', new_record:))

		cl = ProtectRuleAllowsDelete? { EvalProtectRule(@unused) { return false } }
		Assert(cl([], 'test'))

		cl = ProtectRuleAllowsDelete? { EvalProtectRule(@unused) { return '' } }
		Assert(cl([], 'test'))

		cl = ProtectRuleAllowsDelete?
			{ EvalProtectRule(@unused) { return Object(allowDelete: true) } }
		Assert(cl([], 'test'))

		cl = ProtectRuleAllowsDelete?
			{ EvalProtectRule(@unused) { return Object(allowDelete: false) } }
		Assert(cl([], 'test') is: false)

		cl = ProtectRuleAllowsDelete?
			{ EvalProtectRule(@unused) { return Object(abc:, def:) } }
		Assert(cl([], 'test') is: false)
		}
	}