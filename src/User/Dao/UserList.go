package UserDao

type ByAccessId []UserDTO

func (s ByAccessId) Len() int {
	return len(s)
}
func (s ByAccessId) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByAccessId) Less(i, j int) bool {
	return s[i].AccessId < s[j].AccessId
}
