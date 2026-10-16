package registry_test

import (
	"testing"

	"hidden-coast/internal/registry"
)

func TestSeedAndGet(t *testing.T) {
	reg := registry.GetRegistry()

	all := reg.ListAll()
	if len(all) != 15 {
		t.Fatalf("expected 15 SMEs, got %d", len(all))
	}

	sme, ok := reg.Get("mama-aisha")
	if !ok {
		t.Fatalf("expected to find mama-aisha")
	}

	if sme.Name != "Mama Aisha's Kitchen" {
		t.Errorf("unexpected name: %s", sme.Name)
	}

	prod := sme.FindProduct("chicken-biryani")
	if prod == nil {
		t.Fatalf("expected product chicken-biryani")
	}

	if prod.Price != 450 {
		t.Errorf("expected price 450, got %d", prod.Price)
	}
}

func TestSearchByCategoryAndArea(t *testing.T) {
	reg := registry.GetRegistry()

	foodOldTown := reg.SearchByCategory("food", "Old Town")
	if len(foodOldTown) == 0 {
		t.Errorf("expected food SMEs in Old Town")
	}

	crafts := reg.SearchByCategory("crafts", "")
	if len(crafts) == 0 {
		t.Errorf("expected craft SMEs")
	}
}

func TestUpdateStock(t *testing.T) {
	reg := registry.GetRegistry()

	sme, _ := reg.Get("mama-aisha")
	initialStock := sme.FindProduct("chicken-biryani").Stock

	updated := reg.UpdateStock("mama-aisha", "chicken-biryani", -2)
	if !updated {
		t.Fatalf("failed to update stock")
	}

	smeUpdated, _ := reg.Get("mama-aisha")
	newStock := smeUpdated.FindProduct("chicken-biryani").Stock

	if newStock != initialStock-2 {
		t.Errorf("expected stock %d, got %d", initialStock-2, newStock)
	}
}
