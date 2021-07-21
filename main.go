package main

import (
	"blockchain/data"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var Blockchain []data.Block

//use calculateHash to generate a new block
func newBlock(oldBlock data.Block, BPM int) (data.Block, error) {
	var newBlock data.Block

	t := time.Now()

	newBlock.PrevHash = oldBlock.CalculateHash()
	newBlock.Index = 1 + oldBlock.Index
	newBlock.BPM = BPM
	newBlock.Timestamp = t.String()
	newBlock.Hash = newBlock.CalculateHash()

	return newBlock, nil
}

//picking the best chain in the event of disagreeing nodes
func replaceChain(newBlocks []data.Block) error {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
	return nil
}

func runServer() error {
	router := makeMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", os.Getenv("PORT"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handlePostBlockchain).Methods("POST")
	return muxRouter
}

//wrapper function for error handling
func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

//add a new block to server
func handlePostBlockchain(w http.ResponseWriter, r *http.Request) {
	//take in request
	var m data.Message

	//decode request data
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	//make new block, check error, and check new block validity
	newBlock, err := newBlock(Blockchain[len(Blockchain)-1], m.BPM)
	if err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	if newBlock.IsBlockValid(Blockchain[len(Blockchain)-1]) {
		newBlockChain := append(Blockchain, newBlock)
		replaceChain(newBlockChain)
		spew.Dump(Blockchain)
	}

	//respond to requester
	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

//ask for existing blockchain from server
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func main() {
	err := godotenv.Load() //reads items from the .env file so no hard coding is needed
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := data.Block{0, t.String(), 0, " ", " "}
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
	}()

	log.Fatal(runServer())
}
