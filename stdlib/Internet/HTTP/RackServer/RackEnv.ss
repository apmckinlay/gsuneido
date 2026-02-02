// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(socket)
		{
		this['socket'] = socket
		}

	Getter_(member)
		{
		if member is 'body'
			{
			return this['body'] = .Member?('content_length') and
				this['content_length'] isnt 0
				? this['socket'].Read(this['content_length'])
				: ''
			}
		}

	// for test
	Build(@args)
		{
		env = new this(false)
		for m in args.Members()
			env[m] = args[m]
		return env
		}

	Eq?(@args)
		{
		for m in args.Members()
			if this[m] isnt args[m]
				return false
		return true
		}

	ToString()
		{
		s = "RackEnv:"
		for m in this.Members().Sort!()
			s $= "\n\t" $ m $ ": " $ Display(this[m])
		return s
		}
	}