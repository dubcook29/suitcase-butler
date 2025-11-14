package grid

type GridDataBind struct {
	// TODO input and output stream binding
	Src   DataBind
	Desc  DataBind
	Cache []interface{}
}

type DataBind struct {
	WMPId string
	Index string
}

func (g *Grid) bindWriter(wmpid, index string, value []interface{}) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	for i, bind := range g.GridDataBind {
		if bind.Src.WMPId == wmpid && bind.Src.Index == index {
			g.GridDataBind[i].Cache = append(bind.Cache, value...)
		}
	}

	return nil
}

func (g *Grid) bindReader(wmpid, index string) []interface{} {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var data []interface{} = nil
	for _, bind := range g.GridDataBind {
		if bind.Desc.WMPId == wmpid && bind.Desc.Index == index {
			data = append(data, bind.Cache...)
		}
	}
	return data
}
