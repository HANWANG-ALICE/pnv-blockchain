package AccountManager

type NodeID = uint64

type AccountManager struct {
	WorkerNumberSet    map[uint32]string
	VoterNumberSet	   map[uint32]string
	VoterSet           map[string]uint64
	WorkerSet          map[string]uint64
	WorkerCandidateSet map[string]uint64

	WorkerCandidateList []string
}
