#include "MyWebView2Base.h"

class MyAcceleratorKeyPressedHandler : public ICoreWebView2AcceleratorKeyPressedEventHandler {
public:
    MyAcceleratorKeyPressedHandler(MyBrowserObject* pBrowserObject) : refCount(1), pBrowserObject(pBrowserObject) {
        pBrowserObject->AddRef();
    }

    ~MyAcceleratorKeyPressedHandler() {
        pBrowserObject->Release();
    }

    // IUnknown methods
    ULONG STDMETHODCALLTYPE AddRef() override {
        return InterlockedIncrement(&refCount);
    }

    ULONG STDMETHODCALLTYPE Release() override {
        ULONG newRefCount = InterlockedDecrement(&refCount);
        if (newRefCount == 0) {
            delete this;
        }
        return newRefCount;
    }

    HRESULT STDMETHODCALLTYPE QueryInterface(REFIID riid, void** ppvObject) override {
        if (riid == IID_IUnknown || riid == IID_ICoreWebView2AcceleratorKeyPressedEventHandler) {
            *ppvObject = static_cast<ICoreWebView2AcceleratorKeyPressedEventHandler*>(this);
            AddRef();
            return S_OK;
        }
        *ppvObject = nullptr;
        return E_NOINTERFACE;
    }

    // ICoreWebView2AcceleratorKeyPressedEventHandler method
    HRESULT STDMETHODCALLTYPE Invoke(
        ICoreWebView2Controller* sender,
        ICoreWebView2AcceleratorKeyPressedEventArgs* args) override
    {
        COREWEBVIEW2_KEY_EVENT_KIND keyEventKind;
        HRESULT hr = args->get_KeyEventKind(&keyEventKind);
        if (FAILED(hr)) {
            return hr;
        }

        if (keyEventKind == COREWEBVIEW2_KEY_EVENT_KIND_KEY_DOWN || keyEventKind == COREWEBVIEW2_KEY_EVENT_KIND_SYSTEM_KEY_DOWN) {
            UINT key;
            hr = args->get_VirtualKey(&key);
            if (FAILED(hr)) {
                return hr;
            }

            if (pBrowserObject->onAccelPressed(key)) {
                args->put_Handled(TRUE);
            }
        }
        return S_OK;
    }

private:
    LONG refCount;
    MyBrowserObject* pBrowserObject;
};