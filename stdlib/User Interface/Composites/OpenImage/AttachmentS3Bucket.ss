// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		return OptContribution(
			'OpenImageProcessFile', class { S3Bucket() { '' } }).S3Bucket()
		}

	Info()
		{
		c = OptContribution(
			'OpenImageProcessFile', class { S3Bucket() { '' } S3Region() { '' } })
		return Object(bucket: c.S3Bucket(), region: c.S3Region())
		}
	}