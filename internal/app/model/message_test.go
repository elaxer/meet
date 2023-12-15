package model

import "testing"

func TestMessage_Validate(t *testing.T) {
	type fields struct {
		Text string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Empty message",
			fields{"       "},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Text: tt.fields.Text,
			}
			if err := m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Message.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
