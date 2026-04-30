// Package viplite provides Go bindings for the Vivante VIP Lite NPU driver API.
package viplite

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lNBGlinker -lVIPhal
#cgo LDFLAGS: -Wl,-rpath=/usr/lib:.

#include <vip_lite.h>
#include <vip_lite_common.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// =============================================================================
//  Go type aliases for C vip_* typedefs <vip_lite_common.h>
// =============================================================================

type Uint8 = C.vip_uint8_t
type Uint16 = C.vip_uint16_t
type Uint32 = C.vip_uint32_t
type Uint64 = C.vip_uint64_t
type Int8 = C.vip_int8_t
type Int16 = C.vip_int16_t
type Int32 = C.vip_int32_t
type Int64 = C.vip_int64_t
type Char = C.vip_char_t
type Float = C.vip_float_t
type Enum = C.vip_enum
type Ptr = C.vip_ptr
type Float64 = C.vip_float64_t
type Address = C.vip_address_t

// =============================================================================
//  Opaque handle types
// =============================================================================

type Network C.vip_network
type Buffer C.vip_buffer
type Group C.vip_group

// =============================================================================
//  Enum: vip_bool_e
// =============================================================================

const (
	False = C.vip_false_e
	True  = C.vip_true_e
)

// =============================================================================
//  Enum: vip_status_e
// =============================================================================

type Status C.vip_status_e

const (
	StatusErrorCanceled            Status = C.VIP_ERROR_CANCELED
	StatusErrorRecovery            Status = C.VIP_ERROR_RECOVERY
	StatusErrorPowerStop           Status = C.VIP_ERROR_POWER_STOP
	StatusErrorPowerOff            Status = C.VIP_ERROR_POWER_OFF
	StatusErrorFailure             Status = C.VIP_ERROR_FAILURE
	StatusErrorNetworkIncompatible Status = C.VIP_ERROR_NETWORK_INCOMPATIBLE
	StatusErrorNetworkNotPrepared  Status = C.VIP_ERROR_NETWORK_NOT_PREPARED
	StatusErrorMissingInputOutput  Status = C.VIP_ERROR_MISSING_INPUT_OUTPUT
	StatusErrorInvalidNetwork      Status = C.VIP_ERROR_INVALID_NETWORK
	StatusErrorOutOfMemory         Status = C.VIP_ERROR_OUT_OF_MEMORY
	StatusErrorOutOfResource       Status = C.VIP_ERROR_OUT_OF_RESOURCE
	StatusErrorNotSupported        Status = C.VIP_ERROR_NOT_SUPPORTED
	StatusErrorInvalidArguments    Status = C.VIP_ERROR_INVALID_ARGUMENTS
	StatusErrorIO                  Status = C.VIP_ERROR_IO
	StatusErrorTimeout             Status = C.VIP_ERROR_TIMEOUT
	StatusSuccess                  Status = C.VIP_SUCCESS
)

func (s Status) Error() string {
	switch s {
	case StatusErrorCanceled:
		return "network canceled"
	case StatusErrorRecovery:
		return "hardware recovery"
	case StatusErrorPowerStop:
		return "hardware stopped"
	case StatusErrorPowerOff:
		return "hardware power off"
	case StatusErrorFailure:
		return "failure"
	case StatusErrorNetworkIncompatible:
		return "network incompatible"
	case StatusErrorNetworkNotPrepared:
		return "network not prepared"
	case StatusErrorMissingInputOutput:
		return "missing input/output"
	case StatusErrorInvalidNetwork:
		return "invalid network"
	case StatusErrorOutOfMemory:
		return "out of memory"
	case StatusErrorOutOfResource:
		return "out of resource"
	case StatusErrorNotSupported:
		return "not supported"
	case StatusErrorInvalidArguments:
		return "invalid arguments"
	case StatusErrorIO:
		return "IO error"
	case StatusErrorTimeout:
		return "timeout"
	case StatusSuccess:
		return "success"
	default:
		return fmt.Sprintf("unknown status: %d", int32(s))
	}
}

// checkStatus returns an error if s is not VIP_SUCCESS.
func checkStatus(s Status) error {
	if s != StatusSuccess {
		return s
	}
	return nil
}

