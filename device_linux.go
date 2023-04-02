package osmgr

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func ScanDisks() []*Disk {
	d := make([]*Disk, 0)
	var lastDisk *Disk

	// Use the "lsblk" command to list all block devices
	out, err := exec.Command("lsblk", "-o", "NAME,PTTYPE,PTUUID,TYPE,SIZE,MODEL,PARTUUID,FSTYPE,MOUNTPOINT", "-P").Output()
	if err != nil {
		panic(fmt.Sprintf("Error executing lsblk: %v", err))
	}

	// Split the output into lines
	lines := strings.Split(string(out), "\n")

	// Loop through each line
	for _, line := range lines {
		if line == "" {
			continue
		}

		// Split the line into key-value pairs
		pairs := strings.Split(line, "\"")
		if len(pairs) < 6 {
			continue
		}

		block := pairs[1]
		blockTable := pairs[3]
		blockType := pairs[7]
		blockSize := pairs[9]
		blockModel := pairs[11]

		diskID := strings.Replace(pairs[5], "-", "", -1)
		diskIDnum, _ := strconv.ParseInt(diskID, 16, 64)
		diskID = fmt.Sprintf("%d", diskIDnum)

		if blockType == "disk" {
			if lastDisk != nil {
				d = append(d, lastDisk)
			}
			lastDisk = &Disk{
				Table: blockTable,
				Model: blockModel,
				Size:  blockSize,
				Block: block,
				ID:    diskID,
			}
		} else if blockType == "part" {
			partID := pairs[13]
			partFS := pairs[15]
			partMount := pairs[17]

			if blockTable == "dos" {
				partIDnum, err := strconv.ParseInt(partID, 16, 64)
				if err != nil {
					panic(fmt.Sprintf("Failed to parse partition number `%s`for DOS table on disk `%s`: %v", partID, lastDisk.ID, err))
				}
				partID = fmt.Sprintf("%d", partIDnum*512)
			}

			if lastDisk.Partitions == nil {
				lastDisk.Partitions = make([]*Partition, 0)
			}

			lastDisk.Partitions = append(lastDisk.Partitions, &Partition{
				Disk:  lastDisk,
				FS:    partFS,
				Size:  blockSize,
				Block: block,
				Mount: partMount,
				ID:    partID,
			})
		}
	}

	if lastDisk != nil {
		d = append(d, lastDisk)
	}

	return d
}
