// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_StartSpotInstances()
		{
		mock = Mock(AmazonEC2)
		mock.When.postRequest([anyArgs:]).Return('postResult')
		mock.When.StartSpotInstances([anyArgs:]).CallThrough()
		mock.When.makePost([anyArgs:]).CallThrough()

		specs = []
		type = ''
		Assert(mock.Eval(AmazonEC2.StartSpotInstances, specs, type) is:
			'postResult')
		specs = []
		type = 'one-time'

		Assert(mock.Eval(AmazonEC2.StartSpotInstances, specs, type) is:
			'postResult')
		mock.Verify.postRequest('Action=RequestSpotInstances&Type=one-time&'$
			'Version=2016-11-15', region: 'us-east-1')

		specs =  Object('LaunchSpecification.InstanceType': 'AMI-someImage',
			Version: '2016-11-15')
		Assert(mock.Eval(AmazonEC2.StartSpotInstances, specs, type) is:
			'postResult')
		mock.Verify.postRequest('Action=RequestSpotInstances&' $
			'LaunchSpecification.InstanceType=AMI-someImage&Type=one-time&' $
			'Version=2016-11-15', region: 'us-east-1')
		}

	Test_CancelSpotRequest()
		{
		mock = Mock(AmazonEC2)
		mock.When.makePost([anyArgs:]).CallThrough()
		mock.When.postRequest([anyArgs:]).Return('postResult')
		mock.When.CancelSpotRequest([anyArgs:]).CallThrough()

		mock.Eval(AmazonEC2.CancelSpotRequest)
		mock.Verify.Never().postRequest([anyArgs:])

		mock.Eval(AmazonEC2.CancelSpotRequest, '1')
		mock.Verify.postRequest('Action=CancelSpotInstanceRequests&' $
			'SpotInstanceRequestId.1=1&Version=2016-11-15', region: 'us-east-1')

		mock.Eval(AmazonEC2.CancelSpotRequest, '1', '2')
		mock.Verify.postRequest('Action=CancelSpotInstanceRequests&' $
			'SpotInstanceRequestId.1=1&SpotInstanceRequestId.2=2&Version=2016-11-15',
			region: 'us-east-1')

		mock.Eval(AmazonEC2.CancelSpotRequest, '1', '2', '3')
		mock.Verify.postRequest('Action=CancelSpotInstanceRequests&' $
			'SpotInstanceRequestId.1=1&SpotInstanceRequestId.2=2&' $
			'SpotInstanceRequestId.3=3&Version=2016-11-15',
			region: 'us-east-1')
		}

	Test_TerminateInstance()
		{
		mock = Mock(AmazonEC2)
		mock.When.makePost([anyArgs:]).CallThrough()
		mock.When.postRequest([anyArgs:]).Return('postResult')
		mock.When.TerminateInstance([anyArgs:]).CallThrough()

		mock.Eval(AmazonEC2.TerminateInstance)
		mock.Verify.Never().postRequest([anyArgs:])

		mock.Eval(AmazonEC2.TerminateInstance, '1')
		mock.Verify.postRequest('Action=TerminateInstances&' $
			'InstanceId.1=1&Version=2016-11-15', region: 'us-east-1')

		mock.Eval(AmazonEC2.TerminateInstance, '1', '2')
		mock.Verify.postRequest('Action=TerminateInstances&' $
			'InstanceId.1=1&InstanceId.2=2&Version=2016-11-15', region: 'us-east-1')

		mock.Eval(AmazonEC2.TerminateInstance, '1', '2', '3')
		mock.Verify.postRequest('Action=TerminateInstances&' $
			'InstanceId.1=1&InstanceId.2=2&' $
			'InstanceId.3=3&Version=2016-11-15', region: 'us-east-1')
		}

	Test_GetSpotRequestInfo()
		{
		mock = Mock(AmazonEC2)
		mock.When.makePost([anyArgs:]).CallThrough()
		mock.When.postRequest([anyArgs:]).Return('postResult')
		mock.When.buildExtraFilters([anyArgs:]).CallThrough()
		mock.When.GetSpotRequestInfo([anyArgs:]).CallThrough()

		mock.Eval(AmazonEC2.GetSpotRequestInfo)
		mock.Verify.postRequest([anyArgs:])

		mock.Eval(AmazonEC2.GetSpotRequestInfo, '1')
		mock.Verify.postRequest('Action=DescribeSpotInstanceRequests&' $
			'SpotInstanceRequestId.1=1&Version=2016-11-15', region: 'us-east-1')
		mock.Verify.Times(2).postRequest([anyArgs:])

		mock.Eval(AmazonEC2.GetSpotRequestInfo,
			extraFilters: #('status-code': 'fulfilled'))
		mock.Verify.postRequest('Action=DescribeSpotInstanceRequests&' $
			'Filter.1.Name=status-code&Filter.1.Value=fulfilled&Version=2016-11-15',
			region: 'us-east-1')
		}

	Test_GetInstanceInfo()
		{
		mock = Mock(AmazonEC2)
		mock.When.makePost([anyArgs:]).CallThrough()
		mock.When.postRequest([anyArgs:]).Return('postResult')
		mock.When.buildExtraFilters([anyArgs:]).CallThrough()
		mock.When.GetInstanceInfo([anyArgs:]).CallThrough()

		mock.Eval(AmazonEC2.GetInstanceInfo)
		mock.Verify.postRequest([anyArgs:])

		mock.Eval(AmazonEC2.GetInstanceInfo, '1')
		mock.Verify.postRequest('Action=DescribeInstances&' $
			'InstanceId.1=1&Version=2016-11-15', region: 'us-east-1')
		mock.Verify.Times(2).postRequest([anyArgs:])

		mock.Eval(AmazonEC2.GetInstanceInfo,
			extraFilters: #('status-code': 'fulfilled'))
		mock.Verify.postRequest('Action=DescribeInstances&' $
			'Filter.1.Name=status-code&Filter.1.Value=fulfilled&Version=2016-11-15',
			region: 'us-east-1')
		}

	Test_ModifyInstanceMetadataOptions()
		{
		mock = Mock(AmazonEC2)
		mock.When.makePost([anyArgs:]).CallThrough()
		mock.When.postRequest([anyArgs:]).Return(Object(content:
			`<?xml version="1.0" encoding="UTF-8"?>
			<ModifyInstanceMetadataOptionsResponse
				xmlns="http://ec2.amazonaws.com/doc/2016-11-15/">
				<requestId>43382050-f545-4eec-86cc-1a0693ffb250</requestId>
				<instanceMetadataOptions>
					<httpEndpoint>enabled</httpEndpoint>
					<httpProtocolIpv4>enabled</httpProtocolIpv4>
					<httpProtocolIpv6>disabled</httpProtocolIpv6>
					<httpPutResponseHopLimit>2</httpPutResponseHopLimit>
					<httpTokens>required</httpTokens>
					<instanceMetadataTags>disabled</instanceMetadataTags>
					<state>pending</state>
				</instanceMetadataOptions>
				<instanceId>i-09230890c28a6df75</instanceId>
			</ModifyInstanceMetadataOptionsResponse>`))
		options = mock.Eval(AmazonEC2.ModifyInstanceMetadataOptions, 'instance_id')
		Assert(options is:
			Object(instancemetadatatags: "disabled", httpendpoint: "enabled",
				httpputresponsehoplimit: '2', httptokens: "required",
				httpprotocolipv6: "disabled", httpprotocolipv4: "enabled",
				state: "pending"))
		}
	}