// =============================================================================
//  Enum: vip_buffer_format_e
// =============================================================================

type BufferFormat C.vip_buffer_format_e

const (
	BufferFormatFP32   BufferFormat = C.VIP_BUFFER_FORMAT_FP32
	BufferFormatFP16   BufferFormat = C.VIP_BUFFER_FORMAT_FP16
	BufferFormatUint8  BufferFormat = C.VIP_BUFFER_FORMAT_UINT8
	BufferFormatInt8   BufferFormat = C.VIP_BUFFER_FORMAT_INT8
	BufferFormatUint16 BufferFormat = C.VIP_BUFFER_FORMAT_UINT16
	BufferFormatInt16  BufferFormat = C.VIP_BUFFER_FORMAT_INT16
	BufferFormatChar   BufferFormat = C.VIP_BUFFER_FORMAT_CHAR
	BufferFormatBFP16  BufferFormat = C.VIP_BUFFER_FORMAT_BFP16
	BufferFormatInt32  BufferFormat = C.VIP_BUFFER_FORMAT_INT32
	BufferFormatUint32 BufferFormat = C.VIP_BUFFER_FORMAT_UINT32
	BufferFormatInt64  BufferFormat = C.VIP_BUFFER_FORMAT_INT64
	BufferFormatUint64 BufferFormat = C.VIP_BUFFER_FORMAT_UINT64
	BufferFormatFP64   BufferFormat = C.VIP_BUFFER_FORMAT_FP64
	BufferFormatBool8  BufferFormat = C.VIP_BUFFER_FORMAT_BOOL8
)

// =============================================================================
//  Enum: vip_buffer_quantize_format_e
// =============================================================================

type BufferQuantizeFormat C.vip_buffer_quantize_format_e

const (
	BufferQuantizeNone           BufferQuantizeFormat = C.VIP_BUFFER_QUANTIZE_NONE
	BufferQuantizeDynamicFixedPt BufferQuantizeFormat = C.VIP_BUFFER_QUANTIZE_DYNAMIC_FIXED_POINT
	BufferQuantizeTFAsymm        BufferQuantizeFormat = C.VIP_BUFFER_QUANTIZE_TF_ASYMM
	BufferQuantizeMax            BufferQuantizeFormat = C.VIP_BUFFER_QUANTIZE_MAX
)

// =============================================================================
//  Enum: vip_buffer_memory_type_e
// =============================================================================

type BufferMemoryType C.vip_buffer_memory_type_e

const (
	BufferMemoryTypeDefault BufferMemoryType = C.VIP_BUFFER_MEMORY_TYPE_DEFAULT
	BufferMemoryTypeHost    BufferMemoryType = C.VIP_BUFFER_MEMORY_TYPE_HOST
	BufferMemoryTypeDMABuf  BufferMemoryType = C.VIP_BUFFER_MEMORY_TYPE_DMA_BUF
	BufferMemoryTypeMax     BufferMemoryType = C.VIP_BUFFER_MEMORY_TYPE_MAX
)

// =============================================================================
//  Enum: vip_create_network_type_e
// =============================================================================

type CreateNetworkType C.vip_create_network_type_e

const (
	CreateNetworkFromNone   CreateNetworkType = C.VIP_CREATE_NETWORK_FROM_NONE
	CreateNetworkFromFile   CreateNetworkType = C.VIP_CREATE_NETWORK_FROM_FILE
	CreateNetworkFromMemory CreateNetworkType = C.VIP_CREATE_NETWORK_FROM_MEMORY
	CreateNetworkFromFlash  CreateNetworkType = C.VIP_CREATE_NETWORK_FROM_FLASH
	CreateNetworkMax        CreateNetworkType = C.VIP_CREATE_NETWORK_MAX
)

// =============================================================================
//  Enum: vip_dup_network_type_e
// =============================================================================

type DupNetworkType C.vip_dup_network_type_e

