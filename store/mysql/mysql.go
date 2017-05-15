package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oliread/secretshop"
)

type processedReplay struct {
	GameID        uint64
	GameStart     float32
	GameEnd       float32
	StrategyStart float32
	Players       string
	Heroes        string
	FriendlyName  string
}

// Store implementation of secretshop.Store
type Store struct {
	db *sql.DB
}

// NewStore handles creating a store and connecting to a database with information
// from a config file
func NewStore(c *secretshop.Config, data secretshop.ConfigDBInfo) (err error) {
	connString := data.Address
	connPort := strconv.Itoa(data.Port)
	if connPort != "0" {
		connString = strings.Join([]string{data.Address, connPort}, ":")
	}

	connInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s", data.User, data.Pass, connString, data.DB)
	db, err := sql.Open("mysql", connInfo)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	c.Stores["mysql"] = Store{
		db: db,
	}

	return nil
}

// SaveItemPurchase implementation for secretshop.Store
func (s Store) SaveItemPurchase(i *secretshop.ItemPurchase) error {
	stmt, err := s.db.Prepare("INSERT item_purchase SET gameId=?,steamId=?,hero=?,item=?,timestamp=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(i.GameID, i.SteamID, i.Hero, i.Item, i.Timestamp); err != nil {
		return err
	}

	return nil
}

// LoadItemPurchase implementation for secretshop.Store
func (s Store) LoadItemPurchase(filters map[string]interface{}) (i []secretshop.ItemPurchase, err error) {
	query := "SELECT * FROM item_purchase"
	conditions := []string{}
	args := []interface{}{}

	if gameIDs, ok := filters["gameId"]; ok {
		vars := make([]string, len(gameIDs.([]uint64)))
		for i, gameID := range gameIDs.([]uint64) {
			vars[i] = "?"
			args = append(args, gameID)
		}
		gameList := strings.Join(vars, ",")
		conditions = append(conditions, "gameId IN ("+gameList+")")
	}

	if players, ok := filters["player"]; ok {
		vars := make([]string, len(players.([]uint64)))
		for i, player := range players.([]uint64) {
			vars[i] = "?"
			args = append(args, player)
		}
		playerList := strings.Join(vars, ",")
		conditions = append(conditions, "steamId IN ("+playerList+")")
	}

	if heroes, ok := filters["hero"]; ok {
		vars := make([]string, len(heroes.([]string)))
		for i, hero := range heroes.([]string) {
			vars[i] = "?"
			args = append(args, hero)
		}
		heroList := strings.Join(vars, ",")
		conditions = append(conditions, "hero IN ("+heroList+")")
	}

	if items, ok := filters["item"]; ok {
		vars := make([]string, len(items.([]string)))
		for i, item := range items.([]string) {
			vars[i] = "?"
			args = append(args, item)
		}
		itemList := strings.Join(vars, ",")
		conditions = append(conditions, "item IN ("+itemList+")")
	}

	var rows *sql.Rows
	if len(filters) > 0 {
		where := strings.Join(conditions, " AND ")
		query = strings.Join([]string{query, where}, " WHERE ")
		rows, err = s.db.Query(query, args...)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = s.db.Query(query)
		if err != nil {
			return nil, err
		}
	}

	for rows.Next() {
		var (
			gameID    uint64
			steamID   uint64
			hero      string
			item      string
			timestamp float32
		)
		err := rows.Scan(&gameID, &steamID, &hero, &item, &timestamp)
		if err != nil {
			return nil, err
		}
		purchase := secretshop.ItemPurchase{
			GameID:    gameID,
			SteamID:   steamID,
			Hero:      hero,
			Item:      item,
			Timestamp: timestamp,
		}
		i = append(i, purchase)
	}

	return i, nil
}

