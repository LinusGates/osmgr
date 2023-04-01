package osmgr

// OS holds methods to interact with a specific operating system install
type OS interface {
	//Filled in by implementation for OS detection and install steps
	//- Filesystems not yet understood simply won't be available for use
	GetPartitionTables() []string //Usually just gpt and/or dos
	GetBootFilesystems() []string //Return nothing to skip installing/managing a bootloader
	GetBootSize() int64
	GetSystemFilesystems() []string //Return nothing to skip installing a system image
	GetSystemSize() int64
	GetRecoveryFilesystems() []string //Return nothing if recovery files cannot be moved, defaults to system
	GetRecoverySize() int64
	GetHomeFilesystems() []string //Return nothing if home folders cannot be moved, defaults to system
	GetHomeSize() int64
	GetTempFilesystems() []string //Return nothing if temp folders cannot be moved, defaults to system
	GetTempSize() int64
	GetSwap() bool //Return true if swap partition/file support is available

	//Called during the installation phases
	InstallBoot(*Partition)     //Used to install the bootloader, Infer MBR/GPT from *Partition.Disk.Table
	InstallSystem(*Partition)   //Used to install or upgrade the main system image
	InstallRecovery(*Partition) //Used to install or upgrade the recovery image, defaults to system partition
	InstallHome(*Partition)     //Used to migrate or otherwise be aware of home files, defaults to system partition
	InstallTemp(*Partition)     //Defaults to system partition, but useful if alternative temp folders are available
	InstallSwap(*Partition)     //Used to enable a swap partition if one was specified, defaults to temp partition

	//Called to determine if a given disk or a given partition space is valid for install
	//- Calculates a default partition layout to fill in the given space
	//- Specifying a disk should assume the disk's entire space is available for calculations
	//- Specifying a partition should assume the user selected to use this partition's space for the install
	//- Use either *Disk.GetFreePartitionsAvailable() or *Partition.Disk.GetFreePartitionsAvailable() so you don't exceed capacity
	//- Return nil or an empty partition slice to signal an impossibility to install
	GetFreshTableRecommendationDisk(*Disk) []*Partition
	GetFreshTableRecommendationPartition(*Partition) []*Partition
}
