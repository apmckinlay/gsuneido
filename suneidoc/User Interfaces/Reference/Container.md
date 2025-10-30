### Container

Abstract base class for *containers*
i.e. controls that contain other controls.

Controls derived from Container must implement:
GetChildren() => object
: Return a list of the container's child controls.

Has default methods for:
Broadcast(method, args...)
: Calls the method with the supplied arguments on each child control.

SetEnabled(enabled)
: Broadcasts SetEnabled.

GetEnabled()
: Returns True if all child controls are enabled, False otherwise.

SetVisible(visible)
: Broadcasts SetVisible.

SetReadOnly(readOnly)
: Broadcasts SetReadOnly.

GetReadOnly()
: Returns True if all child controls are read-only, False otherwise.

HasFocus?()
: Returns True is any child control has the focus, False otherwise.

Destroy()
: Calls Destroy on all the child controls.