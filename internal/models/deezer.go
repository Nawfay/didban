package models


type Track struct {
	ID     int    `json:"id"`
	TrackToken string `json:"track_token"`
	Title  string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Duration int    `json:"duration"`
	Artist struct {
		Name string `json:"name"`
	} `json:"artist"`
	Album struct {
		Title string `json:"title"`
		Cover string `json:"cover_big"`
	} `json:"album"`
}


// MediaResponse models the v1/get_url response.
type MediaResponse struct {
	Data []struct {
		Media []struct {
			Sources []struct {
				URL      string `json:"url"`
				Provider string `json:"provider"`
			} `json:"sources"`
		} `json:"media"`
		Errors []struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
	} `json:"data"`
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

// getUserDataResponse matches Deezer’s getUserData JSON envelope.
type GetUserDataResponse struct {
	Results struct {
		User struct {
			Options struct {
				LicenseToken string `json:"license_token"`
			} `json:"options"`
		} `json:"user"`
	} `json:"results"`
}
