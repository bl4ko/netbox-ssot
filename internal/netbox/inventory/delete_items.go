package inventory

import "context"

func (nbi *NetboxInventory) DeleteOrphans(ctx context.Context) error {
	// Ensure OrphanObjectPriority and OrphanManager lengths are the same,
	// if not, there are missing entries somewhere and need to be fixed.
	if len(nbi.OrphanManager) != len(nbi.OrphanObjectPriority) {
		panic("len(nbi.OrphanManager) != len(nbi.OrphanObjectPriority). This should not happen. Every orphan managed object must have its corresponding priority")
	}

	for i := 0; i < len(nbi.OrphanObjectPriority); i++ {
		objectAPIPath := nbi.OrphanObjectPriority[i]
		ids := nbi.OrphanManager[objectAPIPath]
		if len(ids) != 0 {
			nbi.Logger.Infof(ctx, "Deleting orphaned objects of type %s", objectAPIPath)
			nbi.Logger.Debugf(ctx, "Ids of objects to be deleted: %v", ids)
			for id := range ids {
				err := nbi.NetboxAPI.DeleteObject(ctx, objectAPIPath, id)
				if err != nil {
					nbi.Logger.Errorf(nbi.Ctx, "delete objects: %s", err)
				}
			}
		}
	}
	return nil
}
