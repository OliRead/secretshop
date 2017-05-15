package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"strconv"

	"encoding/json"

	"strings"

	"github.com/gorilla/mux"
	"github.com/oliread/secretshop"
)

// Handler contains information about a API http handler
type Handler struct {
	conf   secretshop.Config
	Router *mux.Router
}

// NewHandler creates a new http handler for an API instance
func NewHandler(conf secretshop.Config) (h Handler, err error) {
	h = Handler{
		conf: conf,
	}

	h.Router = mux.NewRouter()
	h.Router.Handle("/replay/upload", h.isAuthenticated(http.HandlerFunc(h.replayNewPost))).Methods("POST")
	h.Router.Handle("/replay/friendlyname", h.isAuthenticated(http.HandlerFunc(h.replayFriendlyNamePost))).Methods("POST")

	h.Router.HandleFunc("/replay/info", h.replayInfoGet).Methods("GET")
	h.Router.HandleFunc("/replay/items", h.itemPurchaseGet).Methods("GET")
	h.Router.HandleFunc("/player/info", h.playerInfoGet).Methods("GET")

	return h, nil
}

func (h *Handler) replayNewPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1024 * 1024 * 300); err != nil {
		log.Printf("Error uploading replay: %s", err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	file, handler, err := r.FormFile("replay")
	f, err := os.OpenFile("/tmp/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	if _, err := io.Copy(f, file); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	f.Close()

	log.Printf("Uploaded file [%s] to: /tmp/%s", handler.Filename, handler.Filename)
	replay, err := secretshop.NewReplay("/tmp/" + handler.Filename)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	log.Printf("Parsing Replay [%s]...", handler.Filename)
	if err := replay.Parse(); err != nil {
		log.Printf("Error parsing replay [%s]: %s", handler.Filename, err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	replay.Process()

	log.Printf("Finished parsing Replay [%s], saving item purchases...", handler.Filename)

	for host, store := range h.conf.Stores {
		info, err := store.LoadReplayInfo([]uint64{replay.GameID})
		if err != nil {
			log.Printf("Error loading replay [%d] from store [%s]: %s", replay.GameID, host, err)
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error loading replay [%d] from store [%s]: %s", replay.GameID, host, err)))
			return
		}

		if len(info) >= 1 {
			log.Printf("Could not save replay info or item purchases, replay [%d] has already been parsed by a store [%s]", replay.GameID, host)
			w.WriteHeader(409)
			w.Write([]byte(fmt.Sprintf("Could not save replay info or item purchases, replay [%d] has already been parsed by a store [%s]", replay.GameID, host)))
			return
		}
	}

	count := 0
	for _, purchase := range replay.ItemPurchases {
		for host, store := range h.conf.Stores {
			if err := store.SaveItemPurchase(purchase); err != nil {
				log.Printf("Could not save purchase [%+v] to store [%s]. %s", purchase, host, err)
			}
		}

		count++
	}

	log.Printf("Finished saving itme purchases, saving player info...")
	for _, player := range replay.PlayerInfo {
		for host, store := range h.conf.Stores {
			if err := store.SavePlayerInfo(player); err != nil {
				log.Printf("Could not save player info [%+v] to store [%s]. %s", player, host, err)
			}
		}
	}

	log.Printf("Finished saving player info, saving replay info...")
	log.Printf("%+v", replay.GameID)
	for host, store := range h.conf.Stores {
		if err := store.SaveReplayInfo(replay); err != nil {
			log.Printf("Could not save replay info [%d] to store [%s]. %s", replay.GameID, host, err)
		}
	}

	log.Printf("Succesfully parsed and saved Replay [%s]. Read %d Purchases and saved %d records", handler.Filename, len(replay.ItemPurchases), count)
	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf("Succesfully parsed and saved Replay [%s]. Read %d Purchases and saved %d records", handler.Filename, len(replay.ItemPurchases), count)))
	return
}

func (h *Handler) replayInfoGet(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	gameIdsRaw := r.URL.Query().Get("gameId")
	log.Printf("Grabbing replay [%s] info from store [%s]", gameIdsRaw, host)

	gameIDsData := strings.Split(gameIdsRaw, ",")
	gameIDs := make([]uint64, len(gameIDsData))

	if _, ok := h.conf.Stores[host]; !ok {
		log.Printf("Can't get replay info from store [%s], store does not exist", host)
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("Can't get replay info from store [%s], store does not exist", host)))
		return
	}

	for i, id := range gameIDsData {
		parsed, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			log.Printf("Could not get replay [%s] info: %s", id, err)
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Could not get replay [%s] info: %s", id, err)))
			return
		}
		gameIDs[i] = parsed
	}

	replay, err := h.conf.Stores[host].LoadReplayInfo(gameIDs)
	if err != nil {
		log.Printf("Error loading replay info [%s] from store [%s]: %s", gameIdsRaw, host, err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error loading replay info [%s] from store [%s]: %s", gameIdsRaw, host, err)))
		return
	}

	data, err := json.Marshal(replay)
	if err != nil {
		log.Printf("Error marshalling replay [%+v] to json: %s", replay, err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error marshalling replay [%+v] to json: %s", replay, err)))
		return
	}

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h *Handler) replayFriendlyNamePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	gameID, err := strconv.ParseUint(r.FormValue("gameId"), 10, 64)
	if err != nil {
		log.Printf("Could not set friendly name for replay [%s]: %s", r.FormValue("gameId"), err)
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("Could not set friendly name for replay [%s]: %s", r.FormValue("gameId"), err)))
		return
	}

	friendlyName := r.FormValue("friendlyName")
	if friendlyName == "" {
		log.Printf("Could not set friendly name for replay [%d]: friendlyName cannot be null", gameID)
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("Could not set friendly name for replay [%d]: friendlyName cannot be null", gameID)))
		return
	}

	for host, store := range h.conf.Stores {
		if err := store.SaveReplayInfoFriendlyName(gameID, friendlyName); err != nil {
			log.Printf("Error saving replay [%d] friendly name [%s] to store [%s]: %s", gameID, friendlyName, host, err)
		}
	}

	log.Printf("Saved replay [%d] friendly name [%s]", gameID, friendlyName)
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Saved replay [%d] friendly name [%s]", gameID, friendlyName)))
}

