package registry

import (
	"strings"
	"sync"

	"hidden-coast/internal/models"
)

type Registry struct {
	mu   sync.RWMutex
	smes map[string]models.SME
}

var defaultRegistry *Registry
var once sync.Once

func GetRegistry() *Registry {
	once.Do(func() {
		defaultRegistry = NewRegistry()
		Seed(defaultRegistry)
	})
	return defaultRegistry
}

func NewRegistry() *Registry {
	return &Registry{
		smes: make(map[string]models.SME),
	}
}

func (r *Registry) Add(sme models.SME) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.smes[sme.ID] = sme
}

func (r *Registry) Get(id string) (models.SME, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sme, ok := r.smes[id]
	return sme, ok
}

func (r *Registry) ListAll() []models.SME {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]models.SME, 0, len(r.smes))
	for _, sme := range r.smes {
		list = append(list, sme)
	}
	return list
}

func (r *Registry) SearchByCategory(category, area string) []models.SME {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var results []models.SME
	for _, sme := range r.smes {
		matchCategory := category == "" || strings.EqualFold(sme.Category, category)
		matchArea := area == "" || strings.EqualFold(sme.Location.Area, area) || strings.Contains(strings.ToLower(sme.Location.Area), strings.ToLower(area))
		if matchCategory && matchArea {
			results = append(results, sme)
		}
	}
	return results
}

func (r *Registry) UpdateStock(smeID, productID string, qtyChange int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	sme, ok := r.smes[smeID]
	if !ok {
		return false
	}
	for i := range sme.Products {
		if sme.Products[i].ID == productID {
			sme.Products[i].Stock += qtyChange
			if sme.Products[i].Stock < 0 {
				sme.Products[i].Stock = 0
			}
			r.smes[smeID] = sme
			return true
		}
	}
	return false
}
