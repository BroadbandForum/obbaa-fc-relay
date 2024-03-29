# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: tr477_cpri_message.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='tr477_cpri_message.proto',
  package='tr477_cpri_message.v1',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\n\x18tr477_cpri_message.proto\x12\x15tr477_cpri_message.v1\x1a\x1bgoogle/protobuf/empty.proto\"\x1f\n\rCpriMsgHeader\x12\x0e\n\x06msg_id\x18\x01 \x01(\t\"Z\n\x0fGenericMetadata\x12\x13\n\x0b\x64\x65vice_name\x18\x01 \x01(\t\x12\x18\n\x10\x64\x65vice_interface\x18\x02 \x01(\t\x12\x18\n\x10originating_rule\x18\x03 \x01(\t\":\n\x0bOnuMetadata\x12\x1b\n\x13\x63hannel_termination\x18\x01 \x01(\t\x12\x0e\n\x06onu_id\x18\x02 \x01(\t\"\x0e\n\x0c\x44hcpMetadata\"\x0f\n\rPppoeMetadata\"\xf9\x01\n\x0c\x43priMetaData\x12\x37\n\x07generic\x18\x01 \x01(\x0b\x32&.tr477_cpri_message.v1.GenericMetadata\x12/\n\x03onu\x18\x02 \x01(\x0b\x32\".tr477_cpri_message.v1.OnuMetadata\x12\x33\n\x04\x64hcp\x18\x03 \x01(\x0b\x32#.tr477_cpri_message.v1.DhcpMetadataH\x00\x12\x35\n\x05pppoe\x18\x04 \x01(\x0b\x32$.tr477_cpri_message.v1.PppoeMetadataH\x00\x42\x13\n\x11specific_metadata\"\x87\x01\n\x07\x43priMsg\x12\x34\n\x06header\x18\x01 \x01(\x0b\x32$.tr477_cpri_message.v1.CpriMsgHeader\x12\x36\n\tmeta_data\x18\x02 \x01(\x0b\x32#.tr477_cpri_message.v1.CpriMetaData\x12\x0e\n\x06packet\x18\x03 \x01(\x0c\x62\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_empty__pb2.DESCRIPTOR,])




_CPRIMSGHEADER = _descriptor.Descriptor(
  name='CpriMsgHeader',
  full_name='tr477_cpri_message.v1.CpriMsgHeader',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='msg_id', full_name='tr477_cpri_message.v1.CpriMsgHeader.msg_id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=80,
  serialized_end=111,
)


_GENERICMETADATA = _descriptor.Descriptor(
  name='GenericMetadata',
  full_name='tr477_cpri_message.v1.GenericMetadata',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='device_name', full_name='tr477_cpri_message.v1.GenericMetadata.device_name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='device_interface', full_name='tr477_cpri_message.v1.GenericMetadata.device_interface', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='originating_rule', full_name='tr477_cpri_message.v1.GenericMetadata.originating_rule', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=113,
  serialized_end=203,
)


_ONUMETADATA = _descriptor.Descriptor(
  name='OnuMetadata',
  full_name='tr477_cpri_message.v1.OnuMetadata',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='channel_termination', full_name='tr477_cpri_message.v1.OnuMetadata.channel_termination', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='onu_id', full_name='tr477_cpri_message.v1.OnuMetadata.onu_id', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=205,
  serialized_end=263,
)


_DHCPMETADATA = _descriptor.Descriptor(
  name='DhcpMetadata',
  full_name='tr477_cpri_message.v1.DhcpMetadata',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=265,
  serialized_end=279,
)


_PPPOEMETADATA = _descriptor.Descriptor(
  name='PppoeMetadata',
  full_name='tr477_cpri_message.v1.PppoeMetadata',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=281,
  serialized_end=296,
)


_CPRIMETADATA = _descriptor.Descriptor(
  name='CpriMetaData',
  full_name='tr477_cpri_message.v1.CpriMetaData',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='generic', full_name='tr477_cpri_message.v1.CpriMetaData.generic', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='onu', full_name='tr477_cpri_message.v1.CpriMetaData.onu', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='dhcp', full_name='tr477_cpri_message.v1.CpriMetaData.dhcp', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='pppoe', full_name='tr477_cpri_message.v1.CpriMetaData.pppoe', index=3,
      number=4, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
    _descriptor.OneofDescriptor(
      name='specific_metadata', full_name='tr477_cpri_message.v1.CpriMetaData.specific_metadata',
      index=0, containing_type=None, fields=[]),
  ],
  serialized_start=299,
  serialized_end=548,
)


_CPRIMSG = _descriptor.Descriptor(
  name='CpriMsg',
  full_name='tr477_cpri_message.v1.CpriMsg',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='header', full_name='tr477_cpri_message.v1.CpriMsg.header', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='meta_data', full_name='tr477_cpri_message.v1.CpriMsg.meta_data', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='packet', full_name='tr477_cpri_message.v1.CpriMsg.packet', index=2,
      number=3, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=_b(""),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=551,
  serialized_end=686,
)

