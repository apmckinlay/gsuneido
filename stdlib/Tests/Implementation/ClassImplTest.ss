// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_01_class_without_base()
		{
		class { }
		}
	Test_02_class_with_base()
		{
		Test { }
		}
	Test_03_class_type()
		{
		c = class { }
		Assert(Type(c) is: 'Class')
		}
	Test_04_class_equality_is_identity()
		{
		Assert(class { } isnt: class { })
		c = class { }
		Assert(c is: c)
		}
	Test_05_class_members()
		{
		c = class { M: 123 }
		Assert(c.M is: 123)
		}
	Test_06_classes_are_readonly()
		{
		f = function () { c = class { }; c.Mem = 123 }
		e = Catch(f)
		Assert(e.Has?('readonly') or e.Has?("does not support put"))
		}
	Test_07_class_methods()
		{
		c = class { F() { 123 } }
		Assert(c.F() is: 123)
		}
	Test_08_class_members()
		{
		c = class { M: 123, F() { } }
		Assert(c members: #(F, M))
		}

	Test_09_instance()
		{
		c = class { }
		new c
		}
	Test_10_instance_type()
		{
		c = class { }
		i = new c
		Assert(Type(i) is: 'Instance')
		}
	Test_11_instances_inherit_from_class()
		{
		c = class { M: 123 }
		i = new c
		Assert(i.M is: 123)
		}
	Test_12_memberq_includes_inherited()
		{
		c = class { M: 123 }
		i = new c
		Assert(i.Member?('M'))
		}
	Test_13_method_lookup_starts_in_class()
		{
		c = class { F() { 123 } }
		i = new c
		i.F = 'foo'
		Assert(c.F() is: 123)
		}
	Test_14_instance_equality_is_value()
		{
		c = class { UseDeepEquals: true }
		Assert(new c is: new c)
		}
	Test_15_default_CallClass_is_new()
		{
		c = class { UseDeepEquals: true }
		Assert(c() is: new c)
		}
	Test_16_instances_are_modifiable()
		{
		c = class { }
		i = new c
		i.Mem = 123
		Assert(i.Mem is: 123)
		}
	}