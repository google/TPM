package entrypoints_test

import (
	"crypto/rand"
	"encoding/binary"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/google/TPM/bindings/go/entrypoints"
	"github.com/google/TPM/bindings/go/platform"
)

type mockTimer time.Time

func (t mockTimer) Ticks() time.Duration         { return time.Since(time.Time(t)) }
func (t mockTimer) WasReset() bool               { return false }
func (t mockTimer) WasStopped() bool             { return false }
func (t mockTimer) Adjust(a platform.Adjustment) {}

func init() {
	platform.Current = &platform.Platform{
		Entropy: rand.Reader,
		Timer:   mockTimer(time.Now()),
		Debug:   os.Stderr,
	}
}

func TestHCRTM(t *testing.T) {
	entrypoints.Init()
	data := [][]byte{
		{1, 2, 3},
		make([]byte, 0, 20),
		{5, 5, 5, 5, 5},
		nil,
	}
	if err := entrypoints.HashSequence(slices.Values(data)); err != nil {
		t.Errorf("HashSequence failed: %v", err)
	}
}

func TestEmptyCmd(t *testing.T) {
	entrypoints.Init()
	rsp := entrypoints.ExecuteCommand(nil)
	if len(rsp) != 10 {
		t.Errorf("Expected response of %d bytes, got %d", 10, len(rsp))
	}
}

func TestManufacture(t *testing.T) {
	if err := entrypoints.Manufacture(); err != nil {
		t.Errorf("Manufacture failed with: %v", err)
	}
	if err := entrypoints.Manufacture(); err != nil {
		t.Errorf("Remanufacture failed with: %v", err)
	}
}

func TestCommands(t *testing.T) {
	// Ensure TPM is manufactured and initialized
	entrypoints.Manufacture()
	entrypoints.Init()

	// 1. Startup (Clear)
	startupCmd := []byte{
		0x80, 0x01, // tag: TPM_ST_NO_SESSIONS
		0x00, 0x00, 0x00, 0x0C, // size: 12
		0x00, 0x00, 0x01, 0x44, // code: TPM_CC_Startup
		0x00, 0x00, // startupType: TPM_SU_CLEAR
	}
	rsp := entrypoints.ExecuteCommand(startupCmd)
	t.Logf("Startup response: %x", rsp)
	if len(rsp) < 10 || binary.BigEndian.Uint32(rsp[6:10]) != 0 {
		t.Errorf("Startup failed: %x", rsp)
	}

	// 2. GetRandom (16 bytes)
	getRandomCmd := []byte{
		0x80, 0x01, // tag: TPM_ST_NO_SESSIONS
		0x00, 0x00, 0x00, 0x0C, // size: 12
		0x00, 0x00, 0x01, 0x7B, // code: TPM_CC_GetRandom
		0x00, 0x10, // bytesRequested: 16
	}
	rsp = entrypoints.ExecuteCommand(getRandomCmd)
	t.Logf("GetRandom response: %x", rsp)
	if len(rsp) < 12 || binary.BigEndian.Uint32(rsp[6:10]) != 0 {
		t.Errorf("GetRandom failed: %x", rsp)
	}
	dataSize := binary.BigEndian.Uint16(rsp[10:12])
	if dataSize != 16 || len(rsp) != 12+16 {
		t.Errorf("Expected 16 bytes of random data, got %d", dataSize)
	}

	// 3. Shutdown (Clear)
	shutdownCmd := []byte{
		0x80, 0x01, // tag: TPM_ST_NO_SESSIONS
		0x00, 0x00, 0x00, 0x0C, // size: 12
		0x00, 0x00, 0x01, 0x45, // code: TPM_CC_Shutdown
		0x00, 0x00, // shutdownType: TPM_SU_CLEAR
	}
	rsp = entrypoints.ExecuteCommand(shutdownCmd)
	t.Logf("Shutdown response: %x", rsp)
	if len(rsp) < 10 || binary.BigEndian.Uint32(rsp[6:10]) != 0 {
		t.Errorf("Shutdown failed: %x", rsp)
	}
}
