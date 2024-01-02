// Package tdlib provides CGo wrapper for Telegram TDLib.
//
// It wraps only modern methods for TDLib JSON client:
// https://core.telegram.org/tdlib/docs/td__json__client_8h.html
package tdlib

import "unsafe"

/*
#cgo linux CFLAGS: -I/usr/local/include -I"third-party/td/tdlib/include"
#cgo linux LDFLAGS: -L/usr/local/lib -L"third-party/td/tdlib/lib" -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltdapi -ltddb -ltdsqlite -ltdnet -ltdutils -lstdc++ -lssl -lcrypto -ldl -lz -lm

#include <stdlib.h>
#include <td/telegram/td_json_client.h>
*/
import "C"

// TDCreateClientID
// Returns an opaque identifier of a new TDLib instance.
// The TDLib instance will not send updates until the first request is sent to it.
func TDCreateClientID() int32 {
	return int32(C.td_create_client_id())
}

// TDSend
// Sends request to the TDLib client.
// May be called from any thread.
//
// This might mutate provided `request` slice by setting NULL-byte after content if `cap(request) > len(request)`.
func TDSend(
	clientID int32, // TDLib client identifier.
	request []byte, // JSON-serialized request to TDLib. Doesn't have to be NULL-terminated, the wrapper handles this on its own.
) {
	C.td_send(C.int(clientID), (*C.char)(unsafe.Pointer(&nullTerminated(request)[0])))
}

// TDReceive
// Receives incoming updates and request responses.
// Must not be called simultaneously from two different threads.
// The returned pointer can be used until the next call to td_receive (TDReceive) or td_execute (TDExecute),
// after which it will be deallocated by TDLib.
func TDReceive(
	timeout float64, // The maximum number of seconds allowed for this function to wait for new data.
) []byte {
	result := C.td_receive(C.double(timeout))
	if result == nil {
		return nil
	}

	return C.GoBytes((unsafe.Pointer)(result), strlen(result))
}

// TDExecute
// Synchronously executes a TDLib request.
// A request can be executed synchronously, only if it is documented with "Can be called synchronously".
// The returned pointer can be used until the next call to td_receive (TDReceive) or td_execute (TDExecute),
// after which it will be deallocated by TDLib.
func TDExecute(
	request []byte, // JSON-serialized request to TDLib. Doesn't have to be NULL-terminated, the wrapper handles this on its own.
) []byte {
	result := C.td_execute((*C.char)(unsafe.Pointer(&nullTerminated(request)[0])))
	if result == nil {
		return nil
	}

	return C.GoBytes((unsafe.Pointer)(result), strlen(result))
}

// ----------------------------------------------------------------------------
// Netherworld of C hacks and pointer arithmetics

func nullTerminated(p []byte) []byte {
	if len(p) == 0 {
		return p // empty
	}

	if p[len(p)-1] == 0x00 {
		return p // already null-terminated
	}

	return append(p, 0x00) // extends or reallocates
}

func strlen(p *C.char) C.int {
	if p == nil {
		return 0
	}

	var n C.int
	for ; *p != 0x00; p = (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + 1)) { // -_-
		n++
	}
	return n
}
