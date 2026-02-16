package tui

import (
	"testing"

	"new_era_go/internal/discovery"
	reader18 "new_era_go/internal/protocol/reader18"
	"new_era_go/internal/reader"
)

func TestStartReadingQueuesScanWhenDisconnectedAndNoCandidates(t *testing.T) {
	m := NewModel()
	m.scanning = false
	m.candidates = nil

	next, cmd := m.runControlAction(0)
	if cmd == nil {
		t.Fatal("expected scan/connect command, got nil")
	}

	nm, ok := next.(Model)
	if !ok {
		t.Fatal("unexpected model type")
	}
	if nm.pendingAction != 0 {
		t.Fatalf("expected pendingAction=0, got %d", nm.pendingAction)
	}
	if !nm.scanning {
		t.Fatal("expected scanning=true")
	}
	if !nm.pendingConnect {
		t.Fatal("expected pendingConnect=true")
	}
}

func TestStartReadingQueuesConnectWhenCandidatesExist(t *testing.T) {
	m := NewModel()
	m.scanning = false
	m.pendingConnect = false
	m.candidates = []discovery.Candidate{
		{Host: "192.168.1.200", Port: 6000, Verified: true, Score: 999},
	}

	next, cmd := m.runControlAction(0)
	if cmd == nil {
		t.Fatal("expected connect command, got nil")
	}

	nm := next.(Model)
	if nm.pendingAction != 0 {
		t.Fatalf("expected pendingAction=0, got %d", nm.pendingAction)
	}
	if nm.scanning {
		t.Fatal("expected scanning=false")
	}
	if nm.status == "" {
		t.Fatal("expected status message")
	}
}

func TestOnConnectFinishedRunsPendingStartAction(t *testing.T) {
	m := NewModel()
	m.pendingAction = 0

	next, cmd := m.onConnectFinished(connectFinishedMsg{Endpoint: reader.Endpoint{Host: "10.0.0.5", Port: 6000}, Err: nil})
	if cmd == nil {
		t.Fatal("expected command batch")
	}

	nm := next.(Model)
	if !nm.inventoryRunning {
		t.Fatal("expected inventoryRunning=true")
	}
	if nm.pendingAction != noPendingAction {
		t.Fatalf("expected pendingAction cleared, got %d", nm.pendingAction)
	}
}

func TestOnScanFinishedClearsPendingActionWhenNoCandidates(t *testing.T) {
	m := NewModel()
	m.pendingAction = 2
	m.pendingConnect = true

	next, _ := m.onScanFinished(scanFinishedMsg{Candidates: nil})
	nm := next.(Model)
	if nm.pendingAction != noPendingAction {
		t.Fatalf("expected pendingAction cleared, got %d", nm.pendingAction)
	}
	if nm.pendingConnect {
		t.Fatal("expected pendingConnect=false")
	}
}

func TestSingleInventoryCountsUniqueTagOnlyOnce(t *testing.T) {
	m := NewModel()
	m.inventoryRunning = true

	frame := reader18.Frame{
		Command: reader18.CmdInventorySingle,
		Status:  reader18.StatusNoTag,
		Address: 0x00,
		Data: []byte{
			0x01, 0x01, 0x0C,
			0x30, 0x34, 0x25, 0x7B, 0xF7, 0x19,
			0x4E, 0x40, 0x00, 0x00, 0x00, 0x42,
		},
	}

	m.handleProtocolFrame(frame)
	m.handleProtocolFrame(frame)

	if m.inventoryTagTotal != 1 {
		t.Fatalf("expected unique tag total 1, got %d", m.inventoryTagTotal)
	}
	if len(m.seenTagEPC) != 1 {
		t.Fatalf("expected seen EPC size 1, got %d", len(m.seenTagEPC))
	}
}

func TestLegacyInventoryDoesNotInflateUniqueCounter(t *testing.T) {
	m := NewModel()
	m.inventoryRunning = true
	m.inventoryTagTotal = 0

	frame := reader18.Frame{
		Command: reader18.CmdInventory,
		Status:  reader18.StatusNoTag,
		Data:    []byte{0x01, 0x09},
	}
	m.handleProtocolFrame(frame)

	if m.inventoryTagTotal != 0 {
		t.Fatalf("expected unique tag total unchanged, got %d", m.inventoryTagTotal)
	}
}