const (
	DupNone            DupNetworkType = C.VIP_DUP_NONE
	DupForCmdByNetwork DupNetworkType = C.VIP_DUP_FOR_CMD_BY_NETWORK
	DupForCmdByNBG     DupNetworkType = C.VIP_DUP_FOR_CMD_BY_NBG
	DupFromNBGFile     DupNetworkType = C.VIP_DUP_FROM_NBG_FILE
	DupFromNBGMemory   DupNetworkType = C.VIP_DUP_FROM_NBG_MEMORY
	DupFromNBGFlash    DupNetworkType = C.VIP_DUP_FROM_NBG_FLASH
	DupFromNetwork     DupNetworkType = C.VIP_DUP_FROM_NETWORK
	DupNetworkMax      DupNetworkType = C.VIP_DUP_NETWORK_MAX
)

// =============================================================================
//  Enum: vip_power_property_e
// =============================================================================

type PowerProperty C.vip_power_property_e

const (
	PowerPropertyNone         PowerProperty = C.VIP_POWER_PROPERTY_NONE
	PowerPropertySetFrequency PowerProperty = C.VIP_POWER_PROPERTY_SET_FREQUENCY
	PowerPropertyOff          PowerProperty = C.VIP_POWER_PROPERTY_OFF
	PowerPropertyOn           PowerProperty = C.VIP_POWER_PROPERTY_ON
	PowerPropertyStop         PowerProperty = C.VIP_POWER_PROPERTY_STOP
	PowerPropertyStart        PowerProperty = C.VIP_POWER_PROPERTY_START
	PowerPropertyMax          PowerProperty = C.VIP_POWER_PROPERTY_MAX
)

// =============================================================================
//  Enum: vip_query_hardware_property_e
// =============================================================================

type QueryHWProperty C.vip_query_hardware_property_e

const (
	QueryHWPropCID         QueryHWProperty = C.VIP_QUERY_HW_PROP_CID
	QueryHWPropDeviceCount QueryHWProperty = C.VIP_QUERY_HW_PROP_DEVICE_COUNT
	QueryHWPropCoreCount   QueryHWProperty = C.VIP_QUERY_HW_PROP_CORE_COUNT_EACH_DEVICE
	QueryHWPropMax         QueryHWProperty = C.VIP_QUERY_HW_PROP_MAX
)

// =============================================================================
//  Enum: vip_network_property_e
// =============================================================================

type NetworkProperty C.vip_network_property_e

const (
	NetworkPropLayerCount     NetworkProperty = C.VIP_NETWORK_PROP_LAYER_COUNT
	NetworkPropInputCount     NetworkProperty = C.VIP_NETWORK_PROP_INPUT_COUNT
	NetworkPropOutputCount    NetworkProperty = C.VIP_NETWORK_PROP_OUTPUT_COUNT
	NetworkPropNetworkName    NetworkProperty = C.VIP_NETWORK_PROP_NETWORK_NAME
	NetworkPropAddressInfo    NetworkProperty = C.VIP_NETWORK_PROP_ADDRESS_INFO
	NetworkPropReadRegIRQ     NetworkProperty = C.VIP_NETWORK_PROP_READ_REG_IRQ
	NetworkPropMemoryPoolSize NetworkProperty = C.VIP_NETWORK_PROP_MEMORY_POOL_SIZE
	NetworkPropProfiling      NetworkProperty = C.VIP_NETWORK_PROP_PROFILING
	NetworkPropCoreCount      NetworkProperty = C.VIP_NETWORK_PROP_CORE_COUNT
	// set properties
	NetworkPropChangePPUParam NetworkProperty = C.VIP_NETWORK_PROP_CHANGE_PPU_PARAM
	NetworkPropSetMemoryPool  NetworkProperty = C.VIP_NETWORK_PROP_SET_MEMORY_POOL
	NetworkPropSetDeviceID    NetworkProperty = C.VIP_NETWORK_PROP_SET_DEVICE_ID
	NetworkPropSetPriority    NetworkProperty = C.VIP_NETWORK_PROP_SET_PRIORITY
	NetworkPropSetTimeout     NetworkProperty = C.VIP_NETWORK_PROP_SET_TIME_OUT
	NetworkPropSetCoeffMemory NetworkProperty = C.VIP_NETWORK_PROP_SET_COEFF_MEMORY
)

// =============================================================================
//  Enum: vip_group_property_e
// =============================================================================

type GroupProperty C.vip_group_property_e

