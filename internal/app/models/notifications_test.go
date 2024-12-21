package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotificationType(t *testing.T) {
	tests := []struct {
		name string
		typ  NotificationType
		want string
	}{
		{
			name: "slack type",
			typ:  NotificationTypeSlack,
			want: "slack",
		},
		{
			name: "email type",
			typ:  NotificationTypeEmail,
			want: "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, string(tt.typ))
		})
	}
}

func TestConfigData_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    ConfigData
		wantErr bool
	}{
		{
			name:  "nil input",
			input: nil,
			want:  ConfigData{},
		},
		{
			name:  "empty json",
			input: []byte(`{}`),
			want:  ConfigData{},
		},
		{
			name:  "valid slack config",
			input: []byte(`{"webhook_url": "https://hooks.slack.com/test"}`),
			want: ConfigData{
				"webhook_url": "https://hooks.slack.com/test",
			},
		},
		{
			name:  "valid email config",
			input: []byte(`{"recipients": ["test@example.com"]}`),
			want: ConfigData{
				"recipients": []interface{}{"test@example.com"},
			},
		},
		{
			name:    "invalid json",
			input:   []byte(`{invalid}`),
			wantErr: true,
		},
		{
			name:    "invalid type",
			input:   123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c ConfigData
			err := c.Scan(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, c)
			}
		})
	}
}

func TestConfigData_Value(t *testing.T) {
	tests := []struct {
		name    string
		config  ConfigData
		want    string
		wantErr bool
	}{
		{
			name:   "empty config",
			config: ConfigData{},
			want:   "{}",
		},
		{
			name: "slack config",
			config: ConfigData{
				"webhook_url": "https://hooks.slack.com/test",
			},
			want: `{"webhook_url":"https://hooks.slack.com/test"}`,
		},
		{
			name: "email config",
			config: ConfigData{
				"recipients": []string{"test@example.com"},
			},
			want: `{"recipients":["test@example.com"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.Value()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			gotBytes, ok := got.([]byte)
			assert.True(t, ok)
			assert.JSONEq(t, tt.want, string(gotBytes))
		})
	}
}

func TestConfigData_SlackConfig(t *testing.T) {
	tests := []struct {
		name   string
		config ConfigData
		want   string
	}{
		{
			name:   "empty config",
			config: ConfigData{},
			want:   "",
		},
		{
			name: "valid webhook",
			config: ConfigData{
				"webhook_url": "https://hooks.slack.com/test",
			},
			want: "https://hooks.slack.com/test",
		},
		{
			name: "invalid type",
			config: ConfigData{
				"webhook_url": 123,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.SlackConfig()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfigData_EmailRecipients(t *testing.T) {
	tests := []struct {
		name   string
		config ConfigData
		want   []string
	}{
		{
			name:   "empty config",
			config: ConfigData{},
			want:   []string{},
		},
		{
			name: "single recipient",
			config: ConfigData{
				"recipients": []interface{}{"test@example.com"},
			},
			want: []string{"test@example.com"},
		},
		{
			name: "multiple recipients",
			config: ConfigData{
				"recipients": []interface{}{"test1@example.com", "test2@example.com"},
			},
			want: []string{"test1@example.com", "test2@example.com"},
		},
		{
			name: "invalid type",
			config: ConfigData{
				"recipients": "test@example.com", // string instead of array
			},
			want: []string{},
		},
		{
			name: "mixed types",
			config: ConfigData{
				"recipients": []interface{}{"test@example.com", 123, true},
			},
			want: []string{"test@example.com"}, // only valid emails included
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.EmailRecipients()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNotificationConfig_Integration(t *testing.T) {
	// Test full notification config with database serialization
	config := &NotificationConfig{
		ID:   1,
		Type: NotificationTypeSlack,
		Name: "Test Slack",
		Config: ConfigData{
			"webhook_url": "https://hooks.slack.com/test",
		},
		Enabled:   true,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// Test JSON marshaling/unmarshaling
	jsonBytes, err := json.Marshal(config)
	assert.NoError(t, err)

	var decoded NotificationConfig
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err)

	// Compare original and decoded
	assert.Equal(t, config.ID, decoded.ID)
	assert.Equal(t, config.Type, decoded.Type)
	assert.Equal(t, config.Name, decoded.Name)
	assert.Equal(t, config.Config["webhook_url"], decoded.Config["webhook_url"])
	assert.Equal(t, config.Enabled, decoded.Enabled)
	assert.Equal(t, config.CreatedAt, decoded.CreatedAt)
	assert.Equal(t, config.UpdatedAt, decoded.UpdatedAt)
}
