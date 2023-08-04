package db

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	m "main/models"
	"net/http"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	Db *sql.DB

	Dsn           string
	currentApiKey string

	apiKeys []string

	countForApiKeySwitch int = 0
	requestCount         int = 0

	VolumesResponse   m.VolumesResponse
	VolumeResponse    m.VolumeResponse
	IssueResponse     m.IssuesResponse
	CharacterResponse m.CharactersResponse
)

func InitDb() {
	var envs map[string]string

	envs, err := godotenv.Read(`D:\Dolgaim\Programozás\Golang Learning\Comicbook-Webshop\Backend\.env`)
	if err != nil {
		log.Fatal(err)
	}
	apiKeys = strings.Split(envs["API_KEYS"], ",")

	Dsn = envs["MySQLUsername"] + ":" + envs["MySQLPassword"] + "@tcp(" + envs["MyAddress"] + ":" + envs["MyPort"] + ")/comicbooks"

	Db, err = sql.Open("mysql", Dsn)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = Db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to MySQL")
	LoadComicbooksFromApi2()
	return
}

// func LoadComicbooksFromApi() {
// 	var volumesCount int
// 	row := Db.QueryRow("SELECT COUNT(*) FROM volumes")
// 	err := row.Scan(&volumesCount)
// 	if err != nil {
// 		fmt.Println("Error loading Comicbooks")
// 		return
// 	}
// 	if volumesCount != 0 {
// 		fmt.Println("Comicbooks up to date!")
// 		return
// 	}
// 	fmt.Println("Loading commicbooks...")
// 	fmt.Println("Commicbooks loaded... Saving to database...")
// 	var body []byte
// 	var isOnDelay bool
// 	body, err = GetBodyFromURL("https://comicvine.gamespot.com/api/volumes/")
// 	if err != nil {
// 		fmt.Println("Error at: https://comicvine.gamespot.com/api/volumes/ : ", err.Error())
// 	}
// 	err = HandleBody(body, "vs")
// 	if err != nil {
// 		fmt.Println("Error at handlebody ", err.Error())
// 	}
// 	for i, vUrl := range VolumesResponse.VolumesResults.VolumesUrl {
// 		for {
// 			body, err = GetBodyFromURL(vUrl.VolumeApiDetailUrl)
// 			if err != nil {
// 				fmt.Println("Error at: ", vUrl.VolumeApiDetailUrl, " : ", err.Error())
// 			}
// 			isOnDelay, err = IsOnDelay(body)
// 			if err != nil {
// 				fmt.Println("Error at: ", vUrl.VolumeApiDetailUrl, " : ", err.Error())
// 			}
// 			if isOnDelay {
// 				fmt.Println("DELAY! RequestCount: ", requestCount)
// 				time.Sleep(30 * time.Second)
// 			} else {
// 				err = HandleBody(body, "v")
// 				if err != nil {
// 					fmt.Println(err.Error())
// 				}
// 				InsertVolumeIntoDb(VolumeResponse.VolumeResults[i])
// 				for j, iUrl := range VolumeResponse.VolumeResults[i].IssuesApiDetails {
// 					for {
// 						body, err = GetBodyFromURL(iUrl.IssueApiDetailUrl)
// 						if err != nil {
// 							fmt.Println("Error at: ", iUrl.IssueApiDetailUrl, " : ", err.Error())
// 						}
// 						isOnDelay, err = IsOnDelay(body)
// 						if err != nil {
// 							fmt.Println("Error at: ", iUrl.IssueApiDetailUrl, " : ", err.Error())
// 						}
// 						if isOnDelay {
// 							fmt.Println("DELAY! RequestCount: ", requestCount)
// 							break
// 						} else {
// 							err = HandleBody(body, "i")
// 							if err != nil {
// 								fmt.Println(err.Error())
// 							}
// 							InsertIssueIntoDb(IssueResponse.Issues[j], VolumeResponse.VolumeResults[i].ID)
// 						}
// 					}
// 					for k, cUrl := range VolumeResponse.VolumeResults[i].CharactersApiDetails {
// 						body, err = GetBodyFromURL(cUrl.CharacterApiDetailUrl)
// 						if err != nil {
// 							fmt.Println("Error at: ", cUrl.CharacterApiDetailUrl, " : ", err.Error())
// 						}
// 						isOnDelay, err = IsOnDelay(body)
// 						if err != nil {
// 							fmt.Println("Error at: ", cUrl.CharacterApiDetailUrl, " : ", err.Error())
// 						}
// 						if isOnDelay {
// 							fmt.Println("DELAY! RequestCount: ", requestCount)
// 							time.Sleep(30 * time.Second)
// 						} else {
// 							err = HandleBody(body, "c")
// 							if err != nil {
// 								fmt.Println(err.Error())
// 							}
// 							InsertCharacterIntoDb(CharacterResponse.Characters[k])
// 						}
// 					}
// 				}
// 				break
// 			}
// 		}
// 		fmt.Printf("%d/%d volume saved to database!", i+1, len(VolumesResponse.VolumesResults.VolumesUrl))
// 		fmt.Println()
// 	}
// 	fmt.Println("Comicbooks loaded from API!")

//	}

