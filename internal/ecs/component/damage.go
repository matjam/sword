package component

import "github.com/matjam/sword/internal/ecs"

type DamageRecord struct {
	Amount int
	Source string
}

// Damage records incoming damage and is applied by the injury system.
type Damage struct {
	Records []DamageRecord
}

func (*Damage) ComponentName() ecs.ComponentName {
	return "damage"
}

// RecordDamage records damage to the entity.
func (d *Damage) RecordDamage(amount int, source string) {
	if d.Records == nil {
		d.Records = make([]DamageRecord, 0)
	}

	d.Records = append(d.Records, DamageRecord{amount, source})
}

// ClearDamage clears the damage records.
func (d *Damage) ClearDamage() {
	d.Records = make([]DamageRecord, 0)
}
