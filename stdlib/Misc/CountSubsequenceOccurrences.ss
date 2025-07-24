// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function (list, pattern)
	{
	m = pattern.Size()
	n = list.Size()

	dp = [1].Set_default(0)

	for (i = 0; i < n; i++)
		{
		for (j = m; j > 0; j--)
			{
			if list[i] is pattern[j - 1]
				dp[j] += dp[j - 1]
			}
		}

	return dp[m]
	}
