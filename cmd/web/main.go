package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"sketchdb.cozycole.net/internal/img"
	"sketchdb.cozycole.net/internal/models"

	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
	fileStorage   img.FileStorageInterface
	videos        models.VideoModelInterface
	creators      models.CreatorModelInterface
	people        models.PersonModelInterface
	characters    models.CharacterModelInterface
	search        models.SearchModelInterface
	debugMode     bool
	formDecoder   *form.Decoder
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	debug := flag.Bool("debug", false, "debug mode")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbpool, err := openDB(os.Getenv("DB_URL"))
	if err != nil {
		errorLog.Fatal(err)
	}
	defer dbpool.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	storagePath := os.Getenv("IMG_STORAGE_PATH")
	fileStorage := img.FileStorage{Path: storagePath}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
		formDecoder:   formDecoder,
		fileStorage:   &fileStorage,
		videos:        &models.VideoModel{DB: dbpool, ResultSize: 16},
		creators:      &models.CreatorModel{DB: dbpool},
		people:        &models.PersonModel{DB: dbpool},
		characters:    &models.CharacterModel{DB: dbpool},
		search:        &models.SearchModel{DB: dbpool},
		debugMode:     *debug,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Println("Starting server on", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err = dbpool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return dbpool, nil
}
