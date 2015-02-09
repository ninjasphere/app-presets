package model

import (
	"encoding/json"
	"log"
	"testing"
)

func assert(t *testing.T, description string, assertion func() bool) {
	if !assertion() {
		t.Fatalf("assertion failed: %s\n", description)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	item := &Presets{
		Version: "1.0",
		Scenes: []Scene{
			Scene{
				ID:    "CAFE-BABE-0001",
				Slot:  1,
				Label: "Romantic",
				Scope: "sites/a458dfe3-3a81-43cc-a118-6c42c814f4b3",
				Things: []ThingState{
					ThingState{
						ID: "ba8236f9-a813-11e4-8ab9-7c669d02a706",
						Channels: []ChannelState{
							ChannelState{
								ID:    "on-off",
								State: false,
							},
						},
					},
				},
			},
		},
	}

	serialized, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("marhsalling - %s", err)
	}
	deserialized := &Presets{}
	err = json.Unmarshal(serialized, deserialized)
	if err != nil {
		t.Fatalf("unmarhsalling - %s", err)
	}

	assert(t, "id", func() bool { return deserialized.Scenes[0].ID == item.Scenes[0].ID })
	assert(t, "deserialized.Scenes[0].Things[0].Channels[0].State", func() bool { return deserialized.Scenes[0].Things[0].Channels[0].State == false })

}

func TestDeserialization(t *testing.T) {
	serialized := []byte("{\"version\":\"1.0\",\"scenes\":[{\"uuid\":\"CAFE-BABE-0001\",\"slot\":1,\"label\":\"Romantic\",\"scope\":\"sites/a458dfe3-3a81-43cc-a118-6c42c814f4b3\",\"things\":[{\"id\":\"ba8236f9-a813-11e4-8ab9-7c669d02a706\",\"channels\":[{\"id\":\"on-off\",\"state\":{\"a\":\"b\", \"c\": 2}}]}]}]}")
	deserialized := &Presets{}
	err := json.Unmarshal(serialized, deserialized)
	if err != nil {
		t.Fatalf("unmarhsalling - %s", err)
	}
	log.Printf("%v", deserialized.Scenes[0].Things[0].Channels[0].State)
}
