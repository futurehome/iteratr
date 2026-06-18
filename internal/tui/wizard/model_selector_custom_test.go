package wizard

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// TestModelSelector_CustomMode_EnterSubmits verifies that pressing 'c' enters
// custom mode, and that pressing Enter with a typed model emits a
// ModelSelectedMsg with that model ID. This is the path users need to enter
// a model that isn't in the opencode catalog (e.g. custom/<name> with an
// ITERATR_PROVIDER_URL override).
func TestModelSelector_CustomMode_EnterSubmits(t *testing.T) {
	step := NewModelSelectorStep()
	step.loading = false
	step.allModels = []*ModelInfo{
		{id: "test/model-1", displayName: "test/model-1", providerID: "test"},
	}
	step.buildGroupedList()

	// Initially not in custom mode
	if step.isCustomMode {
		t.Fatal("expected isCustomMode to be false initially")
	}

	// Press 'c' to enter custom mode
	cmd := step.Update(tea.KeyPressMsg{Code: 'c', Text: "c"})
	if cmd == nil {
		t.Fatal("expected a command after pressing 'c'")
	}
	if !step.isCustomMode {
		t.Fatal("expected isCustomMode to be true after pressing 'c'")
	}

	// View should show the custom mode UI
	view := step.View()
	for _, want := range []string{"Enter Custom Model", "Model ID", "ESC cancel"} {
		if !containsString(view, want) {
			t.Errorf("expected view to mention %q; got:\n%s", want, view)
		}
	}

	// Type a custom model directly into the input
	step.customInput.SetValue("custom/qwen3_5-122b")

	// Press Enter — should emit a ModelSelectedMsg
	cmd = step.Update(tea.KeyPressMsg{Code: tea.KeyEnter, Text: "enter"})
	if cmd == nil {
		t.Fatal("expected a command after pressing Enter in custom mode")
	}
	msg := cmd()
	sel, ok := msg.(ModelSelectedMsg)
	if !ok {
		t.Fatalf("expected ModelSelectedMsg, got %T: %v", msg, msg)
	}
	if sel.ModelID != "custom/qwen3_5-122b" {
		t.Errorf("expected ModelID %q, got %q", "custom/qwen3_5-122b", sel.ModelID)
	}
}

// TestModelSelector_CustomMode_EscCancels verifies that pressing Esc in
// custom mode returns to the normal list view and clears the typed value.
func TestModelSelector_CustomMode_EscCancels(t *testing.T) {
	step := NewModelSelectorStep()
	step.loading = false
	step.allModels = []*ModelInfo{
		{id: "test/model-1", displayName: "test/model-1", providerID: "test"},
	}
	step.buildGroupedList()

	// Enter custom mode
	step.Update(tea.KeyPressMsg{Code: 'c', Text: "c"})
	if !step.isCustomMode {
		t.Fatal("expected to be in custom mode")
	}
	step.customInput.SetValue("partial input")

	// Press Esc
	step.Update(tea.KeyPressMsg{Code: tea.KeyEscape, Text: "esc"})

	if step.isCustomMode {
		t.Error("expected to have left custom mode after Esc")
	}
	if got := step.customInput.Value(); got != "" {
		t.Errorf("expected custom input cleared after Esc, got %q", got)
	}

	// View should be back to the normal list view (not "Enter Custom Model")
	view := step.View()
	if containsString(view, "Enter Custom Model") {
		t.Errorf("expected normal list view after Esc; got:\n%s", view)
	}
}

// TestModelSelector_CustomMode_EmptyEnterIsNoop verifies that pressing Enter
// in custom mode with an empty input does not emit a ModelSelectedMsg.
func TestModelSelector_CustomMode_EmptyEnterIsNoop(t *testing.T) {
	step := NewModelSelectorStep()
	step.loading = false
	step.allModels = []*ModelInfo{
		{id: "test/model-1", displayName: "test/model-1", providerID: "test"},
	}
	step.buildGroupedList()

	step.Update(tea.KeyPressMsg{Code: 'c', Text: "c"})
	// Don't set any value; just press Enter
	cmd := step.Update(tea.KeyPressMsg{Code: tea.KeyEnter, Text: "enter"})

	// cmd is allowed to be nil (no model to submit) — but should NOT
	// produce a ModelSelectedMsg. The simplest check: if cmd is non-nil,
	// the resulting message should not be ModelSelectedMsg.
	if cmd != nil {
		if msg := cmd(); msg != nil {
			if _, isSel := msg.(ModelSelectedMsg); isSel {
				t.Errorf("did not expect ModelSelectedMsg on empty Enter, got: %+v", msg)
			}
		}
	}
}

// TestModelSelector_PreferredHeight_CustomMode verifies the modal uses a
// fixed height when in custom mode (independent of the opencode catalog size).
func TestModelSelector_PreferredHeight_CustomMode(t *testing.T) {
	step := NewModelSelectorStep()
	step.loading = false
	step.allModels = []*ModelInfo{
		{id: "a/1", displayName: "a/1", providerID: "a"},
		{id: "b/2", displayName: "b/2", providerID: "b"},
		{id: "c/3", displayName: "c/3", providerID: "c"},
		{id: "d/4", displayName: "d/4", providerID: "d"},
		{id: "e/5", displayName: "e/5", providerID: "e"},
	}
	step.buildGroupedList()

	customHeight := step.PreferredHeight()
	step.Update(tea.KeyPressMsg{Code: 'c', Text: "c"})
	customModeHeight := step.PreferredHeight()
	if customModeHeight != 5 {
		t.Errorf("expected custom mode height to be 5, got %d", customModeHeight)
	}
	if customModeHeight == customHeight {
		t.Logf("note: normal and custom heights happen to match (%d)", customHeight)
	}
}

// containsString is a tiny helper to avoid pulling in strings.Contains
// dependencies just for a couple of test assertions.
func containsString(haystack, needle string) bool {
	if len(needle) == 0 {
		return true
	}
	if len(needle) > len(haystack) {
		return false
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