// SaveReplayInfo implementation for secretshop.Store
func (s Store) SaveReplayInfo(r *secretshop.Replay) error {
	p := processReplay(r)
	stmt, err := s.db.Prepare("INSERT replay_info SET gameId=?,strategyStart=?,gameStart=?,gameEnd=?,players=?,heroes=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(p.GameID, p.StrategyStart, p.GameStart, p.GameEnd, p.Players, p.Heroes); err != nil {
		return err
	}

	return nil
}

// LoadReplayInfo implementation for secretshop.Store
func (s Store) LoadReplayInfo(gameIDs []uint64) (map[uint64]secretshop.Replay, error) {
	r := secretshop.Replay{}
	vars := make([]string, len(gameIDs))
	args := make([]interface{}, len(gameIDs))

	for i := 0; i < len(gameIDs); i++ {
		vars[i] = "?"
		args[i] = gameIDs[i]
	}

	query := "SELECT * FROM replay_info"
	if len(gameIDs) > 0 {
		query = "SELECT * FROM replay_info WHERE gameId IN (" + strings.Join(vars, ",") + ")"
	}

	log.Printf("GameIds: %+v", gameIDs)
	log.Printf("Query: %s", query)
	log.Printf("Args: %+v", args)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	replays := make(map[uint64]secretshop.Replay)
	for rows.Next() {
		var (
			id            uint64
			strategyStart float32
			gameStart     float32
			gameEnd       float32
			players       string
			heroes        string
			friendlyName  string
		)
		rows.Scan(&id, &strategyStart, &gameStart, &gameEnd, &players, &heroes, &friendlyName)
		playerInfo := strings.Split(players, ",")
		heroInfo := strings.Split(heroes, ",")
		r.GameID = id
		r.StrategyStart = strategyStart
		r.GameStart = gameStart
		r.GameEnd = gameEnd
		r.Players = make(map[string]uint64)
		r.FriendlyName = friendlyName
		for i := 0; i < len(playerInfo); i++ {
			player, err := strconv.ParseUint(playerInfo[i], 10, 64)
			if err != nil {
				return nil, err
			}

			r.Players[heroInfo[i]] = player
		}
		replays[id] = r
	}

	return replays, nil
}

// SaveReplayInfoFriendlyName implementation for secretshop.Store
func (s Store) SaveReplayInfoFriendlyName(gameID uint64, friendlyName string) error {
	stmt, err := s.db.Prepare("UPDATE replay_info SET friendlyName=? WHERE gameId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(friendlyName, gameID); err != nil {
		return err
	}

	return nil
}

// SavePlayerInfo implementation for secretshop.Store
func (s Store) SavePlayerInfo(p *secretshop.PlayerInfo) error {
	stmt, err := s.db.Prepare("INSERT player_info SET steamId=?,team=?,name=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(p.SteamID, p.Team, p.Name); err != nil {
		return err
	}

	return nil
}

// LoadPlayerInfo implementation for secretshop.Store
func (s Store) LoadPlayerInfo() (p map[uint64]secretshop.PlayerInfo, err error) {
	playerInfo := make(map[uint64]secretshop.PlayerInfo)
	rows, err := s.db.Query("SELECT * FROM player_info")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			steamID uint64
			team    string
			name    string
		)

		if err := rows.Scan(&steamID, &team, &name); err != nil {
			return nil, err
		}

		player := secretshop.PlayerInfo{
			SteamID: steamID,
			Team:    team,
			Name:    name,
		}

		playerInfo[steamID] = player
	}

	return playerInfo, nil
}

func processReplay(r *secretshop.Replay) (p processedReplay) {
	p = processedReplay{
		GameID:        r.GameID,
		GameStart:     r.GameStart,
		GameEnd:       r.GameEnd,
		StrategyStart: r.StrategyStart,
	}

	index := 0
	heroes := make([]string, len(r.Players))
	players := make([]string, len(r.Players))
	for hero, player := range r.Players {
		heroes[index] = hero
		players[index] = strconv.FormatUint(player, 10)
		index++
	}

	p.Heroes = strings.Join(heroes, ",")
	p.Players = strings.Join(players, ",")
	return p
}
