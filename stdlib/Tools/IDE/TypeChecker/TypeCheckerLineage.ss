class
	{
	CallClass(names)
		{
		if not Object?(names)
			return #()

		entryCache = Object()
		chainCache = Object()
		for name in names
			{
			if chainCache.Member?(name)
				continue
			try
				chainCache[name] = .resolveChain(entryCache, name)
			catch
				chainCache[name] = false
			}
		return chainCache
		}

	resolveChain(entryCache, leaf)
		{
		chain = Object()
		seen  = Object()
		curr  = leaf
		while curr isnt false and not seen.Member?(curr)
			{
			seen[curr] = true
			.ensureEntry(entryCache, curr)
			entry = entryCache[curr]
			if entry is false
				break
			chain.Add(Object(name: entry.name, src: entry.src))
			curr = entry.baseName
			}
		return chain.Empty?() ? false : chain.Reverse!()
		}

	// an object which maps class name to source stubs
	builtinStubs: #(
		SocketServer: #(
			src: `class
				{
				CopyTo(dest, nbytes = false) :number { }
				ManualClose() { }
				Read(nbytes = false) :string|false { }
				Readline() :string { }
				RemoteUser() :string { }
				SetTimeout(seconds) { }
				Write(s :string) { }
				Writeline(s :string) { }
				}`,
			baseName: false
			)
		)

	ensureEntry(entryCache, name)
		{
		if entryCache.Member?(name)
			return
		if .builtinStubs.Member?(name)
			{
			stub = .builtinStubs[name]
			entryCache[name] = Object(:name,
				src: stub.src, baseName: stub.baseName)
			return
			}
		instance = Global(name)
		if Function?(instance)
			{
			entryCache[name] = .functionEntry(name)
			return
			}
		if not Class?(instance)
			{
			entryCache[name] = false
			return
			}
		base = BaseClass(instance)
		src = .normalizeOverride(.safeSourceCode(instance), Name(instance))
		if not String?(src) or src is ''
			{
			entryCache[name] = false
			return
			}
		entryCache[name] = Object(name: Name(instance), :src,
			baseName: base is false ? false : Name(base))
		}

	normalizeOverride(src, name)
		{
		// we need special handling for overridden classes `_ApControl`
		u = '_' $ name
		explicit = 'class[ \t]*:[ \t]*' $ u
		if src =~ explicit
			return src.Replace(explicit, 'class', 1)
		return src.Replace(u, 'class', 1)
		}

	functionEntry(name)
		{
		src = .recordSrc(name)
		if not String?(src) or src is ''
			return false
		return Object(:name, :src, baseName: false)
		}

	// by-name, override order (functions)
	recordSrc(name)
		{
		for lib in Libraries().Reverse!()
			if false isnt src = .libSrc(lib, name)
				return src
		return false
		}

	libSrc(lib, name)
		{
		try
			return Query1(lib, group: -1, :name).text
		catch
			return false
		}

	// SourceCode(fn) starts by trying `fn.Source()`
	// Several stdlib classes (ListComponent, VirtualListGridComponent,
	// ParamsFormatClassForQuery, ...) have Getter_ methods that
	// infinitely recurse when asked for an undefined member, which blows
	// up the stack. We could use SourceCode directlye but dont want the unnecessary
	// logs
	safeSourceCode(instance)
		{
		base = Name(instance).BeforeFirst('.').BeforeFirst(' ')
		s = Display(instance)
		lib = s.AfterFirst('/* ').BeforeFirst(' ')
		recName = lib.Has?('__') ? base $ '__' $ lib.AfterFirst('__') : base
		lib = lib.BeforeFirst('__')
		if lib is "function"
			return ''
		src = .libSrc(lib, recName)
		return src is false ? '' : src
		}
	}
