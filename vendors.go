package main

// should be int but idc
type ItemType string

const (
	ItemTypeMod      = "Mod"
	ItemTypeArchPart = "ArchPart"
	ItemTypeWeapon   = "Weapon"
)

type Vendor struct {
	Name  string
	Items []Item
}

type Item struct {
	Name         string
	Type         ItemType
	StandingCost int
}

var AribiterOfHexis = Vendor{
	"Arbiter of Hexis",
	[]Item{
		{"gilded_truth", ItemTypeMod, 25000},
		{"blade_of_truth", ItemTypeMod, 25000},
		{"avenging_truth", ItemTypeMod, 25000},
		{"stinging_truth", ItemTypeMod, 25000},
		{"seeking_shuriken", ItemTypeMod, 25000},
		{"smoke_shadow", ItemTypeMod, 25000},
		{"fatal_teleport", ItemTypeMod, 25000},
		{"rising_storm", ItemTypeMod, 25000},
		{"endless_lullaby", ItemTypeMod, 25000},
		{"reactive_storm", ItemTypeMod, 25000},
		{"duality", ItemTypeMod, 25000},
		{"calm_and_frenzy", ItemTypeMod, 25000},
		{"peaceful_provocation", ItemTypeMod, 25000},
		{"energy_transfer", ItemTypeMod, 25000},
		{"surging_dash", ItemTypeMod, 25000},
		{"radiant_finish", ItemTypeMod, 25000},
		{"furious_javelin", ItemTypeMod, 25000},
		{"chromatic_blade", ItemTypeMod, 25000},
		{"shattered_storm", ItemTypeMod, 25000},
		{"mending_splinters", ItemTypeMod, 25000},
		{"spectrosiphon", ItemTypeMod, 25000},
		{"mach_crash", ItemTypeMod, 25000},
		{"thermal_transfer", ItemTypeMod, 25000},
		{"cathode_current", ItemTypeMod, 25000},
		{"tribunal", ItemTypeMod, 25000},
		{"warding_thurible", ItemTypeMod, 25000},
		{"lasting_covenant", ItemTypeMod, 25000},
		{"desiccations_curse", ItemTypeMod, 25000},
		{"elemental_sandstorm", ItemTypeMod, 25000},
		{"negation_swarm", ItemTypeMod, 25000},
		{"rift_haven", ItemTypeMod, 25000},
		{"rift_torrent", ItemTypeMod, 25000},
		{"cataclysmic_continuum", ItemTypeMod, 25000},
		{"savior_decoy", ItemTypeMod, 25000},
		{"hushed_invisibility", ItemTypeMod, 25000},
		{"safeguard_switch", ItemTypeMod, 25000},
		{"irradiating_disarm", ItemTypeMod, 25000},
		{"hall_of_malevolence", ItemTypeMod, 25000},
		{"explosive_legerdemain", ItemTypeMod, 25000},
		{"total_eclipse", ItemTypeMod, 25000},
		{"mind_freak", ItemTypeMod, 25000},
		{"pacifying_bolts", ItemTypeMod, 25000},
		{"chaos_sphere", ItemTypeMod, 25000},
		{"assimilate", ItemTypeMod, 25000},
		{"repair_dispensary", ItemTypeMod, 25000},
		{"temporal_erosion", ItemTypeMod, 25000},
		{"intrepid_stand", ItemTypeMod, 25000},
		{"shock_trooper", ItemTypeMod, 25000},
		{"shocking_speed", ItemTypeMod, 25000},
		{"transistor_shield", ItemTypeMod, 25000},
		{"capacitance", ItemTypeMod, 25000},
		{"celestial_stomp", ItemTypeMod, 25000},
		{"enveloping_cloud", ItemTypeMod, 25000},
		{"primal_rage", ItemTypeMod, 25000},
		{"elusive_retribution", ItemTypeMod, 25000},
		{"damage_decoy", ItemTypeMod, 25000},
		{"axios_javelineers", ItemTypeMod, 25000},
		{"warriors_rest", ItemTypeMod, 25000},
		{"coil_recharge", ItemTypeMod, 25000},

		{"decurion_barrel", ItemTypeArchPart, 20000},
		{"phaedra_barrel", ItemTypeArchPart, 20000},
		{"corvas_barrel", ItemTypeArchPart, 20000},
		{"cyngas_barrel", ItemTypeArchPart, 20000},
		{"centaur_aegis", ItemTypeArchPart, 20000},

		{"telos_akbolto", ItemTypeWeapon, 100000},
		{"telos_boltor", ItemTypeWeapon, 125000},
		{"telos_boltace", ItemTypeWeapon, 125000},
	},
}