func LoadComicbooksFromApi2() {
	var volumesCount int

	row := Db.QueryRow("SELECT COUNT(*) FROM `volumes`")
	err := row.Scan(&volumesCount)
	if err != nil {
		fmt.Println("Error loading Comicbooks")
		return
	}
	if volumesCount != 0 {
		fmt.Println("Comicbooks up to date!")
		return
	}
	fmt.Println("Commicbooks loaded... Saving to database...")

	for avoidDuplicates("https://comicvine.gamespot.com/api/volumes/", "vs") {
		fmt.Println("We are on DELAY!")
		time.Sleep(30 * time.Second)
	}
	for i, vUrl := range VolumesResponse.VolumesResults.VolumesUrl {
		for avoidDuplicates(vUrl.VolumeApiDetailUrl, "v") {
			fmt.Println("We are on DELAY! RequestCount: ", requestCount)
			time.Sleep(30 * time.Second)
		}
		InsertVolumeIntoDb(VolumeResponse.VolumeResults[i])
		for j, iUrl := range VolumeResponse.VolumeResults[i].IssuesApiDetails {
			for avoidDuplicates(iUrl.IssueApiDetailUrl, "i") {
				fmt.Println("We are on DELAY! RequestCount: ", requestCount)
				time.Sleep(30 * time.Second)
			}
			InsertIssueIntoDb(IssueResponse.Issues[j], VolumeResponse.VolumeResults[i].ID)
		}
		for k, cUrl := range VolumeResponse.VolumeResults[i].CharactersApiDetails {
			for avoidDuplicates(cUrl.CharacterApiDetailUrl, "c") {
				fmt.Println("We are on DELAY! RequestCount: ", requestCount)
				time.Sleep(30 * time.Second)
			}
			InsertCharacterIntoDb(CharacterResponse.Characters[k])
		}
		fmt.Printf("%d/%d volume saved to database!", i+1, len(VolumesResponse.VolumesResults.VolumesUrl))
		fmt.Println()
	}
	fmt.Println("Comicbooks loaded from API!")
}

func avoidDuplicates(url string, pos string) bool {
	body, err := GetBodyFromURL(url)
	if err != nil {
		fmt.Println("Error at getBody, url: ", url)
		return false
	}
	if IsOnDelay(body) {
		return true
	} else {
		err = HandleBody(body, pos)
		if err != nil {
			fmt.Println("Error at avoidDuplicates>handleBody, error: ", err.Error())
		}
		return false
	}
}

func HandleBody(body []byte, s string) error {
	var err error
	switch s {
	case "vs":
		err = xml.Unmarshal(body, &VolumesResponse)
		break
	case "v":
		err = xml.Unmarshal(body, &VolumeResponse)
		break
	case "i":
		err = xml.Unmarshal(body, &IssueResponse)
		break
	case "c":
		err = xml.Unmarshal(body, &CharacterResponse)
		break
	default:
		err = errors.New("Own misstake at getDataFromUrl()")
	}
	countForApiKeySwitch++
	if countForApiKeySwitch == 100 {
		countForApiKeySwitch = 1
	}
	if err != nil {
		return err
	}
	return nil
}

func IsOnDelay(body []byte) bool {
	var apiError m.ApiError
	err := xml.Unmarshal(body, &apiError)
	if err != nil {
		fmt.Println("Error at unmarshaling isOnDelay, error: ", err)
	}
	if apiError.Error != "OK" {
		return true
	}
	return false
}

func GetBodyFromURL(url string) ([]byte, error) {
	requestCount++
	resp, err := http.Get(url + "?api_key=" + apiKeys[countForApiKeySwitch%2])
	if err != nil {
		fmt.Println("Failed to get data from api")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read body" + err.Error())
		return nil, err
	}
	return body, nil
}

func InsertCharacterIntoDb(c m.Character) {

	if isExistsById(c.ID, "c") {
		return
	}

	stmt, err := Db.Prepare("INSERT INTO `characters` (`id`, `name`, `img`) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID, c.Name, c.Image.OriginalURL)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return
		}
		fmt.Println(err.Error())
	}
}

func InsertIssueIntoDb(i m.Issue, vID int) {
	if isExistsById(i.ID, "i") {
		return
	}

	stmt, err := Db.Prepare("INSERT INTO `issues` (`id`, `volume_id`, `name`, `issue_number`, `img`, `cover_date`, `date_added`) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(i.ID, vID, i.Name, i.IssueNumber, i.Image.OriginalURL, i.CoverDate, i.Date_Added)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func InsertVolumeIntoDb(v m.Volume) {

	if isExistsById(v.ID, "v") {
		return
	}

	stmt, err := Db.Prepare("INSERT INTO `volumes` (`id`, `name`, `img`, `desc`, `publisher`) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(v.ID, v.Name, v.Image.OriginalURL, v.Description, v.Publisher.Name)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func isExistsById(ID int, s string) bool {
	var exists bool
	var table string

	switch s {
	case "c":
		table = `characters`
		break
	case "i":
		table = `issues`
		break
	case "v":
		table = `volumes`
		break
	default:
		table = ``
		return true
	}

	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE `id` = ?)", table)
	err := Db.QueryRow(query, ID).Scan(&exists)
	if err != nil {
		fmt.Println(err.Error())
		return true
	}

	return exists
}
