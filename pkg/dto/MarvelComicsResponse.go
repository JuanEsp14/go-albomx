package dto

type MarvelComicsResponse struct {
	Code            int    `json:"code"`
	Status          string `json:"status"`
	Copyright       string `json:"copyright"`
	AttributionText string `json:"attributionText"`
	AttributionHTML string `json:"attributionHTML"`
	Etag            string `json:"etag"`
	Data            Data   `json:"data"`
}//@Name MarvelComicsResponse

type Data struct {
	Offset  int `json:"offset"`
	Limit   int `json:"limit"`
	Total   int `json:"total"`
	Count   int `json:"count"`
	Results []Result  `json:"results"`
}//@Name Data

type Result struct {
	ID                 int           `json:"id"`
	DigitalID          int           `json:"digitalId"`
	Title              string        `json:"title"`
	IssueNumber        int           `json:"issueNumber"`
	VariantDescription string        `json:"variantDescription"`
	Description        interface{}   `json:"description"`
	Modified           string        `json:"modified"`
	Isbn               string        `json:"isbn"`
	Upc                string        `json:"upc"`
	DiamondCode        string        `json:"diamondCode"`
	Ean                string        `json:"ean"`
	Issn               string        `json:"issn"`
	Format             string        `json:"format"`
	PageCount          int           `json:"pageCount"`
	TextObjects        []interface{} `json:"textObjects"`
	ResourceURI        string        `json:"resourceURI"`
	Urls               []Url  		 `json:"urls"`
	Series 			   Variant       `json:"series"`
	Variants 		   []Variant	 `json:"variants"`
	Dates              []Date 		 `json:"dates"`
	Prices 			   []Price  	 `json:"prices"`
	Thumbnail 		   Image 		 `json:"thumbnail"`
	Images 			   []Image  	 `json:"images"`
	Creators 		   Collaborators `json:"creators"`
	Characters 		   Collaborators `json:"characters"`
	Stories 		   Collaborators `json:"stories"`
	Events 			   Collaborators `json:"events"`
}//@Name Result

type Url struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}//@Name Url

type Variant struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`
}//@Name Variant

type Date struct {
	Type string `json:"type"`
	Date string `json:"date"`
}//@Name Date

type Price struct {
	Type  string  `json:"type"`
	Price float64 `json:"price"`
}//@Name Price

type Image struct {
	Path      string `json:"path"`
	Extension string `json:"extension"`
}//@Name Image

type Collaborators struct {
	Available     int    `json:"available"`
	CollectionURI string `json:"collectionURI"`
	Items         []Item `json:"items"`
	Returned 	  int 	 `json:"returned"`
}//@Name Collaborators

type Item struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`
	Role        string `json:"role"`
}//@Name Item