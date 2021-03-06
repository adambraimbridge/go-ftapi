package ftapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Recommendation struct {
	ID            string  `json:"id"`
	Popularity    float64 `json:"popularity"`
	PublishedDate string  `json:"published"`
	Score         float64 `json:"score"`
	Title         string  `json:"title"`
	URL           string  `json:"url"`
}

type Recommendations struct {
	RawJSON  *[]byte
	Articles []Recommendation `json:"articles"`
	Status   string           `json:"status"`
	Type     string           `json:"type"`
	Version  string           `json:"version"`
}

type RecommendationDocument struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Tags    string `json:"tags"`
}

type RecommendationWeights struct {
	Popularity   float64 `json:"beh.weight.popularity"`
	Freshness    float64 `json:"beh.weight.freshness"`
	MoreLikeThis float64 `json:"beh.weight.moreLikeThis"`
	UserProfile  float64 `json:"beh.weight.userProfile"`
	Covisitation float64 `json:"beh.weight.covisitation"`
}

type RecommendationConfig struct {
	*RecommendationWeights
	Doc *RecommendationDocument `json:"doc"`
}

func (c *Client) RawContextualRecommendationsByUUID(uuid string, count int, recency int) (result *Recommendations, err error) {
	u, err := url.Parse("/recommended-reads-api/recommend/contextual")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("contentid", uuid)
	q.Set("count", strconv.Itoa(count))
	q.Set("sort", "rel")
	q.Set("recency", strconv.Itoa(recency))
	u.RawQuery = q.Encode()
	result = &Recommendations{}
	raw, err := c.FromURL(u.String(), result)
	result.RawJSON = raw
	return result, err
}

func (c *Client) RawContextualRecommendationsByDocument(doc *RecommendationDocument, weights *RecommendationWeights, count int, recency int) (result *Recommendations, err error) {
	u, err := url.Parse("/recommended-reads-api/recommend/contextual/doc")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("count", strconv.Itoa(count))
	q.Set("sort", "rel")
	q.Set("recency", strconv.Itoa(recency))
	u.RawQuery = q.Encode()
	config := RecommendationConfig{weights, doc}
	b, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	result = &Recommendations{}
	raw, err := c.FromURLWithBody(u.String(), b, result)
	result.RawJSON = raw
	return result, err
}

func (c *Client) ContextualRecommendationsByUUID(uuid string, count int, recency int) (result []Recommendation, err error) {
	r, err := c.RawContextualRecommendationsByUUID(uuid, count, recency)
	if err != nil {
		return nil, err
	}
	return r.Articles, nil
}

func (c *Client) ContextualRecommendationsByString(s string, count int, recency int) (result []Recommendation, err error) {
	doc := RecommendationDocument{Content: s}
	r, err := c.RawContextualRecommendationsByDocument(&doc, nil, count, recency)
	if err != nil {
		return nil, err
	}
	return r.Articles, nil
}

func (c *Client) RawBehaviouralRecommendationsByUUID(uuid string, userid string, count int, recency int) (result *Recommendations, err error) {
	u, err := url.Parse("/recommended-reads-api/recommend/behavioural")
	q := u.Query()
	q.Set("contentid", uuid)
	q.Set("userid", userid)
	q.Set("count", strconv.Itoa(count))
	q.Set("sort", "rel")
	q.Set("recency", strconv.Itoa(recency))
	u.RawQuery = q.Encode()
	result = &Recommendations{}
	raw, err := c.FromURL(u.String(), result)
	result.RawJSON = raw
	return result, err
}

func (c *Client) BehaviouralRecommendationsByUUID(uuid string, userid string, count int, recency int) (result []Recommendation, err error) {
	r, err := c.RawBehaviouralRecommendationsByUUID(uuid, userid, count, recency)
	if err != nil {
		return nil, err
	}
	return r.Articles, nil
}

func (c *Client) RawBehaviouralRecommendations(userid string, count int, recency int) (result *Recommendations, err error) {
	u, err := url.Parse("/recommended-reads-api/recommend/behavioural")
	q := u.Query()
	q.Set("userid", userid)
	q.Set("count", strconv.Itoa(count))
	q.Set("sort", "rel")
	q.Set("recency", strconv.Itoa(recency))
	u.RawQuery = q.Encode()
	result = &Recommendations{}
	raw, err := c.FromURL(u.String(), result)
	result.RawJSON = raw
	return result, err
}

func (c *Client) BehaviouralRecommendations(userid string, count int, recency int) (result []Recommendation, err error) {
	r, err := c.RawBehaviouralRecommendations(userid, count, recency)
	if err != nil {
		return nil, err
	}
	return r.Articles, nil
}

func (c *Client) RawPopularRecommendations(count int, recency int) (result *Recommendations, err error) {
	u, err := url.Parse("/recommended-reads-api/recommend/popular")
	q := u.Query()
	q.Set("count", strconv.Itoa(count))
	q.Set("recency", strconv.Itoa(recency))
	u.RawQuery = q.Encode()
	result = &Recommendations{}
	raw, err := c.FromURL(u.String(), result)
	result.RawJSON = raw
	return result, err
}

func (c *Client) PopularRecommendations(count int, recency int) (result []Recommendation, err error) {
	r, err := c.RawPopularRecommendations(count, recency)
	if err != nil {
		return nil, err
	}
	return r.Articles, nil
}