const (
	GroupPropProfiling   GroupProperty = C.VIP_GROUP_PROP_PROFILING
	GroupPropSetDeviceID GroupProperty = C.VIP_GROUP_PROP_SET_DEVICE_ID
	GroupPropSetTimeout  GroupProperty = C.VIP_GROUP_PROP_SET_TIME_OUT
)

// =============================================================================
//  Enum: vip_buffer_property_e
// =============================================================================

type BufferProperty C.vip_buffer_property_e

const (
	BufferPropQuantFormat    BufferProperty = C.VIP_BUFFER_PROP_QUANT_FORMAT
	BufferPropNumOfDimension BufferProperty = C.VIP_BUFFER_PROP_NUM_OF_DIMENSION
	BufferPropSizesOfDim     BufferProperty = C.VIP_BUFFER_PROP_SIZES_OF_DIMENSION
	BufferPropDataFormat     BufferProperty = C.VIP_BUFFER_PROP_DATA_FORMAT
	BufferPropFixedPointPos  BufferProperty = C.VIP_BUFFER_PROP_FIXED_POINT_POS
	BufferPropTFScale        BufferProperty = C.VIP_BUFFER_PROP_TF_SCALE
	BufferPropTFZeroPoint    BufferProperty = C.VIP_BUFFER_PROP_TF_ZERO_POINT
	BufferPropName           BufferProperty = C.VIP_BUFFER_PROP_NAME
)

// =============================================================================
//  Enum: vip_buffer_operation_type_e
// =============================================================================

type BufferOpType C.vip_buffer_operation_type_e

const (
	BufferOpNone       BufferOpType = C.VIP_BUFFER_OPER_TYPE_NONE
	BufferOpFlush      BufferOpType = C.VIP_BUFFER_OPER_TYPE_FLUSH
	BufferOpInvalidate BufferOpType = C.VIP_BUFFER_OPER_TYPE_INVALIDATE
	BufferOpMax        BufferOpType = C.VIP_BUFFER_OPER_TYPE_MAX
)

// =============================================================================
//  Struct: vip_buffer_create_params_t
// =============================================================================

type BufferCreateParams C.vip_buffer_create_params_t

// NewBufferCreateParams initialises a BufferCreateParams with defaults.
func NewBufferCreateParams() *BufferCreateParams {
	p := new(BufferCreateParams)
	*p = BufferCreateParams{}
	return p
}

func (p *BufferCreateParams) NumOfDims() uint32     { return uint32(p.num_of_dims) }
func (p *BufferCreateParams) SetNumOfDims(n uint32) { p.num_of_dims = C.uint(n) }

func (p *BufferCreateParams) Sizes() [6]uint32 {
	var s [6]uint32
	for i := range p.sizes {
		s[i] = uint32(p.sizes[i])
	}
	return s
}
func (p *BufferCreateParams) SetSizes(s []uint32) {
	for i := 0; i < len(s) && i < 6; i++ {
		p.sizes[i] = C.uint(s[i])
	}
}

func (p *BufferCreateParams) DataFormat() BufferFormat     { return BufferFormat(p.data_format) }
func (p *BufferCreateParams) SetDataFormat(f BufferFormat) { p.data_format = C.vip_enum(f) }

func (p *BufferCreateParams) QuantFormat() BufferQuantizeFormat {
	return BufferQuantizeFormat(p.quant_format)
}
func (p *BufferCreateParams) SetQuantFormat(f BufferQuantizeFormat) { p.quant_format = C.vip_enum(f) }

func (p *BufferCreateParams) MemoryType() BufferMemoryType     { return BufferMemoryType(p.memory_type) }
func (p *BufferCreateParams) SetMemoryType(t BufferMemoryType) { p.memory_type = C.uint(t) }

// SetDFPFixedPointPos sets the dynamic-fixed-point position in the quant_data union.
func (p *BufferCreateParams) SetDFPFixedPointPos(pos int32) {
	*(*C.vip_int32_t)(unsafe.Pointer(&p.quant_data[0])) = C.vip_int32_t(pos)
}
func (p *BufferCreateParams) DFPFixedPointPos() int32 {
	return int32(*(*C.vip_int32_t)(unsafe.Pointer(&p.quant_data[0])))
}

