package config

type Config struct {
    TogglConfig TogglConfig `desc:"Toggl config"`
    FirestoreConfig FirestoreConfig `desc:"Firestore config"`
}
type ConfigFile string
type TogglConfig struct {
    APIKey string `desc:"Toggl API key"`
}
type FirestoreConfig struct {
    CollectionName string `desc:"Firestore collection name"`
}
