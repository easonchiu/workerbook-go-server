package schedule

func Start() {
  // gocron.Every(1).Day().At("01:00").Do(SaveAnalysis)
  // gocron.Every(5).Day().Do(SaveAnalysis)
  // gocron.Start()
  SaveAnalysis()
}
