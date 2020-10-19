package API

import (
	"encoding/json"
	"goTestProj/DataModules"
	"log"
	"net/http"
)

// Skatoties uz boilerplate, ko rada go error handlings, programmētājs atcerās DRY likumu un uzraksta kaut ko šādu :)
// Es ieteiktu šādas `general error handling` funkcijas nerakstīt, un apstrādāt kļūdas inline.
//
// Iemesls ir tāds, ka katra kļūda ir individuāli jāapsver.
//
// Šī funkcija tiek lietota tikai vietās, kur programmas inicializācija ir nofeilojusi, un programmu ir jānogalina,
// jo tā nevar nonākt lietojamā stāvoklī.
//
// 90% gadījumu rakstīti funkcijas, kur kļūdas gadījumā tā ir jāatgriež kā viens no return argumentiem, piemēram,
//
// func Search(query string) (string, error) {
//     reqUrl, err := url.Parse("https://google.com")
//     if err != nil {
//         return "", err // return error that is handled somewhere up the stack trace
//     }
//     // logic that makes a search request to google...
// }
//
// Citreiz, kļūdas mums pat neko diži neietekmē, un funkcija var turpināt, neatgriežot kļūdu.
//
// func getImageWidth(r *http.Request, defaultWidth int) int {
//     widthStr := r.URL.Query().Get("width")
//     width, err := strconv.Atoi(widthStr)
//     if err != nil {
//         width = defaultWidth // no biggie, just use the default
//     }
//     return width
// }

func LogErrorIfNeeded(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func RespondWithError(w http.ResponseWriter, status int, error DataModules.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}