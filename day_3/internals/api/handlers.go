package api

import (
	"elastic_study/internals/middleware"
	"elastic_study/internals/view"
	"errors"
	"github.com/dgrijalva/jwt-go/v4"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Place struct {
	ID       int      `json:"id" csv:"id"`
	Name     string   `json:"name" csv:"Name"`
	Address  string   `json:"address" csv:"Address"`
	Phone    string   `json:"phone" csv:"Phone"`
	Location Location `json:"location"`
}

type Location struct {
	Latitude  float64 `json:"lat" csv:"Latitude"`
	Longitude float64 `json:"lon" csv:"Longitude"`
}

type Store interface {
	GetPlaces(limit int, offset int) ([]Place, int, error)
	UploadData(dataFile string, settingsFile string) error
	GetRecommendations(lat float64, lon float64, size int) ([]Place, error)
}

type StorageResponse struct {
	Name        string  `json:"name,omitempty"`
	Total       int     `json:"total,omitempty"`
	Places      []Place `json:"places,omitempty"`
	PrevPage    int     `json:"prev_page,omitempty"`
	NextPage    int     `json:"next_page,omitempty"`
	NumLastPage int     `json:"last_page,omitempty"`
	Error       string  `json:"error,omitempty"`
}

type DBFunc func(r *http.Request, s *Service) (*StorageResponse, error)

type ViewFunc func(w http.ResponseWriter, v interface{}, t *template.Template)

type Service struct {
	PageLimit      int
	NumRecommend   int
	expireDuration time.Duration
	signingKey     []byte
	DB             Store
	mux            *http.ServeMux
	tmpl           *template.Template
}

func NewService(pageLimit int, numRec int, tokenDuration time.Duration, signingKey string, db Store) *Service {
	// Template initialization
	allFiles := []string{"content.tmpl", "footer.tmpl", "header.tmpl", "page.tmpl"}
	var allPaths []string
	for _, tmpl := range allFiles {
		allPaths = append(allPaths, "./templates/"+tmpl)
	}
	tmplFunc := func(a, b int) int {
		return a + b
	}
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{"add": tmplFunc}).ParseFiles(allPaths...))

	// Handlers initialization
	mux := http.NewServeMux()
	service := Service{
		PageLimit:      pageLimit,
		NumRecommend:   numRec,
		expireDuration: tokenDuration,
		signingKey:     []byte(signingKey),
		DB:             db,
		mux:            mux,
		tmpl:           tmpl,
	}

	authMiddleware := middleware.NewCheckAuthHandler(signingKey)
	recommendHandler := http.HandlerFunc(service.getJsonRecommend)
	htmlPlacesHandler := http.HandlerFunc(service.getHtmlPlaces)
	jsonPlacesHandler := http.HandlerFunc(service.getJsonPlaces)
	tokenHandler := http.HandlerFunc(service.getToken)

	service.mux.Handle("/places/", middleware.Logging(htmlPlacesHandler))
	service.mux.Handle("/api/places/", middleware.Logging(jsonPlacesHandler))
	service.mux.Handle("/api/recommend/", middleware.Logging(authMiddleware(recommendHandler)))
	service.mux.Handle("/api/get_token/", middleware.Logging(tokenHandler))
	return &service
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Service) getHtmlPlaces(w http.ResponseWriter, r *http.Request) {
	s.process(w, r, getPage, view.RenderHtml)
}

func (s *Service) getJsonPlaces(w http.ResponseWriter, r *http.Request) {
	s.process(w, r, getPage, view.RenderJSON)
}

func (s *Service) getJsonRecommend(w http.ResponseWriter, r *http.Request) {
	s.process(w, r, getRecommendations, view.RenderJSON)
}

func (s *Service) process(w http.ResponseWriter, r *http.Request, dbF DBFunc, viewF ViewFunc) {
	result, err := dbF(r, s)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		viewF(w, &StorageResponse{Error: err.Error()}, s.tmpl)
	} else {
		viewF(w, result, s.tmpl)
	}
}

func getRecommendations(r *http.Request, s *Service) (*StorageResponse, error) {
	lat, lon, err := getLatLonQuery(r)
	if err != nil {
		return &StorageResponse{}, err
	}
	recommendations, err := s.DB.GetRecommendations(lat, lon, s.NumRecommend)
	if err != nil {
		return &StorageResponse{}, err
	}
	return &StorageResponse{
		Name:   "places",
		Places: recommendations,
	}, nil
}

func getLatLonQuery(r *http.Request) (lat float64, lon float64, err error) {
	keys, ok := r.URL.Query()["lat"]
	if ok {
		lat, err = strconv.ParseFloat(keys[0], 64)
		if err != nil {
			return 0, 0, errors.New("latitude does not have a valid value")
		}
	} else {
		return 0, 0, errors.New("query does not contain a latitude")
	}
	keys, ok = r.URL.Query()["lon"]
	if ok {
		lat, err = strconv.ParseFloat(keys[0], 64)
		if err != nil {
			return 0, 0, errors.New("longitude does not have a valid value")
		}
	} else {
		return 0, 0, errors.New("query does not contain a longitude")
	}
	return
}

func getPage(r *http.Request, s *Service) (*StorageResponse, error) {
	keys, ok := r.URL.Query()["page"]
	pageStr := "1"
	if ok {
		pageStr = keys[0]
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return &StorageResponse{}, errors.New("Invalid 'page' value: " + pageStr)
	}

	places, total, err := s.DB.GetPlaces(s.PageLimit, page-1)
	if err != nil {
		return &StorageResponse{}, err
	}
	return &StorageResponse{
		Name:        "places",
		Total:       total,
		Places:      places,
		PrevPage:    page - 1,
		NextPage:    page + 1,
		NumLastPage: int(math.Ceil(float64(total) / float64(s.PageLimit))),
	}, nil
}

func (s *Service) getToken(w http.ResponseWriter, r *http.Request) {
	type JWT struct {
		Token string `json:"token"`
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: jwt.At(time.Now().Add(s.expireDuration)),
		IssuedAt:  jwt.At(time.Now()),
	})
	signedToken, _ := token.SignedString(s.signingKey)
	view.RenderJSON(w, JWT{signedToken}, nil)
}
