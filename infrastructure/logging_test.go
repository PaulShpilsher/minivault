package infrastructure

import (
	"testing"
	"minivault/mocks"
)

func TestLogger_LogInteraction(t *testing.T) {
	mockLog := &mocks.MockLogger{}
	mockLog.LogInteraction("prompt", "resp")
	if len(mockLog.Interactions) != 1 {
		t.Error("Interaction not logged")
	}
}

func TestLogger_LogError(t *testing.T) {
	mockLog := &mocks.MockLogger{}
	mockLog.LogError("fail", nil)
	if len(mockLog.Errors) != 1 {
		t.Error("Error not logged")
	}
}

func TestLogger_LogWarnAndInfo(t *testing.T) {
	mockLog := &mocks.MockLogger{}
	mockLog.LogWarn("warn1")
	mockLog.LogWarn("warn2")
	mockLog.LogInfo("info1")
	if len(mockLog.Warnings) != 2 || mockLog.Warnings[0] != "warn1" {
		t.Error("Warnings not recorded correctly")
	}
	if len(mockLog.Infos) != 1 || mockLog.Infos[0] != "info1" {
		t.Error("Infos not recorded correctly")
	}
}

func TestLogger_MultipleLogs(t *testing.T) {
	mockLog := &mocks.MockLogger{}
	mockLog.LogInteraction("p1", "r1")
	mockLog.LogInteraction("p2", "r2")
	mockLog.LogError("err1", nil)
	mockLog.LogWarn("warn")
	mockLog.LogInfo("info")
	if len(mockLog.Interactions) != 2 || len(mockLog.Errors) != 1 || len(mockLog.Warnings) != 1 || len(mockLog.Infos) != 1 {
		t.Error("Multiple logs not handled correctly")
	}
}
