// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_set_should_validate()
		{
		r = []
		r.test_simple_pull // so it has deps
		r.name = 'fred' // invalidates test_simple_pull
		r.test_simple_pull = 'joe'
		Assert(r.test_simple_pull is: 'joe')
		}
	Test_explicit_rules()
		{
		r = [name: 'fred', test_amount: 123]
		Assert(r.test_simple_pull is: 'fred and 123')

		r = [].Set_default()
		Assert({ r.test_simple_pull } throws: 'member not found')

		r = [name: 'fred', test_amount: 123]
		r.Set_default()
		r.AttachRule(#rule, function () { 'hello' })
		Assert(r.rule is: 'hello')
		}
	Test_modification_by_observer()
		{
		r = [b: 1, c: 2, e: 3, f: 4].Set_default()
		r.AttachRule('a', function () { .b $ .c })
		r.AttachRule('d', function () { .e $ .f })
		Assert(r.a is: '12')
		Assert(r.d is: '34')
		log = []
		r.Observer({|member| log.Add(member) })
		r.Observer(function (member) { if (member is 'c') .f = 9 })
		r.c = 22
		Assert(log is: #(c, f, a, d))
		}
	}