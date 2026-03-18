package cli

import "testing"

func TestGetCommandUX(t *testing.T) {
	tests := []struct {
		command string
		mode    string
	}{
		{command: cmdSource, mode: "проверка источника данных"},
		{command: cmdCoastline, mode: "анализ реальных данных"},
		{command: cmdParadox, mode: "синтетическая демонстрация"},
		{command: cmdKoch, mode: "синтетическая демонстрация"},
		{command: cmdKochOrganic, mode: "синтетическая демонстрация"},
		{command: cmdDimension, mode: "синтетическая демонстрация"},
		{command: cmdAll, mode: "смешанный сценарий"},
	}

	for _, test := range tests {
		t.Run(test.command, func(t *testing.T) {
			ux := getCommandUX(test.command)
			if ux.Mode != test.mode {
				t.Fatalf("expected mode %q, got %q", test.mode, ux.Mode)
			}
			if ux.Summary == "" {
				t.Fatalf("expected summary for %q", test.command)
			}
			if ux.RuntimeNote == "" {
				t.Fatalf("expected runtime note for %q", test.command)
			}
		})
	}
}
