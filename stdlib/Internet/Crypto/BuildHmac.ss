// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function (message, key, hashFn)
	{
	blocksize = 64
	if key.Size() > blocksize
		key = hashFn(key)
	key = key.RightFill(blocksize, '\x00')
	o_key_pad = StringXor(key, '\x5c')
	i_key_pad = StringXor(key, '\x36')
	return hashFn(o_key_pad $ hashFn(i_key_pad $ message))
	}