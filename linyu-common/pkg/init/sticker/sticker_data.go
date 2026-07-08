package sticker

type StickerData struct {
	ID            string
	Name          string
	IconUrl       string
	Type          string
	IconValue     string
	StickerPackID string
}

func All() []StickerData {
	var all []StickerData
	for _, items := range [][]StickerData{Default(), MiyouRabbit(), LinyuXiaoji()} {
		all = append(all, items...)
	}
	return all
}
