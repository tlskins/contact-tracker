package mongo

type M map[string]interface{}

type PostProcessable interface {
	PostProcess()
}
