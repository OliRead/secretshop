package secretshop

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

var friendlyNames map[string]string

// Config contains details to set up the application
type Config struct {
	BindAddress string                  `toml:"bind"`
	Auth        string                  `toml:"auth"`
	StoreInfo   map[string]ConfigDBInfo `toml:"stores"`
	Stores      map[string]Store
}

// ConfigDBInfo contains details for a database to be used as a store
type ConfigDBInfo struct {
	Address string
	Port    int
	User    string
	Pass    string
	DB      string
}

// ItemPurchase contains information about an individual item purchase
type ItemPurchase struct {
	Item      string      `json:"item"`
	Hero      string      `json:"hero"`
	GameID    uint64      `json:"gameId"`
	SteamID   uint64      `json:"steamId"`
	Timestamp float32     `json:"timestamp"`
	Raw       interface{} `json:"raw,omitempty"`
}

// PlayerInfo contains information about a player and their team and steam account
type PlayerInfo struct {
	SteamID uint64 `json:"steamId"`
	Team    string `json:"team"`
	Name    string `json:"name"`
}

// Store handles interactions with a Database
type Store interface {
	SaveReplayInfo(*Replay) error
	SaveReplayInfoFriendlyName(uint64, string) error
	LoadReplayInfo([]uint64) (map[uint64]Replay, error)
	SavePlayerInfo(*PlayerInfo) error
	LoadPlayerInfo() (map[uint64]PlayerInfo, error)
	SaveItemPurchase(*ItemPurchase) error
	LoadItemPurchase(map[string]interface{}) ([]ItemPurchase, error)
}

