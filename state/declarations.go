// Объявления типов и переменных из пакета state
package state

// StateEntry представляет запись в реестре переменных состояния
type StateEntry struct {
	Login        string `json:"login"`
	Password     string `json:"password"`
	Url          string `json:"url"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiredAt    int64  `json:"expiredAt"`
}

type stateRegistry map[string]*StateEntry

// registry содержит переменные состояния приложения.
// Ключ: строковая метка, значение: запись в реестре.
var registry = make(stateRegistry)

// stateFileName определяет имя файла, в который сериализуется реестр состояния.
var stateFileName = "state.json"
