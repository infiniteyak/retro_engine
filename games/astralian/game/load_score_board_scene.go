package game

import (
	"fmt"
	"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"sort"
    "net/http"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "log"
	"github.com/hajimehoshi/ebiten/v2"
)

type scoreAPIResponse struct {
    Items []scoreEntryAPI `json:"items"`
}

type scoreEntryAPI struct {
    Initials string `json:"initials"`
    Score int `json:"score"`
}

func getScores(pInitials string, pScore int) ([]scoreEntryAPI, bool) {
    if pInitials != "" {
        jsonData, err := json.Marshal(scoreEntryAPI{Initials: pInitials, Score: pScore})
        if err != nil {
            log.Fatal(err)
        }

        _, err = http.Post(
            "https://www.infiniteyak.com/api/collections/astralian_scores/records",
            "application/json",
            bytes.NewBuffer(jsonData),
            )
        if err != nil {
            print(err)
            return nil, true
        }
    }

    response, err := http.Get("https://www.infiniteyak.com/api/collections/astralian_scores/records?sort=-score&perPage=10")
    if err != nil {
        print(err)
        return nil, true
    }
    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        print(err)
        return nil, true
    }

    var responseObject scoreAPIResponse
    if err := json.Unmarshal(responseData, &responseObject); err != nil {
        panic(err)
    }

    return responseObject.Items, false
}

func (this *Game) LoadScoreBoardScene() {
    println("LoadScoreBoardScene")
    this.curScene.SetId(ScoreBoard_sceneId)

    this.GenerateStars(this.screenView)

    titleText := entity.AddTitleText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        18,
        this.screenView,
        "HIGH SCORES",
    )
    titleText.YAlign = entity.Top_fontaligny

    scores, useLocal := getScores(this.curScore.Ident, this.curScore.Value)

    if useLocal {
        this.highScores = append(this.highScores, utility.ScoreEntry{
            Ident: this.curScore.Ident, 
            Value: this.curScore.Value,
        })
        sort.Slice(this.highScores, func(i, j int) bool {
            return this.highScores[i].Value > this.highScores[j].Value
        })
    } else {
        this.highScores = []utility.ScoreEntry{}
        for _, s := range scores {
            this.highScores = append(this.highScores, utility.ScoreEntry{
                Ident: s.Initials, 
                Value: s.Score, 
            })
        }
    }
    scoreCount := len(this.highScores)
    if len(this.highScores) > 10 {
        scoreCount = 10
    }
    blinked := false
    for i := 0; i < scoreCount; i++ {
        match := this.highScores[i].Value == this.curScore.Value && 
                 this.highScores[i].Ident == this.curScore.Ident 
        scoreText := entity.AddNormalText(
            this.ecs, 
            float64(this.screenView.Area.Max.X / 2), 
            float64(40 + i * 15),
            this.screenView,
            "WhiteFont",
            fmt.Sprintf("%02d", i+1) + 
            ". " + 
            this.highScores[i].Ident + 
            " " + 
            fmt.Sprintf("%06d", this.highScores[i].Value),
        )
        scoreText.YAlign = entity.Top_fontaligny
        scoreText.Blink = match && !blinked

        blinked = blinked || match
    }

    this.ResetScore()

    // Advance to the next state when you hit space
    entity.AddInputTrigger(
        this.ecs, 
        ebiten.KeySpace,
        func() {
            this.Transition(Advance_sceneEvent)
        },
    )

    // Start game when you hit Enter
    entity.AddInputTrigger(
        this.ecs, 
        ebiten.KeyEnter,
        func() {
            this.Transition(GameStart_sceneEvent)
        },
    )
}