// SetAffineScale sets the affine scale in the quant_data union.
func (p *BufferCreateParams) SetAffineScale(s float32) {
	*(*C.vip_float_t)(unsafe.Pointer(&p.quant_data[0])) = C.vip_float_t(s)
}
func (p *BufferCreateParams) AffineScale() float32 {
	return float32(*(*C.vip_float_t)(unsafe.Pointer(&p.quant_data[0])))
}

// SetAffineZeroPoint sets the affine zero-point in the quant_data union (at offset 4).
func (p *BufferCreateParams) SetAffineZeroPoint(zp int32) {
	*(*C.vip_int32_t)(unsafe.Pointer(&p.quant_data[4])) = C.vip_int32_t(zp)
}
func (p *BufferCreateParams) AffineZeroPoint() int32 {
	return int32(*(*C.vip_int32_t)(unsafe.Pointer(&p.quant_data[4])))
}

// =============================================================================
//  Struct: vip_power_frequency_t
// =============================================================================

type PowerFrequency C.vip_power_frequency_t

func NewPowerFrequency(pct uint8) *PowerFrequency {
	return &PowerFrequency{fscale_percent: C.uchar(pct)}
}

// =============================================================================
//  Struct: vip_inference_profile_t
// =============================================================================

type InferenceProfile C.vip_inference_profile_t

func (p *InferenceProfile) InferenceTime() uint32 { return uint32(p.inference_time) }
func (p *InferenceProfile) TotalCycle() uint32    { return uint32(p.total_cycle) }

// =============================================================================
//  Struct: vip_ppu_param_t
// =============================================================================

type PPUParam C.vip_ppu_param_t

// =============================================================================
//  Global API
// =============================================================================

// GetVersion returns the VIP Lite driver version.
func GetVersion() uint32 {
	return uint32(C.vip_get_version())
}

// Init initialises the VIP hardware.
// When building with LIBVIP_VERSION_85X, use Init85X(memSize) instead.
func Init() error {
	return checkStatus(Status(C.vip_init()))
}

// Init85X initialises the VIP hardware with a specified video memory size (in bytes).
// Requires LIBVIP_VERSION_85X to be defined during compilation.
// func Init85X(memSize uint32) error {
// 	return checkStatus(Status(C.vip_init(C.uint(memSize))))
// }

// Destroy terminates the VIP Lite driver and shuts down hardware.
func Destroy() error {
	return checkStatus(Status(C.vip_destroy()))
}

// QueryHardware queries hardware capability information.
func QueryHardware(property QueryHWProperty, value unsafe.Pointer, size uint32) error {
	return checkStatus(Status(C.vip_query_hardware(
		C.vip_query_hardware_property_e(property),
		C.uint(size),
		value,
	)))
}

// =============================================================================
//  Buffer API
// =============================================================================

// CreateBuffer creates an input or output buffer.
func CreateBuffer(params *BufferCreateParams, sizeOfParam uint32) (Buffer, error) {
	var buf C.vip_buffer
	s := C.vip_create_buffer(
		(*C.vip_buffer_create_params_t)(unsafe.Pointer(params)),
		C.uint(sizeOfParam),
		&buf,
	)
	return Buffer(buf), checkStatus(Status(s))
}

// CreateBufferFromPhysical creates a buffer from physical addresses.
func CreateBufferFromPhysical(params *BufferCreateParams, physicalTable *Address, sizeTable []Uint32, physicalNum uint32) (Buffer, error) {
	var buf C.vip_buffer
	var psz *C.uint
	if len(sizeTable) > 0 {
		psz = (*C.uint)(unsafe.Pointer(&sizeTable[0]))
	}
	s := C.vip_create_buffer_from_physical(
		(*C.vip_buffer_create_params_t)(unsafe.Pointer(params)),
		(*C.vip_address_t)(unsafe.Pointer(physicalTable)),
		psz,
		C.uint(physicalNum),
		&buf,
	)
	return Buffer(buf), checkStatus(Status(s))
}

// CreateBufferFromHandle creates a buffer from a logical handle.
func CreateBufferFromHandle(params *BufferCreateParams, handle unsafe.Pointer, handleSize uint32) (Buffer, error) {
	var buf C.vip_buffer
	s := C.vip_create_buffer_from_handle(
		(*C.vip_buffer_create_params_t)(unsafe.Pointer(params)),
		C.vip_ptr(handle),
		C.uint(handleSize),
		&buf,
	)
	return Buffer(buf), checkStatus(Status(s))
}

