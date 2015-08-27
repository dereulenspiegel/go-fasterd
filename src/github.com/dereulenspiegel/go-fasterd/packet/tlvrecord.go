package packet

type TLVRecordType uint16

const (
  HandshakeType TLVRecordType = iota
  ReplyCode
  ErrorDetail
  Flags
  Mode
  ProtocolName
  SenderKey
  RecipientKey
  SenderHandshakeKey
  RecipientHandshakeKey
  AuthenticationTag //obsolete
  MTU
  MethodName
  VersionName
  MethodList
  TLVAuthenticationTag
)

type TLVRecord struct {
  Type TLVRecordType
  Length uint16
  Body []byte
}
