package util

func Forever(do func(int) (done bool)) {
  i := 1
  for !do(i) {
    i++
  }
}
