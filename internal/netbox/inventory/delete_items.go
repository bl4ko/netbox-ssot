package inventory

func (nbi *NetBoxInventory) DeleteOrphans() error {
	// Ensure OrphanObjectPriority and OrphanManager lengths are the same,
	// if not, there are missing entries somewhere and need to be fixed.
	if len(nbi.OrphanManager) != len(nbi.OrphanObjectPriority) {
		panic("len(nbi.OrphanManager) != len(nbi.OrhpanObjectPriority). This should not happen. Every orphan managed object must have its corresponding priority")
	}

	for i := 0; i < len(nbi.OrphanObjectPriority); i++ {
		objectAPIPath := nbi.OrphanObjectPriority[i]
		ids := nbi.OrphanManager[objectAPIPath]
		if len(ids) != 0 {
			nbi.Logger.Info("Deleting orphaned objects of type ", objectAPIPath, " with IDs ", ids)
			err := nbi.NetboxApi.BulkDeleteObjects(objectAPIPath, ids)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
