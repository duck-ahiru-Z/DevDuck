package adapter

// Registry は、利用可能なAdapterを管理する。
type Registry struct {
	adapters []Adapter
}

// NewRegistry は、Adapter一覧を受け取ってRegistryを作成する。
func NewRegistry(adapters ...Adapter) *Registry {
	return &Registry{
		adapters: adapters,
	}
}

// Find は、実行されたコマンドに対応するAdapterを探す。
func (r *Registry) Find(
	command string,
	args []string,
) (Adapter, bool) {
	for _, candidate := range r.adapters {
		if candidate.Detect(command, args) {
			return candidate, true
		}
	}

	return nil, false
}
