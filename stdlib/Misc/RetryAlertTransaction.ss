// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function (block)
	{
	try RetryTransaction()
		{|t|
		block(:t)
		}
	catch (err, "RetryTransaction: too many retries")
		{
		AlertTransactionCommitFailed(err)
		return false
		}
	return true
	}