// CreateBufferFromFD creates a buffer from a file descriptor (Linux DMA-BUF).
func CreateBufferFromFD(params *BufferCreateParams, fd uint32, memorySize uint32) (Buffer, error) {
	var buf C.vip_buffer
	s := C.vip_create_buffer_from_fd(
		(*C.vip_buffer_create_params_t)(unsafe.Pointer(params)),
		C.uint(fd),
		C.uint(memorySize),
		&buf,
	)
	return Buffer(buf), checkStatus(Status(s))
}

// DestroyBuffer destroys a buffer object.
func DestroyBuffer(buf Buffer) error {
	return checkStatus(Status(C.vip_destroy_buffer(C.vip_buffer(buf))))
}

// MapBuffer maps a buffer to get a CPU-accessible address.
func MapBuffer(buf Buffer) unsafe.Pointer {
	return C.vip_map_buffer(C.vip_buffer(buf))
}

// UnmapBuffer unmaps a previously mapped buffer.
func UnmapBuffer(buf Buffer) error {
	return checkStatus(Status(C.vip_unmap_buffer(C.vip_buffer(buf))))
}

// GetBufferSize returns the byte size allocated for the buffer.
func GetBufferSize(buf Buffer) uint32 {
	return uint32(C.vip_get_buffer_size(C.vip_buffer(buf)))
}

// FlushBuffer flushes or invalidates a buffer's CPU cache.
func FlushBuffer(buf Buffer, opType BufferOpType) error {
	return checkStatus(Status(C.vip_flush_buffer(
		C.vip_buffer(buf),
		C.vip_buffer_operation_type_e(opType),
	)))
}

// =============================================================================
//  Network API
// =============================================================================

// CreateNetwork creates a network object from binary data (file path or memory).
func CreateNetwork(data unsafe.Pointer, sizeOfData uint32, ntype CreateNetworkType) (Network, error) {
	var net C.vip_network
	s := C.vip_create_network(
		data,
		C.uint(sizeOfData),
		C.vip_create_network_type_e(ntype),
		&net,
	)
	return Network(net), checkStatus(Status(s))
}

// DupNetwork duplicates a network for weight sharing.
func DupNetwork(data unsafe.Pointer, sizeOfData uint32, ntype DupNetworkType, network Network) (Network, error) {
	var dup C.vip_network
	s := C.vip_dup_network(
		data,
		C.uint(sizeOfData),
		C.vip_dup_network_type_e(ntype),
		C.vip_network(network),
		&dup,
	)
	return Network(dup), checkStatus(Status(s))
}

// WeakDupNetwork creates a weak-dup of a network (shares coefficients).
func WeakDupNetwork(network Network) (Network, error) {
	var dup C.vip_network
	s := C.vip_weak_dup_network(C.vip_network(network), &dup)
	return Network(dup), checkStatus(Status(s))
}

// DestroyNetwork destroys a network object and releases all resources.
func DestroyNetwork(network Network) error {
	return checkStatus(Status(C.vip_destroy_network(C.vip_network(network))))
}

// SetNetwork configures a network property.
func SetNetwork(network Network, property NetworkProperty, value unsafe.Pointer) error {
	return checkStatus(Status(C.vip_set_network(
		C.vip_network(network),
		C.vip_enum(property),
		value,
	)))
}

// QueryNetwork queries a network property.
func QueryNetwork(network Network, property NetworkProperty, value unsafe.Pointer) error {
	return checkStatus(Status(C.vip_query_network(
		C.vip_network(network),
		C.vip_enum(property),
		value,
	)))
}

// PrepareNetwork prepares a network for execution on VIP hardware.
func PrepareNetwork(network Network) error {
	return checkStatus(Status(C.vip_prepare_network(C.vip_network(network))))
}

// QueryInput queries a property of a specific network input.
func QueryInput(network Network, index uint32, property BufferProperty, value unsafe.Pointer) error {
	return checkStatus(Status(C.vip_query_input(
		C.vip_network(network),
		C.uint(index),
		C.vip_enum(property),
		value,
	)))
}

