package utils

import "time"

type Movie struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Plot      string    `bson:"plot,omitempty" json:"plot,omitempty"`
	Genres    []string  `bson:"genres,omitempty" json:"genres,omitempty"`
	Runtime   int       `bson:"runtime,omitempty" json:"runtime,omitempty"`
	Cast      []string  `bson:"cast,omitempty" json:"cast,omitempty"`
	Poster    string    `bson:"poster,omitempty" json:"poster,omitempty"`
	Title     string    `bson:"title,omitempty" json:"title,omitempty"`
	FullPlot  string    `bson:"fullplot,omitempty" json:"fullplot,omitempty"`
	Languages []string  `bson:"languages,omitempty" json:"languages,omitempty"`
	Released  time.Time `bson:"released,omitempty" json:"released,omitempty"` // You can also use time.Time if you want to parse dates
	Directors []string  `bson:"directors,omitempty" json:"directors,omitempty"`
	Rated     string    `bson:"rated,omitempty" json:"rated,omitempty"`
	Awards    struct {
		Wins int `json:"wins" bson:"wins"` // Use int for awards.wins
	} `json:"awards" bson:"awards"`
	LastUpdated      string                 `bson:"lastupdated,omitempty" json:"lastupdated,omitempty"`
	Year             interface{}            `bson:"year,omitempty" json:"year,omitempty"`
	IMDB             map[string]interface{} `bson:"imdb,omitempty" json:"imdb,omitempty"`
	Countries        []string               `bson:"countries,omitempty" json:"countries,omitempty"`
	Type             string                 `bson:"type,omitempty" json:"type,omitempty"`
	Tomatoes         map[string]interface{} `bson:"tomatoes,omitempty" json:"tomatoes,omitempty"`
	NumMflixComments int                    `bson:"num_mflix_comments,omitempty" json:"num_mflix_comments,omitempty"`
}

type User struct {
	FirstName string `json:"firstname" bson:"firstname"`
	LastName  string `json:"lastname" bson:"lastname"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	IsAdmin   bool   `json:"isAdmin" bson:"isAdmin"`
}
