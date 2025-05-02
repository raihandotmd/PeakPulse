package db

import (
	"os"

	"github.com/raihandotmd/peakPulse/internal/config"
	supa "github.com/supabase-community/supabase-go"
)

var Supa *supa.Client

func InitSupabase() *supa.Client {
	// Load environment variables
	config.Load()
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")
	var err error
	Supa, err = supa.NewClient(url, key, &supa.ClientOptions{}) // sets up PostgREST, Realtime, Auth, Storage
	if err != nil {
		panic("Failed to initialize Supabase client: " + err.Error())
	}

	return Supa
}

// GetSupabaseClient returns the Supabase client
func GetSupabaseClient() *supa.Client {
	if Supa == nil {
		panic("Supabase client is not initialized")
	}
	return Supa
}
