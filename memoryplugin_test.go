package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createFakeMeminfoFile(t *testing.T, dir, meminfoContent string) {
	t.Helper()
	meminfoPath := filepath.Join(dir, "meminfo")
	if err := os.WriteFile(meminfoPath, []byte(meminfoContent), 0644); err != nil {
		t.Fatalf("failed to write meminfo file: %v", err)
	}
}

func TestFetchMetricsFS(t *testing.T) {
	tmp := t.TempDir()
	meminfoContent := `MemTotal:         985796 kB
MemFree:          512116 kB
MemAvailable:     742312 kB
Buffers:           37324 kB
Cached:           325472 kB
SwapCached:         3628 kB
Active:           118988 kB
Inactive:         252300 kB
Active(anon):       3940 kB
Inactive(anon):    19788 kB
Active(file):     115048 kB
Inactive(file):   232512 kB
Unevictable:       14632 kB
Mlocked:           14632 kB
SwapTotal:       4198396 kB
SwapFree:        4066224 kB
Dirty:                28 kB
Writeback:             0 kB
AnonPages:         21540 kB
Mapped:            33772 kB
Shmem:              9936 kB
KReclaimable:      20612 kB
Slab:              53872 kB
SReclaimable:      20612 kB
SUnreclaim:        33260 kB
KernelStack:        1740 kB
PageTables:         5712 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:     4691292 kB
Committed_AS:     410916 kB
VmallocTotal:   34359738367 kB
VmallocUsed:       17572 kB
VmallocChunk:          0 kB
Percpu:              468 kB
HardwareCorrupted:     0 kB
AnonHugePages:      6144 kB
ShmemHugePages:        0 kB
ShmemPmdMapped:        0 kB
FileHugePages:         0 kB
FilePmdMapped:         0 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
Hugetlb:               0 kB
DirectMap4k:      604012 kB
DirectMap2M:      444416 kB
DirectMap1G:           0 kB`
	createFakeMeminfoFile(t, tmp, meminfoContent)
	plugin := LinuxMemoryPlugin{}
	metrics, err := plugin.FetchMetricsFS(tmp)
	require.NoError(t, err)
	expectedMetrics := map[string]float64{
		"total":       985796 * 1024,
		"kernelstack": 1740 * 1024,
		"vmallocused": 17572 * 1024,
		"pagetables":  5712 * 1024,
		"mapped":      33772 * 1024,
		"anonpages":   21540 * 1024,
		"available":   742312 * 1024,
		"slab":        53872 * 1024,
		"buffers":     37324 * 1024,
		"cached":      325472 * 1024,
		"free":        512116 * 1024,
		"used":        (985796 - 742312) * 1024, // used = total - available
	}
	assert.Equal(t, expectedMetrics, metrics)
}

func TestFetchMetricsFS_MissingFields(t *testing.T) {
	tmp := t.TempDir()
	meminfoContent := `MemTotal:       16281828 kB
MemFree:        11377504 kB
Buffers:          232080 kB
Cached:          3183720 kB
SwapCached:            0 kB
Active:          1571996 kB
Inactive:        1901160 kB
Active(anon):      26200 kB
Inactive(anon):    35220 kB
Active(file):    1545796 kB
Inactive(file):  1865940 kB
Unevictable:       13952 kB
Mlocked:            5776 kB
SwapTotal:       4194300 kB
SwapFree:        4194300 kB
Dirty:                76 kB
Writeback:             0 kB
AnonPages:         70660 kB
Mapped:            15424 kB
Shmem:               804 kB
Slab:            1229304 kB
SReclaimable:    1192260 kB
SUnreclaim:        37044 kB
KernelStack:        4080 kB
PageTables:        10488 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:    12335212 kB
Committed_AS:     205140 kB
VmallocTotal:   34359738367 kB
VmallocUsed:      110512 kB
VmallocChunk:   34359600376 kB
HardwareCorrupted:     0 kB
AnonHugePages:     12288 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:       10240 kB
DirectMap2M:     2070528 kB
DirectMap1G:    14680064 kB`
	createFakeMeminfoFile(t, tmp, meminfoContent)
	plugin := LinuxMemoryPlugin{}
	metrics, err := plugin.FetchMetricsFS(tmp)
	require.NoError(t, err)
	// available is missing in this case
	expectedMetrics := map[string]float64{
		"total":       16281828 * 1024,
		"kernelstack": 4080 * 1024,
		"vmallocused": 110512 * 1024,
		"pagetables":  10488 * 1024,
		"mapped":      15424 * 1024,
		"anonpages":   70660 * 1024,
		"slab":        1229304 * 1024,
		"buffers":     232080 * 1024,
		"cached":      3183720 * 1024,
		"free":        11377504 * 1024,
		"used":        (16281828 - 11377504 - 232080 - 3183720) * 1024, // used = total - free - buffers - cached
	}
	assert.Equal(t, expectedMetrics, metrics)
}

func TestMeminfoValue(t *testing.T) {
	var value uint64 = 12345
	result := meminfoValue(&value)
	expected := float64(12345 * 1024)
	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}

	result = meminfoValue(nil)
	expected = 0
	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}

func TestMeminfoSub(t *testing.T) {
	var total uint64 = 10000
	var sub1 uint64 = 2000
	var sub2 uint64 = 3000

	result := meminfoSub(&total, &sub1, &sub2)
	expected := float64((10000 - 2000 - 3000) * 1024)
	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}

	result = meminfoSub(nil, &sub1, &sub2)
	expected = 0
	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}

	result = meminfoSub(&total, nil, &sub2)
	expected = float64((10000 - 0 - 3000) * 1024)
	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}
