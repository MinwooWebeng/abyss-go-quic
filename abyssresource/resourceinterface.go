package abyssresource

import (
	"time"
)

type IResourceBase interface {
	GetMIME() string
	GetModifyDateUTC() string
	Abandon() //application will never access to this resource after this.
	//this is called by resource manager internally when the application calls ReturnResource(*IResourceBase).
}

type ResourceBase struct { //metadata of resource.
	last_modified time.Time
	//TODO: add more common metadata
}

func (r *ResourceBase) GetModifyDateUTC() time.Time {
	return r.last_modified
}

// type IBlobResource interface {
// 	IResourceBase
// 	GetTotalSize() uint32
// 	GetLoadedSize() uint32
// 	Read(p []byte, offset uint32, length uint32) (n int, err error)
// }
// type IDatagramResource interface {
// 	IResourceBase
// 	GetMaxLength() uint16
// 	Read(p []byte) (n int, err error)
// }
// type IStreamResource interface {
// 	IResourceBase
// 	bufio.Reader
// }
// type IMmapResource interface {
// 	IResourceBase
// 	GetMemorySlice() []uint64
// }
