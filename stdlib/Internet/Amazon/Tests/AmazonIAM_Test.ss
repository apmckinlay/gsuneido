// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_extractCredentials()
		{
		response = '<GetFederationTokenResponse xmlns="https://sts.amazonaws.com/doc/
2011-06-15/">
  <GetFederationTokenResult>
	<Credentials>
	  <SessionToken>
		AQoDYXdzEPT//////////wEXAMPLEtc764bNrC9SAPBSM22wDOk4x4HIZ8j4FZTwdQW
	  </SessionToken>
	  <SecretAccessKey>
	   9lcbEXAMPLEzZFojeqcAwPlTgatkXHa46n0PC
	  </SecretAccessKey>
	  <Expiration>2011-07-15T23:28:33.359Z</Expiration>
	  <AccessKeyId>ASIAEXAMPLE434RO6QZA</AccessKeyId>
	</Credentials>
	<FederatedUser>
	  <Arn>arn:aws:sts::123456789012:federated-user/Bob</Arn>
	  <FederatedUserId>123456789012:Bob</FederatedUserId>
	</FederatedUser>
	<PackedPolicySize>6</PackedPolicySize>
  </GetFederationTokenResult>
  <ResponseMetadata>
	<RequestId>c6104cbe-af31-11e0-8154-cbc7ccf896c7</RequestId>
  </ResponseMetadata>
</GetFederationTokenResponse>'
		results = AmazonIAM.AmazonIAM_extractCredentials(response)
		Assert(results.sessiontoken
			is: 'AQoDYXdzEPT//////////wEXAMPLEtc764bNrC9SAPBSM22wDOk4x4HIZ8j4FZTwdQW')
		Assert(results.secretaccesskey is: '9lcbEXAMPLEzZFojeqcAwPlTgatkXHa46n0PC')
		Assert(results.accesskeyid is: 'ASIAEXAMPLE434RO6QZA')
		}
	}