func init() {
	friendlyNames = make(map[string]string)

	// Basic Items
	friendlyNames["item_aegis"] = "Aegis of the Immortal"
	friendlyNames["item_courier"] = "Animal Courier"
	friendlyNames["item_boots_of_elves"] = "Band of Elvenskin"
	friendlyNames["item_belt_of_strength"] = "Belt of Strength"
	friendlyNames["item_blade_of_alacrity"] = "Blade of Alacrity"
	friendlyNames["item_blades_of_attack"] = "Blades of Attack"
	friendlyNames["item_blight_stone"] = "Blight Stone"
	friendlyNames["item_blink"] = "Blink Dagger"
	friendlyNames["item_boots"] = "Boots of Speed"
	friendlyNames["item_bottle"] = "Bottle"
	friendlyNames["item_broadsword"] = "Broadsword"
	friendlyNames["item_chainmail"] = "Chainmail"
	friendlyNames["item_cheese"] = "Cheese"
	friendlyNames["item_circlet"] = "Circlet"
	friendlyNames["item_clarity"] = "Clarity"
	friendlyNames["item_claymore"] = "Claymore"
	friendlyNames["item_cloak"] = "Cloak"
	friendlyNames["item_demon_edge"] = "Demon Edge"
	friendlyNames["item_dust"] = "Dust of Appearance"
	friendlyNames["item_eagle"] = "Eaglesong"
	friendlyNames["item_enchanted_mango"] = "Enchanted Mango"
	friendlyNames["item_energy_booster"] = "Energy Booster"
	friendlyNames["item_faerie_fire"] = "Faerie Fire"
	friendlyNames["item_flying_courier"] = "Flying Courier"
	friendlyNames["item_gauntlets"] = "Gauntlets of Strength"
	friendlyNames["item_gem"] = "Gem of True Sight"
	friendlyNames["item_ghost"] = "Ghost Sceptre"
	friendlyNames["item_gloves"] = "Gloves of Haste"
	friendlyNames["item_flask"] = "Healing Salve"
	friendlyNames["item_helm_of_iron_will"] = "Helm of Iron Will"
	friendlyNames["item_hyperstone"] = "Hyperstone"
	friendlyNames["item_infused_raindrop"] = "Infused Raindrop"
	friendlyNames["item_branches"] = "Iron Branch"
	friendlyNames["item_javelin"] = "Javelin"
	friendlyNames["item_magic_stick"] = "Magic Stick"
	friendlyNames["item_mantle"] = "Mantle of Intelligence"
	friendlyNames["item_mithril_hammer"] = "Mithril Hammer"
	friendlyNames["item_lifesteal"] = "Morbid Mask"
	friendlyNames["item_mystic_staff"] = "Mystic Staff"
	friendlyNames["item_ward_observer"] = "Observer Ward"
	friendlyNames["item_ogre_axe"] = "Ogre Club"
	friendlyNames["item_orb_of_venom"] = "Orb of Venom"
	friendlyNames["item_platemail"] = "Platemail"
	friendlyNames["item_point_booster"] = "Point Booster"
	friendlyNames["item_quarterstaff"] = "Quarterstaff"
	friendlyNames["item_quelling_blade"] = "Quelling Blade"
	friendlyNames["item_reaver"] = "Reaver"
	friendlyNames["item_ring_of_health"] = "Ring of Health"
	friendlyNames["item_ring_of_protection"] = "Ring of Protection"
	friendlyNames["item_ring_of_regen"] = "Ring of Regen"
	friendlyNames["item_robe"] = "Robe of the Magi"
	friendlyNames["item_relic"] = "Sacred Relic"
	friendlyNames["item_sobi_mask"] = "Sage's Mask"
	friendlyNames["item_ward_sentry"] = "Sentry Ward"
	friendlyNames["item_shadow_amulet"] = "Shadow Amulet"
	friendlyNames["item_slippers"] = "Slippers of Agility"
	friendlyNames["item_smoke_of_deceit"] = "Smoke of Deceit"
	friendlyNames["item_staff_of_wizardry"] = "Staff of Wizardry"
	friendlyNames["item_stout_shield"] = "Stout Shield"
	friendlyNames["item_talisman_of_evasion"] = "Talisman of Evasion"
	friendlyNames["item_tango"] = "Tango"
	friendlyNames["item_tango_single"] = "Tango (Shared)"
	friendlyNames["item_tome_of_knowledge"] = "Tome of Knowledge"
	friendlyNames["item_tpscroll"] = "Town Portal Scroll"
	friendlyNames["item_ultimate_orb"] = "Ultimate Orb"
	friendlyNames["item_vitality_booster"] = "Vitality Booster"
	friendlyNames["item_void_stone"] = "Void Stone"
	friendlyNames["item_wind_lace"] = "Wind Lace"

	// Upgraded Items
	friendlyNames["item_abyssal_blade"] = "Abyssal Blade"
	friendlyNames["item_recipe_abyssal_blade"] = "Abyssal Blade Recipe"
	friendlyNames["item_aether_lens"] = "Aether Lens"
	friendlyNames["item_recipe_aether_lens"] = "Aether Lens Recipe"
	friendlyNames["item_ultimate_scepter"] = "Aghanim's Scepter"
	friendlyNames["item_recipe_ultimate_scepter"] = "Aghanim's Scepter Recipe"
	friendlyNames["item_arcane_boots"] = "Arcane Boots"
	friendlyNames["item_recipe_arcane_boots"] = "Arcane Boots Recipe"
	friendlyNames["item_armlet"] = "Armlet of Mordiggian"
	friendlyNames["item_recipe_armlet"] = "Armlet of Mordiggian Recipe"
	friendlyNames["item_assault"] = "Assault Cuirass"
	friendlyNames["item_recipe_assault"] = "Assault Cuirass Recipe"
	friendlyNames["item_bfury"] = "Battle Fury"
	friendlyNames["item_recipe_bfury"] = "Battle Fury Recipe"
	friendlyNames["item_black_king_bar"] = "Black King Bar"
	friendlyNames["item_recipe_black_king_bar"] = "Black King Bar Recipe"
	friendlyNames["item_blade_mail"] = "Blade Mail"
	friendlyNames["item_recipe_blade_mail"] = "Blade Mail Recipe"
	friendlyNames["item_bloodstone"] = "Bloodstone"
	friendlyNames["item_recipe_bloodstone"] = "Bloodstone Recipe"
	friendlyNames["item_bloodthorn"] = "Bloodthorn"
	friendlyNames["item_recipe_bloodthorn"] = "Bloodthorn Recipe"
	friendlyNames["item_travel_boots"] = "Boots of Travel (Level 1)"
	friendlyNames["item_travel_boots_2"] = "Boots of Travel (Level 2)"
	friendlyNames["item_recipe_travel_boots"] = "Boots of Travel Recipe"
	friendlyNames["item_bracer"] = "Bracer"
	friendlyNames["item_recipe_bracer"] = "Bracer Recipe"
	friendlyNames["item_buckler"] = "Buckler"
	friendlyNames["item_recipe_buckler"] = "Buckler Recipe"
	friendlyNames["item_butterfly"] = "Butterfly"
	friendlyNames["item_recipe_butterfly"] = "Butterfly Recipe"
	friendlyNames["item_crimson_guard"] = "Crimson Guard"
	friendlyNames["item_recipe_crimson_guard"] = "Crimson Guard Recipe"
	friendlyNames["item_lesser_crit"] = "Crystalys"
	friendlyNames["item_recipe_lesser_crit"] = "Crystalys Recipe"
	friendlyNames["item_greater_crit"] = "Daedalus"
	friendlyNames["item_recipe_greater_crit"] = "Daedalus Recipe"
	friendlyNames["item_dagon_#1"] = "Dagon (Level 1)"
	friendlyNames["item_dagon_#2"] = "Dagon (Level 2)"
	friendlyNames["item_dagon_#3"] = "Dagon (Level 3)"
	friendlyNames["item_dagon_#4"] = "Dagon (Level 4)"
	friendlyNames["item_dagon_#5"] = "Dagon (Level 5)"
	friendlyNames["item_recipe_dagon"] = "Dagon Recipe"
	friendlyNames["item_desolator"] = "Desloator"
	friendlyNames["item_recipe_desolator"] = "Desolator Recipe"
	friendlyNames["item_diffusal_blade"] = "Diffusal Blade (Level 1)"
	friendlyNames["item_diffusal_blade_2"] = "Diffusal Blade (Level 2)"
	friendlyNames["item_recipe_diffusal_blade"] = "Diffusal Blade Recipe"
	friendlyNames["item_dragon_lance"] = "Dragon Lance"
	friendlyNames["item_recipe_dragon_lance"] = "Dragon Lance Recipe"
	friendlyNames["item_ancient_janggo"] = "Drum of Endurance"
	friendlyNames["item_recipe_ancient_janggo"] = "Drum of Endurance Recipe"
	friendlyNames["item_echo_sabre"] = "Echo Sabre"
	friendlyNames["item_recipe_echo_sabre"] = "Echo Sabre Recipe"
	friendlyNames["item_ethereal_blade"] = "Ethereal Blade"
	friendlyNames["item_ethereal_blade_recipe"] = "Ethereal Blade Recipe"
	friendlyNames["item_cyclone"] = "Eul's Sceptre of Divinity"
	friendlyNames["item_recipe_cyclone"] = "Eul's Sceptre of Divinity Recipe"
	friendlyNames["item_skadi"] = "Eye of Skadi"
	friendlyNames["item_recipe_skadi"] = "Eye of Skadi Recipe"
	friendlyNames["item_force_staff"] = "Force Staff"
	friendlyNames["item_recipe_force_staff"] = "Force Staff Recipe"
	friendlyNames["item_glimmer_cape"] = "Glimmer Cape"
	friendlyNames["item_recipe_glimmer_cape"] = "Glimmer Cape Recipe"
	friendlyNames["item_guardian_greaves"] = "Guardian Greaves"
	friendlyNames["item_recipe_guardian_greaves"] = "Guardian Greaves Recipe"
	friendlyNames["item_hand_of_midas"] = "Hand of Midas"
	friendlyNames["item_recipe_hand_of_midas"] = "Hand of Midas Recipe"
	friendlyNames["item_headdress"] = "Headdress"
	friendlyNames["item_recipe_headdress"] = "Headdress Recipe"
	friendlyNames["item_heart"] = "Heart of Tarrasque"
	friendlyNames["item_recipe_heart"] = "Heart of Tarrasque Recipe"
	friendlyNames["item_heavens_halberd"] = "Heaven's Halberd"
	friendlyNames["item_recipe_heavens_halberd"] = "Heaven's Halberd Recipe"
	friendlyNames["item_helm_of_the_dominator"] = "Helm of the Dominator"
	friendlyNames["item_recipe_helm_of_the_dominator"] = "Helm of the Dominator Recipe"
	friendlyNames["item_hood_of_defiance"] = "Hood of Defiance"
	friendlyNames["item_recipe_hood_of_defiance"] = "Hood of Defiance Recipe"
	friendlyNames["item_hurricane_pike"] = "Hurricane Pike"
	friendlyNames["item_recipe_hurricane_pike"] = "Hurricane Pike Recipe"
	friendlyNames["item_iron_talon"] = "Iron Talon"
	friendlyNames["item_recipe_iron_talon"] = "Iron Talon Recipe"
	friendlyNames["item_sphere"] = "Linken's Sphere"
	friendlyNames["item_recipe_sphere"] = "Linken's Sphere Recipe"
	friendlyNames["item_lotus_orb"] = "Lotus Orb"
	friendlyNames["item_recipe_lotus_orb"] = "Lotus Orb Recipe"
	friendlyNames["item_maelstrom"] = "Malestrom"
	friendlyNames["item_recipe_maelstrom"] = "Maelstrom Recipe"
	friendlyNames["item_magic_wand"] = "Magic Wand"
	friendlyNames["item_recipe_magic_wand"] = "Magic Wand Recipe"
	friendlyNames["item_manta"] = "Manta Style"
	friendlyNames["item_recipe_manta"] = "Manta Style Recipe"
	friendlyNames["item_mask_of_madness"] = "Mask of Madness"
	friendlyNames["item_recipe_mask_of_madness"] = "Mask of Madness Recipe"
	friendlyNames["item_medallion_of_courage"] = "Medallion of Courage"
	friendlyNames["item_recipe_medallion_of_courage"] = "Medallion of Courage Recipe"
	friendlyNames["item_mekansm"] = "Mekansm"
	friendlyNames["item_recipe_mekansm"] = "Mekansm Recipe"
	friendlyNames["item_mjollnir"] = "Mjollnir"
	friendlyNames["item_recipe_mjollnir"] = "Mjollnir Recipe"
	friendlyNames["item_monkey_king_bar"] = "Monkey King Bar"
	friendlyNames["item_recipe_monkey_king_bar"] = "Monkey King Bar Recipe"
	friendlyNames["item_moon_shard"] = "Moon Shard"
	friendlyNames["item_recipe_moon_shard"] = "Moon Shard Recipe"
	friendlyNames["item_necronomicon"] = "Necronomicon (Level 1)"
	friendlyNames["item_necronomicon_2"] = "Necronomicon (Level 2)"
	friendlyNames["item_necronomicon_3"] = "Necronomicon (Level 3)"
	friendlyNames["item_recipe_necronomicon"] = "Necronomicon Recipe"
	friendlyNames["item_null_talisman"] = "Null Talisman"
	friendlyNames["item_recipe_null_talisman"] = "Null Talisman Recipe"
	friendlyNames["item_oblivion_staff"] = "Oblivion Staff"
	friendlyNames["item_recipe_oblivion_staff"] = "Oblivion Staff Recipe"
	friendlyNames["item_ward_dispenser"] = "Observer and Sentry Wards"
	friendlyNames["item_recipe_ward_dispenser"] = "Observer and Sentry Wards Recipe"
	friendlyNames["item_octarine_core"] = "Octarine Core"
	friendlyNames["item_recipe_octarine_core"] = "Octarine Core Recipe"
	friendlyNames["item_orchid"] = "Orchid Malevolance"
	friendlyNames["item_recipe_orchid"] = "Orchid Malevolance Recipe"
	friendlyNames["item_pers"] = "Perseverance"
	friendlyNames["item_recipe_pers"] = "Perseverance Recipe"
	friendlyNames["item_phase_boots"] = "Phase Boots"
	friendlyNames["item_pipe"] = "Pipe of Insight"
	friendlyNames["item_recipe_pipe"] = "Pipe of Insight Recipe"
	friendlyNames["item_poor_mans_shield"] = "Poor Man's Shield"
	friendlyNames["item_recipe_poor_mans_shield"] = "Poor Man's Shield Recipe"
	friendlyNames["item_power_treads"] = "Power Treads"
	friendlyNames["item_recipe_power_treads"] = "Power Treads Recipe"
	friendlyNames["item_radiance"] = "Radiance"
	friendlyNames["item_recipe_radiance"] = "Radiance Recipe"
	friendlyNames["item_rapier"] = "Divine Rapier"
	friendlyNames["item_recipe_rapier"] = "Divine Rapier Recipe"
	friendlyNames["item_refresher"] = "Refresher Orb"
	friendlyNames["item_recipe_refresher"] = "Refresher Orb Recipe"
	friendlyNames["item_ring_of_aquila"] = "Ring of Aquila"
	friendlyNames["item_recipe_ring_of_aquila"] = "Ring of Aquila Recipe"
	friendlyNames["item_ring_of_basilius"] = "Ring of Basilius"
	friendlyNames["item_recipe_ring_of_basilius"] = "Ring of Basilius Recipe"
	friendlyNames["item_rod_of_atos"] = "Rod of Atos"
	friendlyNames["item_recipe_rod_of_atos"] = "Rod of Atos Recipe"
	friendlyNames["item_sange"] = "Sange"
	friendlyNames["item_recipe_sange"] = "Sange Recipe"
	friendlyNames["item_sange_and_yasha"] = "Sange and Yasha"
	friendlyNames["item_recipe_sange_and_yasha"] = "Sange and Yasha Recipe"
	friendlyNames["item_satanic"] = "Satanic"
	friendlyNames["item_recipe_satanic"] = "Satanic Recipe"
	friendlyNames["item_sheepstick"] = "Scythe of Vyse"
	friendlyNames["item_recipe_sheepstick"] = "Scythe of Vyse Recipe"
	friendlyNames["item_invis_sword"] = "Shadow Blade"
	friendlyNames["item_recipe_invis_sword"] = "Shadow Blade Recipe"
	friendlyNames["item_shivas_guard"] = "Shiva's Guard"
	friendlyNames["item_recipe_shivas_guard"] = "Shiva's Guard Recipe"
	friendlyNames["item_silver_edge"] = "Silver Edge"
	friendlyNames["item_recipe_silver_edge"] = "Silver Edge Recipe"
	friendlyNames["item_basher"] = "Skull Basher"
	friendlyNames["item_recipe_basher"] = "Skull Basher Recipe"
	friendlyNames["item_solar_crest"] = "Solar Crest"
	friendlyNames["item_recipe_solar_crest"] = "Solar Crest Recipe"
	friendlyNames["item_soul_booster"] = "Soul Booster"
	friendlyNames["item_recipe_soul_booster"] = "Soul Booster Recipe"
	friendlyNames["item_soul_ring"] = "Soul Ring"
	friendlyNames["item_recipe_soul_ring"] = "Soul Ring Recipe"
	friendlyNames["item_tranquil_boots"] = "Tranquil Boots"
	friendlyNames["item_recipe_tranquil_boots"] = "Tranquil Boots Recipe"
	friendlyNames["item_urn_of_shadows"] = "Urn of Shadows"
	friendlyNames["item_recipe_urn_of_shadows"] = "Urn of Shadows Recipe"
	friendlyNames["item_vanguard"] = "Vanguard"
	friendlyNames["item_recipe_vanguard"] = "Vanguard Recipe"
	friendlyNames["item_veil_of_discord"] = "Veil of Discord"
	friendlyNames["item_recipe_veil_of_discord"] = "Veil of Discord Recipe"
	friendlyNames["item_vladmir"] = "Vladmir's Offering"
	friendlyNames["item_recipe_vladmir"] = "Vladmir's Offering Recipe"
	friendlyNames["item_wraith_band"] = "Wraith Band"
	friendlyNames["item_recipe_wraith_band"] = "Wraith Band Recipe"
	friendlyNames["item_yasha"] = "Yasha"
	friendlyNames["item_recipe_yasha"] = "Yasha Recipe"

	// Hero Names
	friendlyNames["npc_dota_hero_abaddon"] = "Abaddon"
	friendlyNames["npc_dota_hero_alchemist"] = "Alchemist"
	friendlyNames["npc_dota_hero_antimage"] = "Anti-Mage"
	friendlyNames["npc_dota_hero_ancient_apparition"] = "Ancient Apparition"
	friendlyNames["npc_dota_hero_arcwarden"] = "Arc Warden"
	friendlyNames["npc_dota_hero_axe"] = "Axe"
	friendlyNames["npc_dota_hero_bane"] = "Bane"
	friendlyNames["npc_dota_hero_batrider"] = "Batrider"
	friendlyNames["npc_dota_hero_beastmaster"] = "Beastmaster"
	friendlyNames["npc_dota_hero_bloodseeker"] = "Bloodseeker"
	friendlyNames["npc_dota_hero_bounty_hunter"] = "Bounty Hunter"
	friendlyNames["npc_dota_hero_brewmaster"] = "Brewmaster"
	friendlyNames["npc_dota_hero_bristleback"] = "Bristleback"
	friendlyNames["npc_dota_hero_broodmother"] = "Broodmother"
	friendlyNames["npc_dota_hero_centaur"] = "Centaur Warrunner"
	friendlyNames["npc_dota_hero_chaos_knight"] = "Chaos Knight"
	friendlyNames["npc_dota_hero_chen"] = "Chen"
	friendlyNames["npc_dota_hero_clinkz"] = "Clinkz"
	friendlyNames["npc_dota_hero_rattletrap"] = "Clockwerk"
	friendlyNames["npc_dota_hero_crystal_maiden"] = "Crystal Maiden"
	friendlyNames["npc_dota_hero_dark_seer"] = "Dark Seer"
	friendlyNames["npc_dota_hero_dazzle"] = "Dazzle"
	friendlyNames["npc_dota_hero_death_prophet"] = "Death Prophet"
	friendlyNames["npc_dota_hero_disruptor"] = "Disruptor"
	friendlyNames["npc_dota_hero_doom_bringer"] = "Doom"
	friendlyNames["npc_dota_hero_dragon_knight"] = "Dragon Knight"
	friendlyNames["npc_dota_hero_drow_ranger"] = "Drow Ranger"
	friendlyNames["npc_dota_hero_earth_spirit"] = "Earth Spirit"
	friendlyNames["npc_dota_hero_earthshaker"] = "Earthshaker"
	friendlyNames["npc_dota_hero_elder_titan"] = "Elder Titan"
	friendlyNames["npc_dota_hero_ember_spirit"] = "Ember Spirit"
	friendlyNames["npc_dota_hero_enchantress"] = "Enchantress"
	friendlyNames["npc_dota_hero_enigma"] = "Enigma"
	friendlyNames["npc_dota_hero_faceless_void"] = "Faceless Void"
	friendlyNames["npc_dota_hero_gyrocopter"] = "Gyrocopter"
	friendlyNames["npc_dota_hero_huskar"] = "Huskar"
	friendlyNames["npc_dota_hero_invoker"] = "Invoker"
	friendlyNames["npc_dota_hero_wisp"] = "IO"
	friendlyNames["npc_dota_hero_jakiro"] = "Jakiro"
	friendlyNames["npc_dota_hero_juggernaut"] = "Juggernaut"
	friendlyNames["npc_dota_hero_keeper_of_the_light"] = "Keeper of the Light"
	friendlyNames["npc_dota_hero_kunkka"] = "Kunkka"
	friendlyNames["npc_dota_hero_legion_commander"] = "Legion Commander"
	friendlyNames["npc_dota_hero_leshrac"] = "Leshrac"
	friendlyNames["npc_dota_hero_lich"] = "Lich"
	friendlyNames["npc_dota_hero_life_stealer"] = "Life Stealer"
	friendlyNames["npc_dota_hero_lina"] = "Lina"
	friendlyNames["npc_dota_hero_lion"] = "Lion"
	friendlyNames["npc_dota_hero_lone_druid"] = "Lone Druid"
	friendlyNames["npc_dota_hero_luna"] = "Luna"
	friendlyNames["npc_dota_hero_lycan"] = "Lycan"
	friendlyNames["npc_dota_hero_magnataur"] = "Magnus"
	friendlyNames["npc_dota_hero_medusa"] = "Medusa"
	friendlyNames["npc_dota_hero_meepo"] = "Meepo"
	friendlyNames["npc_dota_hero_mirana"] = "Mirana"
	friendlyNames["npc_dota_hero_morphling"] = "Morphling"
	friendlyNames["npc_dota_hero_monkey_king"] = "Monkey King"
	friendlyNames["npc_dota_hero_naga_siren"] = "Naga Siren"
	friendlyNames["npc_dota_hero_furion"] = "Nature's Prophet"
	friendlyNames["npc_dota_hero_necrolyte"] = "Necrophos"
	friendlyNames["npc_dota_hero_night_stalker"] = "Night Stalker"
	friendlyNames["npc_dota_hero_nyx_assassin"] = "Nyx Assassin"
	friendlyNames["npc_dota_hero_ogre_magi"] = "Ogre Magi"
	friendlyNames["npc_dota_hero_omniknight"] = "Omniknight"
	friendlyNames["npc_dota_hero_oracle"] = "Oracle"
	friendlyNames["npc_dota_hero_obsidian_destroyer"] = "Outworld Devourer"
	friendlyNames["npc_dota_hero_phantom_assassin"] = "Phantom Assassin"
	friendlyNames["npc_dota_hero_phantom_lancer"] = "Phantom Lancer"
	friendlyNames["npc_dota_hero_phoenix"] = "Phoenix"
	friendlyNames["npc_dota_hero_puck"] = "Puck"
	friendlyNames["npc_dota_hero_pudge"] = "Pudge"
	friendlyNames["npc_dota_hero_pugna"] = "Pugna"
	friendlyNames["npc_dota_hero_queenofpain"] = "Queen of Pain"
	friendlyNames["npc_dota_hero_razor"] = "Razor"
	friendlyNames["npc_dota_hero_riki"] = "Riki"
	friendlyNames["npc_dota_hero_rubick"] = "Rubick"
	friendlyNames["npc_dota_hero_sand_king"] = "Sand King"
	friendlyNames["npc_dota_hero_shadow_demon"] = "Shadow Demon"
	friendlyNames["npc_dota_hero_nevermore"] = "Shadow Fiend"
	friendlyNames["npc_dota_hero_shadow_shaman"] = "Shadow Shaman"
	friendlyNames["npc_dota_hero_silencer"] = "Silencer"
	friendlyNames["npc_dota_hero_skywrath_mage"] = "Skywrath Mage"
	friendlyNames["npc_dota_hero_slardar"] = "Slardar"
	friendlyNames["npc_dota_hero_slark"] = "Slark"
	friendlyNames["npc_dota_hero_sniper"] = "Sniper"
	friendlyNames["npc_dota_hero_spectre"] = "Spectre"
	friendlyNames["npc_dota_hero_spirit_breaker"] = "Spirit Breaker"
	friendlyNames["npc_dota_hero_storm_spirit"] = "Storm Sprit"
	friendlyNames["npc_dota_hero_sven"] = "Sven"
	friendlyNames["npc_dota_hero_techies"] = "Techies"
	friendlyNames["npc_dota_hero_templar_assassin"] = "Templar Assassin"
	friendlyNames["npc_dota_hero_terrorblade"] = "Terrorblade"
	friendlyNames["npc_dota_hero_tidehunter"] = "Tidehunter"
	friendlyNames["npc_dota_hero_shredder"] = "Timbersaw"
	friendlyNames["npc_dota_hero_tinker"] = "Tinker"
	friendlyNames["npc_dota_hero_tiny"] = "Tiny"
	friendlyNames["npc_dota_hero_treant"] = "Treant Protector"
	friendlyNames["npc_dota_hero_troll_warlord"] = "Troll Warlord"
	friendlyNames["npc_dota_hero_tusk"] = "Tusk"
	friendlyNames["npc_dota_hero_abyssal_underlord"] = "Underlord"
	friendlyNames["npc_dota_hero_undying"] = "Undying"
	friendlyNames["npc_dota_hero_ursa"] = "Ursa"
	friendlyNames["npc_dota_hero_vengefulspirit"] = "Vengeful Spirit"
	friendlyNames["npc_dota_hero_venomancer"] = "Venomancer"
	friendlyNames["npc_dota_hero_viper"] = "Viper"
	friendlyNames["npc_dota_hero_visage"] = "Visage"
	friendlyNames["npc_dota_hero_warlock"] = "Warlock"
	friendlyNames["npc_dota_hero_weaver"] = "Weaver"
	friendlyNames["npc_dota_hero_windrunner"] = "Windranger"
	friendlyNames["npc_dota_hero_winter_wyvern"] = "Winter Wyvern"
	friendlyNames["npc_dota_hero_witch_doctor"] = "Witch Doctor"
	friendlyNames["npc_dota_hero_skeleton_king"] = "Wraith King"
	friendlyNames["npc_dota_hero_zuus"] = "Zeus"
}

// ReadConfig reads and parses a TOML Config
func ReadConfig(confFile string) (c Config, err error) {
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		return c, err
	}

	if _, err := toml.Decode(string(data), &c); err != nil {
		return c, err
	}

	c.Stores = make(map[string]Store)

	return c, nil
}

// GetFriendlyName returns a friendly name from an internal name
func GetFriendlyName(internalName string) (friendlyName string, err error) {
	if _, ok := friendlyNames[internalName]; !ok {
		return internalName, fmt.Errorf("Internal Name [%s] does not have a friendly name", internalName)
	}

	return friendlyNames[internalName], nil
}
