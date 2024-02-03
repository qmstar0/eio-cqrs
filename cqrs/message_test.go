package cqrs_test

import (
	"context"
	"github.com/qmstar0/eio-cqrs/cqrs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage(t *testing.T) {
	codec := cqrs.NewMessageCodec(cqrs.JsonMarshaler{})

	message1, err1 := codec.EncodeMessage(context.TODO(), &Cmd{Name: "1"})
	message2, err2 := codec.EncodeMessage(context.TODO(), &Cmd{Name: "1"})
	message3, err3 := codec.EncodeMessage(context.TODO(), &Cmd{Name: "2"})
	assert.Equal(t, message1.Payload, message2.Payload)
	assert.NotEqual(t, message3.Payload, message1.Payload)
	assert.NoError(t, err1, err2, err3)

	topic := codec.Topic(&Cmd{})
	assert.Equal(t, "Cmd", topic)

	c := &Cmd{}

	err := codec.DecodeMessage(message1, c)
	assert.NoError(t, err)
	assert.Equal(t, c.Name, "1")

	dataMap := new(map[string]string)
	err = codec.DecodeMessage(message1, dataMap)
	assert.NoError(t, err)

	err = codec.DecodeMessage(message1, 1)
	assert.Error(t, err)
}
