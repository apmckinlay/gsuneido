### PubSub
`Subscribe(event, handler) => handle`
: Registers that handler should be called when event is published. Returns an object that can be used to Unsubscribe

`Publish(event)`
: Calls all the subscribers.

event is normally a string

**Note**: Subscribers are called in the publisher's thread. If the Publish can be from a secondary thread (not the main UI thread) it is the subscriber's responsibility to use [Defer](<../../User Interfaces/Reference/Defer.md>) if they require it.

For example:

``` suneido
New()
	{
	.sub = PubSub.Subscribe(#myevent, .handle_myevent)
	}
handle_myevent()
	{
	...
	}
...
Destroy()
	{
	.sub.Unsubscribe()
	}
```