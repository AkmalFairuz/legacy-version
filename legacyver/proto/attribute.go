package proto

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// Attribute is an entity attribute, that holds specific data such as the health of the entity. Each attribute
// holds a default value, maximum and minimum value, name and its current value.
type Attribute struct {
	protocol.AttributeValue
	// DefaultMin is the default minimum value of the attribute. It's not clear why this field must be sent to
	// the client, but it is required regardless.
	DefaultMin float32
	// DefaultMax is the default maximum value of the attribute. It's not clear why this field must be sent to
	// the client, but it is required regardless.
	DefaultMax float32
	// Default is the default value of the attribute. It's not clear why this field must be sent to the
	// client, but it is required regardless.
	Default float32
	// Modifiers is a slice of AttributeModifiers that are applied to the attribute.
	Modifiers []protocol.AttributeModifier
}

func (x *Attribute) FromLatest(v protocol.Attribute) Attribute {
	x.AttributeValue = v.AttributeValue
	x.DefaultMin = v.DefaultMin
	x.DefaultMax = v.DefaultMax
	x.Default = v.Default
	x.Modifiers = v.Modifiers
	return *x
}

func (x *Attribute) ToLatest() protocol.Attribute {
	return protocol.Attribute{
		AttributeValue: protocol.AttributeValue{
			Name:  x.Name,
			Value: x.Value,
			Max:   x.Max,
			Min:   x.Min,
		},
		DefaultMin: x.DefaultMin,
		DefaultMax: x.DefaultMax,
		Default:    x.Default,
		Modifiers:  x.Modifiers,
	}
}

// Marshal encodes/decodes an Attribute.
func (x *Attribute) Marshal(r protocol.IO) {
	r.Float32(&x.Min)
	r.Float32(&x.Max)
	r.Float32(&x.Value)
	if IsProtoGTE(r, ID729) {
		r.Float32(&x.DefaultMin)
		r.Float32(&x.DefaultMax)
	}
	r.Float32(&x.Default)
	r.String(&x.Name)
	protocol.Slice(r, &x.Modifiers)
}
