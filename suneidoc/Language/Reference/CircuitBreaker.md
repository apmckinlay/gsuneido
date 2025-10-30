### CircuitBreaker

``` suneido
service, runnable, arguments...
```

The default configuration is .threshold = 10 /*times*/, .timeout = 600 /*seconds*/, .timeoutIncrement = 0 /*seconds*/, .failureExpiry = false. The configuration defaults can be overridden using a Contribution for the service.

For example

``` suneido
SampleLibrary_CircuitBreakerConfig
#(
    'AmazonSQS':           (threshold: 20, timeout: 600, timeoutIncrement: 0)
)
```

Implements the basic CircuitBreaker pattern. This class is used to handle calls to services that could be unavailable for an extended period of time. Using this class will help prevent tying up resources involved in making calls to a service that is likely to fail.

Suneido uses the specified **service** name to maintain one CircuitBreaker instance for that service. Later calls to CircuitBreaker for the same service name will use the existing instance.

Once the specified **threshold** for failed calls to the service has been reached, the circuit will be opened which prevents further calls to the service for the duration of the specified **timeout** period. The CircuitBreaker will then try the service call once. If that call fails the circuit goes back to the open state and calls are again prevented for the duration of the timeout period. If the **timeoutIncrement** is specified, the length of time the circuit remains open will be incremented by this amount. The maximum amount of time the circuit will remain open is two hours.

**failureExpiry** (in minutes) can be used to ignore old failures when checking against the **threshehold**

Exceptions in calls are also treated as a failure.

For example:

``` suneido
result = CircuitBreaker('ftp', SampleFtpTransfer, file)
```