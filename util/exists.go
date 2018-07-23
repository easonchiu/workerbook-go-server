package util

// Exists checks if a string is in a set.
func Exists(set []string, find string) bool {
  for _, s := range set {
    if s == find {
      return true
    }
  }
  return false
}