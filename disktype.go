package osmgr

var (
	disks []*Disk
)

func init() {
	disks = ScanDisks()
}

type Disk struct {
	Table      string       `json:"table,omitempty"`
	Model      string       `json:"model,omitempty"`
	Size       string       `json:"size,omitempty"`
	Block      string       `json:"block,omitempty"`
	ID         string       `json:"id,omitempty"`
	Partitions []*Partition `json:"partitions,omitempty"`
}

func GetDisk(diskID string) *Disk {
	for i := 0; i < len(disks); i++ {
		if disks[i].ID == diskID {
			return disks[i]
		}
	}
	return nil
}

func GetDisks() []*Disk {
	return disks
}

func (d *Disk) GetPartition(partID string) *Partition {
	for i := 0; i < len(d.Partitions); i++ {
		if d.Partitions[i].ID == partID {
			return d.Partitions[i]
		}
	}
	return nil
}

func (d *Disk) GetFreePartitionsAvailable() int {
	switch d.Table {
	case "gpt":
		return 128 - len(d.Partitions)
	case "dos":
		return 4 - len(d.Partitions)
	}
	return 0 //We don't know how to operate on this disk yet
}

type Partition struct {
	Disk        *Disk  `json:"-"`
	FS          string `json:"fs,omitempty"`
	Size        string `json:"size,omitempty"`
	OffsetStart int64  `json:"offsetA"`
	OffsetEnd   int64  `json:"offsetB"`
	Block       string `json:"block,omitempty"`
	Mount       string `json:"mount,omitempty"`
	ID          string `json:"id,omitempty"`
}
