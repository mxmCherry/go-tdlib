/*
Implements phone-number auth as described in:
https://core.telegram.org/tdlib/getting-started#user-authorization

Run:

	API_ID=xxx API_HASH=yyy go run cmd/example/main.go
*/
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/mxmCherry/go-tdlib"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// to read user input line(s)
	stdin := bufio.NewReader(os.Stdin)

	clientID := tdlib.TDCreateClientID()

	// https://core.telegram.org/tdlib/docs/classtd_1_1td__api_1_1set_log_verbosity_level.html
	// force WARNs and above, otherwise it's way too verbose
	tdlib.TDExecute([]byte(`{
		"@type": "setLogVerbosityLevel",
		"new_verbosity_level": 2
	}`))

	// https://core.telegram.org/tdlib/docs/classtd_1_1td__api_1_1get_authorization_state.html
	tdlib.TDSend(clientID, []byte(`{
		"@type": "getAuthorizationState"
	}`))

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		result := tdlib.TDReceive(0.100 /* 100 milliseconds */)
		if result == nil {
			continue
		}
		println(string(result))

		var res struct {
			Type               string `json:"@type"`
			AuthorizationState *struct {
				Type string `json:"@type"`
			} `json:"authorization_state,omitempty"`
		}
		if err := json.Unmarshal(result, &res); err != nil {
			panic("unmarshal TDReceive result: " + err.Error())
		}

		switch res.Type {

		// https://core.telegram.org/tdlib/docs/classtd_1_1td__api_1_1update_authorization_state.html
		case "updateAuthorizationState":
			switch res.AuthorizationState.Type {

			case "authorizationStateWaitTdlibParameters":
				// https://core.telegram.org/tdlib/docs/classtd_1_1td__api_1_1set_tdlib_parameters.html
				// https://core.telegram.org/tdlib/docs/classtd_1_1td__api_1_1tdlib_parameters.html
				tdlib.TDSend(clientID, []byte(`{
					"@type": "setTdlibParameters",
					"api_id": "`+os.Getenv("API_ID")+`",
					"api_hash": "`+os.Getenv("API_HASH")+`",
					"system_language_code": "en-US",
					"device_model": "mxmCherry/go-tdlib",
					"application_version": "debug"
				}`))

			case "authorizationStateWaitPhoneNumber":
				_, _ = fmt.Printf("phone number: ")
				phoneNumber, _ := stdin.ReadString('\n')
				phoneNumber = strings.TrimRight(phoneNumber, "\n")
				tdlib.TDSend(clientID, []byte(`{
					"@type": "setAuthenticationPhoneNumber",
					"phone_number": "`+phoneNumber+`"
				}`))

			case "authorizationStateWaitCode":
				_, _ = fmt.Printf("code: ")
				code, _ := stdin.ReadString('\n')
				code = strings.TrimRight(code, "\n")
				tdlib.TDSend(clientID, []byte(`{
					"@type": "checkAuthenticationCode",
					"code": "`+code+`"
				}`))

			case "authorizationStateReady":
				println("authenticated")
				// now you can TDSend/TDExecute your first auth-requiring request

			} // end switch auth state @type

			// other (non-auth-state) updates - real code to have more `case`-s in this result @type switch

		} // end switch result @type
	}
}
