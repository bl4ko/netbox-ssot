package inventory

func (nbi *NetBoxInventory) DeleteOrphans() error {
	for objectAPIPath, ids := range nbi.OrphanManager {
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
