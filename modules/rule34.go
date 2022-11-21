package modules

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/carabelle/alexisbot/commands"
	"github.com/carabelle/alexisbot/utils"
	"github.com/valyala/fasthttp"
	"lukechampine.com/frand"
)

const (
	URL = "https://api.rule34.xxx/index.php?page=dapi&json=1&s=post&q=index&limit=30&tags="
)

var (
	postCache  = make(map[string][]*RPost, 0)
	emptyPost  = &RPost{}
	emptyPosts = make([]*RPost, 0)
)

type RPost struct {
	Score       int                      `json:"score"`
	FileURL     string                   `json:"file_url"`
	Rating      string                   `json:"rating"`
	Color       utils.ImageDominateColor `json:"color"`
	Preview_URL string                   `json:"preview_url"`
}

func decodeTOJSON(tags string) ([]*RPost, error) {
	_, body, err := fasthttp.Get(nil, URL+tags)
	utils.CheckIfError(err)
	json.Unmarshal(body, &emptyPosts)
	return emptyPosts, nil
}

func cacheGetOrElse(tags string) ([]*RPost, error) {
	if get, ok := postCache[tags]; ok {
		return get, nil
	}
	p, err := decodeTOJSON(tags)
	utils.CheckIfError(err)
	postCache[tags] = p
	filteredPosts(tags)
	return postCache[tags], nil
}

func filteredPosts(t string) {
	var wg sync.WaitGroup
	wg.Add(len(postCache[t]))
	for _, post := range postCache[t] {
		go func(post *RPost) {
			defer wg.Done()
			if utils.IsImage(post.FileURL) {
				emptyPosts = append(emptyPosts, post)
				post.Color = utils.GetColorInfo(post.Preview_URL)
			}
		}(post)
	}
	wg.Wait()
	postCache[t] = emptyPosts
	emptyPosts = make([]*RPost, 0)
}

func init() {
	commands.NewCommand("rule34", "Search on rule34.xxx").
		AddOption("tags", "tags to search for", 3, true).
		SetHandler(func(c *commands.Command) {
			arg, _ := c.GetOption("tags")
			string := arg.StringValue()
			string = strings.Replace(string, " ", "_", -1)
			post, err := randomResponse(string)
			c.CheckIfError(err)
			gId, err := c.State.GuildChannel(c.GuildID, c.ChannelID)
			c.CheckIfError(err)
			if !gId.NSFW {
				c.SendEphemeralMessageEmbed(post.GetEmbedModel(arg.StringValue()))
			}
			c.SendInteractionMessageEmbed(post.GetEmbedModel(arg.StringValue()))

		})
}

func randomResponse(tags string) (*RPost, error) {
	postsCache, err := cacheGetOrElse(tags)
	if err != nil || len(postsCache) == 0 {
		return emptyPost, err
	}
	frand.NewSource().Seed(int64(time.Now().Minute()))
	return postsCache[frand.Intn(len(postsCache))], nil
}

func (e *RPost) GetEmbedModel(t string) (u *utils.Embed) {
	t = strings.Replace(t, "_", " ", -1)
	u = utils.NewEmbed().
		SetTitle(strings.Title(t)).
		AddField("Rating", strings.Title(e.Rating)).
		AddField("Score", fmt.Sprintf("%d", e.Score)).
		InlineAllFields().
		SetColor(e.Color.RGB).
		SetImage(e.FileURL).
		SetURL(e.FileURL)
	return
}
