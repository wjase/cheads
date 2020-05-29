package main

import (
	"encoding/json"
	"fmt"
	"time"

	//    "io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/wjase/cheads/pkg/auth"
	"github.com/wjase/cheads/pkg/cheads"
	"github.com/wjase/cheads/pkg/game"
	"github.com/wjase/cheads/pkg/remote"
	"github.com/wjase/cheads/pkg/spa"
	"github.com/wjase/cheads/pkg/websocket"
)

func main() {
	fmt.Println("Celebrity Heads")
	router := mux.NewRouter()

	poolSet := websocket.NePoolSet()
	poolSet.PoolAdded = func(pool remote.ClientPool) {
		engine := cheads.Game{}
		hub := game.NewGameHub(pool)
		hub.AddGame(&engine)
	}

	// usr login
	router.HandleFunc("/signin", auth.Signin)

	// conenct websocket
	router.HandleFunc("/ws", auth.Authorised(poolSet.ServeWsCreateClient()))

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	spa := spa.Handler{StaticPath: "frontend/build", IndexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
