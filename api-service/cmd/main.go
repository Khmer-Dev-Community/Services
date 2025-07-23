package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Khmer-Dev-Community/Services/api-service/utils"
)

func main() {

	cfg, db := InitConfigAndDatabase()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		utils.WarnLog("Received signal: %s", sig.String())
		cancel()
	}()
	services := InitServices(db, &cfg)
	router := InitRoutes(cfg, services)
	go func() {

		log.Println("      \n                   *       \n.         *       *A*       *\n.        *A*     **=**     *A*\n        *\"\"\"*   *|\"\"\"|*   *\"\"\"*\n       *|***|*  *|*+*|*  *|***|*\n*********\"\"\"*___*//+\\\\*___*\"\"\"*********\n@@@@@@@@@@@@@@@@//   \\\\@@@@@@@@@@@@@@@@@\n###############||ព្រះពុទ្ធ||#################\nTTTTTTTTTTTTTTT||ព្រះធម័||TTTTTTTTTTTTTTTTT\nLLLLLLLLLLLLLL//ព្រះសង្ឃ\\\\LLLLLLLLLLLLLLLLL\n៚ សូមប្រោសប្រទានពរឱ្យប្រតិប័ត្តិការណ៍ប្រព្រឹត្តទៅជាធម្មតា ៚ \n៚ ជោគជ័យ   //  ៚សិរីសួរស្តី \\\\   ៚សុវត្តិភាព \n___________//___៚(♨️)៚__\\\\____________\n៚Application Service is Running Port: 80 ")
		if err := StartHTTPServer(cfg.Service.Port, router.(http.Handler)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
		log.Printf("HTTP server successfully started on port %s", cfg.Service.Port)
	}()

	<-ctx.Done()
	utils.InfoLog("Shutting down service...", "")
}
