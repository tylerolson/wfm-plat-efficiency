package main

// should be int but idc
type ItemType string

const (
	Mod      = "Mod"
	ArchPart = "ArchPart"
	Weapon   = "Weapon"
)

type Item struct {
	Name string
	Type ItemType
}

var AribiterOfHexis = []Item{
	{"gilded_truth", Mod},
	{"blade_of_truth", Mod},
	{"avenging_truth", Mod},
	{"stinging_truth", Mod},
	{"seeking_shuriken", Mod},
	{"smoke_shadow", Mod},
	{"fatal_teleport", Mod},
	{"rising_storm", Mod},
	{"endless_lullaby", Mod},
	{"reactive_storm", Mod},
	{"duality", Mod},
	{"calm_and_frenzy", Mod},
	{"peaceful_provocation", Mod},
	{"energy_transfer", Mod},
	{"surging_dash", Mod},
	{"radiant_finish", Mod},
	{"furious_javelin", Mod},
	{"chromatic_blade", Mod},
	{"shattered_storm", Mod},
	{"mending_splinters", Mod},
	{"spectrosiphon", Mod},
	{"mach_crash", Mod},
	{"thermal_transfer", Mod},
	{"cathode_current", Mod},
	{"tribunal", Mod},
	{"warding_thurible", Mod},
	{"lasting_covenant", Mod},
	{"desiccations_curse", Mod},
	{"elemental_sandstorm", Mod},
	{"negation_swarm", Mod},
	{"rift_haven", Mod},
	{"rift_torrent", Mod},
	{"cataclysmic_continuum", Mod},
	{"savior_decoy", Mod},
	{"hushed_invisibility", Mod},
	{"safeguard_switch", Mod},
	{"irradiating_disarm", Mod},
	{"hall_of_malevolence", Mod},
	{"explosive_legerdemain", Mod},
	{"total_eclipse", Mod},
	{"mind_freak", Mod},
	{"pacifying_bolts", Mod},
	{"chaos_sphere", Mod},
	{"assimilate", Mod},
	{"repair_dispensary", Mod},
	{"temporal_erosion", Mod},
	{"intrepid_stand", Mod},
	{"shock_trooper", Mod},
	{"shocking_speed", Mod},
	{"transistor_shield", Mod},
	{"capacitance", Mod},
	{"celestial_stomp", Mod},
	{"enveloping_cloud", Mod},
	{"primal_rage", Mod},
	{"elusive_retribution", Mod},
	{"damage_decoy", Mod},
	{"axios_javelineers", Mod},
	{"warriors_rest", Mod},
	{"coil_recharge", Mod},

	{"decurion_barrel", ArchPart},
	{"phaedra_barrel", ArchPart},
	{"corvas_barrel", ArchPart},
	{"cyngas_barrel", ArchPart},
	{"centaur_aegis", ArchPart},

	{"telos_akbolto", Weapon},
	{"telos_boltor", Weapon},
	{"telos_boltace", Weapon},
}
