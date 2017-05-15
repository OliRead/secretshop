package secretshop

import (
	"fmt"
	"os"

	"github.com/dotabuff/manta"
	"github.com/dotabuff/manta/dota"
)

// Replay holds information about a replay file
type Replay struct {
	File          *os.File          `json:"file,omitempty"`
	StrategyStart float32           `json:"strategyStart"`
	GameStart     float32           `json:"gameStart"`
	GameEnd       float32           `json:"gameEnd"`
	GameID        uint64            `json:"gameId"`
	ItemPurchases []*ItemPurchase   `json:"itemPurchases,omitempty"`
	Players       map[string]uint64 `json:"players"`
	PlayerInfo    []*PlayerInfo     `json:"playerInfo"`
	FriendlyName  string            `json:"friendlyName"`
}

// NewReplay initializes a replay ready to be parsed
func NewReplay(fileName string) (r *Replay, err error) {
	r = &Replay{}
	if r.File, err = os.Open(fileName); err != nil {
		return nil, fmt.Errorf("unable to open file: %s", err)
	}

	r.Players = make(map[string]uint64)

	return r, nil
}

// Parse reads a replay file and pulls out data from it
func (r *Replay) Parse() error {
	defer r.File.Close()

	p, err := manta.NewStreamParser(r.File)
	if err != nil {
		return fmt.Errorf("unable to create parser: %s", err)
	}

	p.Callbacks.OnCDemoFileInfo(func(m *dota.CDemoFileInfo) error {
		data := m.GameInfo.GetDota()
		r.GameID = *data.MatchId
		for _, player := range data.PlayerInfo {
			playerInfo := PlayerInfo{
				SteamID: *player.Steamid,
				Name:    *player.PlayerName,
			}
			r.PlayerInfo = append(r.PlayerInfo, &playerInfo)
			r.Players[*player.HeroName] = *player.Steamid
		}
		return nil
	})

	p.Callbacks.OnCMsgDOTACombatLogEntry(func(m *dota.CMsgDOTACombatLogEntry) error {
		t := m.GetType()

		if t == dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_GAME_STATE {
			v := m.GetValue()
			if v == uint32(dota.DOTA_GameState_DOTA_GAMERULES_STATE_GAME_IN_PROGRESS) {
				r.GameStart = *m.Timestamp
			} else if v == uint32(dota.DOTA_GameState_DOTA_GAMERULES_STATE_STRATEGY_TIME) {
				r.StrategyStart = *m.Timestamp
			} else if v == uint32(dota.DOTA_GameState_DOTA_GAMERULES_STATE_POST_GAME) {
				r.GameEnd = *m.Timestamp
			}

			return nil
		}

		if t == dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_PURCHASE {
			item, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetValue()))
			hero, _ := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetName()))
			timestamp := m.GetTimestamp()

			purchase := ItemPurchase{
				Item:      item,
				Hero:      hero,
				Timestamp: timestamp,
				Raw:       m,
			}

			r.ItemPurchases = append(r.ItemPurchases, &purchase)
			return nil
		}

		return nil
	})

	if err := p.Start(); err != nil {
		return err
	}

	return nil
}

// Process fills in any missing information from a replay after parsing it
func (r *Replay) Process() {
	for _, p := range r.ItemPurchases {
		p.GameID = r.GameID
		p.SteamID = r.Players[p.Hero]
	}
}