func (h *Handler) playerInfoGet(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if _, ok := h.conf.Stores[host]; !ok {
		log.Printf("Can't get player info from store [%s], store does not exist", host)
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("Can't get player info from store [%s], store does not exist", host)))
		return
	}

	playerInfo, err := h.conf.Stores[host].LoadPlayerInfo()
	if err != nil {
		log.Printf("Error loading playerInfo from store [%s]: %s", host, err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error loading playerInfo from store [%s]: %s", host, err)))
		return
	}

	payload, err := json.Marshal(playerInfo)
	if err != nil {
		log.Printf("Error marshalling playerInfo: %s", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error marshalling playerInfo: %s", err)))
		return
	}

	w.WriteHeader(200)
	w.Write(payload)
}

func (h *Handler) itemPurchaseGet(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	log.Printf("Grabbing replay item purchases info from store [%s]", host)

	if _, ok := h.conf.Stores[host]; !ok {
		log.Printf("Can't get replay info from store [%s], store does not exist", host)
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("Can't get replay info from store [%s], store does not exist", host)))
		return
	}

	filters := make(map[string]interface{})
	if filter := r.URL.Query().Get("gameId"); filter != "" {
		log.Printf("Found filter [gameId] in itemPurchaseGet request, parsing...")
		games := strings.Split(filter, ",")
		data := make([]uint64, len(games))
		for i, gameID := range games {
			s, err := strconv.ParseUint(gameID, 10, 64)
			if err != nil {
				log.Printf("Error parsing filter [gameId] in itemPurchaseGet request: %s", err)
				w.WriteHeader(500)
				w.Write([]byte(fmt.Sprintf("Error reading filter [gameId] in itemPurchaseGet request: %s", err)))
				return
			}
			data[i] = s
		}
		filters["gameId"] = data
	}

	if filter := r.URL.Query().Get("player"); filter != "" {
		log.Printf("Found filter [player] in itemPurchaseGet request, parsing...")
		players := strings.Split(filter, ",")
		data := make([]uint64, len(players))
		for i, steamID := range players {
			s, err := strconv.ParseUint(steamID, 10, 64)
			if err != nil {
				log.Printf("Error parsing filter [player] in itemPurchaseGet request: %s", err)
				w.WriteHeader(500)
				w.Write([]byte(fmt.Sprintf("Error reading filter [player] in itemPurchaseGet request: %s", err)))
				return
			}
			data[i] = s
		}
		filters["player"] = data
	}

	if filter := r.URL.Query().Get("hero"); filter != "" {
		log.Printf("Found filter [hero] in itemPurchaseGet request, parsing...")
		heroes := strings.Split(filter, ",")
		filters["hero"] = heroes
	}

	if filter := r.URL.Query().Get("item"); filter != "" {
		log.Printf("Found filter [item] in itemPurchaseGet request, parsing...")
		items := strings.Split(filter, ",")
		filters["item"] = items
	}

	log.Printf("Loading Item Purchases from store [%s] using filters [%+v]", host, filters)
	i, err := h.conf.Stores[host].LoadItemPurchase(filters)
	if err != nil {
		log.Printf("Can't grab item purchases from store [%s]: %s", host, err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Can't grab item purchases from store [%s]: %s", host, err)))
		return
	}

	payload, err := json.Marshal(i)
	if err != nil {
		log.Printf("Error marshalling item purchases as JSON: %s", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	w.Write(payload)
}

func (h *Handler) isAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.conf.Auth == "" {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.URL.Query().Get("auth")
		if auth != h.conf.Auth {
			log.Printf("Failed authentication for client [%s], incorrect auth key supplied: %s", r.Host, auth)
			w.WriteHeader(403)
			w.Write([]byte("403: Forbidden"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
