package schedule

import "github.com/jasonlvhit/gocron"

func Start() {
  gocron.Every(1).Day().At("01:00").Do(SaveAnalysis)
  gocron.Start()
}
