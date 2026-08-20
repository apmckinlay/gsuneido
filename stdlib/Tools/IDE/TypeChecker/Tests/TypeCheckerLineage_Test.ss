// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_builtinStub()
		{
		// SocketServer is a builtin so Global returns it but Class? is false.
		// Resolver should synthesize the stub instead of returning false.
		chain = TypeCheckerLineage(#(SocketServer)).GetDefault(#SocketServer, false)
		Assert(chain isnt: false)
		Assert(chain.Size() is: 1)

		entry = chain[0]
		Assert(entry.name is: #SocketServer)
		Assert(entry.src.Has?(#CopyTo))
		Assert(entry.src.Has?(":number"))
		Assert(entry.src.Has?(#Readline))
		Assert(entry.src.Has?(#Writeline))
		}

	Test_unknownGlobalReturnsFalse()
		{
		name = #DoesNotExist_xyz_123
		chains = TypeCheckerLineage([name])
		Assert(chains.GetDefault(name, #sentinel) is: false)
		}

	Test_nonObjectInputReturnsEmpty()
		{
		Assert(TypeCheckerLineage("not an object") is: #())
		Assert(TypeCheckerLineage(false) is: #())
		}

	Test_functionReturnsSource()
		{
		chains = TypeCheckerLineage(#(NameSplit))
		chain = chains.GetDefault(#NameSplit, false)
		Assert(chain isnt: false)
		Assert(chain.Size() is: 1)
		entry = chain[0]
		Assert(entry.name is: #NameSplit)
		Assert(entry.src.Has?("function (name, split_on = false)"))
		Assert(entry.src.Has?("first:"))
		}

	Test_overrideStackBundlesLowerDefinitions()
		{
		lineage = TypeCheckerLineage
			{
			TypeCheckerLineage_nameOf(x) { return x.name }
			TypeCheckerLineage_baseOf(x) { return x.GetDefault(#base, false) }
			TypeCheckerLineage_libOf(x) { return x.lib }
			TypeCheckerLineage_safeSourceCode(x) { return x.src }
			}
		low = [name: 'Ovr', lib: 'etalib', base: [name: 'Access', lib: 'axonlib'],
			src: `Access
					{
					LowOnly() { }
					}`]
		high = [name: 'Ovr', lib: 'carslib', base: low, src: `_Ovr
					{
					HighOnly() { }
					}`]

		entryCache = Object()
		lineage.TypeCheckerLineage_addOverrideStack(entryCache, 'Ovr', high)

		Assert(entryCache.Members().Sort!() is: #('Ovr', 'Ovr__etalib'))

		leaf = entryCache.Ovr
		Assert(leaf.name is: 'Ovr')
		Assert(leaf.baseName is: 'Ovr__etalib') // the bug: base was dropped
		Assert(leaf.src.Has?('HighOnly'))
		Assert(leaf.src.Has?('Ovr__etalib')) // `_Ovr` retargeted, not stripped
		Assert(leaf.src.Has?('_Ovr') is: false)

		under = entryCache['Ovr__etalib']
		Assert(under.name is: 'Ovr__etalib')
		Assert(under.baseName is: 'Access')
		Assert(under.src.Has?('LowOnly'))
		}
	}
