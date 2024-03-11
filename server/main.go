package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/config"
	"server/internal/infrastructure/migrator"
	"server/internal/modules/auth/models"
	"server/internal/modules/chat/contorller"
	"server/internal/modules/chat/service"
	"server/internal/router"
	"time"
)

func main() {
	time.Sleep(time.Second * 10)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Can't load config from .env")
	}
	//create Config with sens info
	cfg := config.NewConfig()
	//load Config
	cfg.Load()
	//
	//
	//
	//
	//
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.Name,
	)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalln("Can't connect to postgres using sql std driver", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln("Can't ping postgres using sql std driver", err)
	}
	log.Println("Pong from ")
	dbx := sqlx.NewDb(db, "postgres")

	//
	//
	//
	migrator := migrator.NewMigrator(dbx)
	user := &models.User{}
	err = migrator.Migrate(user)
	if err != nil {
		log.Fatalln("Can't create table", err)
	}
	//
	//
	//
	//
	//
	//
	//
	//

	//create new hub who menagge all recived message and created rooms
	hub := service.NewHub()
	go hub.Run()
	// create new contoroller and pass hub into
	ctrl := controller.NewChatController(hub, dbx, []byte(cfg.Key))

	//creage default gin.Engine instance where write all routes and paths
	gin := gin.Default()
	r := router.NewRouter(ctrl)
	r.Route(gin)
	//run listening and serv requests
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Listen, cfg.Port),
		Handler: gin,
	}
	go func() {
		log.Println("Server is startgin on", cfg.Listen, cfg.Port)
		err = srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Can't start listening on %s:%s", cfg.Listen, ":", cfg.Port)
		}

	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Fatalln("Server is shutdown")
	}
	select {
	case <-ctx.Done():
		log.Println("Time out 5 second exists")
	}
	log.Println("Server stoped")
}
