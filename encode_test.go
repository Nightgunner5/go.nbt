package nbt

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

type Player struct {
	FoodSaturationLevel float32 `nbt:"foodSaturationLevel"`
	FoodExhaustionLevel float32 `nbt:"foodExhaustionLevel"`
	FoodTickTimer       uint32  `nbt:"foodTickTimer"`
	FoodLevel           uint32  `nbt:"foodLevel"`

	XpLevel uint32
	XpP     float32
	XpTotal uint32

	Health uint16
	Fire   uint16
	Air    uint16

	AttackTime uint16
	DeathTime  uint16
	HurtTime   uint16

	FallDistance float32

	Sleeping   bool
	SleepTimer uint16

	SpawnX int32
	SpawnY int32
	SpawnZ int32

	OnGround  bool
	Pos       []float64
	Motion    []float64
	Rotation  []float32
	Dimension int32

	GameType  uint32          `nbt:"playerGameType"`
	Abilities PlayerAbilities `nbt:"abilities"`

	Inventory  []InventoryItem
	EnderItems []InventoryItem
}

type PlayerAbilities struct {
	MayFly bool `nbt:"mayfly"`
	Flying bool `nbt:"flying"`

	FlySpeed  float32 `nbt:"flySpeed"`
	WalkSpeed float32 `nbt:"walkSpeed"`

	InstaBuild   bool `nbt:"instabuild"`
	Invulnerable bool `nbt:"invulnerable"`
	MayBuild     bool `nbt:"mayBuild"`
}

type InventoryItem struct {
	Type   uint16 `nbt:"id"`
	Damage uint16
	Count  uint8
	Slot   uint8
}

func TestEncode(t *testing.T) {
	data, err := ioutil.ReadFile("testcases/Nightgunner5.dat")
	if err != nil {
		t.Error(err)
	}

	var reference Player
	err = Unmarshal(GZip, bytes.NewReader(data), &reference)
	if err != nil {
		t.Error(err)
	}

	var encoded bytes.Buffer
	err = Marshal(GZip, &encoded, reference)
	if err != nil {
		t.Error(err)
	}

	var result Player
	err = Unmarshal(GZip, bytes.NewReader(encoded.Bytes()), &result)
	if err != nil {
		t.Error(err)
	}

	left, right := parseStruct(reflect.ValueOf(result)), parseStruct(reflect.ValueOf(reference))
	for field, l := range left {
		r := right[field]

		if !reflect.DeepEqual(l.Interface(), r.Interface()) {
			t.Errorf("Field %s differs:", field)
			t.Logf("Found   : %#v", l)
			t.Logf("Expected: %#v", r)
		}
	}
}
