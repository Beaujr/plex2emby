package process

import (
	"github.com/beaujr/plex2emby/plex"
	"github.com/beaujr/plex2emby/emby"
	"strings"
	"fmt"
	"strconv"
	"time"
)

const embyDateFormat = "20060102150405"
const seriesType = "Series"
const movieType = "Movie"


type Plex2Emby interface {
	Process() error
	processFilms(sectionKey string) error
	processTVShows(sectionKey string) error
}

type Clients struct {
	plex plex.Client
	emby emby.Client
}

func NewPlex2EmbyClient(plexClient plex.Client, embyClient emby.Client) Clients {
	return Clients{plexClient, embyClient}
}

func (c *Clients) Process() error {
	sections, err := c.plex.GetSections()
	if err != nil {
		return err
	}

	for _, section := range sections {
		if strings.Compare(section.Type, "movie") == 0 {
			fmt.Sprintf("PROCESSING SECTION %s", section.Type)
			err = c.processFilms(section.Key)
			if err != nil {
				return err
			}
		}
		if strings.Compare(section.Type, "show") == 0 {
			fmt.Sprintf("PROCESSING SECTION %s", section.Type)
			err = c.processTVShows(section.Key)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Clients) processFilms (sectionKey string) error {
	items, err := c.plex.GetFilmSection(sectionKey)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.LastViewedAt != "" {
			movies, err := c.emby.Search(item.Title, movieType)
			if err != nil {
				return err
			}

			// We've got a match! Always take the first
			// TODO: do some string comparisons
			if len(movies) > 0 {
				movie := movies[0]
				i, err := strconv.ParseInt(item.LastViewedAt, 10, 64)
				if err != nil {
					return err
				}
				tm := time.Unix(i, 0)
				fmt.Println(fmt.Sprintf("%s,%s,watched", item.Title, movie.Name))
				return c.emby.MarkItemAsPlayed(movie.Id, tm.Format(embyDateFormat))

			}
		}
	}
	return nil
}

func(c *Clients) processTVShows (sectionKey string) error {
	items, err := c.plex.GetTVSection(sectionKey)
	if err != nil {
		return err
	}
	for _, item := range items {
		seasons, err :=  c.plex.GetShow(item.Key)
		if err != nil {
			return err
		}

		results, err := c.emby.Search(item.Title, seriesType)
		if err != nil {
			return err
		}
		// We've got a match! Always take the first
		// TODO: do some string comparisons
		if len(results) > 0 {
			embyItem := results[0]
			embyEpisodes, err := c.emby.GetItem(embyItem.Id)
			if err != nil {
				return err
			}
			for _, season := range seasons {
				if strings.Compare(season.Title, "All episodes") != 0 {
					plexEpisodes, err :=  c.plex.GetSeason(season.Key)
					if err != nil {
						return err
					}

					for _, plexEpisode := range plexEpisodes {
						if plexEpisode.ViewCount == "" {
							fmt.Println(fmt.Sprintf("%s,%s,%s,%s,%s", sectionKey, item.Title, season.Title, plexEpisode.EpisodeNumber, "Unwatched"))
						} else {
							fmt.Println(fmt.Sprintf("%s,%s,%s,%s,%s", sectionKey, item.Title, season.Title, plexEpisode.EpisodeNumber, "Watched"))
							for _, embyEpisode := range embyEpisodes {
								plexEpisodeNumber, err := strconv.Atoi(plexEpisode.EpisodeNumber)
								if err != nil {
									return err
								}

								if embyEpisode.IndexNumber == plexEpisodeNumber && strings.Compare(season.Title, embyEpisode.SeasonName) == 0 && embyEpisode.UserData.PlayCount == 0 {
									i, err := strconv.ParseInt(plexEpisode.LastViewedAt, 10, 64)
									if err != nil {
										return err
									}
									tm := time.Unix(i, 0)
									return c.emby.MarkItemAsPlayed(embyEpisode.Id, tm.Format(embyDateFormat))
								}
							}
						}
					}
				}
			}
		} else {
			return fmt.Errorf("%s,%s", item.Title, "NOT FOUND")
		}
	}
	return nil
}