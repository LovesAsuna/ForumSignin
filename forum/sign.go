package forum

type Sign interface {
	Name() string

	BasicUrl() string

	Cookie() string

	Sign() (<-chan string, bool)
}

type nocookieclient struct {
	name string
}

func NewNoCookieClient(name string) Sign {
	return &nocookieclient{
		name,
	}
}

func (client *nocookieclient) Name() string {
	return client.name
}

func (client *nocookieclient) BasicUrl() string {
	return ""
}

func (client *nocookieclient) Cookie() string {
	return ""
}

func (client *nocookieclient) Sign() (<-chan string, bool) {
	c := make(chan string, 1)
	c <- client.name + "Cookie未设置！"
	close(c)
	return c, false
}
