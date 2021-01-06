package transport

import (
	"encoding/json"
	"testing"
)

func TestNodeStatus_UnmarshalJSON(t *testing.T) {
	data := []byte(`
{
  "blockchainHeight": 15246218,
  "stateHeight": 15246218,
  "updatedTimestamp": 1609910738015726800,
  "updatedDate": "2021-01-06T05:25:38.015Z"
}
`)
	var resp NodeStatus
	err := json.Unmarshal(data, &resp)
	if err != nil {
		t.Fatal(err)
	}
}
