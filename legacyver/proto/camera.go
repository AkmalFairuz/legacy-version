package proto

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// CameraPreset represents a basic preset that can be extended upon by more complex instructions.
type CameraPreset struct {
	// Name is the name of the preset. Each preset must have their own unique name.
	Name string
	// Parent is the name of the preset that this preset extends upon. This can be left empty.
	Parent string
	// PosX is the default X position of the camera.
	PosX protocol.Optional[float32]
	// PosY is the default Y position of the camera.
	PosY protocol.Optional[float32]
	// PosZ is the default Z position of the camera.
	PosZ protocol.Optional[float32]
	// RotX is the default pitch of the camera.
	RotX protocol.Optional[float32]
	// RotY is the default yaw of the camera.
	RotY protocol.Optional[float32]
	// RotationSpeed is the speed at which the camera should rotate.
	RotationSpeed protocol.Optional[float32]
	// SnapToTarget determines whether the camera should snap to the target entity or not.
	SnapToTarget protocol.Optional[bool]
	// HorizontalRotationLimit is the horizontal rotation limit of the camera.
	HorizontalRotationLimit protocol.Optional[mgl32.Vec2]
	// VerticalRotationLimit is the vertical rotation limit of the camera.
	VerticalRotationLimit protocol.Optional[mgl32.Vec2]
	// ContinueTargeting determines whether the camera should continue targeting when using aim assist.
	ContinueTargeting protocol.Optional[bool]
	// TrackingRadius is the radius around the camera that the aim assist should track targets.
	TrackingRadius protocol.Optional[float32]
	// ViewOffset is only used in a follow_orbit camera and controls an offset based on a pivot point to the
	// player, causing it to be shifted in a certain direction.
	ViewOffset protocol.Optional[mgl32.Vec2]
	// EntityOffset controls the offset from the entity that the camera should be rendered at.
	EntityOffset protocol.Optional[mgl32.Vec3]
	// Radius is only used in a follow_orbit camera and controls how far away from the player the camera should
	// be rendered.
	Radius protocol.Optional[float32]
	// AudioListener defines where the audio should be played from when using this preset. This is one of the
	// constants above.
	AudioListener protocol.Optional[byte]
	// PlayerEffects is currently unknown.
	PlayerEffects protocol.Optional[bool]
	// AlignTargetAndCameraForward determines whether the camera should align the target and the camera forward
	// or not.
	AlignTargetAndCameraForward protocol.Optional[bool]
	// AimAssist defines the aim assist to use when using this preset.
	AimAssist protocol.Optional[protocol.CameraPresetAimAssist]
}

func (x *CameraPreset) FromLatest(cp protocol.CameraPreset) CameraPreset {
	x.Name = cp.Name
	x.Parent = cp.Parent
	x.PosX = cp.PosX
	x.PosY = cp.PosY
	x.PosZ = cp.PosZ
	x.RotX = cp.RotX
	x.RotY = cp.RotY
	x.RotationSpeed = cp.RotationSpeed
	x.SnapToTarget = cp.SnapToTarget
	x.HorizontalRotationLimit = cp.HorizontalRotationLimit
	x.VerticalRotationLimit = cp.VerticalRotationLimit
	x.ContinueTargeting = cp.ContinueTargeting
	x.TrackingRadius = cp.TrackingRadius
	x.ViewOffset = cp.ViewOffset
	x.EntityOffset = cp.EntityOffset
	x.Radius = cp.Radius
	x.AudioListener = cp.AudioListener
	x.PlayerEffects = cp.PlayerEffects
	x.AlignTargetAndCameraForward = cp.AlignTargetAndCameraForward
	x.AimAssist = cp.AimAssist
	return *x
}

func (x *CameraPreset) ToLatest() protocol.CameraPreset {
	return protocol.CameraPreset{
		Name:                        x.Name,
		Parent:                      x.Parent,
		PosX:                        x.PosX,
		PosY:                        x.PosY,
		PosZ:                        x.PosZ,
		RotX:                        x.RotX,
		RotY:                        x.RotY,
		RotationSpeed:               x.RotationSpeed,
		SnapToTarget:                x.SnapToTarget,
		HorizontalRotationLimit:     x.HorizontalRotationLimit,
		VerticalRotationLimit:       x.VerticalRotationLimit,
		ContinueTargeting:           x.ContinueTargeting,
		TrackingRadius:              x.TrackingRadius,
		ViewOffset:                  x.ViewOffset,
		EntityOffset:                x.EntityOffset,
		Radius:                      x.Radius,
		AudioListener:               x.AudioListener,
		PlayerEffects:               x.PlayerEffects,
		AlignTargetAndCameraForward: x.AlignTargetAndCameraForward,
		AimAssist:                   x.AimAssist,
	}
}

// Marshal encodes/decodes a CameraPreset.
func (x *CameraPreset) Marshal(r protocol.IO) {
	r.String(&x.Name)
	r.String(&x.Parent)
	protocol.OptionalFunc(r, &x.PosX, r.Float32)
	protocol.OptionalFunc(r, &x.PosY, r.Float32)
	protocol.OptionalFunc(r, &x.PosZ, r.Float32)
	protocol.OptionalFunc(r, &x.RotX, r.Float32)
	protocol.OptionalFunc(r, &x.RotY, r.Float32)
	if IsProtoGTE(r, ID729) {
		protocol.OptionalFunc(r, &x.RotationSpeed, r.Float32)
		protocol.OptionalFunc(r, &x.SnapToTarget, r.Bool)
	}
	if IsProtoGTE(r, ID748) {
		protocol.OptionalFunc(r, &x.HorizontalRotationLimit, r.Vec2)
		protocol.OptionalFunc(r, &x.VerticalRotationLimit, r.Vec2)
		protocol.OptionalFunc(r, &x.ContinueTargeting, r.Bool)
	}
	if IsProtoGTE(r, ID766) {
		protocol.OptionalFunc(r, &x.TrackingRadius, r.Float32)
	}
	protocol.OptionalFunc(r, &x.ViewOffset, r.Vec2)
	if IsProtoGTE(r, ID729) {
		protocol.OptionalFunc(r, &x.EntityOffset, r.Vec3)
	}
	protocol.OptionalFunc(r, &x.Radius, r.Float32)
	protocol.OptionalFunc(r, &x.AudioListener, r.Uint8)
	protocol.OptionalFunc(r, &x.PlayerEffects, r.Bool)
	if IsProtoGTE(r, ID748) {
		protocol.OptionalFunc(r, &x.AlignTargetAndCameraForward, r.Bool)
	}
	if IsProtoGTE(r, ID766) {
		protocol.OptionalMarshaler(r, &x.AimAssist)
	}
}