// QueryOutput queries a property of a specific network output.
func QueryOutput(network Network, index uint32, property BufferProperty, value unsafe.Pointer) error {
	return checkStatus(Status(C.vip_query_output(
		C.vip_network(network),
		C.uint(index),
		C.vip_enum(property),
		value,
	)))
}

// SetInput attaches an input buffer to a network.
func SetInput(network Network, index uint32, input Buffer) error {
	return checkStatus(Status(C.vip_set_input(
		C.vip_network(network),
		C.uint(index),
		C.vip_buffer(input),
	)))
}

// SetOutput attaches an output buffer to a network.
func SetOutput(network Network, index uint32, output Buffer) error {
	return checkStatus(Status(C.vip_set_output(
		C.vip_network(network),
		C.uint(index),
		C.vip_buffer(output),
	)))
}

// RunNetwork executes the network on VIP hardware (blocking).
func RunNetwork(network Network) error {
	return checkStatus(Status(C.vip_run_network(C.vip_network(network))))
}

// FinishNetwork finishes using the network for inference.
func FinishNetwork(network Network) error {
	return checkStatus(Status(C.vip_finish_network(C.vip_network(network))))
}

// TriggerNetwork kicks off network execution without waiting for completion.
func TriggerNetwork(network Network) error {
	return checkStatus(Status(C.vip_trigger_network(C.vip_network(network))))
}

// TriggerGroup triggers tasks in a group without waiting.
func TriggerGroup(group Group, num uint32) error {
	return checkStatus(Status(C.vip_trigger_group(C.vip_group(group), C.uint(num))))
}

// WaitNetwork explicitly waits for HW to finish executing the network.
func WaitNetwork(network Network) error {
	return checkStatus(Status(C.vip_wait_network(C.vip_network(network))))
}

// WaitGroup waits for HW to finish executing tasks in a group.
func WaitGroup(group Group) error {
	return checkStatus(Status(C.vip_wait_group(C.vip_group(group))))
}

// CancelNetwork cancels a running network.
func CancelNetwork(network Network) error {
	return checkStatus(Status(C.vip_cancel_network(C.vip_network(network))))
}

// PowerManagement controls VIP core power/frequency.
func PowerManagement(deviceID uint32, property PowerProperty, value unsafe.Pointer) error {
	return checkStatus(Status(C.vip_power_management(
		C.uint(deviceID),
		C.vip_power_property_e(property),
		value,
	)))
}

// CreateGroup creates a group for running multiple tasks without interrupts.
func CreateGroup(count uint32) (Group, error) {
	var g C.vip_group
	s := C.vip_create_group(C.uint(count), &g)
	return Group(g), checkStatus(Status(s))
}

// DestroyGroup destroys a group object.
func DestroyGroup(group Group) error {
	return checkStatus(Status(C.vip_destroy_group(C.vip_group(group))))
}

// SetGroup configures a group property.
func SetGroup(group Group, property GroupProperty, value unsafe.Pointer) error {
	return checkStatus(Status(C.vip_set_group(
		C.vip_group(group),
		C.vip_enum(property),
		value,
	)))
}

// QueryGroup queries a group property.
func QueryGroup(group Group, property GroupProperty, value unsafe.Pointer) error {
	return checkStatus(Status(C.vip_query_group(
		C.vip_group(group),
		C.vip_enum(property),
		value,
	)))
}

// AddNetwork adds a network to a group.
func AddNetwork(group Group, network Network) error {
	return checkStatus(Status(C.vip_add_network(
		C.vip_group(group),
		C.vip_network(network),
	)))
}

// RunGroup executes tasks in a group (blocking).
func RunGroup(group Group, num uint32) error {
	return checkStatus(Status(C.vip_run_group(C.vip_group(group), C.uint(num))))
}

// SetPPUParam changes PPU engine parameters for a network.
func SetPPUParam(network Network, param *PPUParam, index uint32) error {
	return checkStatus(Status(C.vip_set_ppu_param(
		C.vip_network(network),
		(*C.vip_ppu_param_t)(unsafe.Pointer(param)),
		C.uint(index),
	)))
}
