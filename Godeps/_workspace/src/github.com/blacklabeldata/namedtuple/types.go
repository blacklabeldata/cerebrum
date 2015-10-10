package namedtuple

import "hash/fnv"

var syncHash SynchronizedHash = NewHasher(fnv.New32a())

// Empty Tuple
var NIL Tuple = Tuple{}

type TupleType struct {
	Namespace     string // Tuple Namespace
	Name          string // Tuple Name
	NamespaceHash uint32
	Hash          uint32
	versions      [][]Field
	fields        map[string]int
}

type Version struct {
	Num    uint8
	Fields []Field
}

type Field struct {
	Name     string
	Required bool
	Type     FieldType
}

// New creates a new TupleType with the given namespace and type name
func New(namespace string, name string) (t TupleType) {
	hash := syncHash.Hash([]byte(name))
	ns_hash := syncHash.Hash([]byte(namespace))
	t = TupleType{namespace, name, ns_hash, hash, make([][]Field, 0), make(map[string]int)}
	return
}

// AddVersion adds a version to the tuple type
func (t *TupleType) AddVersion(fields ...Field) {
	t.versions = append(t.versions, fields)
	for _, field := range fields {
		t.fields[field.Name] = len(t.fields)
	}
}

// Contains determines is a field exists in the TupleType
func (t *TupleType) Contains(field string) bool {
	_, exists := t.fields[field]
	return exists
}

// Offset determines the numerical offset for the given field
func (t *TupleType) Offset(field string) (offset int, exists bool) {
	offset, exists = t.fields[field]
	return
}

// NumVersions returns the number of version in the tuple type
func (t *TupleType) NumVersions() int {
	return len(t.versions)
}

// Versions returns an array of versions contained in this type
func (t *TupleType) Versions() (vers []Version) {
	vers = make([]Version, t.NumVersions())
	for i := 0; i < t.NumVersions(); i++ {
		vers[i] = Version{uint8(i + 1), t.versions[i]}
	}
	return
}

// Builder creates a new builder from the TupleType
func (t *TupleType) Builder(buffer []byte) TupleBuilder {
	return NewBuilder(*t, buffer)
}

// FieldType is the byte representation of each field type
type FieldType uint8

const (
	Uint8Field FieldType = iota
	Uint8ArrayField
	Int8Field
	Int8ArrayField
	Uint16Field
	Uint16ArrayField
	Int16Field
	Int16ArrayField
	Uint32Field
	Uint32ArrayField
	Int32Field
	Int32ArrayField
	Uint64Field
	Uint64ArrayField
	Int64Field
	Int64ArrayField
	Float32Field
	Float32ArrayField
	Float64Field
	Float64ArrayField
	TimestampField
	TimestampArrayField
	TupleField
	TupleArrayField
	StringField
	StringArrayField
	BooleanField
	BooleanArrayField
)

// TypeCode represents a field type
type TypeCode struct {
	OpCode uint8
	Size   uint8
}

// Field Type constants
var (
	NilCode                  = TypeCode{0, 0}
	TrueCode                 = TypeCode{1, 0}
	FalseCode                = TypeCode{2, 0}
	BooleanArray8Code        = TypeCode{3, 1}
	BooleanArray16Code       = TypeCode{4, 2}
	BooleanArray32Code       = TypeCode{5, 4}
	BooleanArray64Code       = TypeCode{6, 8}
	ByteCode                 = TypeCode{7, 0}
	ByteArray8Code           = TypeCode{8, 1}
	ByteArray16Code          = TypeCode{9, 2}
	ByteArray32Code          = TypeCode{10, 4}
	ByteArray64Code          = TypeCode{11, 8}
	UnsignedByteCode         = TypeCode{12, 0}
	UnsignedByteArray8Code   = TypeCode{13, 1}
	UnsignedByteArray16Code  = TypeCode{14, 2}
	UnsignedByteArray32Code  = TypeCode{15, 4}
	UnsignedByteArray64Code  = TypeCode{16, 8}
	Short8Code               = TypeCode{17, 0}
	Short16Code              = TypeCode{18, 0}
	ShortArray8Code          = TypeCode{19, 1}
	ShortArray16Code         = TypeCode{20, 2}
	ShortArray32Code         = TypeCode{21, 4}
	ShortArray64Code         = TypeCode{22, 8}
	UnsignedShort8Code       = TypeCode{23, 0}
	UnsignedShort16Code      = TypeCode{24, 0}
	UnsignedShortArray8Code  = TypeCode{25, 1}
	UnsignedShortArray16Code = TypeCode{26, 2}
	UnsignedShortArray32Code = TypeCode{27, 4}
	UnsignedShortArray64Code = TypeCode{28, 8}
	Int8Code                 = TypeCode{29, 0}
	Int16Code                = TypeCode{30, 0}
	Int32Code                = TypeCode{31, 0}
	IntArray8Code            = TypeCode{32, 1}
	IntArray16Code           = TypeCode{33, 2}
	IntArray32Code           = TypeCode{34, 4}
	IntArray64Code           = TypeCode{35, 8}
	UnsignedInt8Code         = TypeCode{36, 0}
	UnsignedInt16Code        = TypeCode{37, 0}
	UnsignedInt32Code        = TypeCode{38, 0}
	UnsignedIntArray8Code    = TypeCode{39, 1}
	UnsignedIntArray16Code   = TypeCode{40, 2}
	UnsignedIntArray32Code   = TypeCode{41, 4}
	UnsignedIntArray64Code   = TypeCode{42, 8}
	Long8Code                = TypeCode{43, 0}
	Long16Code               = TypeCode{44, 0}
	Long32Code               = TypeCode{45, 0}
	Long64Code               = TypeCode{46, 0}
	LongArray8Code           = TypeCode{47, 1}
	LongArray16Code          = TypeCode{48, 2}
	LongArray32Code          = TypeCode{49, 4}
	LongArray64Code          = TypeCode{50, 8}
	UnsignedLong8Code        = TypeCode{51, 0}
	UnsignedLong16Code       = TypeCode{52, 0}
	UnsignedLong32Code       = TypeCode{53, 0}
	UnsignedLong64Code       = TypeCode{54, 0}
	UnsignedLongArray8Code   = TypeCode{55, 1}
	UnsignedLongArray16Code  = TypeCode{56, 2}
	UnsignedLongArray32Code  = TypeCode{57, 4}
	UnsignedLongArray64Code  = TypeCode{58, 8}
	DoubleCode               = TypeCode{59, 0}
	DoubleArray8Code         = TypeCode{60, 1}
	DoubleArray16Code        = TypeCode{61, 2}
	DoubleArray32Code        = TypeCode{62, 4}
	DoubleArray64Code        = TypeCode{63, 8}
	FloatCode                = TypeCode{64, 0}
	FloatArray8Code          = TypeCode{65, 1}
	FloatArray16Code         = TypeCode{66, 2}
	FloatArray32Code         = TypeCode{67, 4}
	FloatArray64Code         = TypeCode{68, 8}
	String8Code              = TypeCode{69, 1}
	String16Code             = TypeCode{70, 2}
	String32Code             = TypeCode{71, 4}
	String64Code             = TypeCode{72, 8}
	StringArray8Code         = TypeCode{73, 1}
	StringArray16Code        = TypeCode{74, 2}
	StringArray32Code        = TypeCode{75, 4}
	StringArray64Code        = TypeCode{76, 8}
	TimestampCode            = TypeCode{77, 0}
	TimestampArray8Code      = TypeCode{78, 1}
	TimestampArray16Code     = TypeCode{79, 2}
	TimestampArray32Code     = TypeCode{80, 4}
	TimestampArray64Code     = TypeCode{81, 8}
	Tuple8Code               = TypeCode{82, 1}
	Tuple16Code              = TypeCode{83, 2}
	Tuple32Code              = TypeCode{84, 4}
	Tuple64Code              = TypeCode{85, 8}
	TupleArray8Code          = TypeCode{86, 1}
	TupleArray16Code         = TypeCode{87, 2}
	TupleArray32Code         = TypeCode{88, 4}
	TupleArray64Code         = TypeCode{89, 8}
	// FieldName            = TypeCode{90, 1}
	// VarInt               = TypeCode{91, 0}
	// UVarInt              = TypeCode{92, 0}
)

