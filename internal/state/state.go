package main

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type State struct {
	messageType protoreflect.MessageType
	value       protoreflect.Message
}

func (s *State) Init(mt protoreflect.MessageType) {
	if s.messageType != nil {
		fmt.Printf("state already initialized with message type")
	}

	s.messageType = mt
}

func (s State) MessageDescriptor() protoreflect.MessageDescriptor {
	return s.messageType.Descriptor()
}

func (s State) Value() protoreflect.Message {
	return s.value
}

func (s State) ValueAtPath(path string) (any, error) {
	pathArr := strings.Split(path, ".")

	var curPath string
	var travelledPath []string
	var fd protoreflect.FieldDescriptor
	var cur any
	var next any

	cur = s.value
	for len(pathArr) > 0 {
		curPath = pathArr[0]
		travelledPath = append(travelledPath, pathArr[0])
		pathArr = pathArr[1:]

		switch cur := cur.(type) {
		case protoreflect.Message:
			fds := cur.Descriptor().Fields()

			if !cur.Has(fd) {
				return nil, fmt.Errorf("path %s not present in state", strings.Join(travelledPath, "."))
			}

			fd := fds.ByJSONName(curPath)
			if fd.Kind() == protoreflect.BoolKind {
				next = cur.Get(fd).Bool()
			}
			if fd.Kind() == protoreflect.EnumKind {
				next = cur.Get(fd).Enum()
			}
			if fd.Kind() == protoreflect.Int32Kind ||
				fd.Kind() == protoreflect.Int64Kind ||
				fd.Kind() == protoreflect.Sint32Kind ||
				fd.Kind() == protoreflect.Sint64Kind ||
				fd.Kind() == protoreflect.Sfixed32Kind ||
				fd.Kind() == protoreflect.Sfixed64Kind {
				next = cur.Get(fd).Int()
			}
			if fd.Kind() == protoreflect.Uint32Kind ||
				fd.Kind() == protoreflect.Uint64Kind ||
				fd.Kind() == protoreflect.Fixed32Kind ||
				fd.Kind() == protoreflect.Fixed64Kind {
				next = cur.Get(fd).Uint()
			}
			if fd.Kind() == protoreflect.FloatKind ||
				fd.Kind() == protoreflect.DoubleKind {
				next = cur.Get(fd).Float()
			}
			if fd.Kind() == protoreflect.StringKind {
				next = cur.Get(fd).String()
			}
			if fd.Kind() == protoreflect.BytesKind {
				next = cur.Get(fd).Bytes()
			}
			if fd.Kind() == protoreflect.MessageKind ||
				fd.Kind() == protoreflect.GroupKind {
				next = cur.Get(fd).Message()
			}

			next = cur.Get(fd)
		}

		cur = next
	}

	return nil, nil
}

func (s *State) Update(data []byte) (protoreflect.Message, error) {
	if s.messageType == nil {
		return nil, errors.New("state not initialized")
	}

	m := s.messageType.New().Interface()
	err := protojson.Unmarshal(data, m)
	if err != nil || !m.ProtoReflect().IsValid() {
		return nil, errors.New("data does not match message format")
	}

	return m.ProtoReflect(), nil
}