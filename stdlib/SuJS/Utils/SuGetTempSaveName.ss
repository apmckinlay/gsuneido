// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
function (ext = '')
	{
	EnsureDir('temp')
	return 'temp/' $ Display(Timestamp()).Tr('#.') $ Opt('.', ext)
	}