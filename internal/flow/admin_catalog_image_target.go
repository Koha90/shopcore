package flow

// CatalogImageTargetKind identifies which catalog entity expects image upload.
type CatalogImageTargetKind string

const (
	CatalogImageTargetProduct CatalogImageTargetKind = "product"
	CatalogImageTargetVariant CatalogImageTargetKind = "variant"
)

// CatalogImageInputTarget describes current catalog image pending target.
//
// It is used by transports that can turn uploaded media into local image paths.
type CatalogImageInputTarget struct {
	Kind       CatalogImageTargetKind
	EntityID   int
	EntityCode string
}

// CatalogImageInputTarget returns current product/variant image input target.
func (s *Service) CatalogImageInputTarget(key SessionKey) (CatalogImageInputTarget, bool) {
	if s == nil || s.store == nil {
		return CatalogImageInputTarget{}, false
	}

	session, ok := s.store.Get(key)
	if !ok {
		return CatalogImageInputTarget{}, false
	}

	switch session.Pending.Kind {
	case PendingInputProductImageURL:
		productID, ok := pendingProductID(session.Pending)
		if !ok {
			return CatalogImageInputTarget{}, false
		}

		return CatalogImageInputTarget{
			Kind:       CatalogImageTargetProduct,
			EntityID:   productID,
			EntityCode: session.Pending.Value(PendingValueCode),
		}, true

	case PendingInputVariantImageURL:
		variantID, ok := pendingVariantID(session.Pending)
		if !ok {
			return CatalogImageInputTarget{}, false
		}

		return CatalogImageInputTarget{
			Kind:       CatalogImageTargetVariant,
			EntityID:   variantID,
			EntityCode: session.Pending.Value(PendingValueCode),
		}, true

	default:
		return CatalogImageInputTarget{}, false
	}
}
