package server

import (
	"sync"
	"time"

	agentpb "github.com/thekubefleet/kubefleet/proto"
)

type DataStore struct {
	mu            sync.RWMutex
	agentData     []*agentpb.AgentData
	maxDataPoints int
}

func NewDataStore() *DataStore {
	return &DataStore{
		agentData:     make([]*agentpb.AgentData, 0),
		maxDataPoints: 100, // Keep last 100 data points
	}
}

func (ds *DataStore) StoreAgentData(data *agentpb.AgentData) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// Add timestamp if not present
	if data.Timestamp == 0 {
		data.Timestamp = time.Now().Unix()
	}

	// Add new data
	ds.agentData = append(ds.agentData, data)

	// Keep only the last maxDataPoints
	if len(ds.agentData) > ds.maxDataPoints {
		ds.agentData = ds.agentData[len(ds.agentData)-ds.maxDataPoints:]
	}
}

func (ds *DataStore) GetLatestData() *agentpb.AgentData {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if len(ds.agentData) == 0 {
		return nil
	}

	return ds.agentData[len(ds.agentData)-1]
}

func (ds *DataStore) GetAllData() []*agentpb.AgentData {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]*agentpb.AgentData, len(ds.agentData))
	copy(result, ds.agentData)
	return result
}

func (ds *DataStore) GetDataCount() int {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return len(ds.agentData)
}
