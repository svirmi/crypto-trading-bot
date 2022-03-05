package mongodb

import (
	"fmt"
	"reflect"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type decimal_codec struct{}

func (de decimal_codec) EncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	var decimal_type = reflect.TypeOf(decimal.Zero)
	if !val.IsValid() || val.Type() != decimal_type {
		return bsoncodec.ValueEncoderError{
			Name:     "DecimalEncodeValue",
			Types:    []reflect.Type{decimal_type},
			Received: val}
	}
	return vw.WriteString(val.Interface().(decimal.Decimal).String())
}

func (de decimal_codec) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Kind() != reflect.TypeOf(decimal.Zero).Kind() {
		return bsoncodec.ValueDecoderError{
			Name:     "DecimalDecodeValue",
			Kinds:    []reflect.Kind{reflect.TypeOf(decimal.Zero).Kind()},
			Received: val}
	}

	if vr.Type() != bsontype.String {
		return fmt.Errorf("cannot decode %v into a string type", vr.Type())
	}

	str, err := vr.ReadString()
	if err != nil {
		return err
	}
	decimal, err := decimal.NewFromString(str)
	if err != nil {
		return err
	}
	val.Set(reflect.ValueOf(decimal))
	return nil
}

func build_custom_registry() *bsoncodec.Registry {
	var primitiveCodecs bson.PrimitiveCodecs
	rb := bsoncodec.NewRegistryBuilder()
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	primitiveCodecs.RegisterPrimitiveCodecs(rb)

	rb.RegisterTypeEncoder(reflect.TypeOf(decimal.Zero), decimal_codec{})
	rb.RegisterTypeDecoder(reflect.TypeOf(decimal.Zero), decimal_codec{})
	return rb.Build()
}
