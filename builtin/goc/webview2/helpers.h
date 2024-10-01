
#ifndef suneido_webview2_helpers
#define suneido_webview2_helpers

#include <windows.h>
#include <cwchar>
#include <stdio.h>
#include <iostream>
#include <string>
#include <map>

wchar_t* ConvertCharToWChar(const char* s) {
    size_t newSize = strlen(s) + 1;
    wchar_t* wCharArray = new wchar_t[newSize];
    size_t convertedChars = 0;
    mbstowcs_s(&convertedChars, wCharArray, newSize, s, _TRUNCATE);
    return wCharArray;
}

char * ConvertWcharToChar(LPWSTR s) {
    if (s == nullptr) {
        return nullptr;
    }

    // Determine the length of the resulting string
    int len = WideCharToMultiByte(CP_UTF8, 0, s, -1, nullptr, 0, nullptr, nullptr);
    if (len == 0) {
        return nullptr; // Conversion failed
    }

    // Allocate memory for the converted string
    char* charStr = new char[len];

    // Perform the actual conversion
    WideCharToMultiByte(CP_UTF8, 0, s, -1, charStr, len, nullptr, nullptr);

    return charStr;
}

class MyString {
    public:
        MyString() {}
        ~MyString() {
            Release();
        }
        void Release() {
            if (s) {
                CoTaskMemFree(s);
                s = nullptr;
            }
        }
        LPCWSTR Set(LPCWSTR str) {
            Release();
            if (str) {
                s = MakeString(str);
            }
            return s;
        }
        LPCWSTR Get() {
            return s;
        }
        LPWSTR Copy() {
            if (s) {
                return MakeString(s);
            }
            return nullptr;
        }
    private:
        LPWSTR MakeString(LPCWSTR source) {
            const size_t length = wcslen(source);
            const size_t bytes = (length + 1) * sizeof(*source);

            if (bytes <= length) {
                return nullptr;
            }

            wchar_t* result = reinterpret_cast<wchar_t*>(CoTaskMemAlloc(bytes));

            if (result)
                memcpy(result, source, bytes);
            return result;
         }
        LPWSTR s = nullptr;
};

char hexToChar(const char* hex) {
    char ch = 0;
    for (int i = 0; i < 2; ++i) {
        ch *= 16;
        if (hex[i] >= '0' && hex[i] <= '9') {
            ch += hex[i] - '0';
        } else if (hex[i] >= 'A' && hex[i] <= 'F') {
            ch += hex[i] - 'A' + 10;
        } else if (hex[i] >= 'a' && hex[i] <= 'f') {
            ch += hex[i] - 'a' + 10;
        }
    }
    return ch;
}

void decodeURI(const char* src, char* dest) {
    while (*src) {
        if (*src == '%') {
            if (*(src + 1) && *(src + 2)) {
                *dest = hexToChar(src + 1);
                src += 2;
            }
        } else if (*src == '+') {
            *dest = ' ';  // Convert '+' to space
        } else {
            *dest = *src;  // Copy normal character
        }
        ++src;
        ++dest;
    }
    *dest = '\0';  // Null-terminate the output string
}

// Mapping of file extensions to MIME types
std::map<std::string, std::wstring> mimeMap = {
    {".html", L"text/html"},
    {".htm",  L"text/html"},
    {".css",  L"text/css"},
    {".js",   L"application/javascript"},
    {".json", L"application/json"},
    {".png",  L"image/png"},
    {".jpg",  L"image/jpeg"},
    {".jpeg", L"image/jpeg"},
    {".gif",  L"image/gif"},
    {".svg",  L"image/svg+xml"},
    {".pdf",  L"application/pdf"},
    {".txt",  L"text/plain"},
    // Add more mappings as needed
};

std::wstring getContentTypeFromUri(const char* uri) {
    // Find the last dot in the URI to get the extension
    const char* dot = strrchr(uri, '.');
    if (!dot || dot == uri) {
        return L"text/html"; // Default to binary stream if no extension found
    }

    // Convert the file extension to lowercase (optional)
    std::string extension = dot;
    for (char& c : extension) {
        c = std::tolower(c);
    }

    // Lookup the MIME type in the map
    auto it = mimeMap.find(extension);
    if (it != mimeMap.end()) {
        return it->second;
    } else {
        return L"text/html"; // Default MIME type
    }
}

std::wstring buildContentTypeHeader(const char* uri) {
    std::wstring contentType = getContentTypeFromUri(uri);
    std::wstring header = L"Content-Type: " + contentType + L"\r\n";
    return header;
}

#endif