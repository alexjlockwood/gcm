package gcm

// Message is used by the application server to send a message to
// the GCM server. See the documentation for GCM Architectural
// Overview for more information:
// http://developer.android.com/google/gcm/gcm.html#send-msg
type Message struct {
	RegistrationIDs       []string               `json:"registration_ids"`
	CollapseKey           string                 `json:"collapse_key,omitempty"`
	Data                  map[string]interface{} `json:"data,omitempty"`
	DelayWhileIdle        bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive            int                    `json:"time_to_live,omitempty"`
	RestrictedPackageName string                 `json:"restricted_package_name,omitempty"`
	DryRun                bool                   `json:"dry_run,omitempty"`
}

// NewMessage returns a new Message with the specified payload
// and registration IDs.
func NewMessage(data map[string]interface{}, regIDs ...string) *Message {
	return &Message{RegistrationIDs: regIDs, Data: data}
}

func ProcessGcmResponse(gcmResponse *gcm.Response, gcmMessage *gcm.Message) (map[string]string, map[string]string, map[string]string) {
	invalidDevices := make(map[string]string)
	outdatedDevices := make(map[string]string)
	successDevices := make(map[string]string)

	for i, result := range gcmResponse.Results {
		if result.MessageID != "" {
			outdatedDevices[gcmMessage.RegistrationIDs[i]] = result.MessageID
		}

		// Looking for invalid device
		if result.RegistrationID != "" {
			outdatedDevices[gcmMessage.RegistrationIDs[i]] = result.RegistrationID
		}

		// Looking for invalid device
		if result.Error != "" {
			invalidDevices[gcmMessage.RegistrationIDs[i]] = result.Error
		}
	}

	return invalidDevices, outdatedDevices, successDevices
}