var fieldTypes = map[byte]TypeCode{
	0:  NilCode,
	1:  TrueCode,
	2:  FalseCode,
	3:  BooleanArray8Code,
	4:  BooleanArray16Code,
	5:  BooleanArray32Code,
	6:  BooleanArray64Code,
	7:  ByteCode,
	8:  ByteArray8Code,
	9:  ByteArray16Code,
	10: ByteArray32Code,
	11: ByteArray64Code,
	12: UnsignedByteCode,
	13: UnsignedByteArray8Code,
	14: UnsignedByteArray16Code,
	15: UnsignedByteArray32Code,
	16: UnsignedByteArray64Code,
	17: Short8Code,
	18: Short16Code,
	19: ShortArray8Code,
	20: ShortArray16Code,
	21: ShortArray32Code,
	22: ShortArray64Code,
	23: UnsignedShort8Code,
	24: UnsignedShort16Code,
	25: UnsignedShortArray8Code,
	26: UnsignedShortArray16Code,
	27: UnsignedShortArray32Code,
	28: UnsignedShortArray64Code,
	29: Int8Code,
	30: Int16Code,
	31: Int32Code,
	32: IntArray8Code,
	33: IntArray16Code,
	34: IntArray32Code,
	35: IntArray64Code,
	36: UnsignedInt8Code,
	37: UnsignedInt16Code,
	38: UnsignedInt32Code,
	39: UnsignedIntArray8Code,
	40: UnsignedIntArray16Code,
	41: UnsignedIntArray32Code,
	42: UnsignedIntArray64Code,
	43: Long8Code,
	44: Long16Code,
	45: Long32Code,
	46: Long64Code,
	47: LongArray8Code,
	48: LongArray16Code,
	49: LongArray32Code,
	50: LongArray64Code,
	51: UnsignedLong8Code,
	52: UnsignedLong16Code,
	53: UnsignedLong32Code,
	54: UnsignedLong64Code,
	55: UnsignedLongArray8Code,
	56: UnsignedLongArray16Code,
	57: UnsignedLongArray32Code,
	58: UnsignedLongArray64Code,
	59: DoubleCode,
	60: DoubleArray8Code,
	61: DoubleArray16Code,
	62: DoubleArray32Code,
	63: DoubleArray64Code,
	64: FloatCode,
	65: FloatArray8Code,
	66: FloatArray16Code,
	67: FloatArray32Code,
	68: FloatArray64Code,
	69: String8Code,
	70: String16Code,
	71: String32Code,
	72: String64Code,
	73: StringArray8Code,
	74: StringArray16Code,
	75: StringArray32Code,
	76: StringArray64Code,
	77: TimestampCode,
	78: TimestampArray8Code,
	79: TimestampArray16Code,
	80: TimestampArray32Code,
	81: TimestampArray64Code,
	82: Tuple8Code,
	83: Tuple16Code,
	84: Tuple32Code,
	85: Tuple64Code,
	86: TupleArray8Code,
	87: TupleArray16Code,
	88: TupleArray32Code,
	89: TupleArray64Code,
}
