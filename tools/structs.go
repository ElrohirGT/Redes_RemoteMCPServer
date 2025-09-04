package tools

type MTGCard struct {
	Id       string   `json:"multiverseid"`
	Name     string   `json:"name"`
	ManaCost string   `json:"manaCost"`
	Colors   []string `json:"colors"`
	Type     string   `json:"type"`
	Rarity   string   `json:"rarity"`
	Text     string   `json:"text"`
	ImageUrl string   `json:"imageUrl"`
}