_CPRIMETADATA.fields_by_name['generic'].message_type = _GENERICMETADATA
_CPRIMETADATA.fields_by_name['onu'].message_type = _ONUMETADATA
_CPRIMETADATA.fields_by_name['dhcp'].message_type = _DHCPMETADATA
_CPRIMETADATA.fields_by_name['pppoe'].message_type = _PPPOEMETADATA
_CPRIMETADATA.oneofs_by_name['specific_metadata'].fields.append(
  _CPRIMETADATA.fields_by_name['dhcp'])
_CPRIMETADATA.fields_by_name['dhcp'].containing_oneof = _CPRIMETADATA.oneofs_by_name['specific_metadata']
_CPRIMETADATA.oneofs_by_name['specific_metadata'].fields.append(
  _CPRIMETADATA.fields_by_name['pppoe'])
_CPRIMETADATA.fields_by_name['pppoe'].containing_oneof = _CPRIMETADATA.oneofs_by_name['specific_metadata']
_CPRIMSG.fields_by_name['header'].message_type = _CPRIMSGHEADER
_CPRIMSG.fields_by_name['meta_data'].message_type = _CPRIMETADATA
DESCRIPTOR.message_types_by_name['CpriMsgHeader'] = _CPRIMSGHEADER
DESCRIPTOR.message_types_by_name['GenericMetadata'] = _GENERICMETADATA
DESCRIPTOR.message_types_by_name['OnuMetadata'] = _ONUMETADATA
DESCRIPTOR.message_types_by_name['DhcpMetadata'] = _DHCPMETADATA
DESCRIPTOR.message_types_by_name['PppoeMetadata'] = _PPPOEMETADATA
DESCRIPTOR.message_types_by_name['CpriMetaData'] = _CPRIMETADATA
DESCRIPTOR.message_types_by_name['CpriMsg'] = _CPRIMSG
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

CpriMsgHeader = _reflection.GeneratedProtocolMessageType('CpriMsgHeader', (_message.Message,), {
  'DESCRIPTOR' : _CPRIMSGHEADER,
  '__module__' : 'tr477_cpri_message_pb2'
  # @@protoc_insertion_point(class_scope:tr477_cpri_message.v1.CpriMsgHeader)
  })
_sym_db.RegisterMessage(CpriMsgHeader)

GenericMetadata = _reflection.GeneratedProtocolMessageType('GenericMetadata', (_message.Message,), {
  'DESCRIPTOR' : _GENERICMETADATA,
  '__module__' : 'tr477_cpri_message_pb2'
  # @@protoc_insertion_point(class_scope:tr477_cpri_message.v1.GenericMetadata)
  })
_sym_db.RegisterMessage(GenericMetadata)

OnuMetadata = _reflection.GeneratedProtocolMessageType('OnuMetadata', (_message.Message,), {
  'DESCRIPTOR' : _ONUMETADATA,
  '__module__' : 'tr477_cpri_message_pb2'
  # @@protoc_insertion_point(class_scope:tr477_cpri_message.v1.OnuMetadata)
  })
_sym_db.RegisterMessage(OnuMetadata)

DhcpMetadata = _reflection.GeneratedProtocolMessageType('DhcpMetadata', (_message.Message,), {
  'DESCRIPTOR' : _DHCPMETADATA,
  '__module__' : 'tr477_cpri_message_pb2'
  # @@protoc_insertion_point(class_scope:tr477_cpri_message.v1.DhcpMetadata)
  })
_sym_db.RegisterMessage(DhcpMetadata)

PppoeMetadata = _reflection.GeneratedProtocolMessageType('PppoeMetadata', (_message.Message,), {
  'DESCRIPTOR' : _PPPOEMETADATA,
  '__module__' : 'tr477_cpri_message_pb2'
  # @@protoc_insertion_point(class_scope:tr477_cpri_message.v1.PppoeMetadata)
  })
_sym_db.RegisterMessage(PppoeMetadata)

CpriMetaData = _reflection.GeneratedProtocolMessageType('CpriMetaData', (_message.Message,), {
  'DESCRIPTOR' : _CPRIMETADATA,
  '__module__' : 'tr477_cpri_message_pb2'
  # @@protoc_insertion_point(class_scope:tr477_cpri_message.v1.CpriMetaData)
  })
_sym_db.RegisterMessage(CpriMetaData)

CpriMsg = _reflection.GeneratedProtocolMessageType('CpriMsg', (_message.Message,), {
  'DESCRIPTOR' : _CPRIMSG,
  '__module__' : 'tr477_cpri_message_pb2'
  # @@protoc_insertion_point(class_scope:tr477_cpri_message.v1.CpriMsg)
  })
_sym_db.RegisterMessage(CpriMsg)


# @@protoc_insertion_point(module_scope)
