// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_builtinStub()
		{
		// SocketServer is a builtin so Global returns it but Class? is false.
		// Resolver should synthesize the stub instead of returning false.
		chain = TypeCheckerLineage(#('SocketServer')).GetDefault('SocketServer', false)
		Assert(chain isnt: false)
		Assert(chain.Size() is: 1)

		entry = chain[0]
		Assert(entry.name is: 'SocketServer')
		Assert(entry.src.Has?('CopyTo'))
		Assert(entry.src.Has?(':number'))
		Assert(entry.src.Has?('Readline'))
		Assert(entry.src.Has?('Writeline'))
		}

	Test_unknownGlobalReturnsFalse()
		{
		name = 'DoesNotExist_xyz_123'
		chains = TypeCheckerLineage(Object(name))
		Assert(chains.GetDefault(name, 'sentinel') is: false)
		}

	Test_nonObjectInputReturnsEmpty()
		{
		Assert(TypeCheckerLineage('not an object') is: #())
		Assert(TypeCheckerLineage(false) is: #())
		}
	Test_functionReturnsSource()
		{
		chains = TypeCheckerLineage(#('NameSplit'))
		chain = chains.GetDefault('NameSplit', false)
		Assert(chain isnt: false)
		Assert(chain.Size() is: 1)
		entry = chain[0]
		Assert(entry.name is: 'NameSplit')
		Assert(entry.src.Has?('function (name, split_on = false)'))
		Assert(entry.src.Has?('first:'))
		}
	}
